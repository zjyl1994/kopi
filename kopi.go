package kopi

import (
	"errors"
	"reflect"
	"unicode"
)

type TypeConvFunc func(interface{}) (interface{}, error)

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
				nameMap[opts[i].NameTo] = opts[i].NameFrom
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
		fromName := typeOfDst.Field(i).Name
		dstName := fromName
		if newFromName, ok := nameMap[fromName]; ok {
			fromName = newFromName
		}
		if value, ok := data[fromName]; ok {
			dataType := reflect.TypeOf(value)
			dstFieldType := typeOfDst.Field(i).Type
			if dataType == dstFieldType {
				dstValue := valueOfDst.Elem().FieldByName(dstName)
				if dstValue.CanSet() {
					dstValue.Set(reflect.ValueOf(value))
				}
			} else {
				if meta, ok := typeMap[dataType]; ok {
					if meta.DstType == dstFieldType {
						newValue, err := meta.ConvFunc(value)
						if err != nil {
							return err
						}
						dstValue := valueOfDst.Elem().FieldByName(dstName)
						if dstValue.CanSet() {
							dstValue.Set(reflect.ValueOf(newValue))
						}
					}
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

func T(v interface{}) reflect.Type {
	return reflect.TypeOf(v)
}

func NameOpt(from, to string) Option {
	return Option{
		NameFrom: from,
		NameTo:   to,
	}
}

func TypeOpt(from, to interface{}, conv TypeConvFunc) Option {
	return Option{
		TypeFrom:     T(from),
		TypeTo:       T(to),
		TypeConvFunc: conv,
	}
}

func NewOpt(nameFrom, nameTo string, typeFrom, typeTo interface{}, conv TypeConvFunc) Option {
	return Option{
		NameFrom:     nameFrom,
		NameTo:       nameTo,
		TypeFrom:     T(typeFrom),
		TypeTo:       T(typeTo),
		TypeConvFunc: conv,
	}
}
