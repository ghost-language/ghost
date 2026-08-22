package modules

import (
	"math"

	"ghostlang.org/x/ghost/fault"
	"ghostlang.org/x/ghost/object"
	"ghostlang.org/x/ghost/token"
)

// Ghost has no separate array type, so the array methods read ordinary lists:
// a list of numbers is a vector, and a list of lists of numbers is a matrix.
// That is enough for the work an array library is actually asked to do -
// building ranges, reshaping, and the linear algebra - while leaving every list
// that passes through these methods an ordinary list that the rest of the
// language can index, iterate, and print.

func registerMathArrays() {
	// Building lists of numbers.
	RegisterMethod(MathMethods, "arange", mathArange)
	RegisterMethod(MathMethods, "linspace", mathLinspace)
	RegisterMethod(MathMethods, "zeros", mathZeros)
	RegisterMethod(MathMethods, "ones", mathOnes)
	RegisterMethod(MathMethods, "full", mathFull)
	RegisterMethod(MathMethods, "identity", mathIdentity)

	// Rearranging them.
	RegisterMethod(MathMethods, "reshape", mathReshape)
	RegisterMethod(MathMethods, "flatten", mathFlatten)
	RegisterMethod(MathMethods, "shape", mathShape)
	RegisterMethod(MathMethods, "transpose", mathTranspose)

	// Vectors.
	RegisterMethod(MathMethods, "dot", mathDot)
	RegisterMethod(MathMethods, "matmul", mathDot)
	RegisterMethod(MathMethods, "cross", mathCross)
	RegisterMethod(MathMethods, "outer", mathOuter)
	RegisterMethod(MathMethods, "norm", mathNorm)
	RegisterMethod(MathMethods, "normalize", mathNormalize)
	RegisterMethod(MathMethods, "distance", mathDistance)
	RegisterMethod(MathMethods, "angle", mathAngle)

	// Matrices.
	RegisterMethod(MathMethods, "determinant", mathDeterminant)
	RegisterMethod(MathMethods, "inverse", mathInverse)
	RegisterMethod(MathMethods, "solve", mathSolve)
	RegisterMethod(MathMethods, "trace", mathTrace)
}

// =============================================================================
// Building

// mathArange counts from a start to a stop in steps, stopping before the stop
// value rather than on it. Given whole numbers it answers with whole numbers,
// so it can drive an index directly.
func mathArange(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.arange", tok, args, 1, 3); err != nil {
		return err
	}

	bounds := make([]*object.Number, len(args))

	for index := range args {
		number, err := numberAt("math.arange", tok, args, index)

		if err != nil {
			return err
		}

		bounds[index] = number
	}

	start := 0.0
	stop := bounds[0].Float64()
	step := 1.0
	whole := !bounds[0].IsFloat()

	if len(bounds) > 1 {
		start = bounds[0].Float64()
		stop = bounds[1].Float64()
		whole = whole && !bounds[1].IsFloat()
	}

	if len(bounds) > 2 {
		step = bounds[2].Float64()
		whole = whole && !bounds[2].IsFloat()
	}

	if step == 0 {
		return object.NewError(fault.Value, tok, "`math.arange()` expects a step other than zero")
	}

	count := int(math.Ceil((stop - start) / step))

	if count < 0 {
		count = 0
	}

	elements := make([]object.Object, count)

	for index := 0; index < count; index++ {
		given := start + float64(index)*step

		if whole {
			elements[index] = object.NewInt(int64(given))
		} else {
			elements[index] = object.NewFloat(given)
		}
	}

	return &object.List{Elements: elements}
}

