package inject

import (
	"reflect"
	"testing"
)

type ComplexStruct struct {
	*LessComplexStruct
	Name string
	Slice []int
	privateField bool
	Interface ComplexInterface
}

type LessComplexStruct struct {
	Interface ComplexInterface
	Name string
}

type ComplexInterface interface {
	SomeMethod()
}
type ComplexInterfaceImplementation struct{
	Id string
}

func (c ComplexInterfaceImplementation) SomeMethod() {}

func TestInjector_Initialize(t *testing.T) {
	var toInject ComplexStruct
	injector := Injector{
		InjectableValues: []interface{}{
			"myname", ComplexInterfaceImplementation{"implid"},
		},
	}
	injector.Initialize(&toInject)
	if !reflect.DeepEqual(toInject, ComplexStruct{
		LessComplexStruct: &LessComplexStruct{
			Interface: ComplexInterfaceImplementation{"implid"},
			Name:      "myname",
		},
		Name:              "myname",
		Slice:             nil,
		privateField:      false,
		Interface:         ComplexInterfaceImplementation{"implid"},
	}) {
		t.Fail()
	}
}