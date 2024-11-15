package main

import "fmt"

type Test interface {
	TestLog()
}
type BaseT struct {
	i1 int
}

func (this *BaseT) TestLog() {
	fmt.Println("BaseT log")
}

type DerivedT struct {
	BaseT
	i2 int
}

func (this *DerivedT) TestLog() {
	fmt.Println("DerivedT log")
}

func main() {
	m := make(map[string]Test)
	base := BaseT{}
	derived := DerivedT{}
	m["base"] = &base
	m["derived"] = &derived
	for _, v := range m {
		v.TestLog()
	}
}
