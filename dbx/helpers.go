package dbx

import (
	"database/sql"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

var iScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func isScalar(tValue reflect.Type) bool {
	if tValue.Kind() == reflect.Ptr {
		tValue = tValue.Elem()
	}
	if tValue.String() == "time.Time" {
		return true
	}

	if isScanner(tValue) {
		return true
	}

	switch tValue.Kind() {
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16:
		return true
	case reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16:
		return true
	case reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		return true
	default:
		return false
	}
}

func isScanner(tValue reflect.Type) bool {
	return tValue.Implements(iScanner)
}

func checkType(tValue reflect.Type, expected string, kinds ...reflect.Kind) error {
	var oValue = tValue
	var n = len(kinds) - 1
	for k, knd := range kinds {
		if tValue.Kind() != knd {
			return errors.Errorf("Expected %s. Actual: %s", expected, oValue)
		}
		if k < n {
			tValue = tValue.Elem()
		}
	}

	return nil
}

func checkPtrScalar(tValue reflect.Type) error {
	var oValue = tValue
	if !isPointer(tValue) {
		return errors.Errorf("dest should be a pointer to a scalar. Actual: %s", oValue)
	}
	tValue = tValue.Elem()
	if !isScalar(tValue) {
		return errors.Errorf("dest should be a pointer to a scalar. Actual: %s", oValue)
	}

	return nil
}
func checkPtrSliceScalar(tDest reflect.Type) error {
	var oDest = tDest
	if !isPointer(tDest) {
		return errors.Errorf("dest should be a pointer to a list of scalars. Actual: %s", oDest)
	}
	tDest = tDest.Elem()
	if !isSlice(tDest) {
		return errors.Errorf("dest should be a pointer to a list of scalars. Actual: %s", oDest)
	}
	tDest = tDest.Elem()
	if !isScalar(tDest) {
		return errors.Errorf("dest should be a pointer to a list of scalars. Actual: %s", oDest)
	}

	return nil
}

func isStruct(tValue reflect.Type) bool {
	if tValue.Kind() == reflect.Ptr {
		tValue = tValue.Elem()
	}

	return tValue.Kind() == reflect.Struct
}

func isSlice(tValue reflect.Type) bool {
	if tValue.Kind() == reflect.Ptr {
		tValue = tValue.Elem()
	}

	return tValue.Kind() == reflect.Slice
}

func isPointer(tValue reflect.Type) bool {
	return tValue.Kind() == reflect.Ptr
}

func ptrType(tValue reflect.Type) reflect.Type {
	if tValue.Kind() == reflect.Ptr {
		tValue = tValue.Elem()
	}

	return tValue
}

func ptrTypeAny(value interface{}) reflect.Type {
	return ptrType(reflect.TypeOf(value))
}

func isNotNullableError(err error) bool {
	if err == nil {
		return false
	}
	var msg = err.Error()
	return strings.HasPrefix(msg, "Field ") && strings.HasSuffix(msg, "is not nullable, but scanned a null value")
}

func getScannerValue(scanner sql.Scanner) (interface{}, bool) {
	var vScanner = reflect.ValueOf(scanner).Elem()
	// Ex: NullString
	var name = vScanner.Type().Name()[4:]
	var vValue = vScanner.FieldByName(name).Interface()
	var vValid = vScanner.FieldByName("Valid").Bool()
	return vValue, vValid
}
