package filler

import (
	"math"
	"math/rand"
	"time"
)

// RandomSeed は、RandomGeneratorで使用する乱数のパラメータです
type RandomSeed struct {
	Random         *rand.Rand
	NilRate        float64
	SliceCapacity  int
	SliceMinimum   int
	StringCapacity int
	StringMinimum  int
	MapLength      int
}

func (s *RandomSeed) nilRateDenom() int {
	return int(math.Ceil(1 / s.NilRate))
}

func (s *RandomSeed) sliceCap() int {
	return s.SliceCapacity
}

func (s *RandomSeed) sliceMin() int {
	return s.SliceMinimum
}

func (s *RandomSeed) stringCap() int {
	return s.StringCapacity
}

func (s *RandomSeed) stringMin() int {
	return s.StringMinimum
}

func (s *RandomSeed) mapLen() int {
	return s.MapLength
}

// DefaultSeed は典型的な RandomSeed を取得します
func DefaultSeed() RandomSeed {
	return RandomSeed{
		Random:         rand.New(rand.NewSource(time.Now().Unix())),
		NilRate:        1 / 256.0,
		SliceCapacity:  8,
		StringCapacity: 8,
		SliceMinimum:   0,
		StringMinimum:  0,
		MapLength:      8,
	}
}
