package modules

import (
	"math"
	"sort"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// The reductions read a whole collection and answer with one value. Every one
// of them takes its input through gatherNumbers, so the same method serves all
// three ways a program is likely to hold its numbers:
//
//	math.mean(1, 2, 3)
//	math.mean([1, 2, 3])
//	math.mean([[1, 2], [3, 4]])
//
// The last case flattens, which is what numpy does when no axis is named.

func registerMathStatistics() {
	// Totals and averages.
	registerReduction("sum", reduceSum)
	registerReduction("product", reduceProduct)
	registerReduction("mean", reduceMean)
	registerReduction("median", reduceMedian)
	registerReduction("mode", reduceMode)

	// Spread. The plain forms divide by the number of values, describing the
	// values in hand; the sample forms divide by one less, estimating the
	// population those values were drawn from.
	registerReduction("variance", reduceVariance)
	registerReduction("sampleVariance", reduceSampleVariance)
	registerReduction("standardDeviation", reduceStandardDeviation)
	registerReduction("sampleStandardDeviation", reduceSampleStandardDeviation)

	// Extremes.
	registerReduction("min", reduceMin)
	registerReduction("max", reduceMax)
	registerReduction("argmin", reduceArgmin)
	registerReduction("argmax", reduceArgmax)

	// Running totals, which answer with a list as long as their input.
	registerReduction("cumulativeSum", reduceCumulativeSum)
	registerReduction("cumulativeProduct", reduceCumulativeProduct)

	RegisterMethod(MathMethods, "percentile", mathPercentile)
	RegisterMethod(MathMethods, "quantile", mathQuantile)
	RegisterMethod(MathMethods, "sort", mathSort)
	RegisterMethod(MathMethods, "unique", mathUnique)

	// Whole-number mathematics.
	registerReduction("gcd", reduceGcd)
	registerReduction("lcm", reduceLcm)
	RegisterMethod(MathMethods, "factorial", mathFactorial)
	RegisterMethod(MathMethods, "isPrime", mathIsPrime)
	RegisterMethod(MathMethods, "combinations", mathCombinations)
	RegisterMethod(MathMethods, "permutations", mathPermutations)
}

// reduction is an operation over a flattened run of numbers.
type reduction func(name string, tok token.Token, numbers []*object.Number) object.Object

// registerReduction registers a method that flattens everything it is given
// into one sequence before reducing it.
func registerReduction(name string, operation reduction) {
	qualified := "math." + name

	RegisterMethod(MathMethods, name, func(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
		if err := arityAtLeast(qualified, tok, args, 1); err != nil {
			return err
		}

		numbers, err := gatherNumbers(qualified, tok, args)

		if err != nil {
			return err
		}

		if len(numbers) == 0 {
			return emptyInput(qualified, tok)
		}

		return operation(qualified, tok, numbers)
	})
}

// =============================================================================
// Totals and averages

func reduceSum(name string, tok token.Token, numbers []*object.Number) object.Object {
	return sumOf(numbers)
}

func reduceProduct(name string, tok token.Token, numbers []*object.Number) object.Object {
	total := object.NewInt(1)

	for _, number := range numbers {
		total = total.Mul(number)
	}

	return total
}

func reduceMean(name string, tok token.Token, numbers []*object.Number) object.Object {
	return object.NewFloat(meanOf(toFloats(numbers)))
}

func reduceMedian(name string, tok token.Token, numbers []*object.Number) object.Object {
	values := sortedFloats(numbers)
	middle := len(values) / 2

	if len(values)%2 == 1 {
		return object.NewFloat(values[middle])
	}

	return object.NewFloat((values[middle-1] + values[middle]) / 2)
}

// reduceMode answers with the value that occurs most often, and with the
// smallest of them when several tie, so that the result does not depend on the
// order the values arrived in.
func reduceMode(name string, tok token.Token, numbers []*object.Number) object.Object {
	counts := make(map[float64]int, len(numbers))

	for _, number := range numbers {
		counts[number.Float64()]++
	}

	best := numbers[0]
	bestCount := 0

	for _, number := range numbers {
		count := counts[number.Float64()]

		if count > bestCount || (count == bestCount && number.LessThan(best)) {
			best = number
			bestCount = count
		}
	}

	return best
}

// =============================================================================
// Spread

func reduceVariance(name string, tok token.Token, numbers []*object.Number) object.Object {
	return object.NewFloat(varianceOf(toFloats(numbers), 0))
}

func reduceSampleVariance(name string, tok token.Token, numbers []*object.Number) object.Object {
	if len(numbers) < 2 {
		return object.NewError(fault.Value, tok, "`%s()` expects at least two values", name)
	}

	return object.NewFloat(varianceOf(toFloats(numbers), 1))
}

func reduceStandardDeviation(name string, tok token.Token, numbers []*object.Number) object.Object {
	return object.NewFloat(math.Sqrt(varianceOf(toFloats(numbers), 0)))
}

func reduceSampleStandardDeviation(name string, tok token.Token, numbers []*object.Number) object.Object {
	if len(numbers) < 2 {
		return object.NewError(fault.Value, tok, "`%s()` expects at least two values", name)
	}

	return object.NewFloat(math.Sqrt(varianceOf(toFloats(numbers), 1)))
}

// =============================================================================
// Extremes

func reduceMin(name string, tok token.Token, numbers []*object.Number) object.Object {
	return numbers[indexOfExtreme(numbers, false)]
}

func reduceMax(name string, tok token.Token, numbers []*object.Number) object.Object {
	return numbers[indexOfExtreme(numbers, true)]
}

func reduceArgmin(name string, tok token.Token, numbers []*object.Number) object.Object {
	return object.NewInt(int64(indexOfExtreme(numbers, false)))
}

func reduceArgmax(name string, tok token.Token, numbers []*object.Number) object.Object {
	return object.NewInt(int64(indexOfExtreme(numbers, true)))
}

// =============================================================================
// Running totals

func reduceCumulativeSum(name string, tok token.Token, numbers []*object.Number) object.Object {
	running := make([]*object.Number, len(numbers))
	total := object.NewInt(0)

	for index, number := range numbers {
		total = total.Add(number)
		running[index] = total
	}

	return numberList(running)
}

func reduceCumulativeProduct(name string, tok token.Token, numbers []*object.Number) object.Object {
	running := make([]*object.Number, len(numbers))
	total := object.NewInt(1)

	for index, number := range numbers {
		total = total.Mul(number)
		running[index] = total
	}

	return numberList(running)
}

// =============================================================================
// Order statistics

// mathPercentile returns the value below which the given percentage of the
// input falls, interpolating between neighbours when the percentile lands
// between two values.
func mathPercentile(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return quantileOf("math.percentile", tok, args, 100)
}

// mathQuantile is math.percentile on a 0-1 scale.
func mathQuantile(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return quantileOf("math.quantile", tok, args, 1)
}

func quantileOf(name string, tok token.Token, args []object.Object, scale float64) object.Object {
	if err := arity(name, tok, args, 2); err != nil {
		return err
	}

	numbers, err := gatherNumbers(name, tok, args[:1])

	if err != nil {
		return err
	}

	if len(numbers) == 0 {
		return emptyInput(name, tok)
	}

	position, err := floatAt(name, tok, args, 1)

	if err != nil {
		return err
	}

	if position < 0 || position > scale {
		return object.NewError(fault.Value, tok, "`%s()` expects a position between 0 and %s, got %s", name, object.NewFloat(scale).String(), object.NewFloat(position).String())
	}

	values := sortedFloats(numbers)

	if len(values) == 1 {
		return object.NewFloat(values[0])
	}

	offset := (position / scale) * float64(len(values)-1)
	lower := int(math.Floor(offset))
	upper := int(math.Ceil(offset))

	if lower == upper {
		return object.NewFloat(values[lower])
	}

	return object.NewFloat(values[lower] + (values[upper]-values[lower])*(offset-float64(lower)))
}

// mathSort returns the values in order, ascending unless asked otherwise.
func mathSort(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.sort", tok, args, 1, 2); err != nil {
		return err
	}

	numbers, err := gatherNumbers("math.sort", tok, args[:1])

	if err != nil {
		return err
	}

	descending := false

	if len(args) == 2 {
		given, err := booleanAt("math.sort", tok, args, 1)

		if err != nil {
			return err
		}

		descending = given
	}

	sorted := make([]*object.Number, len(numbers))
	copy(sorted, numbers)

	sort.SliceStable(sorted, func(left int, right int) bool {
		if descending {
			return sorted[right].LessThan(sorted[left])
		}

		return sorted[left].LessThan(sorted[right])
	})

	return numberList(sorted)
}

