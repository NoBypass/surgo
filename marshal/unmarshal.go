package marshal

import (
	"github.com/NoBypass/surgo/v2/errs"
	"reflect"
	"regexp"
	"strconv"
	"time"
)

func (m *Marshaler) Unmarshal(src, dest any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errs.ErrUnmarshal.Withf("type mismatch: %v", r)
		}
	}()
	destVal := reflect.ValueOf(dest)

	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return errs.ErrUnmarshal.Withf("dest must be a non-nil pointer")
	}

	return m.unmarshal(reflect.ValueOf(src), destVal.Elem())
}

func (m *Marshaler) unmarshal(src, dest reflect.Value) error {
	switch src.Kind() {
	case reflect.Bool:
		return m.simpleValueDecoder(src, dest)
	case reflect.String:
		if dest.Type() == reflect.TypeOf(time.Time{}) {
			return m.timeDecoder(src, dest)
		} else if dest.Type().String() == "time.Duration" {
			return m.durationDecoder(src, dest)
		}
		return m.simpleValueDecoder(src, dest)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		if src.Type().AssignableTo(dest.Type()) {
			return m.simpleValueDecoder(src, dest)
		}
		return m.numberDecoder(src, dest)
	case reflect.Slice:
		return m.sliceDecoder(src, dest)
	case reflect.Interface:
		return m.interfaceDecoder(src, dest)
	case reflect.Map:
		return m.mapDecoder(src, dest)
	default:
		return errs.ErrUnmarshal.Withf("cannot unmarshal %s", src.Type())
	}
}

func (m *Marshaler) simpleValueDecoder(src, dest reflect.Value) error {
	dest.Set(src)
	return nil
}

func (m *Marshaler) numberDecoder(src, dest reflect.Value) error {
	switch dest.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		dest.SetInt(src.Convert(dest.Type()).Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		dest.SetUint(src.Convert(dest.Type()).Uint())
	case reflect.Float32, reflect.Float64:
		dest.SetFloat(src.Convert(dest.Type()).Float())
	default:
		return errs.ErrUnmarshal.Withf("unsupported number type: %s", dest.Type())
	}
	return nil
}

func (m *Marshaler) timeDecoder(src, dest reflect.Value) error {
	t, err := time.Parse(time.RFC3339, src.String())
	if err != nil {
		return errs.ErrUnmarshal.Withf("cannot parse time: %w", err)
	}

	dest.Set(reflect.ValueOf(t))
	return nil
}

func (m *Marshaler) durationDecoder(src, dest reflect.Value) error {
	re := regexp.MustCompile(`(\d+)(y|w|d|h|ms|m|s|µs|us|ns)`)
	matches := re.FindAllStringSubmatch(src.String(), -1)

	var duration time.Duration
	for _, match := range matches {
		value, err := strconv.Atoi(match[1])
		if err != nil {
			return errs.ErrUnmarshal.Withf("cannot parse duration: %w", err)
		}
		unit := match[2]

		duration += time.Duration(value) * units[unit]
	}

	dest.Set(reflect.ValueOf(duration))
	return nil
}

func (m *Marshaler) sliceDecoder(src, dest reflect.Value) error {
	slice := reflect.MakeSlice(dest.Type(), src.Len(), src.Len())
	for i := 0; i < src.Len(); i++ {
		if err := m.unmarshal(src.Index(i), slice.Index(i)); err != nil {
			return errs.ErrUnmarshal.Withf("cannot unmarshal slice: %w", err)
		}
	}
	dest.Set(slice)
	return nil
}

func (m *Marshaler) interfaceDecoder(src, dest reflect.Value) error {
	if dest.IsNil() {
		dest.Set(src)
	}
	return nil
}

func (m *Marshaler) mapDecoder(src, dest reflect.Value) error {
	if dest.Kind() == reflect.Ptr {
		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		dest = dest.Elem()
	}

	if dest.Kind() == reflect.Map {
		dest.Set(reflect.MakeMap(dest.Type()))
		for _, key := range src.MapKeys() {
			value := src.MapIndex(key)
			dest.SetMapIndex(key, value)
		}
	} else if dest.Kind() == reflect.Struct {
		for i := 0; i < dest.Type().NumField(); i++ {
			field := dest.Type().Field(i)
			tag := m.tagOf(field)

			mapVal := src.MapIndex(reflect.ValueOf(tag))
			if !mapVal.IsValid() {
				continue
			}

			fieldVal := dest.Field(i)
			if !fieldVal.CanSet() {
				continue
			}

			if err := m.Unmarshal(mapVal.Interface(), fieldVal.Addr().Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}
