package main

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/alinush/go-mcl"
	"github.com/dustin/go-humanize"
	"github.com/icza/gox/timex"
)

func main() {
	testing.Init()
	flag.Parse()
	fmt.Println("Hello, World!")
	mcl.InitFromString("bls12-381")

	BenchmarkExponentiation()
	BenchmarkPairing()
}

func Summary(size uint64, op string, aux string, r *testing.BenchmarkResult) {

	a := time.Duration(r.NsPerOp() / int64(size))
	out := fmt.Sprintf("Time per %s (%d iters%s):", op, r.N, aux)
	fmt.Printf("%-60s %20v\n", out, timex.Round(a, 3))
}

func BenchmarkExponentiation() {

	var size []uint64
	size = []uint64{1000}
	for i := 0; i < len(size); i++ {
		baseG1 := generateG1(size[i])
		baseG2 := generateG2(size[i])
		expoFr := generateFr(size[i])
		fmt.Println("Done generating the data")

		var results testing.BenchmarkResult

		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G1
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.G1Add(&result, &result, &baseG1[j])
				}
			}
		})
		Summary(size[i], "G1Add", "", &results)

		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G1
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(expoFr); j++ {
					mcl.G1Mul(&result, &baseG1[j], &expoFr[j])
				}
			}
		})
		Summary(size[i], "G1Mul", "", &results)

		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G1
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G1MulVec(&result, baseG1, expoFr)
			}
		})
		Summary(1, "G1MulVec", fmt.Sprintf(" size %s", humanize.Comma(int64(size[i]))), &results)
		Summary(size[i], "G1MulVec", fmt.Sprintf(", per exp"), &results)

		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G2
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.G2Add(&result, &result, &baseG2[j])
				}
			}
		})
		Summary(size[i], "G2Add", "", &results)

		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G2
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(expoFr); j++ {
					mcl.G2Mul(&result, &baseG2[j], &expoFr[j])
				}
			}
		})
		Summary(size[i], "G2Mul", "", &results)

		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G2
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G2MulVec(&result, baseG2, expoFr)
			}
		})
		Summary(1, "G2MulVec", fmt.Sprintf(" size %s", humanize.Comma(int64(size[i]))), &results)
		Summary(size[i], "G2MulVec", fmt.Sprintf(", per exp"), &results)
	}
}

func BenchmarkPairing() {

	var size []uint64
	size = []uint64{100}
	var length int
	for i := 0; i < len(size); i++ {
		baseG1 := generateG1(size[i])
		baseG2 := generateG2(size[i])
		baseGT := generateGT(size[i])
		fmt.Println("Done generating the data")

		var results testing.BenchmarkResult
		// =============================================
		// b.Run(fmt.Sprintf("%d/GTMul;", size[i]),
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.GTMul(&result, &result, &baseGT[j])
				}
			}
		})
		Summary(size[i], "GTMul", "", &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/MillerLoop;", size[i]),
		results = testing.Benchmark(func(t *testing.B) {
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.MillerLoop(&baseGT[j], &baseG1[j], &baseG2[j])
				}
			}
		})
		Summary(size[i], "MillerLoop", "", &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/FinalExp;", size[i]),
		results = testing.Benchmark(func(t *testing.B) {
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.FinalExp(&baseGT[j], &baseGT[j])
				}
			}
		})
		Summary(size[i], "FinalExp", "", &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/NaivePairing;", size[i]),
		results = testing.Benchmark(func(t *testing.B) {
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.Pairing(&baseGT[j], &baseG1[j], &baseG2[j])
				}
			}
		})
		Summary(size[i], "Pairing", "", &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/MillerLoopVec;", size[i]),
		length = 2
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.MillerLoopVec(&result, baseG1[:length], baseG2[:length])
				mcl.FinalExp(&result, &result)
			}
		})
		Summary(1, "Multi-Pairing", fmt.Sprintf(" size %s", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "Multi-Pairing", fmt.Sprintf(", per pairing"), &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/MillerLoopVec;", size[i]),
		length = 32
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.MillerLoopVec(&result, baseG1[:length], baseG2[:length])
				mcl.FinalExp(&result, &result)
			}
		})
		Summary(1, "Multi-Pairing", fmt.Sprintf(" size %s", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "Multi-Pairing", fmt.Sprintf(", per pairing"), &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/MillerLoopVec;", size[i]),
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.MillerLoopVec(&result, baseG1, baseG2)
				mcl.FinalExp(&result, &result)
			}
		})
		Summary(1, "Multi-Pairing", fmt.Sprintf(" size %s", humanize.Comma(int64(size[i]))), &results)
		Summary(size[i], "Multi-Pairing", fmt.Sprintf(", per pairing"), &results)
	}
}
