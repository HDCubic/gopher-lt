package lt

import (
	"fmt"
	"reflect"

	lua "github.com/yuin/gopher-lua"
)

//toLValue 转换成LValue并判断是否为空
func toLValue(L *lua.LState, value reflect.Value) (lua.LValue, bool) {
	fmt.Println(value.Kind())
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
		for i := 0; i < value.Len(); i++ {
			lv, _ := toLValue(L, value.Index(i))
			v.RawSet(lua.LNumber(i), lv)
			fmt.Println("33333", lv.String())
		}
		return v, value.Len() == 0
	case reflect.Ptr:
		//return NewLTable(L, value.Interface()), true
		e := value.Elem()
		et := e.Type()
		//ev := reflect.ValueOf(e)
		fmt.Println(et.String())

		v := L.NewTable()
		for i := 0; i < et.NumField(); i++ {
			f := et.Field(i)
			name := f.Tag.Get("lt")
			if lv, yes := toLValue(L, e.Field(i)); !yes {
				v.RawSetString(name, lv)
				fmt.Println("2222", name, lv.String())
			}
		}
		fmt.Println("11111", v.String())
		return v, true
	}
	return NewLTable(L, value.Interface()), true
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
			t.RawSet(lua.LString(name), v)
		}
	}
	return t
}
