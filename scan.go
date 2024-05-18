package surgo

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

func scan(src any, dest any) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return fmt.Errorf("dest must be a non-nil pointer")
	}

	destVal = destVal.Elem()

	switch srcVal.Kind() {
	case reflect.Slice:
		return parseSlice(srcVal, destVal)
	case reflect.Map:
		return parseMap(srcVal, destVal)
	default:
		return parseValue(srcVal, destVal)
	}
}

func parseSlice(srcVal reflect.Value, destVal reflect.Value) error {
	if destVal.Kind() != reflect.Slice {
		return fmt.Errorf("cannot assign slice to %s", destVal.Type())
	}

	destType := destVal.Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(destType), srcVal.Len(), srcVal.Len())

	for i := 0; i < srcVal.Len(); i++ {
		if err := scan(srcVal.Index(i).Interface(), slice.Index(i).Addr().Interface()); err != nil {
			return err
		}
	}

	destVal.Set(slice)
	return nil
}

func parseMap(srcVal reflect.Value, destVal reflect.Value) error {
	destKind := destVal.Kind()
	destType := destVal.Type()

	if destKind == reflect.Ptr {
		if destVal.IsNil() {
			newDest := reflect.New(destType.Elem())
			destVal.Set(newDest)
		}

		destVal = destVal.Elem()
		destType = destVal.Type()
	} else if destKind != reflect.Struct {
		return fmt.Errorf("cannot assign map to %s", destType)
	}

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		tag := field.Tag.Get("db")

		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		mapVal := srcVal.MapIndex(reflect.ValueOf(tag))
		if !mapVal.IsValid() {
			continue
		}

		fieldVal := destVal.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		if err := scan(mapVal.Interface(), fieldVal.Addr().Interface()); err != nil {
			return err
		}
	}

	return nil
}

func parseValue(srcVal reflect.Value, destVal reflect.Value) error {
	if !srcVal.IsValid() {
		return nil
	}

	if !srcVal.Type().AssignableTo(destVal.Type()) {
		if destVal.Kind() == reflect.Int && srcVal.Kind() == reflect.Float64 {
			destVal.SetInt(int64(srcVal.Float()))
			return nil
		} else if destVal.Type() == reflect.TypeOf(time.Time{}) && srcVal.Kind() == reflect.String {
			t, err := time.Parse(time.RFC3339, srcVal.String())
			if err != nil {
				return err
			}

			destVal.Set(reflect.ValueOf(t))
			return nil
		}

		return fmt.Errorf("cannot assign %s to %s", srcVal.Type(), destVal.Type())
	}

	destVal.Set(srcVal)
	return nil
}
