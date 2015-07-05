package turbo

import (
	"math"

	"github.com/wiless/go.matrix"
	"github.com/wiless/vlib"
)

func init() {

	// for j := 0; j < noind; j++ {
	// 	s1 = (f1 * j) % noind
	// 	s2 = (j * j) % noind
	// 	s3 = (s2 * f2) % noind
	// 	l.ILEAVE_SEQ[j] = (s1 + s3) % noind
	// }
}

type LTECodec struct {
	MODULATION  string
	CODERATE    string
	BLOCKLENGTH int
	TAILBITS    int
	ITERATIONS  int
	PUNCTURE    vlib.MatrixF // for rate matching and puncturing
	PUNCBITS    int
	NOENCBITS   int
	CRATE       float64
	ILEAVE_SEQ  vlib.VectorI

	//tr_nextstate [][]int
	tr_nextstate []int
	//Xtr_nextstate [][]int

	// vlib.MatrixF

	// var tr_output vlib.MatrixF
	//tr_output [][]int
	tr_output []int
}

func (l *LTECodec) Init(mod string, coderate string, blocklength int) {
	l.MODULATION = mod
	l.CODERATE = coderate
	l.BLOCKLENGTH = blocklength
	l.TAILBITS = 12
	l.ITERATIONS = 8
	//l.tr_nextstate = [][]int{{0, 4}, {4, 0}, {5, 1}, {1, 5}, {2, 6}, {6, 2}, {7, 3}, {3, 7}}
	// tr_output = vlib.MatrixF{{0, 1}, {0, 1}, {1, 0}, {1, 0}, {1, 0}, {1, 0}, {0, 1}, {0, 1}}
	//l.tr_output = [][]int{{0, 1}, {0, 1}, {1, 0}, {1, 0}, {1, 0}, {1, 0}, {0, 1}, {0, 1}}

	l.tr_nextstate = []int{0, 4, 5, 1, 2, 6, 7, 3, 4, 0, 1, 5, 6, 2, 3, 7}
	l.tr_output = []int{0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 1, 1}
	// l.tr_output = []int{0, 0, 1, 0}, {1, 0}, {1, 0}, {1, 0}, {0, 1}, {0, 1}}
	l.ILEAVE_SEQ = InterleaverLTE(blocklength)

	// l.ctc_ileave_indices(blocklength)
}

func (l LTECodec) encoder(ipdata []int) []int {
	// var tailbits vlib.VectorB
	var next_state, nobits int
	// var bit int

	nobits = len(ipdata)

	paritybits := make([]int, nobits+6)
	next_state = 0
	for i := 0; i < nobits; i++ {

		indx := ipdata[i]*8 + next_state
		next_state = l.tr_nextstate[indx]
		paritybits[i] = l.tr_output[indx]

	}

	tailstate := vlib.NewVectorI(3)

	tailstate[2] = (next_state & 1)
	tailstate[1] = (next_state & 2) >> 1
	tailstate[0] = (next_state & 4) >> 2

	k := nobits
	for i := nobits; i < nobits+3; i++ {
		paritybits[k], paritybits[k+1] = int(tailstate[1]^tailstate[2]), tailstate[0]^tailstate[2]

		k += 2
		tailstate[2] = tailstate[1]
		tailstate[1] = tailstate[0]
		tailstate[0] = 0
	}

	return paritybits
}

func (l *LTECodec) Encode(ipbits []int) []int {

	var LEN int = len(ipbits)
	ileaveipbits := make([]int, LEN)
	for i := 0; i < LEN; i++ {
		ileaveipbits[i] = ipbits[l.ILEAVE_SEQ[i]]
	}
	var parity1 []int
	parity1 = l.encoder(ipbits)
	parity2 := l.encoder(ileaveipbits)

	n := 3*LEN + 12
	codedbits := make([]int, n)
	for bit := 0; bit < LEN; bit++ {
		codedbits[3*bit] = ipbits[bit]
		codedbits[3*bit+1] = parity1[bit]
		codedbits[3*bit+2] = parity2[bit]
	}

	TAILS := 6
	begin := 3 * LEN
	for bit := LEN; bit < (LEN + TAILS); bit++ {
		codedbits[begin] = parity1[bit]
		begin++
	}

	begin = 3*LEN + TAILS
	for bit := LEN; bit < (LEN + TAILS); bit++ {
		codedbits[begin] = parity2[bit]
		begin++
	}

	return codedbits
}

