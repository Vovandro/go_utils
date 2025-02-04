package decode

import (
	"errors"
	"reflect"
	"strconv"
	"time"
)

var (
	ErrorDstNotFound  = errors.New("destination not found")
	ErrorDstNotSet    = errors.New("destination not set")
	ErrorTypeMismatch = errors.New("type mismatch")
)

type DecoderFlag int

const (
	DecoderStrongFoundDst    DecoderFlag = 0x1 << iota // Error if not found destination
	DecoderStrongType                                  // Safe source type or error. Explode inner struct to map in map to map
	DecoderUnwrapStructToMap                           // Unwrap struct to map
)

// Decode копирует данные из источника в назначение, поддерживая различные типы данных
// (структуры, мапы) и их вложенность, используя теги для сопоставления полей.
func Decode(source interface{}, destination interface{}, tag string, flag DecoderFlag) error {
	var sourceVal reflect.Value
	var dstVal reflect.Value

	if v, ok := source.(reflect.Value); ok {
		sourceVal = v
	} else {
		sourceVal = reflect.ValueOf(source)
	}

	if v, ok := destination.(reflect.Value); ok {
		dstVal = v
	} else {
		dstVal = reflect.ValueOf(destination)
	}

	// Проверка, что назначение является указателем и доступно для записи
	if dstVal.Kind() != reflect.Ptr || dstVal.IsNil() {
		return errors.New("destination must be a non-nil pointer")
	}
	dstVal = dstVal.Elem()

	sourceVal = reflect.Indirect(sourceVal)

	return copyValues(sourceVal, dstVal, tag, flag)
}

func copyValues(source reflect.Value, destination reflect.Value, tag string, flag DecoderFlag) error {
	var srcVal reflect.Value
	var dstVal reflect.Value
	sourceIsInterface := source.CanInterface()

	if sourceIsInterface {
		srcVal = reflect.ValueOf(source.Interface())
	} else {
		srcVal = source
	}

	dstVal = destination

	if srcVal.CanAddr() && srcVal.IsNil() {
		zero := reflect.Zero(destination.Type())
		dstVal.Set(zero)
		return nil
	}

	sourceIsPtr := srcVal.Kind() == reflect.Ptr
	dstIsPtr := dstVal.Kind() == reflect.Ptr

	if sourceIsPtr {
		srcVal = srcVal.Elem()
	}

	if dstIsPtr {
		dstVal = dstVal.Elem()
	}

	switch srcVal.Kind() {
	case reflect.Struct:
		switch dstVal.Kind() {
		case reflect.Struct:
			return copyStructToStruct(srcVal, dstVal, tag, flag)

		case reflect.Map:
			return copyStructToMap(srcVal, dstVal, tag, flag)

		case reflect.Interface:
			return copyStructToStruct(srcVal, dstVal, tag, flag)
		}

	case reflect.Map:
		switch dstVal.Kind() {
		case reflect.Struct:
			return copyMapToStruct(srcVal, dstVal, tag, flag)

		case reflect.Map:
			return copyMapToMap(srcVal, dstVal, flag)

		case reflect.Interface:
			return copyMapToMap(srcVal, dstVal, flag)
		}

	case reflect.Slice:
		switch dstVal.Kind() {
		case reflect.Slice:
			dstVal.Set(reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Cap()))

			for i := 0; i < srcVal.Len(); i++ {
				val := srcVal.Index(i)
				for val.Kind() == reflect.Interface {
					val = val.Elem()
				}

				val = reflect.Indirect(val)

				if dstVal.Type().Elem() == val.Type() {
					dstVal.Index(i).Set(val.Convert(dstVal.Type().Elem()))
				} else {
					if flag&DecoderStrongType != 0 {
						return ErrorTypeMismatch
					}

					if converted, err := convertBasicTypes(val, dstVal.Type().Elem()); err == nil {
						dstVal.Index(i).Set(converted)
					} else {
						return err
					}
				}
			}

			return nil

		case reflect.Interface:
			dstVal.Set(reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Cap()))

			for i := 0; i < srcVal.Len(); i++ {
				dstVal.Index(i).Set(reflect.Indirect(srcVal.Index(i)))
			}

			return nil
		}

	default:
		if srcVal.Kind() == dstVal.Kind() {
			dstVal.Set(srcVal.Convert(dstVal.Type()))
			return nil
		} else if flag&DecoderStrongType == 0 {
			if converted, err := convertBasicTypes(srcVal, dstVal.Type()); err == nil {
				dstVal.Set(converted)
				return nil
			} else {
				return ErrorTypeMismatch
			}
		} else {
			return ErrorTypeMismatch
		}
	}

	return ErrorTypeMismatch
}

