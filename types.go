package types

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"

	"github.com/aarondl/null/v8/convert"
	"github.com/friendsofgo/errors"
	"github.com/google/uuid"
)

var (
	_ = time.Time{} // To make sure the time package is imported.
	_ = uuid.UUID{} // To make sure the uuid package is imported.
)

var nullBytes = []byte("null")

func isNullBytes(d []byte) bool {
	return string(d) == string(nullBytes)
}

func ParseFromString(typeAsString, value string) (any, error) {
	switch strings.TrimPrefix(typeAsString, "types.") {

	case "Bool":
		return BoolFromString(value)

	case "Date":
		return DateFromString(value)

	case "Float64":
		return Float64FromString(value)

	case "Int":
		return IntFromString(value)

	case "Int16":
		return Int16FromString(value)

	case "Int64":
		return Int64FromString(value)

	case "JSON":
		return JSONFromString(value)

	case "RichText":
		return RichTextFromString(value)

	case "String":
		return StringFromString(value)

	case "Time":
		return TimeFromString(value)

	case "Timestamp":
		return TimestampFromString(value)

	case "UUID":
		return UUIDFromString(value)

	default:
		return nil, errors.New(fmt.Sprintf("invalid type: %s", typeAsString))
	}
}

func IsEmptyArray(a any) bool {
	switch a.(type) {

	case []Bool:
		return len(a.([]Bool)) == 0

	case []Date:
		return len(a.([]Date)) == 0

	case []Float64:
		return len(a.([]Float64)) == 0

	case []Int:
		return len(a.([]Int)) == 0

	case []Int16:
		return len(a.([]Int16)) == 0

	case []Int64:
		return len(a.([]Int64)) == 0

	case []JSON:
		return len(a.([]JSON)) == 0

	case []RichText:
		return len(a.([]RichText)) == 0

	case []String:
		return len(a.([]String)) == 0

	case []Time:
		return len(a.([]Time)) == 0

	case []Timestamp:
		return len(a.([]Timestamp)) == 0

	case []UUID:
		return len(a.([]UUID)) == 0

	default:
		return false
	}
}

// Bool is used to represent booleans
type Bool struct {
	underlying bool
	isDefined  bool
	isNil      bool
}

// NewBool creates a new Bool object.
func NewBool(underlying bool) Bool {
	return Bool{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewBoolFromPtr creates a new Bool object from a pointer.
func NewBoolFromPtr(underlying *bool) Bool {
	if underlying != nil {
		return NewBool(*underlying)
	}

	return Bool{
		isDefined: true,
		isNil:     true,
	}
}

// NewBoolUndefined creates a new undefined Bool object.
func NewBoolUndefined() Bool {
	return Bool{}
}

func BoolFromStringPtr(strPtr *string) (Bool, error) {
	if strPtr == nil {
		return NewBoolFromPtr(nil), nil
	}

	return BoolFromString(*strPtr)
}

func BoolFromString(str string) (Bool, error) {
	if str == "" {
		return NewBoolFromPtr(nil), nil
	}

	underlying, err := strconv.ParseBool(strings.TrimSpace(str))
	if err != nil {
		return Bool{}, err
	}

	return Bool{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Bool
func (s Bool) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return fmt.Sprintf("%t", s.underlying)
}

// Bool returns the bool value.
func (s Bool) Bool() bool {
	return s.underlying
}

// BoolPtr returns the bool value as a pointer.
func (s Bool) BoolPtr() *bool {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Bool) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Bool) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Bool is nil, which is specifically used by sqlboiler queries
func (s Bool) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Bool, but returns nil if undefined.
func (s Bool) Ptr() *Bool {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Bool-pointer,
// will return an undefined Bool if the pointer is nil.
func (s *Bool) Val() Bool {
	if s == nil {
		return NewBoolFromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Bool) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Bool) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Bool) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = false
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Bool) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying, nil
}

// Date is used to represent dates according to the ISO 8601 standard.
type Date struct {
	underlying time.Time
	isDefined  bool
	isNil      bool
}

// NewDate creates a new Date object.
func NewDate(underlying time.Time) Date {
	return Date{
		underlying: underlyingTime(underlying, "2006-01-02"),
		isDefined:  true,
		isNil:      false,
	}
}

