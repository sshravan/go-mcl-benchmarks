package main

import (
	"flag"
	"fmt"
	"testing"

	"github.com/alinush/go-mcl"
	"github.com/dustin/go-humanize"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	// a := time.Duration(r.NsPerOp() / int64(size))
	// out := fmt.Sprintf("Time per %s (%d iters%s):", op, r.N, aux)
	// fmt.Printf("%-60s %20v\n", out, a)

	p := message.NewPrinter(language.English)
	a := float64(r.NsPerOp()/int64(size)) / float64(1000) // Convert ns to us
	out := fmt.Sprintf("Time per %s (%s%d iters):", op, aux, r.N)
	p.Printf("%-60s %20.3f us\n", out, a)
}

func BenchmarkExponentiation() {

	var size []uint64
	size = []uint64{1000}
	var length int
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
					mcl.G1Neg(&result, &baseG1[j])
				}
			}
		})
		Summary(size[i], "G1Neg", "", &results)
		// =============================================
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
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.G1Sub(&result, &result, &baseG1[j])
				}
			}
		})
		Summary(size[i], "G1Sub", "", &results)
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
		length = 2
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G1
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G1MulVec(&result, baseG1[:length], expoFr[:length])
			}
		})
		Summary(1, "G1MulVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "G1MulVec", fmt.Sprintf("per exp, "), &results)

		// =============================================
		length = 32
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G1
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G1MulVec(&result, baseG1[:length], expoFr[:length])
			}
		})
		Summary(1, "G1MulVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "G1MulVec", fmt.Sprintf("per exp, "), &results)

		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G1
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G1MulVec(&result, baseG1, expoFr)
			}
		})
		Summary(1, "G1MulVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		Summary(size[i], "G1MulVec", fmt.Sprintf("per exp, "), &results)

		fmt.Println("=============================================")
		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G2
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG2); j++ {
					mcl.G2Neg(&result, &baseG2[j])
				}
			}
		})
		Summary(size[i], "G2Neg", "", &results)
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
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.G2Sub(&result, &result, &baseG2[j])
				}
			}
		})
		Summary(size[i], "G2Sub", "", &results)

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
		length = 2
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G2
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G2MulVec(&result, baseG2[:length], expoFr[:length])
			}
		})
		Summary(1, "G2MulVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "G2MulVec", fmt.Sprintf("per exp, "), &results)

		// =============================================
		length = 32
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G2
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G2MulVec(&result, baseG2[:length], expoFr[:length])
			}
		})
		Summary(1, "G2MulVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "G2MulVec", fmt.Sprintf("per exp, "), &results)

		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.G2
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.G2MulVec(&result, baseG2, expoFr)
			}
		})
		Summary(1, "G2MulVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		Summary(size[i], "G2MulVec", fmt.Sprintf("per exp, "), &results)
		fmt.Println("=============================================")
		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.Fr
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(expoFr); j++ {
					mcl.FrNeg(&result, &expoFr[j])
				}
			}
		})
		Summary(size[i], "FrNeg", "", &results)
		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.Fr
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(expoFr); j++ {
					mcl.FrInv(&result, &expoFr[j])
				}
			}
		})
		Summary(size[i], "FrInv", "", &results)
		fmt.Println("=============================================")
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
		expoFr := generateFr(size[i])
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
		// b.Run(fmt.Sprintf("%d/GTMul;", size[i]),
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			result.SetString("1", 10)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				for j := 0; j < len(baseG1); j++ {
					mcl.GTPow(&result, &baseGT[j], &expoFr[j])
				}
			}
		})
		Summary(size[i], "GTPow", "", &results)
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
		// b.Run(fmt.Sprintf("%d/MillerLoopVec;", size[i]),
		length = 2
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.MillerLoopVec(&result, baseG1[:length], baseG2[:length])
			}
		})
		Summary(1, "MillerLoopVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "MillerLoopVec", fmt.Sprintf("per MillerLoop, "), &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/MillerLoopVec;", size[i]),
		length = 32
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.MillerLoopVec(&result, baseG1[:length], baseG2[:length])
			}
		})
		Summary(1, "MillerLoopVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "MillerLoopVec", fmt.Sprintf("per MillerLoop, "), &results)

		// =============================================
		// b.Run(fmt.Sprintf("%d/MillerLoopVec;", size[i]),
		results = testing.Benchmark(func(t *testing.B) {
			var result mcl.GT
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				mcl.MillerLoopVec(&result, baseG1, baseG2)
			}
		})
		Summary(1, "MillerLoopVec", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		Summary(size[i], "MillerLoopVec", fmt.Sprintf("per MillerLoop, "), &results)
		fmt.Println("=============================================")

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
		Summary(1, "Multi-Pairing", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "Multi-Pairing", fmt.Sprintf("per pairing, "), &results)

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
		Summary(1, "Multi-Pairing", fmt.Sprintf("size %s, ", humanize.Comma(int64(length))), &results)
		Summary(uint64(length), "Multi-Pairing", fmt.Sprintf("per pairing, "), &results)

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
		Summary(1, "Multi-Pairing", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		Summary(size[i], "Multi-Pairing", fmt.Sprintf("per pairing, "), &results)
		fmt.Println("=============================================")
		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var a mcl.Fr
			a.Random()
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				a.IsEqual(&expoFr[0])
				for j := 0; j < len(expoFr)-1; j++ {
					expoFr[j].IsEqual(&expoFr[j+1])
				}
			}
		})
		Summary(size[i], "FrIsEqual", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var a mcl.G1
			a.Random()
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				a.IsEqual(&baseG1[0])
				for j := 0; j < len(baseG1)-1; j++ {
					baseG1[j].IsEqual(&baseG1[j+1])
				}
			}
		})
		Summary(size[i], "G1IsEqual", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var a mcl.G2
			a.Random()
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				a.IsEqual(&baseG2[0])
				for j := 0; j < len(baseG2)-1; j++ {
					baseG2[j].IsEqual(&baseG2[j+1])
				}
			}
		})
		Summary(size[i], "G2IsEqual", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		// =============================================
		results = testing.Benchmark(func(t *testing.B) {
			var a mcl.GT
			a.SetInt64(1)
			t.ResetTimer()
			for i := 0; i < t.N; i++ {
				a.IsEqual(&baseGT[0])
				for j := 0; j < len(baseGT)-1; j++ {
					baseGT[j].IsEqual(&baseGT[j+1])
				}
			}
		})
		Summary(size[i], "GTIsEqual", fmt.Sprintf("size %s, ", humanize.Comma(int64(size[i]))), &results)
		fmt.Println("=============================================")
	}
}
