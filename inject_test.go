package inject

import (
	"errors"
	"reflect"
	"testing"
)

type ComplexStruct struct {
	*LessComplexStruct
	Name         string
	Slice        []int
	privateField bool
	Interface    ComplexInterface
}

type LessComplexStruct struct {
	Interface ComplexInterface
	Name      string
}

type ComplexInterface interface {
	SomeMethod()
}
type ComplexInterfaceImplementation struct {
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
		Name:         "myname",
		Slice:        nil,
		privateField: false,
		Interface:    ComplexInterfaceImplementation{"implid"},
	}) {
		t.Fail()
	}
}

type OuterStruct struct {
	A *MediumStruct1
	B *MediumStruct2
	C *MediumStruct3
}

type MediumStruct1 struct {
	Inner *InnerStruct
}

type MediumStruct2 struct {
	Inner *InnerStruct
}

type MediumStruct3 struct {
	Inner *InnerStruct
}

type InnerStruct struct {
	initialized bool
}

func (i *InnerStruct) Init() error {
	if i.initialized {
		return errors.New("already initialized")
	}
	i.initialized = true
	return nil
}

func TestInjector_Initialize_CheckSingleInitialization(t *testing.T) {
	var toInject OuterStruct
	injector := Injector{}
	if err := injector.Initialize(&toInject); err != nil {
		t.Error(err)
	}
}

type OuterFailedStruct struct {
	Middle *MiddleFailedStruct
	Inner  *InnerFailedStruct
}

type MiddleFailedStruct struct {
	Inner *InnerFailedStruct
}

type InnerFailedStruct struct{}

var initializeCounter int

func (i *InnerFailedStruct) Init() error {
	initializeCounter++
	return errors.New("failed")
}

func TestInjector_Initialize_FailedWithoutCritical(t *testing.T) {
	var toInject OuterFailedStruct
	injector := Injector{}
	if err := injector.Initialize(&toInject); err != nil {
		t.Error(err)
	}
	if initializeCounter != 1 {
		t.Fatal("Initialize() called incorrect number of times", initializeCounter)
	}
	if toInject.Middle == nil {
		t.Fatal("Middle struct not initialized")
	}
	if toInject.Inner != nil {
		t.Fatal("First inner struct initialized")
	}
	if toInject.Middle.Inner != nil {
		t.Fatal("Second inner struct initialized")
	}
}

func TestInjector_Initialize_FailedWithCritical(t *testing.T) {
	var toInject OuterFailedStruct
	injector := Injector{
		FailOnInitializationError: true,
	}
	if err := injector.Initialize(&toInject); err == nil {
		t.Error("did not fail initialization")
	}
}