// NewDateFromPtr creates a new Date object from a pointer.
func NewDateFromPtr(underlying *time.Time) Date {
	if underlying != nil {
		return NewDate(*underlying)
	}

	return Date{
		isDefined: true,
		isNil:     true,
	}
}

// NewDateUndefined creates a new undefined Date object.
func NewDateUndefined() Date {
	return Date{}
}

func DateFromStringPtr(strPtr *string) (Date, error) {
	if strPtr == nil {
		return NewDateFromPtr(nil), nil
	}

	return DateFromString(*strPtr)
}

func DateFromString(str string) (Date, error) {
	if str == "" {
		return NewDateFromPtr(nil), nil
	}

	layouts := []string{
		"2006-01-02",  // YYYY-MM-DD
		"01-02-06",    // MM-DD-YY, US format short.. Apparently what excel makes dates into.
		"02-01-06",    // DD-MM-YY, Reverse order from Excelize
		"06-01-02",    // YY-MM-DD, Can only happen if Year is > 31 so the above check DD-MM-YY has failed
		"01-02-2006",  // MM-DD-YYYY, US format
		"02-Jan-2006", // DD-MMM-YYYY, old style Oracle
		"02-Jan-06",   // DD-MMM-YY, old style Oracle
	}

	var underlying time.Time
	var err error
	for _, layout := range layouts {
		underlying, err = time.Parse(layout, str)
		if err == nil {
			break
		}
		err = errors.New("invalid date format: " + str)
	}

	if err != nil {
		return Date{}, err
	}

	return Date{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Date
func (s Date) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return s.underlying.Format("2006-01-02")
}

// Date returns the time.Time value.
func (s Date) Date() time.Time {
	return s.underlying
}

// DatePtr returns the time.Time value as a pointer.
func (s Date) DatePtr() *time.Time {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Date) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Date) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Date is nil, which is specifically used by sqlboiler queries
func (s Date) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Date, but returns nil if undefined.
func (s Date) Ptr() *Date {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Date-pointer,
// will return an undefined Date if the pointer is nil.
func (s *Date) Val() Date {
	if s == nil {
		return NewDateFromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Date) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying.Format("2006-01-02"))
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Date) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	var str string
	err := json.Unmarshal(d, &str)
	if err != nil {
		return err
	}

	s.underlying, err = time.Parse("2006-01-02", str)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Date) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = time.Time{}
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Date) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying, nil
}

// Float64 is used to represent 64-bit floating point numbers.
type Float64 struct {
	underlying float64
	isDefined  bool
	isNil      bool
}

// NewFloat64 creates a new Float64 object.
func NewFloat64(underlying float64) Float64 {
	return Float64{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewFloat64FromPtr creates a new Float64 object from a pointer.
func NewFloat64FromPtr(underlying *float64) Float64 {
	if underlying != nil {
		return NewFloat64(*underlying)
	}

	return Float64{
		isDefined: true,
		isNil:     true,
	}
}

// NewFloat64Undefined creates a new undefined Float64 object.
func NewFloat64Undefined() Float64 {
	return Float64{}
}

func Float64FromStringPtr(strPtr *string) (Float64, error) {
	if strPtr == nil {
		return NewFloat64FromPtr(nil), nil
	}

	return Float64FromString(*strPtr)
}

func Float64FromString(str string) (Float64, error) {
	if str == "" {
		return NewFloat64FromPtr(nil), nil
	}

	underlying, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		return Float64{}, err
	}

	return Float64{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Float64
func (s Float64) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	// First, format the float with two decimal places
	formatted := fmt.Sprintf("%.2f", s.underlying)

	// Convert to float to trim unnecessary zeros,
	// ignore the error since we know it shouldn't fail.
	floatVal, _ := strconv.ParseFloat(formatted, 64)

	// Reformat the float without unnecessary zeros
	formatted = fmt.Sprintf("%g", floatVal)

	// Replace dot with comma
	return strings.Replace(formatted, ".", ",", 1)
}

// Float64 returns the float64 value.
func (s Float64) Float64() float64 {
	return s.underlying
}

// Float64Ptr returns the float64 value as a pointer.
func (s Float64) Float64Ptr() *float64 {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Float64) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Float64) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Float64 is nil, which is specifically used by sqlboiler queries
func (s Float64) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Float64, but returns nil if undefined.
func (s Float64) Ptr() *Float64 {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Float64-pointer,
// will return an undefined Float64 if the pointer is nil.
func (s *Float64) Val() Float64 {
	if s == nil {
		return NewFloat64FromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Float64) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Float64) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Float64) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = 0
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Float64) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying, nil
}

