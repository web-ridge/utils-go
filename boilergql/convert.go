package boilergql

import (
	"strconv"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/iancoleman/strcase"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/types/pgeo"

	"github.com/volatiletech/sqlboiler/v4/types"
)

const IDSeparator = "-"

type RemovedID struct {
	ID uint
}

type RemovedUint64ID struct {
	ID uint64
}

type RemovedStringID struct {
	ID string
}

func RemovedIDsToBoiler(removedIds []RemovedID) []uint {
	uintIDs := make([]uint, len(removedIds))
	for index, id := range removedIds {
		uintIDs[index] = id.ID
	}
	return uintIDs
}

func RemovedIDsToBoilerUint(removedIds []RemovedID) []uint {
	return RemovedIDsToBoiler(removedIds)
}

func RemovedIDsToBoilerInt(removedIds []RemovedID) []int {
	uintIDs := make([]int, len(removedIds))
	for index, id := range removedIds {
		uintIDs[index] = int(id.ID)
	}
	return uintIDs
}

func RemovedIDsToBoilerInt64(removedIds []RemovedID) []int64 {
	intIDs := make([]int64, len(removedIds))
	for index, id := range removedIds {
		intIDs[index] = int64(id.ID)
	}
	return intIDs
}

func RemovedIDsToBoilerString(removedIds []RemovedStringID) []string {
	sIDs := make([]string, len(removedIds))
	for index, id := range removedIds {
		sIDs[index] = id.ID
	}
	return sIDs
}

func RemovedUint64IDsToBoiler(removedIds []RemovedUint64ID) []uint64 {
	uintIDs := make([]uint64, len(removedIds))
	for index, id := range removedIds {
		uintIDs[index] = id.ID
	}
	return uintIDs
}

func IntsToInterfaces(ints []int) []interface{} {
	interfaces := make([]interface{}, len(ints))
	for index, number := range ints {
		interfaces[index] = number
	}
	return interfaces
}

func StringsToInterfaces(strings []string) []interface{} {
	interfaces := make([]interface{}, len(strings))
	for index, v := range strings {
		interfaces[index] = v
	}
	return interfaces
}

func FloatsToInterfaces(fs []float64) []interface{} {
	interfaces := make([]interface{}, len(fs))
	for index, number := range fs {
		interfaces[index] = number
	}
	return interfaces
}

func IDsToBoilerInterfaces(ids []string) []interface{} {
	interfaces := make([]interface{}, len(ids))
	for index, id := range ids {
		interfaces[index] = IDToBoiler(id)
	}
	return interfaces
}

func StringIDsToBoilerString(ids []string) []string {
	sa := make([]string, len(ids))
	for index, stringID := range ids {
		sa[index] = StringIDToBoilerString(stringID)
	}
	return sa
}

func StringIDToBoilerString(id string) string {
	splitID := strings.SplitN(id, IDSeparator, 2)
	if len(splitID) == 2 {
		return splitID[1]
	}
	return ""
}

func IDsToBoiler(ids []string) []uint {
	ints := make([]uint, len(ids))
	for index, stringID := range ids {
		ints[index] = IDToBoiler(stringID)
	}
	return ints
}

func IDToBoiler(id string) uint {
	splitID := strings.SplitN(id, IDSeparator, 2)
	if len(splitID) == 2 {
		// nolint: errcheck
		i, _ := strconv.ParseUint(splitID[1], 10, 64)
		return uint(i)
	}
	return 0
}

func IDsToBoilerUint(ids []string) []uint {
	return IDsToBoiler(ids)
}

func IDToBoilerUint(id string) uint {
	return IDToBoiler(id)
}

func IDToBoilerInt(id string) int {
	return int(IDToBoiler(id))
}

func IDsToBoilerInt(ids []string) []int {
	ints := make([]int, len(ids))
	for index, stringID := range ids {
		ints[index] = IDToBoilerInt(stringID)
	}
	return ints
}

func IDToNullBoiler(id string) null.Uint {
	uintID := IDToBoiler(id)
	if uintID == 0 {
		return null.NewUint(0, false)
	}
	return null.Uint{
		Uint:  uintID,
		Valid: true,
	}
}

func NullUintToNullInt(u null.Uint) null.Int {
	return null.Int{
		Int:   int(u.Uint),
		Valid: u.Valid,
	}
}

