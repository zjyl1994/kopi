package kopi

import (
	"errors"
	"reflect"
	"unicode"
)

type KopiOption struct {
}

var (
	ErrTypeNotPtr    = errors.New("type not ptr")
	ErrTypeNotStruct = errors.New("type not struct")
)

func Kopi(dst interface{}, src interface{}, opts ...KopiOption) error {
	dataMap, err := struct2map(src)
	if err != nil {
		return err
	}
	return map2struct(dataMap, dst)
}

func struct2map(v interface{}) (data map[string]interface{}, err error) {
	data = make(map[string]interface{})
	typeOfStruct := reflect.TypeOf(v)
	if typeOfStruct.Kind() != reflect.Struct {
		return nil, ErrTypeNotStruct
	}
	valueOfStruct := reflect.ValueOf(v)
	for i := 0; i < typeOfStruct.NumField(); i++ {
		fieldName := typeOfStruct.Field(i).Name
		if checkNameExported(fieldName) {
			data[fieldName] = valueOfStruct.FieldByName(fieldName).Interface()
		}
	}
	return data, nil
}

func map2struct(data map[string]interface{}, dst interface{}) error {
	typeOfDst := reflect.TypeOf(dst)
	if typeOfDst.Kind() != reflect.Ptr {
		return ErrTypeNotPtr
	}
	valueOfDst := reflect.ValueOf(dst)
	typeOfDst = valueOfDst.Elem().Type()
	for i := 0; i < typeOfDst.NumField(); i++ {
		fieldName := typeOfDst.Field(i).Name
		if value, ok := data[fieldName]; ok {
			if reflect.TypeOf(value) == typeOfDst.Field(i).Type {
				dstValue := valueOfDst.Elem().FieldByName(fieldName)
				if dstValue.CanSet() {
					dstValue.Set(reflect.ValueOf(value))
				}
			}
		}
	}
	return nil
}

func checkNameExported(name string) bool {
	if name == "" {
		return false
	}
	runes := []rune(name)
	return unicode.IsUpper(runes[0])
}
