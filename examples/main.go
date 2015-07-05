package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/oldmonkABA/turbo"
	"github.com/wiless/vlib"
)

func bpskmod(ipbits []int) (modsyms []int) {
	n := len(ipbits)
	modsyms = make([]int, n)
	for i := 0; i < n; i++ {
		modsyms[i] = 2*ipbits[i] - 1

	}
	return modsyms
}

func bpskawgndemod(recsym []int, variance float64) (llr []float64) {
	n := len(recsym)
	llr = make([]float64, n)
	for i := 0; i < n; i++ {
		llr[i] = float64(4*recsym[i]) / variance

	}
	return llr
}

func main() {
	seedvalue := time.Now().Unix()
	rand.Seed(seedvalue)
	runtime.GOMAXPROCS(0)
	BLOCKLENGTH := 40
	NOBLOCKS := 1000
	fmt.Printf("BLOCKLENGTH = %d \t NOBLOCKS = %d\n", BLOCKLENGTH, NOBLOCKS)
	//var output vlib.VectorB
	startTime := time.Now()
	var codec turbo.LTECodec
	berc := 0
	codec.Init("BPSK", "_13_", BLOCKLENGTH)
	variance := 1.0
	for i := 1; i <= NOBLOCKS; i++ {
		fmt.Printf("\rProcessing block no : %d", i)
		input := vlib.Randsrc(BLOCKLENGTH, 2)
		output := codec.Encode(input)
		mod := bpskmod(output)
		llr := bpskawgndemod(mod, variance)
		dec_output := codec.Decode(llr)
		//fmt.Print("Input bits = %v", input)
		//fmt.Print("\nOuput bits = %v", dec_output)
		for j := 0; j < codec.BLOCKLENGTH; j++ {
			berc = berc + input[j] ^ dec_output[j]

		}

	} //
	fmt.Printf("\nNumber of error = %d\t", berc)
	fmt.Printf("BER = %f\n", float64(berc)/(float64(BLOCKLENGTH)*float64(NOBLOCKS)))
	fmt.Printf("Time Elapsed : %v\n", time.Since(startTime))

}