func IDToGraphQL(id uint, tableName string) string {
	return strcase.ToLowerCamel(tableName) + IDSeparator + strconv.Itoa(int(id))
}

func StringIDToGraphQL(id string, tableName string) string {
	return strcase.ToLowerCamel(tableName) + IDSeparator + id
}

func StringIDsToGraphQL(ids []string, tableName string) []string {
	stringIDs := make([]string, len(ids))
	for index, id := range ids {
		stringIDs[index] = StringIDToGraphQL(id, tableName)
	}
	return stringIDs
}

func IntIDToGraphQL(id int, tableName string) string {
	return IDToGraphQL(uint(id), tableName)
}

func IDsToGraphQL(ids []uint, tableName string) []string {
	stringIDs := make([]string, len(ids))
	for index, id := range ids {
		stringIDs[index] = IDToGraphQL(id, tableName)
	}
	return stringIDs
}

func UintIDsToGraphQL(ids []uint, tableName string) []string {
	return IDsToGraphQL(ids, tableName)
}

func IntIDsToGraphQL(ids []int, tableName string) []string {
	stringIDs := make([]string, len(ids))
	for index, id := range ids {
		stringIDs[index] = IntIDToGraphQL(id, tableName)
	}
	return stringIDs
}

func NullDotBoolToPointerBool(v null.Bool) *bool {
	return v.Ptr()
}

func PointerBoolToBool(v *bool) bool {
	if v == nil {
		return false
	}
	return *v
}

func NullDotStringToPointerString(v null.String) *string {
	return v.Ptr()
}

func NullDotTimeToInt(v null.Time) int {
	if !v.Valid {
		return 0
	}
	u := int(v.Time.Unix())
	return u
}

func NullDotTimeToPointerInt(v null.Time) *int {
	if !v.Valid {
		return nil
	}
	u := int(v.Time.Unix())
	return &u
}

func TimeDotTimeToPointerInt(v time.Time) *int {
	u := int(v.Unix())
	return &u
}

func TimeDotTimeToInt(v time.Time) int {
	return int(v.Unix())
}

func IntToTimeDotTime(v int) time.Time {
	return time.Unix(int64(v), 0)
}

func NullDotStringToString(v null.String) string {
	if !v.Valid {
		return ""
	}

	return v.String
}

func NullDotUintToPointerInt(v null.Uint) *int {
	if !v.Valid {
		return nil
	}
	u := int(v.Uint)
	return &u
}

func PointerStringToString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func PointerIntToNullDotTime(v *int) null.Time {
	if v == nil {
		return null.Time{
			Time:  time.Time{},
			Valid: false,
		}
	}
	return null.TimeFrom(time.Unix(int64(*v), 0))
}

func StringToNullDotString(v string) null.String {
	return null.StringFrom(v)
}

func PointerStringToNullDotString(v *string) null.String {
	return null.StringFromPtr(v)
}

func PointerBoolToNullDotBool(v *bool) null.Bool {
	return null.BoolFromPtr(v)
}

func TypesNullDecimalToFloat64(v types.NullDecimal) float64 {
	if v.Big == nil {
		return 0
	}
	f, _ := v.Float64()
	return f
}

func Float64ToTypesNullDecimal(v float64) types.NullDecimal {
	d := new(decimal.Big)
	d.SetFloat64(v)
	return types.NewNullDecimal(d)
}

func TypesDecimalToFloat64(v types.Decimal) float64 {
	if v.Big == nil {
		return 0
	}
	f, _ := v.Float64()
	return f
}

func Float64ToTypesDecimal(v float64) types.Decimal {
	d := new(decimal.Big)
	d.SetFloat64(v)
	return types.NewDecimal(d)
}

func PointerFloat64ToTypesDecimal(v *float64) types.Decimal {
	if v == nil {
		return types.NewDecimal(decimal.New(0, 0))
	}
	d := new(decimal.Big)
	d.SetFloat64(*v)
	return types.NewDecimal(d)
}

func PointerFloat64ToTypesNullDecimal(v *float64) types.NullDecimal {
	if v == nil {
		return types.NewNullDecimal(nil)
	}
	return Float64ToTypesNullDecimal(*v)
}

func TypesNullDecimalToPointerString(v types.NullDecimal) *string {
	if v.Big == nil {
		return nil
	}
	s := v.String()
	if s == "" {
		return nil
	}
	return &s
}

