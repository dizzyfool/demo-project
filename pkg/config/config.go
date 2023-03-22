package config

import (
	"fmt"
	"reflect"

	"github.com/joeshaw/envdecode"
)

func Env(cfg interface{}) interface{} {
	val := reflect.ValueOf(cfg)
	if val.Kind() != reflect.Struct {
		panic(fmt.Errorf("DebugConfig must be a struct"))
	}

	new := reflect.New(val.Type())
	new.Elem().Set(val)
	if err := parse(new.Interface()); err != nil {
		panic(err)
	}
	val = new.Elem()
	values := []reflect.Value{val}
	types := []reflect.Type{val.Type()}
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		if valueField.Kind() == reflect.Struct {
			values = append(values, valueField)
			types = append(types, valueField.Type())
		}
	}
	fn := reflect.FuncOf(nil, types, false)
	return reflect.MakeFunc(fn, func(_ []reflect.Value) (results []reflect.Value) {
		return values
	}).Interface()
}

func parse(cfg interface{}) error {
	return envdecode.Decode(cfg)
}
