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

	if destKind == reflect.Map {
		if destType.Key().Kind() != reflect.String {
			return fmt.Errorf("cannot assign to map with non-string key type")
		}

		destVal.Set(reflect.MakeMap(destType))

		for _, keyVal := range srcVal.MapKeys() {
			srcElemVal := srcVal.MapIndex(keyVal)
			destElemVal := reflect.New(destType.Elem()).Elem()

			if err := scan(srcElemVal.Interface(), destElemVal.Addr().Interface()); err != nil {
				return err
			}

			destVal.SetMapIndex(keyVal, destElemVal)
		}

		return nil
	} else if destKind == reflect.Ptr {
		if destVal.IsNil() {
			destVal.Set(reflect.New(destType.Elem()))
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
		if destVal.Type() == reflect.TypeOf(ID{}) && srcVal.Kind() == reflect.String {
			destVal.Set(reflect.ValueOf(stringToID(strings.Split(srcVal.String(), ":")[1])))
			return nil
		} else if strings.Split(destVal.Type().String(), "[")[0] == "surgo.Record" {
			if srcVal.Kind() == reflect.String {
				destVal.FieldByName("ID").Set(reflect.ValueOf(stringToID(srcVal.String())))
				return nil
			} else {
				underlyingType := destVal.FieldByName("Data").Elem().Type()
				newInstance := reflect.New(underlyingType)
				if err := scan(srcVal.Interface(), newInstance.Interface()); err != nil {
					return err
				}
				destVal.FieldByName("Data").Set(newInstance.Elem())
				return nil
			}
		} else if destVal.Kind() == reflect.Int && srcVal.Kind() == reflect.Float64 {
			destVal.SetInt(int64(srcVal.Float()))
			return nil
		} else if destVal.Type() == reflect.TypeOf(time.Time{}) && srcVal.Kind() == reflect.String {
			t, err := stringToTime(srcVal.String())
			if err != nil {
				return err
			}

			destVal.Set(reflect.ValueOf(t))
			return nil
		} else if destVal.Type().String() == "time.Duration" && srcVal.Kind() == reflect.String {
			d, err := stringToDuration(srcVal.String())
			if err != nil {
				return err
			}

			destVal.Set(reflect.ValueOf(d))
			return nil
		}

		destVal.Set(reflect.Zero(destVal.Type()))
		return nil
	}

	destVal.Set(srcVal)
	return nil
}
