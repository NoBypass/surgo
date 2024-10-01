package marshal

import (
	"github.com/NoBypass/surgo/v2/errs"
	"reflect"
	"time"
)

func (m *Marshaler) Unmarshal(src, dest any) error {
	srcVal := reflect.ValueOf(src)
	destVal := reflect.ValueOf(dest)

	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return errs.ErrUnmarshal.Withf("dest must be a non-nil pointer")
	}

	destVal = destVal.Elem()

	switch srcVal.Kind() {
	case reflect.Slice:
		return m.parseSlice(srcVal, destVal)
	case reflect.Map:
		return m.parseMap(srcVal, destVal)
	default:
		return parseValue(srcVal, destVal)
	}
}

func (m *Marshaler) parseSlice(srcVal reflect.Value, destVal reflect.Value) error {
	if destVal.Kind() != reflect.Slice {
		return errs.ErrUnmarshal.Withf("cannot assign slice to %s", destVal.Type())
	}

	destType := destVal.Type().Elem()
	slice := reflect.MakeSlice(reflect.SliceOf(destType), srcVal.Len(), srcVal.Len())

	for i := 0; i < srcVal.Len(); i++ {
		if err := m.Unmarshal(srcVal.Index(i).Interface(), slice.Index(i).Addr().Interface()); err != nil {
			return err
		}
	}

	destVal.Set(slice)
	return nil
}

func (m *Marshaler) parseMap(srcVal reflect.Value, destVal reflect.Value) error {
	destKind := destVal.Kind()
	destType := destVal.Type()

	if destKind == reflect.Map {
		if destType.Key().Kind() != reflect.String {
			return errs.ErrUnmarshal.Withf("cannot assign to map with non-string key type")
		}

		destVal.Set(reflect.MakeMap(destType))

		for _, keyVal := range srcVal.MapKeys() {
			srcElemVal := srcVal.MapIndex(keyVal)
			destElemVal := reflect.New(destType.Elem()).Elem()

			if err := m.Unmarshal(srcElemVal.Interface(), destElemVal.Addr().Interface()); err != nil {
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
		return errs.ErrUnmarshal.Withf("cannot assign map to %s", destType)
	}

	for i := 0; i < destType.NumField(); i++ {
		field := destType.Field(i)
		dbTag := field.Tag.Get("db")
		if dbTag == "" {
			dbTag = field.Tag.Get(string(*m))
		}
		if dbTag == "" {
			dbTag = field.Name
		}

		mapVal := srcVal.MapIndex(reflect.ValueOf(dbTag))
		if !mapVal.IsValid() {
			continue
		}

		fieldVal := destVal.Field(i)
		if !fieldVal.CanSet() {
			continue
		}

		if err := m.Unmarshal(mapVal.Interface(), fieldVal.Addr().Interface()); err != nil {
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
			t, err := stringToTime(srcVal.String())
			if err != nil {
				return errs.ErrUnmarshal.Withf("cannot parse time: %w", err)
			}

			destVal.Set(reflect.ValueOf(t))
			return nil
		} else if destVal.Type().String() == "time.Duration" && srcVal.Kind() == reflect.String {
			d, err := stringToDuration(srcVal.String())
			if err != nil {
				return errs.ErrUnmarshal.Withf("cannot parse duration: %w", err)
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