func TypesNullDecimalToPointerFloat64(v types.NullDecimal) *float64 {
	if v.Big == nil {
		return nil
	}
	s, ok := v.Float64()
	if !ok {
		return nil
	}
	return &s
}

func PointerStringToTypesNullDecimal(v *string) types.NullDecimal {
	if v == nil {
		return types.NewNullDecimal(nil)
	}
	d := new(decimal.Big)
	if _, ok := d.SetString(*v); !ok {
		nd := types.NewNullDecimal(nil)
		if err := d.Context.Err(); err != nil {
			return nd
		}
		// TODO: error handling maybe write log line here
		// https://github.com/volatiletech/sqlboiler/blob/master/types/decimal.go#L156
		return nd
	}

	return types.NewNullDecimal(d)
}

func PointerIntToInt(v *int) int {
	if v == nil {
		return 0
	}
	return *v
}

func PointerIntToUint(v *int) uint {
	if v == nil {
		return 0
	}
	return uint(*v)
}

func PointerIntToInt8(v *int) int8 {
	if v == nil {
		return 0
	}
	return int8(*v)
}

func PointerIntToNullDotInt(v *int) null.Int {
	return null.IntFromPtr((v))
}

func PointerIntToNullDotUint(v *int) null.Uint {
	if v == nil {
		return null.UintFromPtr(nil)
	}
	uv := *v
	return null.UintFrom(uint(uv))
}

func NullDotIntToPointerInt(v null.Int) *int {
	return v.Ptr()
}

func IntToInt8(v int) int8 {
	return int8(v)
}

func Int8ToInt(v int8) int {
	return int(v)
}

func NullDotFloat64ToPointerFloat64(v null.Float64) *float64 {
	return v.Ptr()
}

func PointerFloat64ToNullDotFloat64(v *float64) null.Float64 {
	return null.Float64FromPtr(v)
}

func IntToUint(v int) uint {
	return uint(v)
}

func UintToInt(v uint) int {
	return int(v)
}

func Int16ToInt(v int16) int {
	return int(v)
}

func IntToInt16(v int) int16 {
	return int16(v)
}

func PointerIntToInt16(v *int) int16 {
	if v != nil {
		return int16(*v)
	}

	return 0
}

func BoolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func IntToBool(v int) bool {
	return v == 1
}

func NullDotBoolToPointerInt(v null.Bool) *int {
	if !v.Valid {
		return nil
	}

	if v.Bool {
		i := 1
		return &i
	}
	i := 0
	return &i
}

func PointerIntToNullDotBool(v *int) null.Bool {
	if v == nil {
		return null.Bool{
			Valid: false,
		}
	}
	return null.Bool{
		Valid: true,
		Bool:  *v == 1,
	}
}

func NullDotIntToUint(v null.Int) uint {
	return uint(v.Int)
}

func NullDotUintToUint(v null.Uint) uint {
	return v.Uint
}

func NullDotIntIsFilled(v null.Int) bool {
	return !v.IsZero()
}

func NullDotUintIsFilled(v null.Uint) bool {
	return !v.IsZero()
}

func UintIsFilled(v uint) bool {
	return v != 0
}

func IntIsFilled(v int) bool {
	return v != 0
}

func NullDotStringIsFilled(v null.String) bool {
	return !v.IsZero()
}

func StringIsFilled(v string) bool {
	return v != ""
}

func PointerIntToTimeDotTime(v *int) time.Time {
	if v == nil {
		return time.Time{}
	}

	return time.Unix(int64(*v), 0)
}

func PointerFloat64ToNullDotFloat32(v *float64) null.Float32 {
	val := null.Float32{}
	if v != nil {
		val.SetValid(float32(*v))
	}
	return val
}

func Float32ToFloat64(v float32) float64 {
	return float64(v)
}

func Float64ToFloat32(v float64) float32 {
	return float32(v)
}

func PointerFloat32ToFloat64(v *float32) *float64 {
	if v == nil {
		return nil
	}
	result := float64(*v)
	return &result
}

func PointerFloat64ToFloat32(v *float64) *float32 {
	if v == nil {
		return nil
	}
	result := float32(*v)
	return &result
}

func Float32ToNullFloat64(v float32) null.Float64 {
	return null.Float64From(float64(v))
}

func Float64ToNullFloat32(v float64) null.Float32 {
	return null.Float32From(float32(v))
}

