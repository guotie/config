package config

import (
	"fmt"
	//"reflect"
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

func TestScan(t *testing.T) {
	type Is struct {
		Strin string `json:"str"`
		Iin   int64
	}
	var s struct {
		Inner Is
		Str   string
		I     int
		Slc   []int
	}
	var (
		m1 = make(map[string]int)
		m2 = make(map[string]string)
	)

	err := Scan("ts", &s)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(s)
	//println(s.Inner.Strin, s.Inner.Iin, s.Str, s.I)
	err = Scan("tm1", &m1)
	fmt.Println(m1)

	err = Scan("tm2", &m2)
	fmt.Println(m2)

}

func testScan2(t *testing.T) {
	m3 := make(map[string]map[string]int)

	err := Scan("tm3", &m3)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(m3)
}

func testScan3(t *testing.T) {
	type IE3 struct {
		Ie1 int
		Ie2 string
		Ie3 float32
	}
	m4 := make(map[string]IE3)
	err := Scan("tm4", &m4)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("m4:", m4)

	var m5 map[string]IE3
	err = Scan("tm4", &m5)
	fmt.Println("m5:", m5)
}

func testScan4(t *testing.T) {
	var s1, is1 []int
	var s2, is2 []string

	err := Scan("sl1", &s1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("s1:", s1)

	is1 = make([]int, 0)
	err = Scan("sl1", &is1)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("is1:", is1)

	err = Scan("sl2", &s2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("s2:", s2)

	is2 = make([]string, 0)
	err = Scan("sl2", &is2)
	if err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println("is2:", is2)
}

func TestScan5(t *testing.T) {
	type St struct {
		I   int
		Str string
	}
	var (
		st1 []St
		st2 []map[string]int
		st3 [][]int
	)

	Scan("ss1", &st1)
	fmt.Println("st1", st1)
	Scan("ss2", &st2)
	fmt.Println("st2", st2)
	Scan("ss3", &st3)
	fmt.Println("st3", st3)
}
