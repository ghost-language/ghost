package decimal

// This is a modified and minimized fork of shopspring/decimal.
// (https://github.com/shopspring/decimal)
// Release under the MIT license.

import (
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Decimal represents a fixed-point decimal. It is immutable.
// number = value * 10 ^ exp
type Decimal struct {
	value *big.Int
	exp   int32
}

// DivisionPrecision is the number of decimal places in the result when it
// doesn't divide exactly.
var DivisionPrecision = 16

// Zero constant, to make computations faster.Zero.rescale()
var Zero = New(0, 1)

var tenInt = big.NewInt(10)

// New returns a new fixed-point decimal, value * 10 ^ exp.
func New(value int64, exp int32) Decimal {
	return Decimal{
		value: big.NewInt(value),
		exp:   exp,
	}
}

// NewFromString returns a new Decimal from a string representation.
// Trailing zeros are not trimmed.
func NewFromString(value string) (Decimal, error) {
	originalInput := value
	var intString string
	var exp int64

	eIndex := strings.IndexAny(value, "Ee")

	if eIndex != -1 {
		expInt, err := strconv.ParseInt(value[eIndex+1:], 10, 32)

		if err != nil {
			if e, ok := err.(*strconv.NumError); ok && e.Err == strconv.ErrRange {
				return Decimal{}, fmt.Errorf("can't convert %s to decimal: fractional part too long", value)
			}

			return Decimal{}, fmt.Errorf("can't convert %s to decimal: exponent is not numeric", value)
		}

		value = value[:eIndex]
		exp = expInt
	}

	parts := strings.Split(value, ".")

	if len(parts) == 1 {
		// There is not decimal point, we can just parse the original string as
		// an int
		intString = value
	} else if len(parts) == 2 {
		intString = parts[0] + parts[1]
		expInt := -len(parts[1])
		exp += int64(expInt)
	} else {
		return Decimal{}, fmt.Errorf("can't convert %s to decimal: too many .s", value)
	}

	dValue := new(big.Int)

	_, ok := dValue.SetString(intString, 10)

	if !ok {
		return Decimal{}, fmt.Errorf("can't convert %s to decimal", value)
	}

	if exp < math.MinInt32 || exp > math.MaxInt32 {
		return Decimal{}, fmt.Errorf("can't convert %s to decimal: fractional part too long", originalInput)
	}

	return Decimal{
		value: dValue,
		exp:   int32(exp),
	}, nil
}

// NewFromInt converts an int64 to Decimal.
func NewFromInt(value int64) Decimal {
	return Decimal{
		value: big.NewInt(value),
		exp:   0,
	}
}

// String returns the string representation of the decimal
// with the fixed point.
func (d Decimal) String() string {
	return d.string(true)
}

// IntPart returns the integer component of the decimal.
func (d Decimal) IntPart() int64 {
	scaleD := d.rescale(0)

	return scaleD.value.Int64()
}

// Neg returns -d.
func (d Decimal) Neg() Decimal {
	d.ensureInitialized()

	val := new(big.Int).Neg(d.value)

	return Decimal{
		value: val,
		exp:   d.exp,
	}
}

// Add returns d + d2.
func (d Decimal) Add(d2 Decimal) Decimal {
	rd, rd2 := RescalePair(d, d2)

	d3Value := new(big.Int).Add(rd.value, rd2.value)

	return Decimal{
		value: d3Value,
		exp:   rd.exp,
	}
}

// Sub returns d - d2.
func (d Decimal) Sub(d2 Decimal) Decimal {
	rd, rd2 := RescalePair(d, d2)

	d3Value := new(big.Int).Sub(rd.value, rd2.value)

	return Decimal{
		value: d3Value,
		exp:   rd.exp,
	}
}

// Mul returns d * d2.
func (d Decimal) Mul(d2 Decimal) Decimal {
	d.ensureInitialized()
	d2.ensureInitialized()

	expInt64 := int64(d.exp) + int64(d2.exp)

	if expInt64 > math.MaxInt32 || expInt64 < math.MinInt32 {
		panic(fmt.Sprintf("exponent %v overflows an int32!", expInt64))
	}

	d3Value := new(big.Int).Mul(d.value, d2.value)

	return Decimal{
		value: d3Value,
		exp:   int32(expInt64),
	}
}

// Div returns d / d2.
func (d Decimal) Div(d2 Decimal) Decimal {
	return d.DivRound(d2, int32(DivisionPrecision))
}

// Mod returns d % d2.
func (d Decimal) Mod(d2 Decimal) Decimal {
	quo := d.Div(d2).Truncate(0)

	return d.Sub(d2.Mul(quo))
}

// Abs returns the absolute value of the decimal.
func (d Decimal) Abs() Decimal {
	d.ensureInitialized()
	d2Value := new(big.Int).Abs(d.value)

	return Decimal{
		value: d2Value,
		exp:   d.exp,
	}
}

// Cmp compares the numbers represented by d and d2 and return:
//
//    -1 if d <  d2
//     0 if d == d2
//    +1 if d >  d2
func (d Decimal) Cmp(d2 Decimal) int {
	d.ensureInitialized()
	d2.ensureInitialized()

	if d.exp == d2.exp {
		return d.value.Cmp(d2.value)
	}

	rd, rd2 := RescalePair(d, d2)

	return rd.value.Cmp(rd2.value)
}

// Equal returns whether the numbers represented by d and d2 are equal.
func (d Decimal) Equal(d2 Decimal) bool {
	return d.Cmp(d2) == 0
}

// GreaterThan returns true when d is greater than d2.
func (d Decimal) GreaterThan(d2 Decimal) bool {
	return d.Cmp(d2) == 1
}

// GreaterThanOrEqual returns true when d is greater than or equal to d2.
func (d Decimal) GreaterThanOrEqual(d2 Decimal) bool {
	cmp := d.Cmp(d2)

	return cmp == 1 || cmp == 0
}

// LessThan returns true when d is greater than d2.
func (d Decimal) LessThan(d2 Decimal) bool {
	return d.Cmp(d2) == -1
}

// LessThanOrEqual returns true when d is greater than or equal to d2.
func (d Decimal) LessThanOrEqual(d2 Decimal) bool {
	cmp := d.Cmp(d2)

	return cmp == -1 || cmp == 0
}

// Sign returns:
//
//    -1 if d <  0
//     0 if d == 0
//    +1 if d >  0
func (d Decimal) Sign() int {
	if d.value == nil {
		return 0
	}

	return d.value.Sign()
}

// Truncate truncates off digits from the number, without rounding.
func (d Decimal) Truncate(precision int32) Decimal {
	d.ensureInitialized()

	if precision >= 0 && -precision > d.exp {
		return d.rescale(-precision)
	}

	return d
}

// DivRound divides and rounds to a given precision.
func (d Decimal) DivRound(d2 Decimal, precision int32) Decimal {
	q, r := d.QuoRem(d2, precision)

	var rv2 big.Int

	rv2.Abs(r.value)
	rv2.Lsh(&rv2, 1)

	r2 := Decimal{value: &rv2, exp: r.exp + precision}

	var c = r2.Cmp(d2.Abs())

	if c < 0 {
		return q
	}

	if d.value.Sign()*d2.value.Sign() < 0 {
		return q.Sub(New(1, -precision))
	}

	return q.Add(New(1, -precision))
}

// QuoRem does division with remainder.
func (d Decimal) QuoRem(d2 Decimal, precision int32) (Decimal, Decimal) {
	d.ensureInitialized()
	d2.ensureInitialized()

	if d2.value.Sign() == 0 {
		panic("decimal division by 0")
	}

	scale := -precision
	e := int64(d.exp - d2.exp - scale)

	if e > math.MaxInt32 || e < math.MinInt32 {
		panic("overflow in decimal QuoRem")
	}

	var aa, bb, expo big.Int
	var scalerest int32

	if e < 0 {
		aa = *d.value
		expo.SetInt64(-e)
		bb.Exp(tenInt, &expo, nil)
		bb.Mul(d2.value, &bb)
		scalerest = d.exp
	} else {
		expo.SetInt64(e)
		aa.Exp(tenInt, &expo, nil)
		aa.Mul(d.value, &aa)
		bb = *d2.value
		scalerest = scale + d2.exp
	}

	var q, r big.Int
	q.QuoRem(&aa, &bb, &r)
	dq := Decimal{value: &q, exp: scale}
	dr := Decimal{value: &r, exp: scalerest}

	return dq, dr
}

// RescalePair rescales two decimals to common exponential values.
// (minimal exp of both decimals)
func RescalePair(d1 Decimal, d2 Decimal) (Decimal, Decimal) {
	d1.ensureInitialized()
	d2.ensureInitialized()

	if d1.exp == d2.exp {
		return d1, d2
	}

	baseScale := min(d1.exp, d2.exp)

	if baseScale != d1.exp {
		return d1.rescale(baseScale), d2
	}

	return d1, d2.rescale(baseScale)
}

//

func (d *Decimal) ensureInitialized() {
	if d.value == nil {
		d.value = new(big.Int)
	}
}

func min(x int32, y int32) int32 {
	if x >= y {
		return y
	}

	return x
}

func (d Decimal) string(trimTrailingZeros bool) string {
	if d.exp >= 0 {
		return d.rescale(0).value.String()
	}

	abs := new(big.Int).Abs(d.value)
	str := abs.String()

	var intPart, fractionalPart string

	dExpInt := int(d.exp)

	if len(str) > -dExpInt {
		intPart = str[:len(str)+dExpInt]
		fractionalPart = str[len(str)+dExpInt:]
	} else {
		intPart = "0"

		num0s := -dExpInt - len(str)
		fractionalPart = strings.Repeat("0", num0s) + str
	}

	if trimTrailingZeros {
		i := len(fractionalPart) - 1

		for ; i >= 0; i-- {
			if fractionalPart[i] != '0' {
				break
			}
		}

		fractionalPart = fractionalPart[:i+1]
	}

	number := intPart

	if len(fractionalPart) > 0 {
		number += "." + fractionalPart
	}

	if d.value.Sign() < 0 {
		return "-" + number
	}

	return number
}

func (d Decimal) rescale(exp int32) Decimal {
	d.ensureInitialized()

	if d.exp == exp {
		return Decimal{
			new(big.Int).Set(d.value),
			d.exp,
		}
	}

	diff := math.Abs(float64(exp) - float64(d.exp))
	value := new(big.Int).Set(d.value)

	expScale := new(big.Int).Exp(tenInt, big.NewInt(int64(diff)), nil)

	if exp > d.exp {
		value = value.Quo(value, expScale)
	} else if exp < d.exp {
		value = value.Mul(value, expScale)
	}

	return Decimal{
		value: value,
		exp:   exp,
	}
}