// mathUnique drops repeats, keeping the order the values were first seen in.
func mathUnique(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityAtLeast("math.unique", tok, args, 1); err != nil {
		return err
	}

	numbers, err := gatherNumbers("math.unique", tok, args)

	if err != nil {
		return err
	}

	seen := make(map[float64]bool, len(numbers))
	unique := make([]*object.Number, 0, len(numbers))

	for _, number := range numbers {
		if seen[number.Float64()] {
			continue
		}

		seen[number.Float64()] = true
		unique = append(unique, number)
	}

	return numberList(unique)
}

// =============================================================================
// Whole-number mathematics

func reduceGcd(name string, tok token.Token, numbers []*object.Number) object.Object {
	divisor := absInteger(numbers[0].Int64())

	for _, number := range numbers[1:] {
		divisor = greatestCommonDivisor(divisor, absInteger(number.Int64()))
	}

	return object.NewInt(divisor)
}

func reduceLcm(name string, tok token.Token, numbers []*object.Number) object.Object {
	multiple := absInteger(numbers[0].Int64())

	for _, number := range numbers[1:] {
		next := absInteger(number.Int64())

		if multiple == 0 || next == 0 {
			multiple = 0

			continue
		}

		divisor := greatestCommonDivisor(multiple, next)
		product, ok := multiplyChecked(multiple/divisor, next)

		if !ok {
			return object.NewError(fault.Value, tok, "`%s()` overflowed; the result is too large to represent", name)
		}

		multiple = product
	}

	return object.NewInt(multiple)
}