func (l *LTECodec) demuxsoftval(softval []float64) (sysllhd1 []float64, lg1 []float64, sysllhd2 []float64, lg2 []float64) {
	bl := l.BLOCKLENGTH + 4
	sysllhd1 = make([]float64, bl)
	lg1 = make([]float64, bl)
	lg2 = make([]float64, bl)
	for i := 0; i < l.BLOCKLENGTH; i++ {
		sysllhd1[i+1] = softval[3*i]
		lg1[i+1] = softval[3*i+1]
		lg2[i+1] = softval[3*i+2]

	}
	sysllhd2 = make([]float64, bl)
	for i := 0; i < l.BLOCKLENGTH; i++ {
		sysllhd2[i+1] = sysllhd1[l.ILEAVE_SEQ[i]+1]
	}
	begin := 3 * l.BLOCKLENGTH
	for i := l.BLOCKLENGTH; i < (l.BLOCKLENGTH + 3); i++ {

		sysllhd1[i+1] = softval[begin]
		begin++
		lg1[i+1] = softval[begin]
		begin++
	}

	begin = 3*l.BLOCKLENGTH + 6
	for i := l.BLOCKLENGTH; i < (l.BLOCKLENGTH + 3); i++ {

		sysllhd2[i+1] = softval[begin]
		begin++
		lg2[i+1] = softval[begin]
		begin++
	}

	return sysllhd1, lg1, sysllhd2, lg2
}

func (l *LTECodec) branch_metric(sysllhd []float64, parllhd []float64, le *matrix.DenseMatrix, lgr *matrix.DenseMatrix) {

	LEN := len(sysllhd)

	var value float64
	for i := 1; i < LEN; i++ {

		for j := 0; j < 4; j++ {

			if j < 2 {
				value = sysllhd[i] + le.Get(i, 0) + float64((j+1)%2)*parllhd[i]

			} else {

				value = le.Get(i, 1) + float64((j+1)%2)*parllhd[i]

			}
			lgr.Set(i, j, value)
		}
	}

}

func (l *LTECodec) forward_metric(lgr *matrix.DenseMatrix, la *matrix.DenseMatrix) {

	alpha_t := [][]int{{0, 1}, {3, 2}, {4, 5}, {7, 6}, {1, 0}, {2, 3}, {5, 4}, {6, 7}}
	br_metric := [][]int{{0, 3}, {1, 2}, {1, 2}, {0, 3}, {0, 3}, {1, 2}, {1, 2}, {0, 3}}
	xtemp := vlib.NewVectorF(2)

	var min_m float64
	LEN := lgr.Rows() - 1
	la.Set(0, 0, 0)

	for i := 1; i < 8; i++ {
		la.Set(0, i, 1e8)
	}

	for k := 1; k <= LEN; k++ {
		for i := 0; i < 8; i++ {

			for j := 0; j < 2; j++ {
				xtemp[j] = la.Get(k-1, alpha_t[i][j]) + lgr.Get(k, br_metric[i][j])

				minval := xtemp[0]
				if xtemp[1] < minval {
					minval = xtemp[1]
				}
				// la.Set(k, i, vlib.Min(xtemp))
				la.Set(k, i, minval)
			}
		}

		// min_m = vlib.Min(la.GetRowVector(k).Array())
		Ncols := la.Cols()
		min_m = 10.0e8
		for i := 0; i < Ncols; i++ {
			x := la.Get(k, i)
			if min_m > x {
				min_m = x
			}
		}

		for i := 0; i < 8; i++ {
			num := 2

			la.Set(k, i, la.Get(k, i)-min_m)
			if k > LEN-3 {

				diff := LEN - k
				if i >= (int)(num*diff) {
					la.Set(k, i, 1e8)
				}

			}
		}

	}

}