func PointerFloat64ToFloat64(v *float64) float64 {
	if v == nil {
		return 0
	}
	return *v
}

func NullDotFloat32ToPointerFloat64(v null.Float32) *float64 {
	if v.IsZero() {
		return nil
	}
	val := new(float64)
	*val = float64(v.Float32)
	return val
}

func PointerIntToNullDotInt16(v *int) null.Int16 {
	val := null.Int16{}
	if v != nil {
		val.SetValid(int16(*v))
	}

	return val
}

func NullDotInt16ToPointerInt(v null.Int16) *int {
	if v.IsZero() {
		return nil
	}

	val := new(int)

	*val = int(v.Int16)
	return val
}

func PgeoPointToGeoPoint(v pgeo.Point) GeoPoint {
	return GeoPoint{
		X: v.X,
		Y: v.Y,
	}
}

func GeoPointToPgeoPoint(v GeoPoint) pgeo.Point {
	return pgeo.NewPoint(v.X, v.Y)
}

func PointerGeoPointToPgeoPoint(v *GeoPoint) pgeo.Point {
	if v == nil {
		return pgeo.NewPoint(0, 0)
	}

	return pgeo.NewPoint(v.X, v.Y)
}

func TimeDotTimeToPointerTimeDotTime(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}

	val := new(time.Time)
	*val = v

	return val
}

func NullDotTimeToPointerTimeDotTime(v null.Time) *time.Time {
	if !v.Valid {
		return nil
	}

	return TimeDotTimeToPointerTimeDotTime(v.Time)
}

func PointerTimeDotTimeToNullDotTime(v *time.Time) null.Time {
	return null.TimeFromPtr(v)
}

func PointerTimeToTimeDotTime(v *time.Time) time.Time {
	if v == nil {
		return time.Time{}
	}
	return *v
}

func IntToString(v int) string {
	return strconv.Itoa(v)
}

func StringToInt(v string) int {
	i, _ := strconv.Atoi(v) //nolint:errcheck
	return i
}

func StringToUint(v string) uint {
	return uint(StringToInt(v))
}

func PointerIntToString(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}

func PointerStringToInt(v *string) int {
	if v == nil {
		return 0
	}
	i, _ := strconv.Atoi(*v) //nolint:errcheck
	return i
}

func PointerStringToUint(v *string) uint {
	return uint(PointerStringToInt(v))
}

func PointerStringToNullDotUint(v *string) null.Uint {
	if v == nil {
		return null.UintFromPtr(nil)
	}
	return null.UintFrom(PointerStringToUint(v))
}

func StringToByteSlice(v string) []byte {
	return []byte(v)
}

func ByteSliceToString(v []byte) string {
	return string(v)
}

func PointerStringToByteSlice(v *string) []byte {
	if v == nil {
		return nil
	}
	return []byte(*v)
}

func Int64ToInt(v int64) int {
	return int(v)
}

func NullDotUintToPointerString(v null.Uint) *string {
	if !v.Valid {
		return nil
	}
	u := UintToString(v.Uint)
	return &u
}

func NullDotUint64ToPointerString(v null.Uint64) *string {
	if !v.Valid {
		return nil
	}
	u := Uint64ToString(v.Uint64)
	return &u
}

func UintToString(v uint) string {
	u := strconv.FormatUint(uint64(v), 10)
	return u
}

func IDToBoilerUint64(v string) uint64 {
	return uint64(IDToBoilerUint(v))
}

func Uint64ToUint(v uint64) uint {
	return uint(v)
}

func IDsToBoilerUint64(a []string) []uint64 {
	ids := IDsToBoilerUint(a)
	ids2 := make([]uint64, len(ids))
	for i, id := range ids {
		ids2[i] = uint64(id)
	}

	return ids2
}

func NullDotInt64ToPointerInt(v null.Int64) *int {
	intV := int(v.Int64)
	return &intV
}

func Uint64ToString(v uint64) string {
	return IntToString(int(v))
}

func IntToInt64(v int) int64 {
	return int64(v)
}

func PointerIntToInt64(v *int) int64 {
	if v == nil {
		return 0
	}
	return int64(*v)
}

func PointerIntToNullDotInt64(v *int) null.Int64 {
	if v == nil {
		return null.Int64FromPtr(nil)
	}
	return null.Int64From(int64(*v))
}
