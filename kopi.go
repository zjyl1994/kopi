package kopi

import (
	"errors"
	"reflect"
	"unicode"
)

type TypeConvFunc func(interface{}) interface{}

type Option struct {
	NameFrom string
	NameTo   string

	TypeFrom     reflect.Type
	TypeTo       reflect.Type
	TypeConvFunc TypeConvFunc
}

type typeConvMeta struct {
	DstType  reflect.Type
	ConvFunc TypeConvFunc
}

var (
	ErrTypeNotPtr    = errors.New("type not ptr")
	ErrTypeNotStruct = errors.New("type not struct")
	ErrInvalidOption = errors.New("invalid options")
)

func Kopi(dst interface{}, src interface{}, opts ...Option) error {
	nameMap := make(map[string]string)
	typeMap := make(map[reflect.Type]typeConvMeta)
	for i := 0; i < len(opts); i++ {
		if opts[i].NameFrom != "" && checkNameExported(opts[i].NameFrom) {
			if opts[i].NameTo != "" && checkNameExported(opts[i].NameTo) {
				nameMap[opts[i].NameFrom] = opts[i].NameTo
			} else {
				return ErrInvalidOption
			}
		}
		if opts[i].TypeFrom != nil {
			if opts[i].TypeTo != nil && opts[i].TypeConvFunc != nil {
				typeMap[opts[i].TypeFrom] = typeConvMeta{
					DstType:  opts[i].TypeTo,
					ConvFunc: opts[i].TypeConvFunc,
				}
			} else {
				return ErrInvalidOption
			}
		}
	}
	dataMap, err := struct2map(src)
	if err != nil {
		return err
	}
	return map2struct(dataMap, dst, nameMap, typeMap)
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

func map2struct(data map[string]interface{}, dst interface{},
	nameMap map[string]string, typeMap map[reflect.Type]typeConvMeta) error {
	typeOfDst := reflect.TypeOf(dst)
	if typeOfDst.Kind() != reflect.Ptr {
		return ErrTypeNotPtr
	}
	valueOfDst := reflect.ValueOf(dst)
	typeOfDst = valueOfDst.Elem().Type()
	for i := 0; i < typeOfDst.NumField(); i++ {
		fieldName := typeOfDst.Field(i).Name
		if value, ok := data[fieldName]; ok {
			dataType := reflect.TypeOf(value)
			dstFieldType := typeOfDst.Field(i).Type
			if dataType == dstFieldType {
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