// mathLinspace divides the span between two values into a fixed number of
// points, including both ends. Where arange is told the step and works out the
// count, linspace is told the count and works out the step.
func mathLinspace(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.linspace", tok, args, 3); err != nil {
		return err
	}

	start, err := floatAt("math.linspace", tok, args, 0)

	if err != nil {
		return err
	}

	stop, err := floatAt("math.linspace", tok, args, 1)

	if err != nil {
		return err
	}

	count, err := integerAt("math.linspace", tok, args, 2)

	if err != nil {
		return err
	}

	if count < 0 {
		return object.NewError(fault.Value, tok, "`math.linspace()` expects a count of zero or greater, got %d", count)
	}

	values := make([]float64, count)

	if count == 1 {
		values[0] = start
	}

	if count > 1 {
		step := (stop - start) / float64(count-1)

		for index := int64(0); index < count; index++ {
			values[index] = start + float64(index)*step
		}

		values[count-1] = stop
	}

	return floatList(values)
}

func mathZeros(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return filled("math.zeros", tok, args, object.NewInt(0))
}

func mathOnes(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	return filled("math.ones", tok, args, object.NewInt(1))
}

// mathFull builds a list, or a list of lists, of a repeated value.
func mathFull(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.full", tok, args, 2, 3); err != nil {
		return err
	}

	fill, err := numberAt("math.full", tok, args, len(args)-1)

	if err != nil {
		return err
	}

	return filled("math.full", tok, args[:len(args)-1], fill)
}

// filled is the shape half of zeros, ones, and full: one argument gives a list,
// two give a list of rows.
func filled(name string, tok token.Token, args []object.Object, fill *object.Number) object.Object {
	if err := arityRange(name, tok, args, 1, 2); err != nil {
		return err
	}

	dimensions := make([]int64, len(args))

	for index := range args {
		size, err := integerAt(name, tok, args, index)

		if err != nil {
			return err
		}

		if size < 0 {
			return object.NewError(fault.Value, tok, "`%s()` expects sizes of zero or greater, got %d", name, size)
		}

		dimensions[index] = size
	}

	if len(dimensions) == 1 {
		return repeated(fill, dimensions[0])
	}

	rows := make([]object.Object, dimensions[0])

	for index := range rows {
		rows[index] = repeated(fill, dimensions[1])
	}

	return &object.List{Elements: rows}
}

func repeated(fill *object.Number, count int64) *object.List {
	elements := make([]object.Object, count)

	for index := range elements {
		elements[index] = fill
	}

	return &object.List{Elements: elements}
}

// mathIdentity builds the square matrix that leaves any matrix it multiplies
// unchanged.
func mathIdentity(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.identity", tok, args, 1); err != nil {
		return err
	}

	size, err := integerAt("math.identity", tok, args, 0)

	if err != nil {
		return err
	}

	if size < 0 {
		return object.NewError(fault.Value, tok, "`math.identity()` expects a size of zero or greater, got %d", size)
	}

	rows := make([]object.Object, size)

	for row := int64(0); row < size; row++ {
		columns := make([]object.Object, size)

		for column := int64(0); column < size; column++ {
			if row == column {
				columns[column] = object.NewInt(1)
			} else {
				columns[column] = object.NewInt(0)
			}
		}

		rows[row] = &object.List{Elements: columns}
	}

	return &object.List{Elements: rows}
}

// =============================================================================
// Rearranging

// mathReshape lays the values out in new dimensions. One dimension may be given
// as -1, in which case it is worked out from how many values there are.
func mathReshape(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.reshape", tok, args, 2, 3); err != nil {
		return err
	}

	numbers, err := gatherNumbers("math.reshape", tok, args[:1])

	if err != nil {
		return err
	}

	dimensions := make([]int64, len(args)-1)

	for index := range dimensions {
		size, argumentErr := integerAt("math.reshape", tok, args, index+1)

		if argumentErr != nil {
			return argumentErr
		}

		dimensions[index] = size
	}

	total := int64(len(numbers))
	inferred := -1
	known := int64(1)

	for index, size := range dimensions {
		if size == -1 {
			if inferred >= 0 {
				return object.NewError(fault.Value, tok, "`math.reshape()` can infer only one dimension")
			}

			inferred = index

			continue
		}

		if size < 0 {
			return object.NewError(fault.Value, tok, "`math.reshape()` expects sizes of zero or greater, or -1 to infer, got %d", size)
		}

		known *= size
	}

	if inferred >= 0 {
		if known == 0 || total%known != 0 {
			return object.NewError(fault.Value, tok, "`math.reshape()` cannot infer a dimension for %d values", total)
		}

		dimensions[inferred] = total / known
		known = total
	}

	if known != total {
		return object.NewError(fault.Value, tok, "`math.reshape()` cannot fit %d values into the requested shape", total)
	}

	if len(dimensions) == 1 {
		return numberList(numbers)
	}

	rows := make([]object.Object, dimensions[0])

	for row := int64(0); row < dimensions[0]; row++ {
		rows[row] = numberList(numbers[row*dimensions[1] : (row+1)*dimensions[1]])
	}

	return &object.List{Elements: rows}
}

