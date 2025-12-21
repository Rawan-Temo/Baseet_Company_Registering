// utils/updater.go
package utils

import (
	"reflect"
)

// UpdateStruct updates only non-nil fields from source to destination
// src is the DTO and dest is the model
func UpdateStruct(dest, src interface{}) {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src).Elem()
	srcType := srcVal.Type()

	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Field(i)
		srcFieldType := srcType.Field(i)

		// Skip unexported fields
		if !srcFieldType.IsExported() {
			continue
		}

		// Check if field is a pointer and not nil
		if srcField.Kind() == reflect.Ptr && !srcField.IsNil() {
			destField := destVal.FieldByName(srcFieldType.Name)
			if destField.IsValid() && destField.CanSet() {
				destField.Set(srcField.Elem())
			}
		}
	}
}