// Int is used to represent integers.
type Int struct {
	underlying int
	isDefined  bool
	isNil      bool
}

// NewInt creates a new Int object.
func NewInt(underlying int) Int {
	return Int{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewIntFromPtr creates a new Int object from a pointer.
func NewIntFromPtr(underlying *int) Int {
	if underlying != nil {
		return NewInt(*underlying)
	}

	return Int{
		isDefined: true,
		isNil:     true,
	}
}

// NewIntUndefined creates a new undefined Int object.
func NewIntUndefined() Int {
	return Int{}
}

func IntFromStringPtr(strPtr *string) (Int, error) {
	if strPtr == nil {
		return NewIntFromPtr(nil), nil
	}

	return IntFromString(*strPtr)
}

func IntFromString(str string) (Int, error) {
	if str == "" {
		return NewIntFromPtr(nil), nil
	}

	parsed, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	underlying := int(parsed)

	if err != nil {
		return Int{}, err
	}

	return Int{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Int
func (s Int) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return fmt.Sprintf("%d", s.underlying)
}

// Int returns the int value.
func (s Int) Int() int {
	return s.underlying
}

// IntPtr returns the int value as a pointer.
func (s Int) IntPtr() *int {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Int) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Int) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Int is nil, which is specifically used by sqlboiler queries
func (s Int) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Int, but returns nil if undefined.
func (s Int) Ptr() *Int {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Int-pointer,
// will return an undefined Int if the pointer is nil.
func (s *Int) Val() Int {
	if s == nil {
		return NewIntFromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Int) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Int) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Int) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = 0
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Int) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return int64(s.underlying), nil
}

// Int16 is used to represent 16-bit integers.
type Int16 struct {
	underlying int16
	isDefined  bool
	isNil      bool
}

// NewInt16 creates a new Int16 object.
func NewInt16(underlying int16) Int16 {
	return Int16{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewInt16FromPtr creates a new Int16 object from a pointer.
func NewInt16FromPtr(underlying *int16) Int16 {
	if underlying != nil {
		return NewInt16(*underlying)
	}

	return Int16{
		isDefined: true,
		isNil:     true,
	}
}

// NewInt16Undefined creates a new undefined Int16 object.
func NewInt16Undefined() Int16 {
	return Int16{}
}

func Int16FromStringPtr(strPtr *string) (Int16, error) {
	if strPtr == nil {
		return NewInt16FromPtr(nil), nil
	}

	return Int16FromString(*strPtr)
}

func Int16FromString(str string) (Int16, error) {
	if str == "" {
		return NewInt16FromPtr(nil), nil
	}

	parsed, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	underlying := int16(parsed)

	if err != nil {
		return Int16{}, err
	}

	return Int16{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Int16
func (s Int16) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return fmt.Sprintf("%d", s.underlying)
}

// Int16 returns the int16 value.
func (s Int16) Int16() int16 {
	return s.underlying
}

// Int16Ptr returns the int16 value as a pointer.
func (s Int16) Int16Ptr() *int16 {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Int16) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Int16) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Int16 is nil, which is specifically used by sqlboiler queries
func (s Int16) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Int16, but returns nil if undefined.
func (s Int16) Ptr() *Int16 {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Int16-pointer,
// will return an undefined Int16 if the pointer is nil.
func (s *Int16) Val() Int16 {
	if s == nil {
		return NewInt16FromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Int16) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Int16) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Int16) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = 0
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Int16) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return int64(s.underlying), nil
}

// Int64 is used to represent 64-bit integers.
type Int64 struct {
	underlying int64
	isDefined  bool
	isNil      bool
}

