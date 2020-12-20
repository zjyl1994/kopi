package main

import (
	"fmt"

	"github.com/zjyl1994/kopi"
)

type Enum uint8

type A struct {
	ID         int
	Name       string
	Value      Enum
	ExtraField string
	nonExp     string
}

type B struct {
	ID             int
	Name           string
	Value          Enum
	AnotherField   string
	notExportFIeld string
}

func main() {
	src := A{
		ID:   23,
		Name: "This is Name",
	}
	var dst B
	fmt.Printf("A:%#v\n", src)
	fmt.Printf("Kopi Result:%v\n", kopi.Kopi(&dst, src))
	fmt.Printf("B:%#v\n", dst)
}