func (l *LTECodec) backward_metric(lgr *matrix.DenseMatrix, lb *matrix.DenseMatrix) {

	beta_t := [][]int{{0, 4}, {4, 0}, {5, 1}, {1, 5}, {2, 6}, {6, 2}, {7, 3}, {3, 7}}
	br_metric := [][]int{{0, 3}, {0, 3}, {1, 2}, {1, 2}, {1, 2}, {1, 2}, {0, 3}, {0, 3}}
	temp := vlib.NewVectorF(2)
	var min_m float64
	LEN := lgr.Rows() - 1
	var num int
	lb.Set(LEN, 0, 0)
	for i := 1; i < 8; i++ {
		lb.Set(LEN, i, 1e8)
	}
	for k := LEN - 1; k >= 0; k-- {
		for i := 0; i < 8; i++ {
			for j := 0; j < 2; j++ {
				temp[j] = lb.Get(k+1, beta_t[i][j]) + lgr.Get(k+1, br_metric[i][j])
			}
			lb.Set(k, i, vlib.Min(temp))
		}

		Ncols := lb.Cols()
		min_m = 10.0e8
		for i := 0; i < Ncols; i++ {
			x := lb.Get(k, i)
			if min_m > x {
				min_m = x
			}
		}

		// xx := lb.GetRowVector(k)
		// if xx == nil {
		// 	log.Panicf("Received Empty GetROW VECTOR")
		// }
		// yy := xx.Array()

		// if yy == nil {
		// 	log.Panicf("Received Empty Array VECTOR")
		// }
		// min_m = vlib.Min(yy)

		for i := 0; i < 8; i++ {
			num = 2
			lb.Set(k, i, lb.Get(k, i)-min_m)
			if k < 3 {
				diff := float64(2 - k)
				num = num * int(math.Pow(2, diff))
				if i%num != 0 {
					lb.Set(k, i, 1e8)
				}
			}
		}

	}
}
func (l *LTECodec) extrinsic_info(alpha *matrix.DenseMatrix, beta *matrix.DenseMatrix, parllhd []float64, le *matrix.DenseMatrix) {

	LE_0 := [][]int{{0, 0, 0}, {1, 4, 0}, {2, 5, 1}, {3, 1, 1}, {4, 2, 1}, {5, 6, 1}, {6, 7, 0}, {7, 3, 0}}
	LE_1 := [][]int{{0, 4, 1}, {1, 0, 1}, {2, 1, 0}, {3, 5, 0}, {4, 6, 0}, {5, 2, 0}, {6, 3, 1}, {7, 7, 1}}
	LEN := alpha.Rows() - 4
	temp_0, temp_1 := vlib.NewVectorF(8), vlib.NewVectorF(8)

	for k := 1; k <= LEN; k++ {
		for i := 0; i < 8; i++ {
			temp_0[i] = alpha.Get(k-1, LE_0[i][0]) + beta.Get(k, LE_0[i][1]) + float64((LE_0[i][2]+1)%2)*parllhd[k]
			temp_1[i] = alpha.Get(k-1, LE_1[i][0]) + beta.Get(k, LE_1[i][1]) + float64((LE_1[i][2]+1)%2)*parllhd[k]
		}
		le.Set(k, 0, temp_0.Min())
		le.Set(k, 1, temp_1.Min())
	}
}

func (l *LTECodec) ctc_ileave(ipmat *matrix.DenseMatrix, opmat *matrix.DenseMatrix) {
	ROWS := ipmat.Rows()
	COLS := ipmat.Cols()
	for j := 0; j < ROWS-4; j++ {
		rowvec := ipmat.GetRowVector(l.ILEAVE_SEQ[j] + 1)
		for k := 0; k < COLS; k++ {
			opmat.Set(j+1, k, rowvec.Get(0, k))
		}
	}
}

func (l *LTECodec) ctc_deileave(ipmat *matrix.DenseMatrix, opmat *matrix.DenseMatrix) {

	ROWS := ipmat.Rows()
	COLS := ipmat.Cols()

	for j := 0; j < ROWS-4; j++ {
		rowvec := ipmat.GetRowVector(j + 1)
		for k := 0; k < COLS; k++ {
			opmat.Set(l.ILEAVE_SEQ[j]+1, k, rowvec.Get(0, k))
		}
	}

}

