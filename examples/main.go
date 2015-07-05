package main

import (
	"../"
	"fmt"
	"math/rand"
	"runtime"
	"time"

	//"github.com/wiless/turbo"
	"github.com/wiless/vlib"
)

func bpskmod(ipbits []int) (modsyms []float64) {
	n := len(ipbits)
	modsyms = make([]float64, n)

	for i := 0; i < n; i++ {
		modsyms[i] = 2.0*float64(ipbits[i]) - 1.0
	}
	return modsyms

}

func bpskawgndemod(recsym []float64, variance float64) (llr []float64) {
	n := len(recsym)
	llr = make([]float64, n)
	scale := 4.0 / variance

	// now := time.Now()
	for i := 0; i < n; i++ {
		llr[i] = scale * recsym[i]
	}
	// t1 := time.Since(now)
	// fmt.Println("Elapsed ", t1)
	// fmt.Println("output ", llr[0:10])
	// var v blas64.Vector
	// v.Inc = 1
	// v.Data = recsym
	// now = time.Now()
	// blas64.Scal(n, scale, v)
	// t2 := time.Since(now)
	// if t1 < t2 {
	// 	fmt.Println("OLD - NEW  ", t1, t2, t1-t2)
	// }

	// fmt.Println("output ", v.Data[0:10])

	return llr
}

func main() {
	seedvalue := time.Now().Unix()
	rand.Seed(seedvalue)
	CORES := 1
	runtime.GOMAXPROCS(16)
	BLOCKLENGTH := 6144
	NOBLOCKS := 500
	fmt.Printf("BLOCKLENGTH = %d \t NOBLOCKS = %d\n", BLOCKLENGTH, NOBLOCKS)
	//var output vlib.VectorB
	startTime := time.Now()
	variance := 1.0
	ch := make(chan bool, CORES)
	// var wg sync.WaitGroup
	for c := 0; c < CORES; c++ {

		// wg.Add(1)
		go func(c int, ch chan bool) {

			var berc int
			berc = 0
			symbol := []byte{'|', '+', '-', 'X'}
			// _ = symbol
			var codec turbo.LTECodec
			codec.Init("BPSK", "_13_", BLOCKLENGTH)

			for i := 0; i < NOBLOCKS/CORES; i++ {
				input := vlib.Randsrc(BLOCKLENGTH, 2)
				// fmt.Printf("\rCORE %d %c   :  BLOCK  : %d", c, symbol[c%4], i)
				// input := vlib.Randsrc(BLOCKLENGTH, 2)

				output := codec.Encode(input)
				mod := bpskmod(output)
				llr := bpskawgndemod(mod, variance)
				dec_output := codec.Decode(llr)
				//fmt.Print("Input bits = %v", input)
				//fmt.Print("\nOuput bits = %v", dec_output)

				for j := 0; j < codec.BLOCKLENGTH; j++ {
					berc = berc + input[j] ^ dec_output[j]

				}

			}
			fmt.Printf("\nCORE %d %c   :  BLOCK  : %f", c, symbol[c%4], float64(berc)/(float64(BLOCKLENGTH)*float64(NOBLOCKS)))
			// wg.Done()
			ch <- true
		}(c, ch)
	} //
	// wg.Wait()
	for i := 0; i < CORES; i++ {
		fmt.Print(<-ch)
	}
	berc := 0
	fmt.Printf("\nNumber of error = %d\t", berc)
	fmt.Printf("BER = %f\n", float64(berc)/(float64(BLOCKLENGTH)*float64(NOBLOCKS)))
	fmt.Printf("Time Elapsed : %v\n", time.Since(startTime))

}
