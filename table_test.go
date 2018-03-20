package lt

import (
	"fmt"
	"testing"

	lua "github.com/yuin/gopher-lua"
	luajson "layeh.com/gopher-json"
)

type SValue struct {
	Name string   `json:"name" lt:"name"`
	List []string `json:"list" lt:"list"`
}
type Demo struct {
	ResourceID string  `json:"resource_id" lt:"resource_id"`
	Value      int64   `json:"value" lt:"value"`
	IsEnable   bool    `json:"is_enable" lt:"is_enable"`
	FValue     float64 `json:"f_value" lt:"f_value"`
	SValue     *SValue `json:"s_value" lt:"s_value"`
}

// Build 构造lua.LTable
func Build(L *lua.LState) int {
	value := Demo{
		ResourceID: "3121412",
		Value:      123214,
		IsEnable:   true,
		FValue:     1.23,
		SValue: &SValue{
			Name: "dasd",
			List: []string{"1", "2", "3", "4", "5"},
		},
	}
	table := NewLTable(L, value)
	L.Push(table)
	return 1
}

// Parse 解析lua.LTable
func Parse(L *lua.LState) int {
	t := L.ToTable(1)
	value := Demo{}
	err := FromLTable(t, &value)
	fmt.Println(value, err)
	return 0
}

func TestNewLTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	luajson.Preload(L)
	L.SetGlobal("build", L.NewFunction(Build))
	L.SetGlobal("parse", L.NewFunction(Parse))
	if err := L.DoString(`
	json = require("json")
	t = build()
	for k, v in pairs(t) do
		print(k, v)
	end
	for k, v in pairs(t.s_value) do
		print(k, v)
	end
	for k, v in pairs(t.s_value.list) do
		print(k, v)
	end
	te = json.encode(t)
	print(te)
	td = json.decode(te)
	print(td)
	parse(td)
	`); err != nil {
		panic(err)
	}
}
