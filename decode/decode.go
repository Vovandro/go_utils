package decode

import (
	"errors"
	"reflect"
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

		default:
			return ErrorTypeMismatch
		}

	case reflect.Map:
		switch dstVal.Kind() {
		case reflect.Struct:
			return copyMapToStruct(srcVal, dstVal, tag, flag)

		case reflect.Map:
			return copyMapToMap(srcVal, dstVal, flag)

		case reflect.Interface:
			return copyMapToMap(srcVal, dstVal, flag)

		default:
			return ErrorTypeMismatch
		}

	case reflect.Slice:
		switch dstVal.Kind() {
		case reflect.Slice:
			dstVal.Set(reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Cap()))

			for i := 0; i < srcVal.Len(); i++ {
				dstVal.Index(i).Set(reflect.Indirect(srcVal.Index(i)))
			}

			return nil

		case reflect.Interface:
			dstVal.Set(reflect.MakeSlice(dstVal.Type(), srcVal.Len(), srcVal.Cap()))

			for i := 0; i < srcVal.Len(); i++ {
				dstVal.Index(i).Set(reflect.Indirect(srcVal.Index(i)))
			}

			return nil

		default:
			return ErrorTypeMismatch
		}

	default:
		if srcVal.Type().ConvertibleTo(dstVal.Type()) {
			dstVal.Set(srcVal.Convert(dstVal.Type()))
			return nil
		} else {
			return ErrorTypeMismatch
		}
	}

	return ErrorTypeMismatch
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

		if dstField.Type().ConvertibleTo(destination.Field(dstTags[sourceFieldName]).Type()) {
			destination.Field(dstTags[sourceFieldName]).Set(dstField)
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

		if dstField.Type().ConvertibleTo(destination.Field(dstTags[sourceFieldName]).Type()) {
			destination.Field(dstTags[sourceFieldName]).Set(dstField)
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
		if sourceValue.Type().ConvertibleTo(destination.Type().Elem()) {
			destination.SetMapIndex(key, sourceValue.Convert(destination.Type().Elem()))
		} else {
			return ErrorTypeMismatch
		}
	}
	return nil
}
