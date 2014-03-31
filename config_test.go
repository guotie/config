package config

import (
	"testing"
)

func init() {
	ReadCfg("./test.json")
}

func TestString(t *testing.T) {
	if GetString("redisProto") != "tcp" {
		t.Fatal("get redisProto failed!")
	}
}

func TestInt(t *testing.T) {
	maxNews, ok := GetInt("maxNews")
	if (!ok) || (maxNews != 5000) {
		t.Fatal("get maxnews failed:", ok, maxNews)
	}
}

func TestFloat(t *testing.T) {
	ft, ok := GetFloat("floatTest")
	if (!ok) || (ft != 1.67) {
		t.Fatal("get float failed:", ok, ft)
	}
}