func (l *LTECodec) lambda_calc(alpha *matrix.DenseMatrix, beta *matrix.DenseMatrix, gamma *matrix.DenseMatrix, lambda *matrix.DenseMatrix) {
	La_0 := [][]int{{0, 0, 0}, {1, 4, 0}, {2, 5, 1}, {3, 1, 1}, {4, 2, 1}, {5, 6, 1}, {6, 7, 0}, {7, 3, 0}}
	La_1 := [][]int{{0, 4, 3}, {1, 0, 3}, {2, 1, 2}, {3, 5, 2}, {4, 6, 2}, {5, 2, 2}, {6, 3, 3}, {7, 7, 3}}
	temp_0, temp_1 := vlib.NewVectorF(8), vlib.NewVectorF(8)
	LEN := alpha.Rows() - 1

	for k := 1; k <= LEN; k++ {
		for i := 0; i < 8; i++ {
			temp_0[i] = alpha.Get(k-1, La_0[i][0]) + beta.Get(k, La_0[i][1]) + gamma.Get(k, La_0[i][2])
			temp_1[i] = alpha.Get(k-1, La_1[i][0]) + beta.Get(k, La_1[i][1]) + gamma.Get(k, La_1[i][2])
		}
		lambda.Set(k, 0, vlib.Min(temp_0))
		lambda.Set(k, 1, vlib.Min(temp_1))
	}
}

func (l *LTECodec) Decode(softval []float64) (decodedbits []int) {
	TRE_LENGTH := l.BLOCKLENGTH + 4
	la := matrix.Zeros(TRE_LENGTH, 8)
	lb := matrix.Zeros(TRE_LENGTH, 8)
	lmda := matrix.Zeros(TRE_LENGTH, 2)
	lgr := matrix.Zeros(TRE_LENGTH, 4)
	lei1 := matrix.Zeros(TRE_LENGTH, 2)
	lei2 := matrix.Zeros(TRE_LENGTH, 2)

	leo1 := matrix.Zeros(TRE_LENGTH, 2)
	leo2 := matrix.Zeros(TRE_LENGTH, 2)
	var sub float64
	sysl, lg1, sysil, lg2 := l.demuxsoftval(softval)
	decintbits := make([]int, l.BLOCKLENGTH)
	//fmt.Printf("\nSys1 = %v", sysl)
	//fmt.Printf("\nSys2 = %v ", sysil)
	//fmt.Printf("\nlg1 = %v", lg1)
	//fmt.Printf("\nlg2 = %v\n", lg2)
	//fmt.Printf("%v\n", lgr)
	//fmt.Print(la)
	//fmt.Print(lb)
	//fmt.Print(leo1)
	//fmt.Print(lei2)
	for it := 0; it < l.ITERATIONS; it++ {
		// Decoder -1
		l.branch_metric(sysl, lg1, lei1, lgr)
		l.forward_metric(lgr, la)
		l.backward_metric(lgr, lb)
		l.extrinsic_info(la, lb, lg1, leo1)
		l.ctc_ileave(leo1, lei2)
		// Decoder -2
		l.branch_metric(sysil, lg2, lei2, lgr)
		l.forward_metric(lgr, la)
		l.backward_metric(lgr, lb)
		l.extrinsic_info(la, lb, lg2, leo2)
		l.ctc_deileave(leo2, lei1)
	}
	l.lambda_calc(la, lb, lgr, lmda)
	//fmt.Print(lmda)
	for i := 0; i < l.BLOCKLENGTH; i++ {
		sub = lmda.Get(i+1, 0) - lmda.Get(i+1, 1)
		if sub > 0 {
			decintbits[i] = 1
		}
	}
	//fmt.Print(decintbits)
	decodedbits = make([]int, l.BLOCKLENGTH)
	for i := 0; i < l.BLOCKLENGTH; i++ {
		decodedbits[l.ILEAVE_SEQ[i]] = decintbits[i]
	}
	return decodedbits
}
