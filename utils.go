package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/alinush/go-mcl"
)

func generateG1(count uint64) []mcl.G1 {
	base := make([]mcl.G1, count)
	for i := uint64(0); i < count; i++ {
		base[i].Random()
	}
	return base
}

func generateG2(count uint64) []mcl.G2 {
	base := make([]mcl.G2, count)
	for i := uint64(0); i < count; i++ {
		base[i].Random()
	}
	return base
}

func generateFr(count uint64) []mcl.Fr {
	base := make([]mcl.Fr, count)
	for i := uint64(0); i < count; i++ {
		base[i].Random()
	}
	return base
}

func generateGT(count uint64) []mcl.GT {
	rand.Seed(time.Now().UnixNano())
	N := int64(math.MaxInt64)
	var v int64
	base := make([]mcl.GT, count)
	for i := uint64(0); i < count; i++ {
		v = rand.Int63n(N)
		base[i].SetInt64(v)
	}
	return base
}

func getKeyValues(db map[string]float64) ([]string, []float64) {

	keys := make([]string, 0, len(db))
	values := make([]float64, 0, len(db))
	for k, v := range db {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
