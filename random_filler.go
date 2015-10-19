package filler

import (
	crand "crypto/rand"
	"encoding/base32"
	"io"
	"reflect"
	"strings"
)

// RandomFiller は、プリミティブ型の値をランダムに設定する Filler を取得します
func RandomFiller(seed *RandomSeed) *Filler {
	gen := &Filler{}
	if seed == nil {
		def := DefaultSeed()
		seed = &def
	}
	gen.RegisterKind(reflect.Ptr, func() (interface{}, FactoryState) {
		if seed.Random.Intn(seed.nilRateDenom()) == 0 {
			return nil, Done
		}
		return nil, None
	})
	gen.RegisterKind(reflect.Bool, func() (interface{}, FactoryState) {
		return seed.Random.Intn(2) == 0, Done
	})
	gen.RegisterKind(reflect.Int, func() (interface{}, FactoryState) {
		return seed.Random.Int(), Done
	})
	gen.RegisterKind(reflect.Int8, func() (interface{}, FactoryState) {
		return int8(seed.Random.Intn(256) - 128), Done
	})
	gen.RegisterKind(reflect.Int16, func() (interface{}, FactoryState) {
		return int16(seed.Random.Intn(65536) - 32768), Done
	})
	gen.RegisterKind(reflect.Int32, func() (interface{}, FactoryState) {
		return int32(seed.Random.Int63n(65536*65536) - 65536*32768), Done
	})
	gen.RegisterKind(reflect.Int64, func() (interface{}, FactoryState) {
		return seed.Random.Int63() * int64(seed.Random.Intn(2)*2-1), Done
	})
	gen.RegisterKind(reflect.Uint, func() (interface{}, FactoryState) {
		return seed.Random.Uint32(), Done
	})
	gen.RegisterKind(reflect.Uint8, func() (interface{}, FactoryState) {
		return uint8(seed.Random.Intn(256)), Done
	})
	gen.RegisterKind(reflect.Uint16, func() (interface{}, FactoryState) {
		return uint16(seed.Random.Intn(65536)), Done
	})
	gen.RegisterKind(reflect.Uint32, func() (interface{}, FactoryState) {
		return uint32(seed.Random.Int63n(int64(65536)*65536 - 1)), Done
	})
	gen.RegisterKind(reflect.Uint64, func() (interface{}, FactoryState) {
		return uint64(seed.Random.Int63()), Done
	})
	gen.RegisterKind(reflect.Float32, func() (interface{}, FactoryState) {
		return seed.Random.Float32(), Done
	})
	gen.RegisterKind(reflect.Float64, func() (interface{}, FactoryState) {
		return seed.Random.Float64(), Done
	})
	gen.RegisterKind(reflect.Complex64, func() (interface{}, FactoryState) {
		return complex(seed.Random.Float32(), seed.Random.Float32()), Done
	})
	gen.RegisterKind(reflect.Complex128, func() (interface{}, FactoryState) {
		return complex(seed.Random.Float64(), seed.Random.Float64()), Done
	})
	gen.RegisterKind(reflect.String, func() (interface{}, FactoryState) {
		b := make([]byte, seed.Random.Intn(seed.stringCap()-seed.stringMin())+seed.stringMin())
		_, err := io.ReadFull(crand.Reader, b)
		if err != nil {
			return "", None
		}
		return strings.TrimRight(base32.StdEncoding.EncodeToString(b), "="), Done
	})

	gen.RegisterFunc(func(typ reflect.Type) (interface{}, FactoryState) {
		if typ.Kind() == reflect.Slice {
			cap := seed.sliceCap()
			return reflect.MakeSlice(typ, seed.Random.Intn(cap-seed.sliceMin())+seed.sliceMin(), cap).Interface(), Init
		}
		if typ.Kind() == reflect.Map {
			kT := typ.Key()
			vT := typ.Elem()
			val := reflect.MakeMap(typ)

			for i := 0; i < seed.Random.Intn(seed.mapLen()); i++ {
				kV := gen.Make(kT)
				vV := gen.Make(vT)
				val.SetMapIndex(reflect.ValueOf(kV), reflect.ValueOf(vV))
			}
			return val.Interface(), Init
		}
		return nil, None
	})
	return gen
}
