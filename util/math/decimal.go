package decimal

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"

	"github.com/sirupsen/logrus"

	"github.com/pingcap/types"
)

var (
	Zero = New(0, 0)
	One  = New(1, 0)
)

type Decimal struct {
	value *types.MyDecimal
}

func NewFromString(value string) (*Decimal, error) {
	dec := new(types.MyDecimal)
	if err := dec.FromString([]byte(value)); err != nil && err != types.ErrTruncated {
		return nil, err
	}

	return &Decimal{value: dec}, nil
}

func NewFromFloat64(value float64) (*Decimal, error) {
	dec := new(types.MyDecimal)
	if err := dec.FromFloat64(value); err != nil {
		return nil, err
	}

	return &Decimal{value: dec}, nil
}

func NewFromUint64(value uint64) *Decimal {
	dec := new(types.MyDecimal)

	return &Decimal{value: dec.FromUint(value)}
}

func New(value int64, exp int) *Decimal {
	dec := types.NewDecFromInt(value)
	if exp < 0 {
		pow := types.NewDecFromInt(1)
		if err := pow.Shift(-exp); err != nil {
			logError(err, "decimal shift error")
		}

		if err := types.DecimalDiv(dec, pow, dec, 10); err != nil {
			logError(err, "decimal divide error")
		}
	} else {
		if err := dec.Shift(exp); err != nil {
			logError(err, "decimal shift error")
		}
	}
	return &Decimal{value: dec}
}

func (d *Decimal) Abs() *Decimal {
	coefficient := New(1, 0)
	if d.value.IsNegative() {
		coefficient = New(-1, 0)
	}
	return d.Mul(coefficient)
}

func (d *Decimal) Add(d2 *Decimal) *Decimal {
	result := new(types.MyDecimal)
	if err := types.DecimalAdd(d.value, d2.value, result); err != nil {
		logError(err, "decimal add error")
	}

	return &Decimal{value: result}
}

func (d *Decimal) Sub(d2 *Decimal) *Decimal {
	result := new(types.MyDecimal)
	if err := types.DecimalSub(d.value, d2.value, result); err != nil {
		logError(err, "decimal subtract error")
	}

	return &Decimal{value: result}
}

func (d *Decimal) Mul(d2 *Decimal) *Decimal {
	result := new(types.MyDecimal)
	if err := types.DecimalMul(d.value, d2.value, result); err != nil {
		logError(err, "decimal multiply error")
	}

	return &Decimal{value: result}
}

func (d *Decimal) Div(d2 *Decimal) *Decimal {
	result := new(types.MyDecimal)
	if err := types.DecimalDiv(d.value, d2.value, result, 10); err != nil {
		logError(err, "decimal divide error")
	}

	return &Decimal{value: result}
}

func (d *Decimal) Min(d2 *Decimal) *Decimal {
	if d.LessThan(d2) {
		return d
	}
	return d2
}

func (d *Decimal) Max(d2 *Decimal) *Decimal {
	if d.GreaterThan(d2) {
		return d
	}
	return d2
}

func (d *Decimal) IsZero() bool {
	return d.value.IsZero()
}

func (d *Decimal) IsNegative() bool {
	return d.value.IsNegative()
}

func (d *Decimal) Float64() float64 {
	result, err := d.value.ToFloat64()
	if err != nil {
		logError(err, "decimal to float64 error")
	}

	return result
}

func (d *Decimal) Int64() int64 {
	result, err := d.value.ToInt()
	if err != nil {
		logError(err, "decimal to int64 error")
	}

	return result
}

func (d *Decimal) BigInt() *big.Int {
	result, ok := big.NewInt(0).SetString(d.StringTruncateFixed(0), 10)
	if !ok {
		logError(errors.New("decimal to big int error"))
	}

	return result
}

