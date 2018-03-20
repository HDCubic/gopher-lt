package lt

import (
	"fmt"
	"reflect"
	"strconv"

	lua "github.com/yuin/gopher-lua"
)

//toLValue 转换成LValue并判断是否为空
func toLValue(L *lua.LState, value reflect.Value) (lua.LValue, bool) {
	switch value.Kind() {
	case reflect.String:
		return lua.LString(value.String()), value.Len() == 0
	case reflect.Bool:
		return lua.LBool(value.Bool()), !value.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return lua.LNumber(value.Int()), value.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return lua.LNumber(value.Uint()), value.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return lua.LNumber(value.Float()), value.Float() == 0
	case reflect.Slice:
		v := L.NewTable()
		//vs := []lua.LValue{}
		for i := 0; i < value.Len(); i++ {
			lv, _ := toLValue(L, value.Index(i))
			//vs = append(vs, lv)
			v.RawSetInt(i+1, lv)
			//v.RawSet(lua.LNumber(i), lv)
		}
		return v, value.Len() == 0
	case reflect.Ptr:
		//return NewLTable(L, value.Interface()), true
		e := value.Elem()
		et := e.Type()

		v := L.NewTable()
		for i := 0; i < et.NumField(); i++ {
			f := et.Field(i)
			name := f.Tag.Get("lt")
			if lv, yes := toLValue(L, e.Field(i)); !yes {
				v.RawSetString(name, lv)
			}
		}
		return v, false
	}
	return NewLTable(L, value.Interface()), false
}

// NewLTable 从value生成一个lua table
func NewLTable(L *lua.LState, value interface{}) *lua.LTable {
	t := L.NewTable()
	vt := reflect.TypeOf(value)
	vv := reflect.ValueOf(value)
	for i := 0; i < vt.NumField(); i++ {
		f := vt.Field(i)
		name := f.Tag.Get("lt")
		if v, yes := toLValue(L, vv.Field(i)); !yes {
			//t.RawSet(lua.LString(name), v)
			t.RawSetString(name, v)
		}
	}
	return t
}

func fromLVlaue(value lua.LValue) interface{} {
	fmt.Println(value.Type())
	switch value.Type() {
	case lua.LTNumber:
		fmt.Println("11111111", value)
		return value
	case lua.LTString:
		fmt.Println("22222222", value)
		return value
	case lua.LTBool:
		fmt.Println("33333333", value)
		return value
	case lua.LTNil:
		fmt.Println("44444444", value)
	case lua.LTTable:
		fmt.Println("5555555", value)
	}
	return nil
}

// FromLTable 从lua table 绑定object
func FromLTable(t *lua.LTable, value interface{}) error {
	vv := reflect.ValueOf(value)
	e := vv.Elem()
	et := e.Type()
	fmt.Println("AAAAAAAAAAAAAAAAAAAAAa", vv, e, et)
	for i := 0; i < et.NumField(); i++ {
		f := et.Field(i)
		name := f.Tag.Get("lt")
		v := t.RawGetString(name)
		switch f.Type.Kind() {
		case reflect.Int, reflect.Int32, reflect.Int64:
			tv, err := strconv.ParseInt(v.String(), 10, 64)
			if err != nil {
				return err
			}
			e.Field(i).SetInt(tv)
		case reflect.Uint, reflect.Uint32, reflect.Uint64:
			tv, err := strconv.ParseUint(v.String(), 10, 64)
			if err != nil {
				return err
			}
			e.Field(i).SetUint(tv)
		case reflect.Float32, reflect.Float64:
			tv, err := strconv.ParseFloat(v.String(), 64)
			if err != nil {
				return err
			}
			e.Field(i).SetFloat(tv)
		case reflect.Bool:
			tv, err := strconv.ParseBool(v.String())
			if err != nil {
				return err
			}
			e.Field(i).SetBool(tv)
		case reflect.Ptr:
			fmt.Println("aaaaaaaaaaaaaaaa", v.Type(), name, f.Name)
			switch v.Type() {
			case lua.LTNil:
			case lua.LTTable:
				tv := e.Field(i).Elem()
				fmt.Println(tv.Kind())
				FromLTable(v.(*lua.LTable), &tv)
			}
		}
		//fmt.Println(f.Name, name, v, f.Type.Kind())
		//fromLVlaue(v)
	}
	//switch vt.Kind() {
	//}
	//vv := reflect.ValueOf(value)
	//for i := 0; i < vt.NumField(); i++ {
	//	fmt.Println(i)
	//	//f := vt.Field(i)
	//	//name := f.Tag.Get("lt")
	//	//v := t.RawGetString(name)
	//	//fmt.Println(name, v.String())
	//}
	return nil
}