// mathFlatten collapses any nesting into a single list of numbers.
func mathFlatten(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityAtLeast("math.flatten", tok, args, 1); err != nil {
		return err
	}

	numbers, err := gatherNumbers("math.flatten", tok, args)

	if err != nil {
		return err
	}

	return numberList(numbers)
}

// mathShape reports the dimensions of a nested list, outermost first. It stops
// at the first level that is ragged, which is the level at which the values
// stop describing a rectangular array.
func mathShape(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.shape", tok, args, 1); err != nil {
		return err
	}

	if _, err := listAt("math.shape", tok, args, 0); err != nil {
		return err
	}

	return integerList(dimensionsOf(args[0]))
}

func dimensionsOf(obj object.Object) []int64 {
	list, ok := obj.(*object.List)

	if !ok {
		return nil
	}

	dimensions := []int64{int64(len(list.Elements))}

	if len(list.Elements) == 0 {
		return dimensions
	}

	inner := dimensionsOf(list.Elements[0])

	if inner == nil {
		return dimensions
	}

	for _, element := range list.Elements[1:] {
		if !sameDimensions(dimensionsOf(element), inner) {
			return dimensions
		}
	}

	return append(dimensions, inner...)
}

func sameDimensions(left []int64, right []int64) bool {
	if len(left) != len(right) {
		return false
	}

	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}

	return true
}

// mathTranspose swaps the rows and columns of a matrix. A plain vector has no
// rows to swap, and comes back unchanged.
func mathTranspose(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.transpose", tok, args, 1); err != nil {
		return err
	}

	rows, isVector, err := toMatrix("math.transpose", tok, args, 0)

	if err != nil {
		return err
	}

	if isVector {
		return floatList(rows[0])
	}

	return matrixList(transposeOf(rows))
}

// =============================================================================
// Vectors

// mathDot multiplies two vectors, a matrix and a vector, or two matrices. Which
// one it is doing follows from the shapes it is given: two vectors produce a
// number, a matrix and a vector produce a vector, and two matrices produce a
// matrix.
func mathDot(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.dot", tok, args, 2); err != nil {
		return err
	}

	left, leftIsVector, err := toMatrix("math.dot", tok, args, 0)

	if err != nil {
		return err
	}

	right, rightIsVector, err := toMatrix("math.dot", tok, args, 1)

	if err != nil {
		return err
	}

	// A vector on the right stands for a column, one on the left for a row.
	if rightIsVector {
		right = transposeOf(right)
	}

	if len(left[0]) != len(right) {
		return object.NewError(fault.Value, tok, "`math.dot()` expects the width of the first operand to match the height of the second, got %d and %d", len(left[0]), len(right))
	}

	product := multiplyMatrices(left, right)

	switch {
	case leftIsVector && rightIsVector:
		return object.NewFloat(product[0][0])
	case leftIsVector:
		return floatList(product[0])
	case rightIsVector:
		return floatList(transposeOf(product)[0])
	}

	return matrixList(product)
}