func (d *Decimal) Uint64() (uint64, error) {
	result, err := d.value.ToUint()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (d *Decimal) IsInt64() bool {
	_, err := d.value.ToInt()
	if err != nil && err != types.ErrTruncated {
		return false
	}
	return true
}

func (d *Decimal) Cmp(d2 *Decimal) int {
	return d.value.Compare(d2.value)
}

func (d *Decimal) Equal(d2 *Decimal) bool {
	return d.Cmp(d2) == 0
}

func (d *Decimal) LessThan(d2 *Decimal) bool {
	return d.Cmp(d2) < 0
}

func (d *Decimal) LessOrEquals(d2 *Decimal) bool {
	return d.Cmp(d2) <= 0
}

func (d *Decimal) GreaterThan(d2 *Decimal) bool {
	return d.Cmp(d2) > 0
}

func (d *Decimal) GreaterOrEquals(d2 *Decimal) bool {
	return d.Cmp(d2) >= 0
}

// TruncateFixed  get the truncate price of decimal , (8.061).TruncateFixed(2) ->  "8.06"
func (d *Decimal) TruncateFixed(places int) *Decimal {
	result := new(types.MyDecimal)
	if err := d.value.Round(result, places, types.ModeTruncate); err != nil {
		logError(err, "decimal string truncate fixed error")
	}

	return &Decimal{value: result}
}

// RoundFixed  get the round price of decimal , (8.06).RoundFixed(1) ->  "8.1"
func (d *Decimal) RoundFixed(places int) *Decimal {
	result := new(types.MyDecimal)
	if err := d.value.Round(result, places, types.ModeHalfEven); err != nil {
		logError(err, "decimal string round fixed error")
	}

	return &Decimal{value: result}
}

// CeilFixed get the ceil price of decimal ,  (8.06).CeilFixed -> "9"
func (d *Decimal) CeilFixed(places int) *Decimal {
	fixedRate := d.RoundFixed(places)
	if fixedRate.Cmp(d) == -1 {
		fixedRate = fixedRate.Add(New(1, -places))
	}

	return fixedRate.RoundFixed(places)
}

func (d *Decimal) String() string {
	return d.value.String()
}

// StringWithoutZero return decimal to string without zero
func (d *Decimal) StringWithoutZero() string {
	return d.value.StringWithoutZero()
}

// StringRoundFixed  get the round price of decimal , (8.06).StringRoundFixed(1) ->  "8.1"
func (d *Decimal) StringRoundFixed(places int) string {
	return d.RoundFixed(places).String()
}

// StringTruncateFixed  get the truncate price of decimal , (8.06).StringTruncateFixed(2) ->  "8.00"
func (d *Decimal) StringTruncateFixed(places int) string {
	return d.TruncateFixed(places).String()
}

// StringCeilFixed get the ceil price of decimal ,  (8.06).StringCeilFixed -> "9"
func (d *Decimal) StringCeilFixed(places int) string {
	return d.CeilFixed(places).String()
}

func logError(err error, args ...interface{}) {
	if err != types.ErrTruncated {
		logrus.WithField("err", err).Panic(args...)
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Decimal) UnmarshalJSON(decimalBytes []byte) error {
	if string(decimalBytes) == "null" {
		return nil
	}

	str, err := trimQuote(decimalBytes)
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", decimalBytes, err)
	}

	value := &types.MyDecimal{}
	if err := value.UnmarshalJSON([]byte(str)); err != nil {
		return err
	}

	d.value = value
	return nil
}

// MarshalJSONWithoutQuotes should be set to true if you want the decimal to
// be JSON marshaled as a number, instead of as a string.
// WARNING: this is dangerous for decimals with many digits, since many JSON
// unmarshallers (ex: Javascript's) will unmarshal JSON numbers to IEEE 754
// double-precision floating point numbers, which means you can potentially
// silently lose precision.
var MarshalJSONWithoutQuotes = false

// MarshalJSON implements the json.Marshaler interface.
func (d Decimal) MarshalJSON() ([]byte, error) {
	if MarshalJSONWithoutQuotes {
		return d.value.MarshalJSON()
	}

	value, _ := d.value.MarshalJSON()
	str := "\"" + string(value) + "\""
	return []byte(str), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for XML
// deserialization.
func (d *Decimal) UnmarshalText(text []byte) error {
	str := string(text)
	dec, err := NewFromString(str)
	d.value = dec.value
	if err != nil {
		return fmt.Errorf("error decoding string '%s': %s", str, err)
	}

	return nil
}

// MarshalText implements the encoding.TextMarshaler interface for XML
// serialization.
func (d Decimal) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// Scan implements the sql.Scanner interface for database deserialization
func (d *Decimal) Scan(value interface{}) error {
	// first try to see if the data is stored in database as a Numeric datatype
	var err error
	var dec *Decimal
	switch v := value.(type) {
	case float64:
		// numeric in sqlite3 sends us float64
		if dec, err = NewFromFloat64(v); err != nil {
			return err
		}
	case int64:
		// at least in sqlite3 when the value is 0 in db, the data is sent
		// to us as an int64 instead of a float64 ...
		dec = New(int64(v), 0)
	default:
		// default is trying to interpret value stored as string
		str, err := trimQuote(v)
		if err != nil {
			return err
		}

		if dec, err = NewFromString(str); err != nil {
			return err
		}
	}

	d.value = dec.value
	return nil
}

// Value implements the sql.Valuer interface for database serialization
func (d Decimal) Value() (driver.Value, error) {
	if d.value == nil {
		return "0", nil
	}

	return d.value.String(), nil
}

func trimQuote(value interface{}) (string, error) {
	var bytes []byte

	switch v := value.(type) {
	case float32:
		bytes = []byte(fmt.Sprint(value))
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	default:
		return "", fmt.Errorf("could not convert value '%+v' to byte array of type '%T'",
			value, value)
	}

	// If the amount is quoted, strip the quotes
	if len(bytes) > 2 && bytes[0] == '"' && bytes[len(bytes)-1] == '"' {
		bytes = bytes[1 : len(bytes)-1]
	}
	return string(bytes), nil
}
