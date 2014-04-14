package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"
)

var (
	cfg         = make(map[string]interface{})
	KeyNotExist = errors.New("key not exist")
)

// read config file
//
func ReadCfg(fn string) {
	if fn == "" {
		fn = "./config.json"
	}
	content, err := ioutil.ReadFile(fn)
	if err != nil {
		panic(err.Error())
	}
	err = json.Unmarshal(content, &cfg)
	if err != nil {
		panic(err.Error())
	}
}

// Get
// if the key do not exist, return nil
func Get(key string) interface{} {
	if v, ok := cfg[key]; ok {
		return v
	}

	return nil
}

// GetMust
// if the key do not exist, panic
func GetMust(key string) interface{} {
	v, ok := cfg[key]
	if !ok {
		panic("Not found " + key + " in config.")
	}

	return v
}

// GetString
// Get string type by key
func GetString(key string) string {
	v := Get(key)
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}

// GetString
// Get string type by key, if not exist, panic
func GetStringMust(key string) string {
	v := Get(key)
	if v == nil {
		panic("Not found key " + key)
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}

	panic("Cannot convert key " + key + " to string.")
	return ""
}

// GetString
// Get string type by key, if key not exist, return the default
func GetStringDefault(key, def string) string {
	s := GetString(key)
	if s == "" {
		return def
	}
	return s
}

// GetInt
func GetInt(key string) (int64, bool) {
	v := Get(key)
	if v == nil {
		return 0, false
	}
	if i, ok := v.(float64); ok {
		return int64(i), true
	}

	return 0, false
}

func GetInt64Default(key string, dv int64) int64 {
	v := Get(key)
	if v == nil {
		return dv
	}
	if i, ok := v.(float64); ok {
		return int64(i)
	}

	return dv
}

func GetIntDefault(key string, dv int) int {
	v := Get(key)
	if v == nil {
		return dv
	}
	if i, ok := v.(float64); ok {
		return int(i)
	}

	return dv
}

func GetFloat(key string) (float64, bool) {
	v := Get(key)
	if v == nil {
		return 0, false
	}
	if i, ok := v.(float64); ok {
		return i, true
	}

	return 0.0, false
}

// Scan
//
func Scan(key string, dest interface{}) (err error) {
	val, ok := cfg[key]
	if !ok {
		return KeyNotExist
	}

	rv := reflect.ValueOf(dest)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("param must be pointer and not nil")
	}

	switch reflect.Indirect(rv).Kind() {
	case reflect.Map:
		err = decodeMap(val, rv)

	case reflect.Struct:
		err = decodeStruct(val, rv)

	case reflect.Slice:
		err = decodeSlice(val, rv)

	case reflect.Array:
		err = decodeArray(val, rv)

	default:
		return fmt.Errorf("Unsupport type: %s", rv.Kind().String())
	}

	return
}