// mathCross returns the cross product. Three-dimensional vectors give the
// vector perpendicular to both; two-dimensional vectors give the single number
// that is all their cross product can hold.
func mathCross(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.cross", tok, args, 2); err != nil {
		return err
	}

	left, err := toVector("math.cross", tok, args, 0)

	if err != nil {
		return err
	}

	right, err := toVector("math.cross", tok, args, 1)

	if err != nil {
		return err
	}

	if len(left) != len(right) || (len(left) != 2 && len(left) != 3) {
		return object.NewError(fault.Value, tok, "`math.cross()` expects two vectors of the same length, either both of 2 or both of 3")
	}

	if len(left) == 2 {
		return object.NewFloat(left[0]*right[1] - left[1]*right[0])
	}

	return floatList([]float64{
		left[1]*right[2] - left[2]*right[1],
		left[2]*right[0] - left[0]*right[2],
		left[0]*right[1] - left[1]*right[0],
	})
}

// mathOuter multiplies every element of one vector by every element of another,
// giving a matrix as tall as the first and as wide as the second. It is the
// shape a weight gradient takes, and the reason it is a method rather than a
// composition of transpose and dot is that spelling it out that way obscures
// what is being asked for.
func mathOuter(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.outer", tok, args, 2); err != nil {
		return err
	}

	left, err := toVector("math.outer", tok, args, 0)

	if err != nil {
		return err
	}

	right, err := toVector("math.outer", tok, args, 1)

	if err != nil {
		return err
	}

	rows := make([][]float64, len(left))

	for row, scale := range left {
		rows[row] = make([]float64, len(right))

		for column, value := range right {
			rows[row][column] = scale * value
		}
	}

	return matrixList(rows)
}

// mathNorm returns the length of a vector. The default is the ordinary
// straight-line length; a second argument asks for another p-norm, where 1 sums
// the absolute values and math.infinity takes the largest of them.
func mathNorm(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityRange("math.norm", tok, args, 1, 2); err != nil {
		return err
	}

	values, err := gatherFloats("math.norm", tok, args[:1])

	if err != nil {
		return err
	}

	order := 2.0

	if len(args) == 2 {
		given, argumentErr := floatAt("math.norm", tok, args, 1)

		if argumentErr != nil {
			return argumentErr
		}

		order = given
	}

	if order <= 0 {
		return object.NewError(fault.Value, tok, "`math.norm()` expects an order greater than zero")
	}

	if math.IsInf(order, 1) {
		largest := 0.0

		for _, given := range values {
			largest = math.Max(largest, math.Abs(given))
		}

		return object.NewFloat(largest)
	}

	total := 0.0

	for _, given := range values {
		total += math.Pow(math.Abs(given), order)
	}

	return object.NewFloat(math.Pow(total, 1/order))
}

// mathNormalize scales a vector to a length of one, keeping its direction.
func mathNormalize(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.normalize", tok, args, 1); err != nil {
		return err
	}

	values, err := toVector("math.normalize", tok, args, 0)

	if err != nil {
		return err
	}

	length := 0.0

	for _, given := range values {
		length += given * given
	}

	length = math.Sqrt(length)

	if length == 0 {
		return object.NewError(fault.Value, tok, "`math.normalize()` cannot normalize a vector of length zero")
	}

	scaled := make([]float64, len(values))

	for index, given := range values {
		scaled[index] = given / length
	}

	return floatList(scaled)
}

// mathDistance returns the straight-line distance between two points. The
// points may be written out as coordinates or passed as vectors, in as many
// dimensions as they have: math.distance(0, 0, 3, 4) and
// math.distance([0, 0], [3, 4]) are the same call.
func mathDistance(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityAtLeast("math.distance", tok, args, 2); err != nil {
		return err
	}

	values, err := gatherFloats("math.distance", tok, args)

	if err != nil {
		return err
	}

	if len(values) < 2 || len(values)%2 != 0 {
		return object.NewError(fault.Value, tok, "`math.distance()` expects two points of the same number of coordinates, got %d coordinates", len(values))
	}

	half := len(values) / 2
	total := 0.0

	for index := 0; index < half; index++ {
		difference := values[half+index] - values[index]
		total += difference * difference
	}

	return object.NewFloat(math.Sqrt(total))
}

