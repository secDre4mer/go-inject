package inject

import (
	"fmt"
	"reflect"
)

type Injector struct {
	InjectableValues   []interface{}
	FailOnUninjectable bool
	FailOnInitializationError bool
}

type Initializable interface {
	Init() error
}

func (i *Injector) Initialize(object interface{}) error {
	value := reflect.ValueOf(object)
	if value.Kind() != reflect.Ptr {
		panic("Must pass pointer to injectable object to injector")
	}
	i.initializeObject(value.Elem())
	return nil
}

func (i *Injector) initializeObject(object reflect.Value) error {
	// Check if object itself is a fixed injectable value
	for _, injectable := range i.InjectableValues {
		if reflect.TypeOf(injectable).AssignableTo(object.Type()) {
			object.Set(reflect.ValueOf(injectable))
			return nil
		}
	}
	// We have no fixed value for Object
	switch object.Kind() {
	case reflect.Ptr:
		newObject := reflect.New(object.Type().Elem())
		if err := i.initializeObject(newObject.Elem()); err != nil {
			return err
		}
		var initializationFailed bool
		if initializable, isInitializable := newObject.Interface().(Initializable); isInitializable {
			if err := initializable.Init(); err != nil {
				if i.FailOnInitializationError {
					return err
				} else {
					initializationFailed = true
				}
			}
		}
		if !initializationFailed {
			object.Set(newObject)
			i.InjectableValues = append(i.InjectableValues, newObject.Interface())
		}
		return nil
	case reflect.Struct:
		for fi := 0; fi < object.NumField(); fi++ {
			fieldType := object.Type().Field(fi)
			if fieldType.PkgPath != "" { // unexported field
				continue
			}
			if err := i.initializeObject(object.Field(fi)); err != nil {
				return fmt.Errorf("could not inject object of type %v: %w", object.Type().String(), err)
			}
		}
		return nil
	}
	if i.FailOnUninjectable {
		return fmt.Errorf("could not inject object of type %v", object.Type().String())
	}
	return nil
}