// NewInt64 creates a new Int64 object.
func NewInt64(underlying int64) Int64 {
	return Int64{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewInt64FromPtr creates a new Int64 object from a pointer.
func NewInt64FromPtr(underlying *int64) Int64 {
	if underlying != nil {
		return NewInt64(*underlying)
	}

	return Int64{
		isDefined: true,
		isNil:     true,
	}
}

// NewInt64Undefined creates a new undefined Int64 object.
func NewInt64Undefined() Int64 {
	return Int64{}
}

func Int64FromStringPtr(strPtr *string) (Int64, error) {
	if strPtr == nil {
		return NewInt64FromPtr(nil), nil
	}

	return Int64FromString(*strPtr)
}

func Int64FromString(str string) (Int64, error) {
	if str == "" {
		return NewInt64FromPtr(nil), nil
	}

	parsed, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	underlying := int64(parsed)

	if err != nil {
		return Int64{}, err
	}

	return Int64{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Int64
func (s Int64) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return fmt.Sprintf("%d", s.underlying)
}

// Int64 returns the int64 value.
func (s Int64) Int64() int64 {
	return s.underlying
}

// Int64Ptr returns the int64 value as a pointer.
func (s Int64) Int64Ptr() *int64 {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Int64) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Int64) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Int64 is nil, which is specifically used by sqlboiler queries
func (s Int64) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Int64, but returns nil if undefined.
func (s Int64) Ptr() *Int64 {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Int64-pointer,
// will return an undefined Int64 if the pointer is nil.
func (s *Int64) Val() Int64 {
	if s == nil {
		return NewInt64FromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Int64) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Int64) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Int64) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = 0
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Int64) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return int64(s.underlying), nil
}

// JSON is used to represent JSON data.
type JSON struct {
	underlying json.RawMessage
	isDefined  bool
	isNil      bool
}

// NewJSON creates a new JSON object.
func NewJSON(underlying json.RawMessage) JSON {
	return JSON{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewJSONFromPtr creates a new JSON object from a pointer.
func NewJSONFromPtr(underlying *json.RawMessage) JSON {
	if underlying != nil {
		return NewJSON(*underlying)
	}

	return JSON{
		isDefined: true,
		isNil:     true,
	}
}

// NewJSONUndefined creates a new undefined JSON object.
func NewJSONUndefined() JSON {
	return JSON{}
}

func JSONFromStringPtr(strPtr *string) (JSON, error) {
	if strPtr == nil {
		return NewJSONFromPtr(nil), nil
	}

	return JSONFromString(*strPtr)
}

func JSONFromString(str string) (JSON, error) {
	if str == "" {
		return NewJSONFromPtr(nil), nil
	}

	underlying, err := json.Marshal(strings.TrimSpace(str))
	if err != nil {
		return JSON{}, err
	}

	return JSON{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output JSON
func (s JSON) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return string(s.underlying)
}

// JSON returns the json.RawMessage value.
func (s JSON) RawMessage() json.RawMessage {
	return s.underlying
}

// JSONPtr returns the json.RawMessage value as a pointer.
func (s JSON) RawMessagePtr() *json.RawMessage {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s JSON) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s JSON) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if JSON is nil, which is specifically used by sqlboiler queries
func (s JSON) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for JSON, but returns nil if undefined.
func (s JSON) Ptr() *JSON {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a JSON-pointer,
// will return an undefined JSON if the pointer is nil.
func (s *JSON) Val() JSON {
	if s == nil {
		return NewJSONFromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s JSON) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *JSON) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *JSON) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	source, ok := value.([]byte)
	if !ok {
		return errors.New("incompatible type for json")
	}

	s.underlying = append((s.underlying)[0:0], source...)

	return nil
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s JSON) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return []byte(s.underlying), nil
}

func (s *JSON) Marshal(obj interface{}) error {
	res, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	s.isDefined = true
	s.isNil = isNullBytes(res)
	s.underlying = res
	return nil
}

// RichText is used to represent rich text.
type RichText struct {
	underlying string
	isDefined  bool
	isNil      bool
}