// mathAngle returns the angle in radians from one point to another, measured
// from the positive x axis and covering every quadrant.
func mathAngle(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arityAtLeast("math.angle", tok, args, 2); err != nil {
		return err
	}

	values, err := gatherFloats("math.angle", tok, args)

	if err != nil {
		return err
	}

	if len(values) != 4 {
		return object.NewError(fault.Value, tok, "`math.angle()` expects two points of two coordinates, got %d coordinates", len(values))
	}

	return object.NewFloat(math.Atan2(values[3]-values[1], values[2]-values[0]))
}

// =============================================================================
// Matrices

func mathTrace(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rows, err := squareMatrix("math.trace", tok, args)

	if err != nil {
		return err
	}

	total := 0.0

	for index := range rows {
		total += rows[index][index]
	}

	return object.NewFloat(total)
}

// mathDeterminant returns the factor by which a matrix scales area or volume. A
// determinant of zero means the matrix collapses its input, and so cannot be
// inverted.
func mathDeterminant(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rows, err := squareMatrix("math.determinant", tok, args)

	if err != nil {
		return err
	}

	_, determinant, singular := eliminate(copyMatrix(rows), nil)

	if singular {
		return object.NewFloat(0)
	}

	return object.NewFloat(determinant)
}

// mathInverse returns the matrix that undoes this one.
func mathInverse(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	rows, err := squareMatrix("math.inverse", tok, args)

	if err != nil {
		return err
	}

	solution, _, singular := eliminate(copyMatrix(rows), identityMatrix(len(rows)))

	if singular {
		return object.NewError(fault.Value, tok, "`math.inverse()` cannot invert a singular matrix")
	}

	return matrixList(solution)
}

// mathSolve finds the x that satisfies a·x = b, without forming the inverse.
// The right-hand side may be a vector or a matrix of several right-hand sides.
func mathSolve(scope *object.Scope, tok token.Token, args ...object.Object) object.Object {
	if err := arity("math.solve", tok, args, 2); err != nil {
		return err
	}

	rows, err := squareMatrix("math.solve", tok, args[:1])

	if err != nil {
		return err
	}

	right, rightIsVector, err := toMatrix("math.solve", tok, args, 1)

	if err != nil {
		return err
	}

	if rightIsVector {
		right = transposeOf(right)
	}

	if len(right) != len(rows) {
		return object.NewError(fault.Value, tok, "`math.solve()` expects the right-hand side to have %d rows, got %d", len(rows), len(right))
	}

	solution, _, singular := eliminate(copyMatrix(rows), copyMatrix(right))

	if singular {
		return object.NewError(fault.Value, tok, "`math.solve()` cannot solve a singular system")
	}

	if rightIsVector {
		return floatList(transposeOf(solution)[0])
	}

	return matrixList(solution)
}

// eliminate runs Gauss-Jordan elimination with partial pivoting, reducing left
// to the identity while applying the same operations to right. It returns the
// transformed right-hand side, the determinant of left, and whether the matrix
// turned out to be singular. Passing a nil right-hand side asks for the
// determinant alone.
func eliminate(left [][]float64, right [][]float64) ([][]float64, float64, bool) {
	size := len(left)
	determinant := 1.0

	for column := 0; column < size; column++ {
		pivot := column

		for row := column + 1; row < size; row++ {
			if math.Abs(left[row][column]) > math.Abs(left[pivot][column]) {
				pivot = row
			}
		}

		if left[pivot][column] == 0 {
			return nil, 0, true
		}

		if pivot != column {
			left[pivot], left[column] = left[column], left[pivot]
			determinant = -determinant

			if right != nil {
				right[pivot], right[column] = right[column], right[pivot]
			}
		}

		scale := left[column][column]
		determinant *= scale

		for index := range left[column] {
			left[column][index] /= scale
		}

		if right != nil {
			for index := range right[column] {
				right[column][index] /= scale
			}
		}

		for row := 0; row < size; row++ {
			if row == column || left[row][column] == 0 {
				continue
			}

			factor := left[row][column]

			for index := range left[row] {
				left[row][index] -= factor * left[column][index]
			}

			if right != nil {
				for index := range right[row] {
					right[row][index] -= factor * right[column][index]
				}
			}
		}
	}

	return right, determinant, false
}