func decodeMap(val interface{}, rv reflect.Value) (err error) {
	mval := val.(map[string]interface{})
	//fmt.Println("decode map:", mval)

	rv = reflect.Indirect(rv)
	if rv.IsNil() {
		//rv = reflect.MakeMap(rv.Type())
		rv.Set(reflect.MakeMap(rv.Type()))
	}

	keytyp := rv.Type().Key().Kind()
	if keytyp != reflect.String {
		return fmt.Errorf("map key's type must be string: %s",
			keytyp.String())
	}

	eletyp := rv.Type().Elem().Kind()
	//fmt.Println(rv.Type().String(), eletyp.String(), rv.Type().Elem())

	for key, valfd := range mval {
		switch eletyp {
		case reflect.Map:
			tmp := reflect.MakeMap(rv.Type().Elem())
			err = decodeMap(valfd, tmp)
			//fmt.Println("map entry:", tmp.Interface())
			rv.SetMapIndex(reflect.ValueOf(key),
				reflect.ValueOf(tmp.Interface()))

		case reflect.Struct:
			tmp := reflect.New(rv.Type().Elem())
			err = decodeStruct(valfd, tmp)
			rv.SetMapIndex(reflect.ValueOf(key),
				reflect.ValueOf(reflect.Indirect(tmp).Interface()))

		case reflect.Slice:
			tmp := reflect.MakeSlice(rv.Type().Elem(), 0, 0)
			err = decodeSlice(valfd, tmp)
			rv.SetMapIndex(reflect.ValueOf(key),
				reflect.ValueOf(reflect.Indirect(tmp).Interface()))

		case reflect.Array:
			tmp := reflect.New(reflect.TypeOf(eletyp))
			err = decodeArray(valfd, tmp)
			rv.SetMapIndex(reflect.ValueOf(key),
				reflect.ValueOf(reflect.Indirect(tmp).Interface()))

		case reflect.Int:
			tmp := int(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Int16:
			tmp := int16(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Int32:
			tmp := int32(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Int64:
			tmp := int64(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Int8:
			tmp := int8(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Float32:
			tmp := float32(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Float64:
			tmp := valfd.(float64)
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Uint:
			tmp := uint64(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Uint16:
			tmp := uint64(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Uint32:
			tmp := uint64(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Uint64:
			tmp := uint64(valfd.(float64))
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.String:
			tmp := valfd.(string)
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))

		case reflect.Bool:
			tmp := valfd.(bool)
			rv.SetMapIndex(reflect.ValueOf(key), reflect.ValueOf(tmp))
		}
		if err != nil {
			return err
		}
	}

	//fmt.Println("decode result:", rv.Interface())
	return nil
}

func decodeStruct(val interface{}, rv reflect.Value) (err error) {
	mval := val.(map[string]interface{})

	if rv.Kind() == reflect.Ptr && rv.IsNil() {
		rv.Set(reflect.New(rv.Type()))
	}

	rv = reflect.Indirect(rv)

	typ := rv.Type()

	var valfd interface{}
	for i := 0; i < rv.NumField(); i++ {
		fd := rv.Field(i)
		tag := typ.Field(i).Tag.Get("json")
		if tag == "-" {
			continue
		}
		if tag != "" {
			valfd = mval[tag]
		} else {
			valfd = mval[typ.Field(i).Name]
			if valfd == nil {
				valfd = mval[strings.ToLower(typ.Field(i).Name)]
				if valfd == nil {
					continue
				}
			}
		}

		switch fd.Kind() {
		case reflect.Map:
			err = decodeMap(valfd, fd)
		case reflect.Struct:
			err = decodeStruct(valfd, fd)
		case reflect.Slice:
			err = decodeSlice(valfd, fd)
		case reflect.Array:
			err = decodeArray(valfd, fd)

		case reflect.Int:
			fallthrough
		case reflect.Int16:
			fallthrough
		case reflect.Int32:
			fallthrough
		case reflect.Int64:
			fallthrough
		case reflect.Int8:
			fallthrough
		case reflect.String:
			fallthrough
		case reflect.Float32:
			fallthrough
		case reflect.Float64:
			fallthrough
		case reflect.Uint:
			fallthrough
		case reflect.Uint16:
			fallthrough
		case reflect.Uint32:
			fallthrough
		case reflect.Uint64:
			fallthrough
		case reflect.Bool:
			err = decodePrimary(valfd, fd)
		default:
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func decodeSlice(val interface{}, rv reflect.Value) error {
	sval := val.([]interface{})

	rv = reflect.Indirect(rv)
	if rv.IsNil() {
		//fmt.Println("decodeSlice: nil slice, create new slice.", rv.Type().String())
		rv.Set(reflect.New(rv.Type()).Elem())
		//rv = reflect.MakeSlice(rv.Type(), 0, 0)
	}

	eletyp := rv.Type().Elem().Kind()
	//fmt.Println("decode slice:", rv.Type().String(), eletyp.String(), rv.Type().Elem())

	var (
		err error
		tmp reflect.Value
	)
	for i, valfd := range sval {
		if i >= rv.Cap() {
			newcap := rv.Cap() + rv.Cap()/2
			if newcap < 4 {
				newcap = 4
			}
			newv := reflect.MakeSlice(rv.Type(), rv.Len(), newcap)
			reflect.Copy(newv, rv)
			rv.Set(newv)
		}
		rv.SetLen(i + 1)

		switch eletyp {
		case reflect.Map:
			//reflect.Make
			tmp = reflect.MakeMap(rv.Type().Elem())
			err = decodeMap(valfd, tmp)
			rv.Index(i).Set(reflect.ValueOf(tmp.Interface()))

		case reflect.Struct:
			tmp = reflect.New(rv.Type().Elem())
			err = decodeStruct(valfd, tmp)
			rv.Index(i).Set(reflect.ValueOf(tmp.Elem().Interface()))

		case reflect.Slice:
			tmp = reflect.New(rv.Type().Elem())
			err = decodeSlice(valfd, tmp)
			rv.Index(i).Set(reflect.ValueOf(tmp.Elem().Interface()))

		case reflect.Array:
			tmp = reflect.New(reflect.TypeOf(eletyp))
			err = decodeArray(valfd, tmp)
			rv.Index(i).Set(reflect.ValueOf(tmp.Interface()))

		case reflect.Int:
			pv := int(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Int16:
			pv := int16(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Int32:
			pv := int32(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Int64:
			pv := int64(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Int8:
			pv := int8(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Float32:
			pv := float32(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Float64:
			pv := valfd.(float64)
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Uint:
			pv := uint64(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Uint16:
			pv := uint64(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Uint32:
			pv := uint64(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Uint64:
			pv := uint64(valfd.(float64))
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.String:
			pv := valfd.(string)
			rv.Index(i).Set(reflect.ValueOf(pv))

		case reflect.Bool:
			pv := valfd.(bool)
			rv.Index(i).Set(reflect.ValueOf(pv))
		}

		if err != nil {
			return err
		}
		//rv = reflect.Append(rv,
		//	reflect.ValueOf(tmp.Interface()))
	}

	//fmt.Println("decode slice result:", rv.Interface())
	return nil
}

func decodeArray(val interface{}, rv reflect.Value) error {
	panic("decode Array Not implement!")
	return nil
}

func decodePrimary(val interface{}, rv reflect.Value) error {
	irv := reflect.Indirect(rv)
	switch irv.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Int8:
		rv.SetInt(int64(val.(float64)))

	case reflect.Uint:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		rv.SetUint(uint64(val.(float64)))

	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		rv.SetFloat(val.(float64))

	case reflect.String:
		rv.SetString(val.(string))

	case reflect.Bool:
		rv.SetBool(val.(bool))

	default:
		panic("No type: " + irv.Kind().String())
	}

	return nil
}