// NewRichText creates a new RichText object.
func NewRichText(underlying string) RichText {
	return RichText{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewRichTextFromPtr creates a new RichText object from a pointer.
func NewRichTextFromPtr(underlying *string) RichText {
	if underlying != nil {
		return NewRichText(*underlying)
	}

	return RichText{
		isDefined: true,
		isNil:     true,
	}
}

// NewRichTextUndefined creates a new undefined RichText object.
func NewRichTextUndefined() RichText {
	return RichText{}
}

func RichTextFromStringPtr(strPtr *string) (RichText, error) {
	if strPtr == nil {
		return NewRichTextFromPtr(nil), nil
	}

	return RichTextFromString(*strPtr)
}

func RichTextFromString(str string) (RichText, error) {
	if str == "" {
		return NewRichTextFromPtr(nil), nil
	}

	var err error
	underlying := strings.TrimSpace(str)

	if err != nil {
		return RichText{}, err
	}

	return RichText{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output RichText
func (s RichText) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return string(s.underlying)
}

// RichText returns the string value.
func (s RichText) RichText() string {
	return s.underlying
}

// RichTextPtr returns the string value as a pointer.
func (s RichText) RichTextPtr() *string {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s RichText) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s RichText) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if RichText is nil, which is specifically used by sqlboiler queries
func (s RichText) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for RichText, but returns nil if undefined.
func (s RichText) Ptr() *RichText {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a RichText-pointer,
// will return an undefined RichText if the pointer is nil.
func (s *RichText) Val() RichText {
	if s == nil {
		return NewRichTextFromPtr(nil)
	}

	return *s
}

// RichTextToLower returns the underlying value of RichText in lower case.
func RichTextToLower(s RichText) RichText {
	if !s.IsNil() {
		s.underlying = strings.ToLower(s.underlying)
	}

	return s
}

// RichTextsToLower returns the underlying values of RichTexts in lower case.
func RichTextsToLower(s []RichText) []RichText {
	for i := range s {
		s[i].underlying = strings.ToLower(s[i].underlying)
	}

	return s
}

// Text returns the plain text value of the rich text.
//
// The method basically converts HTML content to plain text,
// removing all HTML tags and unescaping HTML entities.
//
// For example, "<p>Hello my &lt;b&gt;friend&lt;/b&gt;</p>" becomes "Hello my <b>friend</b>".
func (s RichText) Text() (string, error) {
	doc, err := html.Parse(strings.NewReader(s.underlying))
	if err != nil {
		return "", err
	}

	// walkNodes recursively traverses the HTML node tree and extracts text from text nodes
	var walkNodes func(b *bytes.Buffer, n *html.Node) error

	walkNodes = func(b *bytes.Buffer, n *html.Node) error {
		if n.Type == html.TextNode {
			_, err := b.WriteString(n.Data)
			if err != nil {
				return err
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := walkNodes(b, c); err != nil {
				return err
			}
		}

		// Add double newlines for specific closing tags unless it's the last node
		if n.Type == html.ElementNode {
			switch n.Data {
			case "p", "h1", "h2", "h3", "pre", "ul", "ol":
				if n.NextSibling != nil || n.Parent.NextSibling != nil {
					_, err := b.WriteString("\n\n")
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}

	var b bytes.Buffer
	if err := walkNodes(&b, doc); err != nil {
		return "", err
	}

	return strings.TrimSuffix(b.String(), "\n\n"), nil
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s RichText) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	text, err := s.Text()
	if err != nil {
		return nil, errors.Wrap(err, "cannot convert to text "+s.underlying)
	}

	richText := struct {
		Content string `json:"content"`
		Text    string `json:"text"`
	}{
		Content: s.underlying,
		Text:    text,
	}

	jsonBytes, err := json.Marshal(richText)
	if err != nil {
		return nil, errors.Wrap(err, s.underlying)
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *RichText) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	richText := struct {
		Content string `json:"content"` // We only care about the content
		Text    string `json:"-"`
	}{}

	err := json.Unmarshal(d, &richText)
	if err != nil {
		return err
	}

	s.underlying = richText.Content

	s.underlying = strings.TrimSpace(s.underlying)

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *RichText) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s RichText) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying, nil
}

// String is used to represent strings.
type String struct {
	underlying string
	isDefined  bool
	isNil      bool
}

// NewString creates a new String object.
func NewString(underlying string) String {
	return String{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewStringFromPtr creates a new String object from a pointer.
func NewStringFromPtr(underlying *string) String {
	if underlying != nil {
		return NewString(*underlying)
	}

	return String{
		isDefined: true,
		isNil:     true,
	}
}

// NewStringUndefined creates a new undefined String object.
func NewStringUndefined() String {
	return String{}
}

func StringFromStringPtr(strPtr *string) (String, error) {
	if strPtr == nil {
		return NewStringFromPtr(nil), nil
	}

	return StringFromString(*strPtr)
}

func StringFromString(str string) (String, error) {
	if str == "" {
		return NewStringFromPtr(nil), nil
	}

	var err error
	underlying := strings.TrimSpace(str)

	if err != nil {
		return String{}, err
	}

	return String{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String returns the string value.
func (s String) String() string {
	return s.underlying
}

// StringPtr returns the string value as a pointer.
func (s String) StringPtr() *string {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s String) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s String) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if String is nil, which is specifically used by sqlboiler queries
func (s String) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for String, but returns nil if undefined.
func (s String) Ptr() *String {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a String-pointer,
// will return an undefined String if the pointer is nil.
func (s *String) Val() String {
	if s == nil {
		return NewStringFromPtr(nil)
	}

	return *s
}

// StringToLower returns the underlying value of String in lower case.
func StringToLower(s String) String {
	if !s.IsNil() {
		s.underlying = strings.ToLower(s.underlying)
	}

	return s
}

// StringsToLower returns the underlying values of Strings in lower case.
func StringsToLower(s []String) []String {
	for i := range s {
		s[i].underlying = strings.ToLower(s[i].underlying)
	}

	return s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s String) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *String) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	s.underlying = strings.TrimSpace(s.underlying)

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *String) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = ""
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s String) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying, nil
}

// Time is used to represent a times by the format "HH:MM"
type Time struct {
	underlying time.Time
	isDefined  bool
	isNil      bool
}

// NewTime creates a new Time object.
func NewTime(underlying time.Time) Time {
	return Time{
		underlying: underlyingTime(underlying, "15:04"),
		isDefined:  true,
		isNil:      false,
	}
}

// NewTimeFromPtr creates a new Time object from a pointer.
func NewTimeFromPtr(underlying *time.Time) Time {
	if underlying != nil {
		return NewTime(*underlying)
	}

	return Time{
		isDefined: true,
		isNil:     true,
	}
}

// NewTimeUndefined creates a new undefined Time object.
func NewTimeUndefined() Time {
	return Time{}
}

func TimeFromStringPtr(strPtr *string) (Time, error) {
	if strPtr == nil {
		return NewTimeFromPtr(nil), nil
	}

	return TimeFromString(*strPtr)
}

func TimeFromString(str string) (Time, error) {
	if str == "" {
		return NewTimeFromPtr(nil), nil
	}

	underlying, err := time.Parse("15:04", strings.TrimSpace(str))
	if err != nil {
		return Time{}, err
	}

	return Time{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Time
func (s Time) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return s.underlying.Format("15:04")
}

// Time returns the time.Time value.
func (s Time) Time() time.Time {
	return s.underlying
}

// TimePtr returns the time.Time value as a pointer.
func (s Time) TimePtr() *time.Time {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Time) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Time) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Time is nil, which is specifically used by sqlboiler queries
func (s Time) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Time, but returns nil if undefined.
func (s Time) Ptr() *Time {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Time-pointer,
// will return an undefined Time if the pointer is nil.
func (s *Time) Val() Time {
	if s == nil {
		return NewTimeFromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Time) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying.Format("15:04"))
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Time) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	var str string
	err := json.Unmarshal(d, &str)
	if err != nil {
		return err
	}

	s.underlying, err = time.Parse("15:04", str)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Time) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	val, ok := value.(string)
	if ok {
		t, err := time.Parse("15:04:05", val)
		if err == nil {
			s.underlying = t
			return nil
		}
	}
	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Time) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying, nil
}

// Timestamp is used to represent a timestamps according to the RFC3339 format.
type Timestamp struct {
	underlying time.Time
	isDefined  bool
	isNil      bool
}

// NewTimestamp creates a new Timestamp object.
func NewTimestamp(underlying time.Time) Timestamp {
	return Timestamp{
		underlying: underlyingTime(underlying, "2006-01-02T15:04:05Z07:00"),
		isDefined:  true,
		isNil:      false,
	}
}

// NewTimestampFromPtr creates a new Timestamp object from a pointer.
func NewTimestampFromPtr(underlying *time.Time) Timestamp {
	if underlying != nil {
		return NewTimestamp(*underlying)
	}

	return Timestamp{
		isDefined: true,
		isNil:     true,
	}
}

// NewTimestampUndefined creates a new undefined Timestamp object.
func NewTimestampUndefined() Timestamp {
	return Timestamp{}
}

func TimestampFromStringPtr(strPtr *string) (Timestamp, error) {
	if strPtr == nil {
		return NewTimestampFromPtr(nil), nil
	}

	return TimestampFromString(*strPtr)
}

func TimestampFromString(str string) (Timestamp, error) {
	if str == "" {
		return NewTimestampFromPtr(nil), nil
	}

	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05Z",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04",
		"2006-01-02 15:04",
		"2006-01-02",
		"1/2/06 15:04",
		"1/2/06 15:04:05",
		"1/2/2006 15:04",
		"1/2/2006 15:04:05",
	}

	for _, format := range formats {
		underlying, err := time.Parse(format, strings.TrimSpace(str))
		if err == nil {
			return Timestamp{
				underlying: underlying,
				isDefined:  true,
				isNil:      false,
			}, nil
		}
	}

	underlying, err := time.Parse("2006-01-02T15:04:05Z07:00", strings.TrimSpace(str))
	if err != nil {
		return Timestamp{}, err
	}

	return Timestamp{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

// String output Timestamp
func (s Timestamp) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return s.underlying.Format("2006-01-02T15:04:05Z07:00")
}

// Timestamp returns the time.Time value.
func (s Timestamp) Timestamp() time.Time {
	return s.underlying
}

// TimestampPtr returns the time.Time value as a pointer.
func (s Timestamp) TimestampPtr() *time.Time {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s Timestamp) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s Timestamp) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if Timestamp is nil, which is specifically used by sqlboiler queries
func (s Timestamp) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for Timestamp, but returns nil if undefined.
func (s Timestamp) Ptr() *Timestamp {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a Timestamp-pointer,
// will return an undefined Timestamp if the pointer is nil.
func (s *Timestamp) Val() Timestamp {
	if s == nil {
		return NewTimestampFromPtr(nil)
	}

	return *s
}

func (t Timestamp) After(other Timestamp) bool {
	return t.Timestamp().After(other.Timestamp())
}

func (t Timestamp) Before(other Timestamp) bool {
	return t.Timestamp().Before(other.Timestamp())
}

func (t Timestamp) Equal(other Timestamp) bool {
	return t.Timestamp().Equal(other.Timestamp())
}

// MinutesUntil returns the minutes until the given timestamp
func (from Timestamp) MinutesUntil(to Timestamp) int {
	return int(to.Timestamp().Sub(from.Timestamp()).Minutes())
}

func (t Timestamp) Date() Date {
	return NewDate(t.Timestamp())
}

// Returns a new Timestamp with the time set to the start of the day.
func (s Timestamp) StartOfDay(location *time.Location) Timestamp {
	return NewTimestamp(time.Date(
		s.underlying.Year(),
		s.underlying.Month(),
		s.underlying.Day(),
		0,
		0,
		0,
		0,
		location,
	))
}

// Returns a new Timestamp with the time set to the end of the day.
func (s Timestamp) EndOfDay(location *time.Location) Timestamp {
	return NewTimestamp(time.Date(
		s.underlying.Year(),
		s.underlying.Month(),
		s.underlying.Day(),
		23,
		59,
		59,
		0,
		location,
	))
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s Timestamp) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying.Format("2006-01-02T15:04:05Z07:00"))
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *Timestamp) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	var str string
	err := json.Unmarshal(d, &str)
	if err != nil {
		return err
	}

	s.underlying, err = time.Parse("2006-01-02T15:04:05Z07:00", str)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *Timestamp) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s Timestamp) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying, nil
}

