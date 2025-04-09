package filler

import (
	"reflect"
	"time"
)

func FillStruct[T any]() T {
	var res T
	t := reflect.TypeOf(res)
	if t.Kind() == reflect.Struct {
		v := reflect.New(t).Elem()
		fillValue(v)
		return v.Interface().(T)
	}
	if t.Kind() == reflect.Ptr {
		v := reflect.New(t.Elem())
		if t.Elem().Kind() == reflect.Struct {
			fillValue(v.Elem())
		}
		return v.Interface().(T)
	}
	return res
}

func fillValue(v reflect.Value) {
	if !v.CanSet() {
		return
	}

	switch v.Kind() {
	case reflect.Bool:
		var x bool
		v.SetBool(x)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var x int64
		v.SetInt(x)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var x uint64
		v.SetUint(x)
	case reflect.Float32, reflect.Float64:
		var x float64
		v.SetFloat(x)
	case reflect.Complex64, reflect.Complex128:
		var x complex128
		v.SetComplex(x)
	case reflect.Array:
		for i := range v.Len() {
			fillValue(v.Index(i))
		}
	case reflect.Chan:
		// closed channel by default
		c := reflect.MakeChan(v.Type(), 0)
		c.Close()
		v.Set(c)
	case reflect.Func:
		if v.IsNil() {
			fType := v.Type()
			f := reflect.MakeFunc(fType, func(args []reflect.Value) []reflect.Value {
				ret := make([]reflect.Value, fType.NumOut())
				for i := range fType.NumOut() {
					outType := fType.Out(i)
					ret[i] = reflect.Zero(outType)
				}
				return ret
			})
			v.Set(f)
		}
	case reflect.Interface:
		if v.IsNil() {
			if v.NumMethod() == 0 {
				v.Set(reflect.ValueOf(struct{}{}))
			}
		}
	case reflect.Map:
		m := reflect.MakeMap(v.Type())
		kType, vType := v.Type().Key(), v.Type().Elem()
		key, val := reflect.New(kType).Elem(), reflect.New(vType).Elem()
		fillValue(key)
		fillValue(val)
		m.SetMapIndex(key, val)
		v.Set(m)
	case reflect.Ptr:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		fillValue(v.Elem())
	case reflect.Slice:
		// one-element slice by default
		s := reflect.MakeSlice(v.Type(), 1, 1)
		fillValue(s.Index(0))
		v.Set(s)
	case reflect.String:
		var x string
		v.SetString(x)
	case reflect.Struct:
		// special case for time.Time
		if v.Type() == reflect.TypeOf(time.Time{}) {
			v.Set(reflect.ValueOf(time.Now()))
			return
		}
		for i := range v.NumField() {
			fillValue(v.Field(i))
		}
	default:
		// reflect.Invalid | reflect.Uintptr | reflect.UnsafePointer
	}
}
