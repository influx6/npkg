package nreflect

import (
	"reflect"
	"unsafe"

	"github.com/influx6/npkg/nerror"
	"github.com/jinzhu/copier"
)

// DeepCopy copies giving value returning a new distinct value,
// else an error due to failure.
//
// DeepCopy will never attempt to duplicate the following types, as it will simply
// return the value it received.
//
// 1. Channels
// 2. Interface
// 3. unsafe.Pointer
// 4. Func
//
// For structs, we use github.com/jinzhu/copier for copying, so all rules in that library
// applies here also.
func DeepCopy(elem interface{}) (interface{}, error) {
	var err error
	var refValue reflect.Value

	switch bm := elem.(type) {
	case reflect.Type:
		var base = ValueOf(bm)
		if bm.Kind() == reflect.Ptr {
			return base, nil
		}
		return IndirectValue(base), nil
	case reflect.Value:
		refValue, err = deepCopyValue(IndirectValue(bm), bm.Kind() == reflect.Ptr)
	default:
		var base = ValueOf(elem)
		refValue, err = deepCopyValue(IndirectValue(base), base.Kind() == reflect.Ptr)
	}

	if err != nil {
		return nil, err
	}
	return refValue.Interface(), nil
}

// DeepCopyValue uses the internal deepCopyValue function to make a copy of giving
// reflect value type adequately ensuring to return a pointer if a pointer value
// was giving to it else returning a value.
//
// DeepCopy will never attempt to duplicate the following types, as it will simply
// return the value it received.
//
// 1. Channels
// 2. Interface
// 3. unsafe.Pointer
// 4. Func
//
// For structs, we use github.com/jinzhu/copier for copying, so all rules in that library
// applies here also.
func DeepCopyValue(baseValue reflect.Value) (interface{}, error) {
	var refValue, err = deepCopyValue(IndirectValue(baseValue), baseValue.Kind() == reflect.Ptr)
	if err != nil {
		return nil, err
	}
	return refValue.Interface(), nil
}

// deepCopyValue copies giving reflect.Value returning a new distinct value containing
// copy of data, else an error due to failure.
//
// DeepCopy will never attempt to duplicate the following types, as it will simply
// return the value it received.
//
// 1. Channels
// 2. Interface
// 3. unsafe.Pointer
// 4. Func
//
// For structs, we use github.com/jinzhu/copier for copying, so all rules in that library
// applies here also.
func deepCopyValue(baseValue reflect.Value, pointerType bool) (reflect.Value, error) {
	if !baseValue.IsValid() {
		return baseValue, nil
	}

	var valueType = baseValue.Type()

	switch baseValue.Kind() {
	case reflect.Func:
		return baseValue, nil
	case reflect.Interface:
		return baseValue, nil
	case reflect.Chan:
		return baseValue, nil
	case reflect.UnsafePointer:
		return baseValue, nil
	case reflect.Map:
		var newPtrVal = reflect.MakeMap(baseValue.Type())
		var newValue = IndirectValue(newPtrVal)
		if err := copyMap(&newValue, baseValue); err != nil {
			return reflect.Value{}, err
		}
		if pointerType {
			return newPtrVal, nil
		}
		return newValue, nil
	case reflect.Struct:
		var newPtrVal = reflect.New(valueType)
		var newValue = IndirectValue(newPtrVal)
		if err := copier.Copy(newPtrVal.Interface(), baseValue.Interface()); err != nil {
			return reflect.Value{}, nerror.Wrap(err, "Failed to copy struct type %T", baseValue.Interface())
		}
		if pointerType {
			return newPtrVal, nil
		}
		return newValue, nil
	case reflect.Array:
		var newPtrVal = reflect.New(valueType)
		var newValue = IndirectValue(newPtrVal)
		var copied = reflect.Copy(newValue, baseValue)
		if copied != baseValue.Len() {
			return reflect.Value{}, nerror.New("copied has different length to original")
		}
		if pointerType {
			return newPtrVal, nil
		}
		return newValue, nil
	case reflect.Slice:
		var newPtrVal = reflect.MakeSlice(valueType, baseValue.Len(), baseValue.Cap())
		var newValue = IndirectValue(newPtrVal)
		var copied = reflect.Copy(newValue, baseValue)
		if copied != baseValue.Len() {
			return reflect.Value{}, nerror.New("copied has different length to original")
		}
		if pointerType {
			return newPtrVal, nil
		}
		return newValue, nil
	}

	var newPtrValue = reflect.New(valueType)
	if newPtrValue.CanSet() {
		return reflect.Value{}, nerror.New("unable to set value for type %q", baseValue.Type())
	}

	var newValue = IndirectValue(newPtrValue)
	newValue.Set(baseValue)
	if pointerType {
		return newPtrValue, nil
	}
	return newValue, nil
}

func copyMap(dest *reflect.Value, src reflect.Value) error {
	if dest.Kind() != reflect.Map && src.Kind() != reflect.Map {
		return nerror.New("unable to copy in between map and another type")
	}

	var keys = src.MapKeys()
	for _, key := range keys {
		newKey, err := deepCopyValue(IndirectValue(key), key.Kind() == reflect.Ptr)
		if err != nil {
			return err
		}

		var mapValue = src.MapIndex(key)
		newValue, err := deepCopyValue(IndirectValue(mapValue), mapValue.Kind() == reflect.Ptr)
		if err != nil {
			return err
		}
		dest.SetMapIndex(newKey, newValue)
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
func CopyStringKeyMap(src map[string]interface{}) (map[string]interface{}, error) {
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