// mathFactorial returns n!, exactly while the answer fits in a whole number and
// as a float beyond that.
func mathFactorial(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.factorial", tok, args, 1); err != nil {
		return err
	}

	given, err := integerAt("math.factorial", tok, args, 0)

	if err != nil {
		return err
	}

	if given < 0 {
		return object.NewError(fault.Value, tok, "`math.factorial()` expects a value of zero or greater, got %d", given)
	}

	result := int64(1)

	for step := int64(2); step <= given; step++ {
		product, ok := multiplyChecked(result, step)

		if !ok {
			return object.NewFloat(math.Gamma(float64(given) + 1))
		}

		result = product
	}

	return object.NewInt(result)
}

func mathIsPrime(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.isPrime", tok, args, 1); err != nil {
		return err
	}

	number, err := numberAt("math.isPrime", tok, args, 0)

	if err != nil {
		return err
	}

	if !isInteger(number) {
		return toBoolean(false)
	}

	return toBoolean(isPrime(number.Int64()))
}

// mathCombinations returns how many ways k things can be chosen from n when
// order does not matter.
func mathCombinations(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	total, chosen, err := choiceArguments("math.combinations", tok, args)

	if err != nil {
		return err
	}

	if chosen > total-chosen {
		chosen = total - chosen
	}

	result := int64(1)

	for step := int64(0); step < chosen; step++ {
		product, ok := multiplyChecked(result, total-step)

		if !ok {
			return object.NewFloat(math.Round(math.Exp(logGamma(float64(total)+1) - logGamma(float64(chosen)+1) - logGamma(float64(total-chosen)+1))))
		}

		result = product / (step + 1)
	}

	return object.NewInt(result)
}

// mathPermutations returns how many ways k things can be chosen from n when
// order matters.
func mathPermutations(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	total, chosen, err := choiceArguments("math.permutations", tok, args)

	if err != nil {
		return err
	}

	result := int64(1)

	for step := int64(0); step < chosen; step++ {
		product, ok := multiplyChecked(result, total-step)

		if !ok {
			return object.NewFloat(math.Round(math.Exp(logGamma(float64(total)+1) - logGamma(float64(total-chosen)+1))))
		}

		result = product
	}

	return object.NewInt(result)
}

func choiceArguments(name string, tok token.Token, args []object.Object) (int64, int64, *object.Error) {
	if err := arity(name, tok, args, 2); err != nil {
		return 0, 0, err
	}

	total, err := integerAt(name, tok, args, 0)

	if err != nil {
		return 0, 0, err
	}

	chosen, err := integerAt(name, tok, args, 1)

	if err != nil {
		return 0, 0, err
	}

	if total < 0 || chosen < 0 {
		return 0, 0, object.NewError(fault.Value, tok, "`%s()` expects values of zero or greater", name)
	}

	if chosen > total {
		return 0, 0, object.NewError(fault.Value, tok, "`%s()` cannot choose %d from %d", name, chosen, total)
	}

	return total, chosen, nil
}

// =============================================================================
// Shared helpers

func sumOf(numbers []*object.Number) *object.Number {
	total := object.NewInt(0)

	for _, number := range numbers {
		total = total.Add(number)
	}

	return total
}

func meanOf(values []float64) float64 {
	total := 0.0

	for _, given := range values {
		total += given
	}

	return total / float64(len(values))
}

// varianceOf divides by len(values)-correction, which is 0 for the population
// variance and 1 for the sample variance.
func varianceOf(values []float64, correction int) float64 {
	average := meanOf(values)
	total := 0.0

	for _, given := range values {
		total += (given - average) * (given - average)
	}

	return total / float64(len(values)-correction)
}

func sortedFloats(numbers []*object.Number) []float64 {
	values := toFloats(numbers)

	sort.Float64s(values)

	return values
}

// indexOfExtreme finds the first largest or smallest value, so that min, max,
// argmin, and argmax all agree on which one they mean when values tie.
func indexOfExtreme(numbers []*object.Number, largest bool) int {
	best := 0

	for index, number := range numbers[1:] {
		if largest && number.GreaterThan(numbers[best]) {
			best = index + 1
		}

		if !largest && number.LessThan(numbers[best]) {
			best = index + 1
		}
	}

	return best
}

func greatestCommonDivisor(left int64, right int64) int64 {
	for right != 0 {
		left, right = right, left%right
	}

	return left
}

func absInteger(given int64) int64 {
	if given < 0 {
		return -given
	}

	return given
}

func isPrime(given int64) bool {
	if given < 2 {
		return false
	}

	if given%2 == 0 {
		return given == 2
	}

	for divisor := int64(3); divisor*divisor <= given; divisor += 2 {
		if given%divisor == 0 {
			return false
		}
	}

	return true
}

func emptyInput(name string, tok token.Token) *object.Error {
	return object.NewError(fault.Value, tok, "`%s()` expects at least one value", name)
}
