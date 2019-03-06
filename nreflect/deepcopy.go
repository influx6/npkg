package nreflect

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/gokit/npkg/nerror"
	"github.com/jinzhu/copier"
)

// DeepCopy copies giving value returning a new distinct value,
// else an error due to failure.
func DeepCopy(elem interface{}) (interface{}, error) {
	var baseValue = IndirectValue(ValueOf(elem))

	var baseType reflect.Type
	var newType reflect.Value

	switch baseValue.Kind() {
	case reflect.Map:
		baseType = reflect.MapOf(baseValue.Type(), baseValue.Elem().Type())
		newType = reflect.MakeMap(baseType)
	case reflect.Array:
		baseType = baseValue.Type()
		newType = reflect.MakeSlice(baseValue.Type(), baseValue.Len(), baseValue.Cap())
		fmt.Printf("Array: %#v -> %#v =? %q\n", newType, baseValue, newType.Kind())
	case reflect.Slice:
		baseType = baseValue.Type()
		newType = reflect.MakeSlice(baseValue.Type(), baseValue.Len(), baseValue.Cap())
		fmt.Printf("Slice: %#v -> %#v = %q\n", baseValue, newType, newType.Kind())
	default:
		baseType = TypeOf(elem)
		newType = reflect.New(baseType)
		fmt.Printf("Def: %#v -> %#v = %q \n", baseValue, newType, baseType)
	}

	//newType = IndirectValue(newType)

	switch baseValue.Kind() {
	case reflect.Map:
		if err := copyMap(newType, baseValue); err != nil {
			return nil, err
		}
	case reflect.Struct:
		copier.Copy(newType.Interface(), elem)
	case reflect.UnsafePointer:
	case reflect.Slice:
		reflect.Copy(newType, baseValue)
	case reflect.String:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		fmt.Printf("Value: %#v -> %#v \n", newType, baseValue)
		newType.Set(baseValue)
	case reflect.Float64, reflect.Float32:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		newType.Set(baseValue)
	case reflect.Bool:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		newType.Set(baseValue)
	case reflect.Chan:
		return nil, nerror.New("unable to copy channel type")
	case reflect.Array:
		reflect.Copy(newType, baseValue)
	case reflect.Complex64:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		newType.Set(baseValue)
	case reflect.Complex128:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		newType.Set(baseValue)
	case reflect.Uintptr:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		newType.Set(baseValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		newType.Set(baseValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if newType.CanSet() {
			return nil, nerror.New("unable to set value for type")
		}
		newType.Set(baseValue)
	}
	return newType.Interface(), nil
}

func copyMap(dest reflect.Value, src reflect.Value) error {
	if dest.Kind() != reflect.Map && src.Kind() != reflect.Map {
		return nerror.New("unable to copy in between map and another type")
	}

	var keys = dest.MapKeys()
	for _, key := range keys {
		newKey, err := DeepCopy(key)
		if err != nil {
			return err
		}

		newValue, err := DeepCopy(dest.MapIndex(key))
		if err != nil {
			return err
		}
		src.SetMapIndex(ValueOf(newKey), ValueOf(newValue))
	}
	return nil
}

// func IndirectValue returns the reflect.Value if the provided value is a pointer kind.
func IndirectValue(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

// func IndirectType returns the reflect.Type if the provided value is a pointer kind.
func IndirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

// ValueOf returns giving value of elem, if elem is a reflect.Type returns
// a new reflect.Value for giving type.
func ValueOf(elem interface{}) reflect.Value {
	if telem, ok := elem.(reflect.Value); ok {
		return telem
	}
	if telem, ok := elem.(reflect.Type); ok {
		return reflect.New(telem)
	}
	return reflect.ValueOf(elem)
}

// TypeOf returns giving type of elem.
func TypeOf(elem interface{}) reflect.Type {
	if telem, ok := elem.(reflect.Type); ok {
		return telem
	}
	if telem, ok := elem.(reflect.Value); ok {
		return telem.Type()
	}
	return reflect.TypeOf(elem)
}

// CopyStringMap returns a new copy of a giving string map.
func CopyStringMap(src map[string]string) map[string]string {
	var dest = make(map[string]string, len(src))
	for key, value := range src {
		dest[key] = value
	}
	return dest
}

// CopyStringKeyMap returns a new copy of a giving string keyed map.
func copyStringKeyMap(src map[string]interface{}) (map[string]interface{}, error) {
	var dest = make(map[string]interface{}, len(src))
	for key, value := range src {
		var myValue, err = DeepCopy(value)
		if err != nil {
			return dest, err
		}
		dest[key] = myValue
	}
	return dest, nil
}

// CopyInterfaceKeyMap returns a new copy of a giving interface keyed map.
func CopyInterfaceKeyMap(src map[interface{}]interface{}) (map[interface{}]interface{}, error) {
	var dest = make(map[interface{}]interface{}, len(src))
	for key, value := range src {
		myValue, err := DeepCopy(value)
		if err != nil {
			return dest, err
		}
		myKey, err := DeepCopy(key)
		if err != nil {
			return dest, err
		}
		dest[myKey] = myValue
	}
	return dest, nil
}

// CopyBytes returns a new copy of giving byte slice.
func CopyBytes(bu []byte) []byte {
	var cu = make([]byte, len(bu))
	copy(cu, bu)
	return cu
}

//*****************************************************
// unsafe methods
//*****************************************************

// byte2String converts a byte slice into a string.
func bytes2String(bc []byte) string {
	return *(*string)(unsafe.Pointer(&bc))
}

// string2Bytes converts a string into a byte slice.
func string2Bytes(bc string) []byte {
	return *(*[]byte)(unsafe.Pointer(&bc))
}