// UUID is used to represent a UUID.
type UUID struct {
	underlying uuid.UUID
	isDefined  bool
	isNil      bool
}

// NewRandomUUID generates a new UUID object.
func NewRandomUUID() UUID {
	return UUID{
		underlying: uuid.New(),
		isDefined:  true,
		isNil:      false,
	}
}

// NewUUID creates a new UUID object.
func NewUUID(underlying uuid.UUID) UUID {
	return UUID{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}
}

// NewUUIDFromPtr creates a new UUID object from a pointer.
func NewUUIDFromPtr(underlying *uuid.UUID) UUID {
	if underlying != nil {
		return NewUUID(*underlying)
	}

	return UUID{
		isDefined: true,
		isNil:     true,
	}
}

// NewUUIDUndefined creates a new undefined UUID object.
func NewUUIDUndefined() UUID {
	return UUID{}
}

func UUIDFromStringPtr(strPtr *string) (UUID, error) {
	if strPtr == nil {
		return NewUUIDFromPtr(nil), nil
	}

	return UUIDFromString(*strPtr)
}

func UUIDFromString(str string) (UUID, error) {
	if str == "" {
		return NewUUIDFromPtr(nil), nil
	}

	underlying, err := uuid.Parse(strings.TrimSpace(str))
	if err != nil {
		return UUID{}, err
	}

	return UUID{
		underlying: underlying,
		isDefined:  true,
		isNil:      false,
	}, nil
}

