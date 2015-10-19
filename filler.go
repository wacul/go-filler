package filler

import (
	"reflect"
	"sync"
)

// Filler はオブジェクトの内容を生成する
type Filler struct {
	names       map[string]Factory
	kinds       map[reflect.Kind]Factory
	types       []TypedFactory
	initOnce    sync.Once
	FieldFilter func(structType reflect.Type, field reflect.StructField, value reflect.Value) bool
}

// FactoryState は、オブジェクトを生成する関数の戻り値
type FactoryState int

const (
	// None は、オブジェクトが生成されなかったことを示す
	None FactoryState = iota
	// Init は、オブジェクトが生成されて初期化されたことを示す
	Init
	// Done は、オブジェクトが生成されて全ての値がFillされたことを示す
	Done
)

// Factory はオブジェクトを生成する関数の定義
type Factory func() (interface{}, FactoryState)

// TypedFactory は型の指定に基づきオブジェクトを生成する関数の定義
type TypedFactory func(typ reflect.Type) (interface{}, FactoryState)

func (g *Filler) init() {
	g.initOnce.Do(func() {
		g.names = make(map[string]Factory)
		g.kinds = make(map[reflect.Kind]Factory)
		g.types = []TypedFactory{
			g.genName,
			g.genKind,
		}
	})
}

// RegisterName は pkgPath: パッケージパス、typeName: 型名 に紐づく con: 生成関数を登録します
func (g *Filler) RegisterName(pkgPath string, typeName string, con Factory) {
	g.init()
	g.names[fullName(pkgPath, typeName)] = con
}

// RegisterKind は kind: 型種別に紐づく con: 生成関数を登録します
func (g *Filler) RegisterKind(kind reflect.Kind, con Factory) {
	g.init()
	g.kinds[kind] = con
}

// RegisterFunc はオブジェクトの生成関数を登録します
func (g *Filler) RegisterFunc(con TypedFactory) {
	g.init()
	g.types = append(g.types, con)
}

func (g *Filler) genName(typ reflect.Type) (interface{}, FactoryState) {
	typeName := fullName(typ.PkgPath(), typ.Name())
	if typeName == "" || g.names == nil {
		return nil, None
	}
	con, ok := g.names[typeName]
	if !ok {
		return nil, None
	}

	return con()
}

func (g *Filler) genKind(typ reflect.Type) (interface{}, FactoryState) {
	kind := typ.Kind()
	if kind == reflect.Invalid || g.kinds == nil {
		return nil, None
	}
	con, ok := g.kinds[kind]
	if !ok {
		return nil, None
	}

	return con()
}

func (g *Filler) genType(typ reflect.Type) (interface{}, FactoryState) {
	for _, gen := range g.types {
		obj, state := gen(typ)
		if state != None {
			return obj, state
		}
	}
	return nil, None
}

// Fill は、指定されたオブジェクトに再帰的に値を設定します
func (g *Filler) Fill(object interface{}) {
	objVal := reflect.ValueOf(object).Elem()
	g.fillType(objVal)
}

// Make は、指定された型に再帰的に値を設定したオブジェクトを取得します
func (g *Filler) Make(typ reflect.Type) interface{} {
	objVal := reflect.New(typ)
	g.fillValue(typ, objVal.Elem())
	return objVal.Elem().Interface()
}

func (g *Filler) fillType(val reflect.Value) {
	g.fillValue(val.Type(), val)
}

func (g *Filler) fillValue(typ reflect.Type, val reflect.Value) {
	if !val.CanSet() {
		return
	}

	obj, state := g.genType(typ)
	if state != None {
		got := reflect.ValueOf(obj)
		if (got.Kind() != reflect.Ptr || !got.IsNil()) && (val.Kind() != reflect.Ptr || !val.IsNil()) {
			con := got.Convert(typ)
			val.Set(con)
		}
	}
	if state == Done {
		return
	}

	switch typ.Kind() {
	case reflect.Ptr:
		t := typ.Elem()     // ポインタ内容の型
		v := reflect.New(t) // メモリ領域の確保
		val.Set(v)
		if !val.IsNil() {
			g.fillType(val.Elem())
		}
	case reflect.Array:
		for i := 0; i < val.Len(); i++ {
			g.fillType(val.Index(i))
		}
	case reflect.Slice:
		if !val.IsNil() {
			for i := 0; i < val.Len(); i++ {
				g.fillType(val.Index(i))
			}
		}
	case reflect.Map:
		if !val.IsNil() {
			v := typ.Elem()
			for _, key := range val.MapKeys() {
				g.fillValue(v, val.MapIndex(key))
			}
		}
	case reflect.Struct:
		t := typ
		v := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			if g.FieldFilter == nil || g.FieldFilter(typ, typ.Field(i), v.Field(i)) {
				g.fillType(v.Field(i))
			}
		}
		val.Set(v)
	}
}

func fullName(pkgPath, typeName string) string {
	if typeName != "" {
		if pkgPath == "" {
			return typeName
		}
		return pkgPath + "." + typeName
	}
	return ""
}
