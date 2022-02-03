package common

import (
	"runtime/debug"
	"time"

	"github.com/pkg/errors"
)

// RecoverError
func RecoverError(err interface{}, place string) error {
	if err == nil {
		return nil
	}
	var stack = string(debug.Stack())
	return errors.Errorf("Recovering from %s: %+v\r\nStack:\r\n%s", place, err, stack)
}

// PtrBool method
func PtrBool(value bool) *bool {
	return &value
}

// PtrString method
func PtrString(value string) *string {
	return &value
}

// PtrInt method
func PtrInt(value int) *int {
	return &value
}

// PtrByte method
func PtrByte(value byte) *byte {
	return &value
}

// PtrInt16 method
func PtrInt16(value int16) *int16 {
	return &value
}

// PtrInt32 method
func PtrInt32(value int32) *int32 {
	return &value
}

// PtrInt64 method
func PtrInt64(value int64) *int64 {
	return &value
}

// PtrUint method
func PtrUint(value uint) *uint {
	return &value
}

// PtrUint16 method
func PtrUint16(value uint16) *uint16 {
	return &value
}

// PtrUint32 method
func PtrUint32(value uint32) *uint32 {
	return &value
}

// PtrUint64 method
func PtrUint64(value uint64) *uint64 {
	return &value
}

// PtrFloat32 method
func PtrFloat32(value float32) *float32 {
	return &value
}

// PtrFloat64 method
func PtrFloat64(value float64) *float64 {
	return &value
}

// PtrTime method
func PtrTime(value time.Time) *time.Time {
	return &value
}
