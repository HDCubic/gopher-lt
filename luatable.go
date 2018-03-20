package lt

import (
	"fmt"
	"reflect"

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

//fromLValue 转换成实际值并判断是否为空
func fromLValue(L *lua.LState, value *lua.LValue) interface{} {
	return nil
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

// FromLValue 从lua table 绑定object
func FromLValue(L *lua.LState, t *lua.LTable, value interface{}) {
	vt := reflect.TypeOf(value)
	//vv := reflect.ValueOf(value)
	fmt.Println(vt.String(), vt.Kind())
	switch vt.Kind() {
	case reflect.String:
	case reflect.Bool:
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
	case reflect.Float32, reflect.Float64:
	case reflect.Slice:
	case reflect.Ptr:
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
}
