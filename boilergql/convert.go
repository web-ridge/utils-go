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

type RemovedID struct {
	ID uint
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

func RemovedIDsToBoilerString(removedIds []RemovedStringID) []string {
	sIDs := make([]string, len(removedIds))
	for index, id := range removedIds {
		sIDs[index] = id.ID
	}
	return sIDs
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

func IDsToBoiler(ids []string) []uint {
	ints := make([]uint, len(ids))
	for index, stringID := range ids {
		ints[index] = IDToBoiler(stringID)
	}
	return ints
}

func IDToBoiler(ID string) uint {
	splitted := strings.Split(ID, "-")
	if len(splitted) > 1 {
		// nolint: errcheck
		i, _ := strconv.ParseUint(splitted[1], 10, 64)
		return uint(i)
	}
	return 0
}

func IDsToBoilerUint(ids []string) []uint {
	return IDsToBoiler(ids)
}

func IDToBoilerUint(ID string) uint {
	return IDToBoiler(ID)
}

func IDToBoilerInt(ID string) int {
	return int(IDToBoiler(ID))
}

func IDsToBoilerInt(ids []string) []int {
	ints := make([]int, len(ids))
	for index, stringID := range ids {
		ints[index] = IDToBoilerInt(stringID)
	}
	return ints
}

func IDToNullBoiler(ID string) null.Uint {
	uintID := IDToBoiler(ID)
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
	return strcase.ToLowerCamel(tableName) + "-" + strconv.Itoa(int(id))
}

func IntIDToGraphQL(id int, tableName string) string {
	return strcase.ToLowerCamel(tableName) + "-" + strconv.Itoa(id)
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

func NullDotTimeToPointerInt(v null.Time) *int {
	if !v.Valid {
		return nil
	}
	u := int(v.Time.Unix())
	return &u
}

func TimeTimeToInt(v time.Time) int {
	return int(v.Unix())
}

func IntToTimeTime(v int) time.Time {
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
	s := v.String()
	if s == "" {
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
		Valid: v != nil,
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

func PointerIntToTimeTime(v *int) time.Time {
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
