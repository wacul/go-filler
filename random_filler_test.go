package filler

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/mgo.v2/bson"
)

func TestRandomFiller(t *testing.T) {

	var object struct {
		ID      bson.ObjectId
		Name    string
		Number  int
		Complex complex64

		Array [2]string
		Slice []string
		Map   map[string]string

		Object struct {
			ID     bson.ObjectId
			Name   string
			Number int

			StringP *string
			ObjectP *struct {
				ID     bson.ObjectId
				Name   string
				Number int

				Array [2]string
				Slice []string
				Map   map[string]string
			}

			ArrayP *[2]string
			SliceP *[]string
			MapP   *map[string]string

			StringPP **string
			ObjectPP **struct {
				ID     bson.ObjectId
				Name   string
				Number int
				Array  [2]string
			}

			ArrayPP **[2]string
			SlicePP **[]string
			MapPP   **map[string]string

			PArray  [2]*string
			PPArray [2]**string
			PSlice  []*string
			PPSlice []**string
			PMap    map[string]*string
			PPMap   map[string]**string
		}
	}
	src := DefaultSeed()
	src.NilRate = 1.0 / 256
	gen := RandomFiller(&src)
	gen.RegisterName(
		"gopkg.in/mgo.v2/bson",
		"ObjectId",
		func() (interface{}, FactoryState) { return bson.NewObjectId(), Done },
	)
	gen.Fill(&object)
	spew.Dump(object)
}
