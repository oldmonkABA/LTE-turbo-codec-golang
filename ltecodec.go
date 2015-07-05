package turbo

import (
	"log"
	"math"

	"github.com/wiless/go.matrix"
	"github.com/wiless/vlib"
)

func init() {

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

	l.ctc_ileave_indices(blocklength)

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

		// paritybits[k] = int(tailstate[0] ^ tailstate[2])
		k += 2
		tailstate[2] = tailstate[1]
		tailstate[1] = tailstate[0]
		tailstate[0] = 0
	}

	return paritybits
}

func (l *LTECodec) ctc_ileave_indices(noind int) {
	l.ILEAVE_SEQ.Resize(noind)
	var f1, f2 int
	//FILE *inter;
	var s1, s2, s3 int
	//implementation of interleaver for LTE..Table[x] of 5.1.3.2.3
	switch noind {
	case 40:
		{
			f1 = 3
			f2 = 10
			break
		}
	case 48:
		{
			f1 = 7
			f2 = 12
			break
		}
	case 56:
		{
			f1 = 19
			f2 = 42
			break
		}
	case 64:
		{
			f1 = 7
			f2 = 16
			break
		}
	case 72:
		{
			f1 = 7
			f2 = 18
			break
		}
	case 80:
		{
			f1 = 11
			f2 = 20
			break
		}
	case 88:
		{
			f1 = 5
			f2 = 22
			break
		}
	case 96:
		{
			f1 = 11
			f2 = 24
			break
		}
	case 104:
		{
			f1 = 7
			f2 = 26
			break
		}
	case 112:
		{
			f1 = 41
			f2 = 84
			break
		}
	case 120:
		{
			f1 = 103
			f2 = 90
			break
		}
	case 128:
		{
			f1 = 15
			f2 = 32
			break
		}
	case 136:
		{
			f1 = 9
			f2 = 34
			break
		}
	case 144:
		{
			f1 = 17
			f2 = 108
			break
		}
	case 152:
		{
			f1 = 9
			f2 = 38
			break
		}
	case 160:
		{
			f1 = 21
			f2 = 120
			break
		}
	case 168:
		{
			f1 = 101
			f2 = 84
			break
		}
	case 176:
		{
			f1 = 21
			f2 = 44
			break
		}
	case 184:
		{
			f1 = 57
			f2 = 46
			break
		}
	case 192:
		{
			f1 = 23
			f2 = 48
			break
		}
	case 200:
		{
			f1 = 13
			f2 = 50
			break
		}
	case 208:
		{
			f1 = 27
			f2 = 52
			break
		}
	case 216:
		{
			f1 = 11
			f2 = 36
			break
		}
	case 224:
		{
			f1 = 27
			f2 = 56
			break
		}
	case 232:
		{
			f1 = 85
			f2 = 58
			break
		}
	case 240:
		{
			f1 = 29
			f2 = 60
			break
		}
	case 248:
		{
			f1 = 33
			f2 = 62
			break
		}
	case 256:
		{
			f1 = 15
			f2 = 32
			break
		}
	case 264:
		{
			f1 = 17
			f2 = 198
			break
		}
	case 272:
		{
			f1 = 33
			f2 = 68
			break
		}
	case 280:
		{
			f1 = 103
			f2 = 210
			break
		}
	case 288:
		{
			f1 = 19
			f2 = 36
			break
		}
	case 296:
		{
			f1 = 19
			f2 = 74
			break
		}
	case 304:
		{
			f1 = 37
			f2 = 76
			break
		}
	case 312:
		{
			f1 = 19
			f2 = 78
			break
		}
	case 320:
		{
			f1 = 21
			f2 = 120
			break
		}
	case 328:
		{
			f1 = 21
			f2 = 82
			break
		}
	case 336:
		{
			f1 = 115
			f2 = 84
			break
		}
	case 344:
		{
			f1 = 193
			f2 = 86
			break
		}
	case 352:
		{
			f1 = 21
			f2 = 44
			break
		}
	case 360:
		{
			f1 = 133
			f2 = 90
			break
		}
	case 368:
		{
			f1 = 81
			f2 = 46
			break
		}
	case 376:
		{
			f1 = 45
			f2 = 94
			break
		}
	case 384:
		{
			f1 = 23
			f2 = 48
			break
		}
	case 392:
		{
			f1 = 243
			f2 = 98
			break
		}
	case 400:
		{
			f1 = 151
			f2 = 40
			break
		}
	case 408:
		{
			f1 = 155
			f2 = 102
			break
		}
	case 416:
		{
			f1 = 25
			f2 = 52
			break
		}
	case 424:
		{
			f1 = 51
			f2 = 106
			break
		}
	case 432:
		{
			f1 = 47
			f2 = 72
			break
		}
	case 440:
		{
			f1 = 91
			f2 = 110
			break
		}
	case 448:
		{
			f1 = 29
			f2 = 168
			break
		}
	case 456:
		{
			f1 = 29
			f2 = 114
			break
		}
	case 464:
		{
			f1 = 247
			f2 = 58
			break
		}
	case 472:
		{
			f1 = 29
			f2 = 118
			break
		}
	case 480:
		{
			f1 = 89
			f2 = 180
			break
		}
	case 488:
		{
			f1 = 91
			f2 = 122
			break
		}
	case 496:
		{
			f1 = 157
			f2 = 62
			break
		}
	case 504:
		{
			f1 = 55
			f2 = 84
			break
		}
	case 512:
		{
			f1 = 31
			f2 = 64
			break
		}
	case 528:
		{
			f1 = 17
			f2 = 66
			break
		}
	case 544:
		{
			f1 = 35
			f2 = 68
			break
		}
	case 560:
		{
			f1 = 227
			f2 = 420
			break
		}
	case 576:
		{
			f1 = 65
			f2 = 96
			break
		}
	case 592:
		{
			f1 = 19
			f2 = 74
			break
		}
	case 608:
		{
			f1 = 37
			f2 = 76
			break
		}
	case 624:
		{
			f1 = 41
			f2 = 234
			break
		}
	case 640:
		{
			f1 = 39
			f2 = 80
			break
		}
	case 656:
		{
			f1 = 185
			f2 = 82
			break
		}
	case 672:
		{
			f1 = 43
			f2 = 252
			break
		}
	case 688:
		{
			f1 = 21
			f2 = 86
			break
		}
	case 704:
		{
			f1 = 155
			f2 = 44
			break
		}
	case 720:
		{
			f1 = 79
			f2 = 120
			break
		}
	case 736:
		{
			f1 = 139
			f2 = 92
			break
		}
	case 752:
		{
			f1 = 23
			f2 = 94
			break
		}
	case 768:
		{
			f1 = 217
			f2 = 48
			break
		}
	case 784:
		{
			f1 = 25
			f2 = 98
			break
		}
	case 800:
		{
			f1 = 17
			f2 = 80
			break
		}
	case 816:
		{
			f1 = 127
			f2 = 102
			break
		}
	case 832:
		{
			f1 = 25
			f2 = 52
			break
		}
	case 848:
		{
			f1 = 239
			f2 = 106
			break
		}
	case 864:
		{
			f1 = 17
			f2 = 48
			break
		}
	case 880:
		{
			f1 = 137
			f2 = 110
			break
		}
	case 896:
		{
			f1 = 215
			f2 = 112
			break
		}
	case 912:
		{
			f1 = 29
			f2 = 114
			break
		}
	case 928:
		{
			f1 = 15
			f2 = 58
			break
		}
	case 944:
		{
			f1 = 147
			f2 = 118
			break
		}
	case 960:
		{
			f1 = 29
			f2 = 60
			break
		}
	case 976:
		{
			f1 = 59
			f2 = 122
			break
		}
	case 992:
		{
			f1 = 65
			f2 = 124
			break
		}
	case 1008:
		{
			f1 = 55
			f2 = 84
			break
		}
	case 1024:
		{
			f1 = 31
			f2 = 64
			break
		}
	case 1056:
		{
			f1 = 17
			f2 = 66
			break
		}
	case 1088:
		{
			f1 = 171
			f2 = 204
			break
		}
	case 1120:
		{
			f1 = 67
			f2 = 140
			break
		}
	case 1152:
		{
			f1 = 35
			f2 = 72
			break
		}
	case 1184:
		{
			f1 = 19
			f2 = 74
			break
		}
	case 1216:
		{
			f1 = 39
			f2 = 76
			break
		}
	case 1248:
		{
			f1 = 19
			f2 = 78
			break
		}
	case 1280:
		{
			f1 = 199
			f2 = 240
			break
		}
	case 1312:
		{
			f1 = 21
			f2 = 82
			break
		}
	case 1344:
		{
			f1 = 211
			f2 = 252
			break
		}
	case 1376:
		{
			f1 = 21
			f2 = 86
			break
		}
	case 1408:
		{
			f1 = 43
			f2 = 88
			break
		}
	case 1440:
		{
			f1 = 149
			f2 = 60
			break
		}
	case 1472:
		{
			f1 = 45
			f2 = 92
			break
		}
	case 1504:
		{
			f1 = 49
			f2 = 846
			break
		}
	case 1536:
		{
			f1 = 71
			f2 = 48
			break
		}
	case 1568:
		{
			f1 = 13
			f2 = 28
			break
		}
	case 1600:
		{
			f1 = 17
			f2 = 80
			break
		}
	case 1632:
		{
			f1 = 25
			f2 = 102
			break
		}
	case 1664:
		{
			f1 = 183
			f2 = 104
			break
		}
	case 1696:
		{
			f1 = 55
			f2 = 954
			break
		}
	case 1728:
		{
			f1 = 127
			f2 = 96
			break
		}
	case 1760:
		{
			f1 = 27
			f2 = 110
			break
		}
	case 1792:
		{
			f1 = 29
			f2 = 112
			break
		}
	case 1824:
		{
			f1 = 29
			f2 = 114
			break
		}
	case 1856:
		{
			f1 = 57
			f2 = 116
			break
		}
	case 1888:
		{
			f1 = 45
			f2 = 354
			break
		}
	case 1920:
		{
			f1 = 31
			f2 = 120
			break
		}
	case 1952:
		{
			f1 = 59
			f2 = 610
			break
		}
	case 1984:
		{
			f1 = 185
			f2 = 124
			break
		}
	case 2016:
		{
			f1 = 113
			f2 = 420
			break
		}
	case 2048:
		{
			f1 = 31
			f2 = 64
			break
		}
	case 2112:
		{
			f1 = 17
			f2 = 66
			break
		}
	case 2176:
		{
			f1 = 171
			f2 = 136
			break
		}
	case 2240:
		{
			f1 = 209
			f2 = 420
			break
		}
	case 2304:
		{
			f1 = 253
			f2 = 216
			break
		}
	case 2368:
		{
			f1 = 367
			f2 = 444
			break
		}
	case 2432:
		{
			f1 = 265
			f2 = 456
			break
		}
	case 2496:
		{
			f1 = 181
			f2 = 468
			break
		}
	case 2560:
		{
			f1 = 39
			f2 = 80
			break
		}
	case 2624:
		{
			f1 = 27
			f2 = 164
			break
		}
	case 2688:
		{
			f1 = 127
			f2 = 504
			break
		}
	case 2752:
		{
			f1 = 143
			f2 = 172
			break
		}
	case 2816:
		{
			f1 = 43
			f2 = 88
			break
		}
	case 2880:
		{
			f1 = 29
			f2 = 300
			break
		}
	case 2944:
		{
			f1 = 45
			f2 = 92
			break
		}
	case 3008:
		{
			f1 = 157
			f2 = 188
			break
		}
	case 3072:
		{
			f1 = 47
			f2 = 96
			break
		}
	case 3136:
		{
			f1 = 13
			f2 = 28
			break
		}
	case 3200:
		{
			f1 = 111
			f2 = 240
			break
		}
	case 3264:
		{
			f1 = 443
			f2 = 204
			break
		}
	case 3328:
		{
			f1 = 51
			f2 = 104
			break
		}
	case 3392:
		{
			f1 = 51
			f2 = 212
			break
		}
	case 3456:
		{
			f1 = 451
			f2 = 192
			break
		}
	case 3520:
		{
			f1 = 257
			f2 = 220
			break
		}
	case 3584:
		{
			f1 = 57
			f2 = 336
			break
		}
	case 3648:
		{
			f1 = 313
			f2 = 228
			break
		}
	case 3712:
		{
			f1 = 271
			f2 = 232
			break
		}
	case 3776:
		{
			f1 = 179
			f2 = 236
			break
		}
	case 3840:
		{
			f1 = 331
			f2 = 120
			break
		}
	case 3904:
		{
			f1 = 363
			f2 = 244
			break
		}
	case 3968:
		{
			f1 = 375
			f2 = 248
			break
		}
	case 4032:
		{
			f1 = 127
			f2 = 168
			break
		}
	case 4096:
		{
			f1 = 31
			f2 = 64
			break
		}
	case 4160:
		{
			f1 = 33
			f2 = 130
			break
		}
	case 4224:
		{
			f1 = 43
			f2 = 264
			break
		}
	case 4288:
		{
			f1 = 33
			f2 = 134
			break
		}
	case 4352:
		{
			f1 = 477
			f2 = 408
			break
		}
	case 4416:
		{
			f1 = 35
			f2 = 138
			break
		}
	case 4480:
		{
			f1 = 233
			f2 = 280
			break
		}
	case 4544:
		{
			f1 = 357
			f2 = 142
			break
		}
	case 4608:
		{
			f1 = 337
			f2 = 480
			break
		}
	case 4672:
		{
			f1 = 37
			f2 = 146
			break
		}
	case 4736:
		{
			f1 = 71
			f2 = 144
			break
		}
	case 4800:
		{
			f1 = 71
			f2 = 120
			break
		}
	case 4864:
		{
			f1 = 37
			f2 = 152
			break
		}
	case 4928:
		{
			f1 = 39
			f2 = 462
			break
		}
	case 4992:
		{
			f1 = 127
			f2 = 234
			break
		}
	case 5056:
		{
			f1 = 39
			f2 = 158
			break
		}
	case 5120:
		{
			f1 = 39
			f2 = 80
			break
		}
	case 5184:
		{
			f1 = 31
			f2 = 96
			break
		}
	case 5248:
		{
			f1 = 113
			f2 = 902
			break
		}
	case 5312:
		{
			f1 = 41
			f2 = 166
			break
		}
	case 5376:
		{
			f1 = 251
			f2 = 336
			break
		}
	case 5440:
		{
			f1 = 43
			f2 = 170
			break
		}
	case 5504:
		{
			f1 = 21
			f2 = 86
			break
		}
	case 5568:
		{
			f1 = 43
			f2 = 174
			break
		}
	case 5632:
		{
			f1 = 45
			f2 = 176
			break
		}
	case 5696:
		{
			f1 = 45
			f2 = 178
			break
		}
	case 5760:
		{
			f1 = 161
			f2 = 120
			break
		}
	case 5824:
		{
			f1 = 89
			f2 = 182
			break
		}
	case 5888:
		{
			f1 = 323
			f2 = 184
			break
		}
	case 5952:
		{
			f1 = 47
			f2 = 186
			break
		}
	case 6016:
		{
			f1 = 23
			f2 = 94
			break
		}
	case 6080:
		{
			f1 = 47
			f2 = 190
			break
		}
	case 6144:
		{
			f1 = 263
			f2 = 480
			break
		}
	default:
		{
			log.Panicf("ctc_ileave_indices(...):unsupported block size for no HARQ CTC,%d\n", noind)
			//exit(0)
		}
	}
	for j := 0; j < noind; j++ {
		s1 = (f1 * j) % noind
		s2 = (j * j) % noind
		s3 = (s2 * f2) % noind
		l.ILEAVE_SEQ[j] = (s1 + s3) % noind
	}
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

				la.Set(k, i, vlib.Min(xtemp))
			}
		}
		min_m = vlib.Min(la.GetRowVector(k).Array())
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

		min_m = vlib.Min(lb.GetRowVector(k).Array())

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
		le.Set(k, 0, vlib.Min(temp_0))
		le.Set(k, 1, vlib.Min(temp_1))
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