func convertBasicTypes(source reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	if source.Kind() == reflect.Interface {
		source = source.Elem()
	}

	switch targetType.Kind() {
	case reflect.Interface:
		return source, nil

	case reflect.String:
		switch source.Kind() {
		case reflect.String:
			return reflect.ValueOf(source.String()), nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(strconv.FormatInt(source.Int(), 10)), nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(strconv.FormatFloat(source.Float(), 'f', -1, 64)), nil
		case reflect.Bool:
			return reflect.ValueOf(strconv.FormatBool(source.Bool())), nil
		}
		return reflect.Value{}, ErrorTypeMismatch

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if _, ok := targetType.MethodByName("Nanoseconds"); ok {
			switch source.Kind() {
			case reflect.String:
				d, _ := time.ParseDuration(source.String())
				return reflect.ValueOf(d).Convert(targetType), nil
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return reflect.ValueOf(source.Int()).Convert(targetType), nil
			case reflect.Float32, reflect.Float64:
				return reflect.ValueOf(source.Float()).Convert(targetType), nil
			}
		} else {
			switch source.Kind() {
			case reflect.String:
				newValue := reflect.New(targetType)

				if m := newValue.Elem().MethodByName("Parse"); m.IsValid() && m.Type().NumIn() == 1 && m.Type().In(0).Kind() == reflect.String &&
					m.Type().NumOut() == 1 && m.Type().Out(0) == targetType {
					res := newValue.Elem().MethodByName("Parse").Call([]reflect.Value{reflect.ValueOf(source.String())})
					return res[0], nil
				} else if m := newValue.MethodByName("Parse"); m.IsValid() && m.Type().NumIn() == 1 && m.Type().In(0).Kind() == reflect.String {
					m.Call([]reflect.Value{reflect.ValueOf(source.String())})
					return newValue.Elem(), nil
				} else {
					if intValue, err := strconv.ParseInt(source.String(), 10, 64); err == nil {
						return reflect.ValueOf(intValue).Convert(targetType), nil
					} else {
						return reflect.Value{}, err
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				return reflect.ValueOf(source.Int()).Convert(targetType), nil
			case reflect.Float32, reflect.Float64:
				return reflect.ValueOf(source.Float()).Convert(targetType), nil
			case reflect.Bool:
				return reflect.ValueOf(source.Int()).Convert(targetType), nil
			}
		}
		return reflect.Value{}, ErrorTypeMismatch

	case reflect.Float32, reflect.Float64:
		switch source.Kind() {
		case reflect.String:
			newValue := reflect.New(targetType)

			if m := newValue.Elem().MethodByName("Parse"); m.IsValid() && m.Type().NumIn() == 1 && m.Type().In(0).Kind() == reflect.String &&
				m.Type().NumOut() == 1 && m.Type().Out(0) == targetType {
				res := newValue.Elem().MethodByName("Parse").Call([]reflect.Value{reflect.ValueOf(source.String())})
				return res[0], nil
			} else if m := newValue.MethodByName("Parse"); m.IsValid() && m.Type().NumIn() == 1 && m.Type().In(0).Kind() == reflect.String {
				m.Call([]reflect.Value{reflect.ValueOf(source.String())})
				return newValue.Elem(), nil
			} else {
				if intValue, err := strconv.ParseFloat(source.String(), 64); err == nil {
					return reflect.ValueOf(intValue).Convert(targetType), nil
				} else {
					return newValue.Elem(), err
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(source.Int()).Convert(targetType), nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(source.Float()).Convert(targetType), nil
		case reflect.Bool:
			return reflect.ValueOf(source.Int()).Convert(targetType), nil
		}
		return reflect.Value{}, ErrorTypeMismatch

	case reflect.Bool:
		switch source.Kind() {
		case reflect.String:
			if intValue, err := strconv.ParseBool(source.String()); err == nil {
				return reflect.ValueOf(intValue).Convert(targetType), nil
			} else {
				return reflect.Value{}, err
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return reflect.ValueOf(source.Int()).Convert(targetType), nil
		case reflect.Float32, reflect.Float64:
			return reflect.ValueOf(source.Float()).Convert(targetType), nil
		case reflect.Bool:
			return reflect.ValueOf(source.Bool()).Convert(targetType), nil
		}
		return reflect.Value{}, ErrorTypeMismatch
	}

	return reflect.Value{}, ErrorTypeMismatch
}

func copyStructToMap(source reflect.Value, destination reflect.Value, tag string, flag DecoderFlag) error {
	if destination.IsNil() {
		destination.Set(reflect.MakeMap(destination.Type()))
	}

	typeOfSource := source.Type()
	for i := 0; i < source.NumField(); i++ {
		srcField := source.Field(i)
		var fieldName string

		if tag != "" {
			fieldName = typeOfSource.Field(i).Tag.Get(tag)
			if fieldName == "" {
				continue
			}
		} else {
			fieldName = typeOfSource.Field(i).Name
		}

		var data reflect.Value
		if destination.Type().Elem().Kind() == reflect.Interface {
			if flag&DecoderUnwrapStructToMap != 0 && srcField.Kind() == reflect.Struct {
				data = reflect.MakeMap(reflect.TypeOf(map[string]interface{}{}))
			} else if srcField.Kind() == reflect.Interface {
				if reflect.ValueOf(srcField.Interface()).Kind() == reflect.Ptr {
					data = reflect.New(reflect.TypeOf(srcField.Interface()).Elem())
				} else {
					data = reflect.New(reflect.TypeOf(srcField.Interface())).Elem()
				}
			} else if srcField.Kind() == reflect.Ptr {
				data = reflect.New(srcField.Elem().Type())
			} else {
				data = reflect.New(srcField.Type()).Elem()
			}
		} else {
			data = reflect.New(destination.Type().Elem()).Elem()
		}

		err := copyValues(srcField, data, tag, flag)
		if err != nil {
			return err
		}

		destination.SetMapIndex(reflect.ValueOf(fieldName), data)
	}
	return nil
}

func copyStructToStruct(source reflect.Value, destination reflect.Value, tag string, flag DecoderFlag) error {
	sourceType := source.Type()
	dstType := destination.Type()

	dstTags := make(map[string]int)

	for i := 0; i < destination.NumField(); i++ {
		if tag != "" {
			fieldTag := dstType.Field(i).Tag.Get(tag)
			if fieldTag == "" {
				continue
			}

			dstTags[fieldTag] = i
		} else {
			dstTags[dstType.Field(i).Name] = i
		}
	}

	for i := 0; i < source.NumField(); i++ {
		srcField := source.Field(i)

		var sourceFieldName string

		if tag != "" {
			sourceFieldName = sourceType.Field(i).Tag.Get(tag)
			if sourceFieldName == "" {
				continue
			}
		} else {
			sourceFieldName = sourceType.Field(i).Name
		}

		if _, ok := dstTags[sourceFieldName]; !ok {
			if flag&DecoderStrongFoundDst != 0 {
				return ErrorDstNotFound
			}
			continue
		}

		dstField := destination.Field(dstTags[sourceFieldName])
		if !dstField.IsValid() || !dstField.CanSet() {
			return ErrorDstNotSet
		}

		if dstField.Kind() == reflect.Interface {
			if flag&DecoderUnwrapStructToMap != 0 && srcField.Kind() == reflect.Struct {
				dstField = reflect.MakeMap(reflect.TypeOf(map[string]interface{}{}))
			} else if srcField.Kind() == reflect.Interface {
				if reflect.ValueOf(srcField.Interface()).Kind() == reflect.Ptr {
					dstField = reflect.New(reflect.TypeOf(srcField.Interface()).Elem())
				} else {
					dstField = reflect.New(reflect.TypeOf(srcField.Interface())).Elem()
				}
			} else if srcField.Kind() == reflect.Ptr {
				dstField = reflect.New(srcField.Elem().Type())
			} else {
				dstField = reflect.New(srcField.Type()).Elem()
			}
		}

		err := copyValues(srcField, dstField, tag, flag)
		if err != nil {
			return err
		}

		if dstField.Kind() == destination.Field(dstTags[sourceFieldName]).Kind() {
			destination.Field(dstTags[sourceFieldName]).Set(dstField)
		} else if flag&DecoderStrongType == 0 {
			if converted, err := convertBasicTypes(dstField, destination.Field(dstTags[sourceFieldName]).Type()); err == nil {
				destination.Field(dstTags[sourceFieldName]).Set(converted)
				return nil
			} else {
				return ErrorTypeMismatch
			}
		} else {
			return ErrorTypeMismatch
		}
	}
	return nil
}

func copyMapToStruct(source reflect.Value, destination reflect.Value, tag string, flag DecoderFlag) error {
	dstType := destination.Type()
	dstTags := make(map[string]int)

	for i := 0; i < destination.NumField(); i++ {
		if tag != "" {
			fieldTag := dstType.Field(i).Tag.Get(tag)
			if fieldTag == "" {
				continue
			}

			dstTags[fieldTag] = i
		} else {
			dstTags[dstType.Field(i).Name] = i
		}
	}
	for _, key := range source.MapKeys() {
		sourceFieldName := key.String()

		if _, ok := dstTags[sourceFieldName]; !ok {
			if flag&DecoderStrongFoundDst != 0 {
				return ErrorDstNotFound
			}
			continue
		}

		srcField := source.MapIndex(key)

		dstField := destination.Field(dstTags[sourceFieldName])
		if !dstField.IsValid() || !dstField.CanSet() {
			return ErrorDstNotSet
		}

		if dstField.Kind() == reflect.Interface {
			if flag&DecoderUnwrapStructToMap != 0 && srcField.Kind() == reflect.Struct {
				dstField = reflect.MakeMap(reflect.TypeOf(map[string]interface{}{}))
			} else if srcField.Kind() == reflect.Interface {
				if reflect.ValueOf(srcField.Interface()).Kind() == reflect.Ptr {
					dstField = reflect.New(reflect.TypeOf(srcField.Interface()).Elem())
				} else {
					dstField = reflect.New(reflect.TypeOf(srcField.Interface())).Elem()
				}
			} else if srcField.Kind() == reflect.Ptr {
				dstField = reflect.New(srcField.Elem().Type())
			} else {
				dstField = reflect.New(srcField.Type()).Elem()
			}
		}

		err := copyValues(srcField, dstField, tag, flag)
		if err != nil {
			return err
		}

		if dstField.Kind() == destination.Field(dstTags[sourceFieldName]).Kind() {
			destination.Field(dstTags[sourceFieldName]).Set(dstField)
		} else if flag&DecoderStrongType == 0 {
			if converted, err := convertBasicTypes(dstField, destination.Field(dstTags[sourceFieldName]).Type()); err == nil {
				destination.Field(dstTags[sourceFieldName]).Set(converted)
				return nil
			} else {
				return ErrorTypeMismatch
			}
		} else {
			return ErrorTypeMismatch
		}
	}
	return nil
}

func copyMapToMap(source reflect.Value, destination reflect.Value, flag DecoderFlag) error {
	if destination.IsNil() {
		destination.Set(reflect.MakeMap(destination.Type()))
	}

	for _, key := range source.MapKeys() {
		sourceValue := source.MapIndex(key)
		if sourceValue.Kind() == destination.Type().Elem().Kind() {
			destination.SetMapIndex(key, sourceValue.Convert(destination.Type().Elem()))
		} else if flag&DecoderStrongType == 0 {
			if converted, err := convertBasicTypes(sourceValue, destination.Type().Elem()); err == nil {
				destination.SetMapIndex(key, converted)
				return nil
			} else {
				return ErrorTypeMismatch
			}
		} else {
			return ErrorTypeMismatch
		}
	}
	return nil
}