func UUIDsFromStrings(strings []string) []UUID {
	uuids := make([]UUID, len(strings))
	for i := range strings {
		uuids[i] = NewUUID(uuid.MustParse(strings[i]))
	}
	return uuids
}

func UUIDsToStrings(uuids []UUID) []string {
	strings := make([]string, len(uuids))
	for i := range uuids {
		strings[i] = uuids[i].String()
	}
	return strings
}

// String output UUID
func (s UUID) String() string {
	// If the value is nil we return an empty string
	if s.IsNil() {
		return ""
	}

	return s.underlying.String()
}

// UUID returns the uuid.UUID value.
func (s UUID) UUID() uuid.UUID {
	return s.underlying
}

// UUIDPtr returns the uuid.UUID value as a pointer.
func (s UUID) UUIDPtr() *uuid.UUID {
	if s.IsNil() {
		return nil
	}
	return &s.underlying
}

// IsDefined returns true if the value was defined in the JSON input or was scanned from the database.
func (s UUID) IsDefined() bool {
	return s.isDefined
}

// IsNil returns true if the value is nil or undefined.
func (s UUID) IsNil() bool {
	// if the value is undefined, it is nil even though "isNil" will be set to false
	if !s.isDefined {
		return true
	}

	return s.isNil
}

// IsZero checks if UUID is nil, which is specifically used by sqlboiler queries
func (s UUID) IsZero() bool { return s.IsNil() }