// =============================================================================
// Shared helpers

// toVector reads an argument as a flat list of numbers.
func toVector(name string, tok token.Token, args []object.Object, index int) ([]float64, *object.Error) {
	list, err := listAt(name, tok, args, index)

	if err != nil {
		return nil, err
	}

	values := make([]float64, len(list.Elements))

	for position, element := range list.Elements {
		number, ok := element.(*object.Number)

		if !ok {
			return nil, object.NewError(fault.Argument, tok, "`%s()` expects argument %d to be a list of numbers, got %s at index %d", name, index+1, object.TypeName(element), position)
		}

		values[position] = number.Float64()
	}

	return values, nil
}

// toMatrix reads an argument as a rectangular list of lists. A flat list of
// numbers is read as a single row, and the second result says so, which is how
// the vector and matrix cases of dot and solve tell themselves apart.
func toMatrix(name string, tok token.Token, args []object.Object, index int) ([][]float64, bool, *object.Error) {
	list, err := listAt(name, tok, args, index)

	if err != nil {
		return nil, false, err
	}

	if len(list.Elements) == 0 {
		return nil, false, object.NewError(fault.Argument, tok, "`%s()` expects argument %d to be a non-empty list", name, index+1)
	}

	if _, ok := list.Elements[0].(*object.Number); ok {
		values, err := toVector(name, tok, args, index)

		if err != nil {
			return nil, false, err
		}

		return [][]float64{values}, true, nil
	}

	rows := make([][]float64, len(list.Elements))
	width := -1

	for position, element := range list.Elements {
		row, ok := element.(*object.List)

		if !ok {
			return nil, false, object.NewError(fault.Argument, tok, "`%s()` expects argument %d to be a list of numbers or a list of rows, got %s at index %d", name, index+1, object.TypeName(element), position)
		}

		values, err := toVector(name, tok, []object.Object{row}, 0)

		if err != nil {
			return nil, false, err
		}

		if width >= 0 && len(values) != width {
			return nil, false, object.NewError(fault.Argument, tok, "`%s()` expects every row of argument %d to be the same length, got %d and %d", name, index+1, width, len(values))
		}

		width = len(values)
		rows[position] = values
	}

	return rows, false, nil
}

// squareMatrix reads the single matrix argument shared by the matrix methods.
func squareMatrix(name string, tok token.Token, args []object.Object) ([][]float64, *object.Error) {
	if err := arity(name, tok, args, 1); err != nil {
		return nil, err
	}

	rows, isVector, err := toMatrix(name, tok, args, 0)

	if err != nil {
		return nil, err
	}

	if isVector || len(rows) != len(rows[0]) {
		return nil, object.NewError(fault.Argument, tok, "`%s()` expects a square matrix", name)
	}

	return rows, nil
}

func multiplyMatrices(left [][]float64, right [][]float64) [][]float64 {
	product := make([][]float64, len(left))

	for row := range left {
		product[row] = make([]float64, len(right[0]))

		for column := range right[0] {
			total := 0.0

			for index := range right {
				total += left[row][index] * right[index][column]
			}

			product[row][column] = total
		}
	}

	return product
}

func transposeOf(rows [][]float64) [][]float64 {
	transposed := make([][]float64, len(rows[0]))

	for column := range rows[0] {
		transposed[column] = make([]float64, len(rows))

		for row := range rows {
			transposed[column][row] = rows[row][column]
		}
	}

	return transposed
}

func copyMatrix(rows [][]float64) [][]float64 {
	copied := make([][]float64, len(rows))

	for index, row := range rows {
		copied[index] = make([]float64, len(row))
		copy(copied[index], row)
	}

	return copied
}

func identityMatrix(size int) [][]float64 {
	rows := make([][]float64, size)

	for row := range rows {
		rows[row] = make([]float64, size)
		rows[row][row] = 1
	}

	return rows
}

func matrixList(rows [][]float64) *object.List {
	elements := make([]object.Object, len(rows))

	for index, row := range rows {
		elements[index] = floatList(row)
	}

	return &object.List{Elements: elements}
}
