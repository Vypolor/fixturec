package gen

import (
	"github.com/Vypolor/without_external/mypackage1"
	"github.com/Vypolor/without_external/mypackage2"
)

type ToGen struct {
	myType1 mypackage1.MyType1
	myType2 mypackage2.MyType2
}

func New(myType1 mypackage1.MyType1, myType2 mypackage2.MyType2) *ToGen {
	return &ToGen{
		myType1: myType1,
		myType2: myType2,
	}
}

func (t *ToGen) Sum() int {
	return t.myType1.Call1() + t.myType2.Call2()
}
