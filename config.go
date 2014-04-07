package config

import (
	"encoding/json"
	"io/ioutil"
)

var (
	cfg = make(map[string]interface{})
)

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

func Get(key string) interface{} {
	if v, ok := cfg[key]; ok {
		return v
	}

	return nil
}

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

func GetStringDefault(key, def string) string {
	s := GetString(key)
	if s == "" {
		return def
	}
	return s
}

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

func GetBoolean(key string) (bool, bool) {
	v := Get(key)
	if v == nil {
		return false, false
	}

	if i, ok := v.(bool); ok {
		return i, true
	}

	return false, false
}

func GetBooleanDefault(key string, dv bool) bool {
	v := Get(key)
	if v == nil {
		return dv
	}
	if i, ok := v.(bool); ok {
		return i
	}

	return dv
}