// Ptr returns the pointer for UUID, but returns nil if undefined.
func (s UUID) Ptr() *UUID {
	if !s.isDefined {
		return nil
	}

	return &s
}

// Val returns the value of a UUID-pointer,
// will return an undefined UUID if the pointer is nil.
func (s *UUID) Val() UUID {
	if s == nil {
		return NewUUIDFromPtr(nil)
	}

	return *s
}

// MarshalJSON implements the json Marshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Marshaler
func (s UUID) MarshalJSON() ([]byte, error) {
	if s.IsNil() {
		return nullBytes, nil
	}

	jsonBytes, err := json.Marshal(s.underlying)
	if err != nil {
		return nil, errors.Wrap(err, s.String())
	}

	return jsonBytes, nil
}

// UnmarshalJSON implements the json Unmarshaler interface.
//
// See: https://pkg.go.dev/encoding/json#Unmarshaler
func (s *UUID) UnmarshalJSON(d []byte) error {
	s.isNil = isNullBytes(d)
	s.isDefined = true

	if s.isNil {
		return nil
	}

	err := json.Unmarshal(d, &s.underlying)
	if err != nil {
		return err
	}

	return nil
}

// Scan assigns a value from a database driver and implements the sql Scanner interface.
//
// See https://pkg.go.dev/database/sql#Scanner
func (s *UUID) Scan(value interface{}) error {
	s.isNil = (nil == value)
	s.isDefined = true

	if s.isNil {
		s.underlying = uuid.Nil
		return nil
	}

	return convert.ConvertAssign(&s.underlying, value)
}

// Value implements the driver Valuer interface.
//
// See https://pkg.go.dev/database/sql/driver#Valuer
func (s UUID) Value() (driver.Value, error) {
	if s.IsNil() {
		return nil, nil
	}
	return s.underlying.String(), nil
}

func underlyingTime(t time.Time, format string) time.Time {
	t, _ = time.Parse(format, t.Format(format))
	return t.UTC()
}

// Types is an interface which can be used for generated code to force package dependency
type Types interface{}
