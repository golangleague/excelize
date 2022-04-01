package excelize

import (
	"container/list"
	"math"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/efp"
)

func prepareCalcData(cellData [][]interface{}) *File {
	f := NewFile()
	for r, row := range cellData {
		for c, value := range row {
			cell, _ := CoordinatesToCellName(c+1, r+1)
			_ = f.SetCellValue("Sheet1", cell, value)
		}
	}
	return f
}

func TestCalcCellValue(t *testing.T) {
	cellData := [][]interface{}{
		{1, 4, nil, "Month", "Team", "Sales"},
		{2, 5, nil, "Jan", "North 1", 36693},
		{3, nil, nil, "Jan", "North 2", 22100},
		{0, nil, nil, "Jan", "South 1", 53321},
		{nil, nil, nil, "Jan", "South 2", 34440},
		{nil, nil, nil, "Feb", "North 1", 29889},
		{nil, nil, nil, "Feb", "North 2", 50090},
		{nil, nil, nil, "Feb", "South 1", 32080},
		{nil, nil, nil, "Feb", "South 2", 45500},
	}
	mathCalc := map[string]string{
		"=2^3":            "8",
		"=1=1":            "TRUE",
		"=1=2":            "FALSE",
		"=1<2":            "TRUE",
		"=3<2":            "FALSE",
		"=1<\"-1\"":       "TRUE",
		"=\"-1\"<1":       "FALSE",
		"=\"-1\"<\"-2\"":  "TRUE",
		"=2<=3":           "TRUE",
		"=2<=1":           "FALSE",
		"=1<=\"-1\"":      "TRUE",
		"=\"-1\"<=1":      "FALSE",
		"=\"-1\"<=\"-2\"": "TRUE",
		"=2>1":            "TRUE",
		"=2>3":            "FALSE",
		"=1>\"-1\"":       "FALSE",
		"=\"-1\">-1":      "TRUE",
		"=\"-1\">\"-2\"":  "FALSE",
		"=2>=1":           "TRUE",
		"=2>=3":           "FALSE",
		"=1>=\"-1\"":      "FALSE",
		"=\"-1\">=-1":     "TRUE",
		"=\"-1\">=\"-2\"": "FALSE",
		"=1&2":            "12",
		"=15%":            "0.15",
		"=1+20%":          "1.2",
		`="A"="A"`:        "TRUE",
		`="A"<>"A"`:       "FALSE",
		// Engineering Functions
		// BESSELI
		"=BESSELI(4.5,1)": "15.3892227537359",
		"=BESSELI(32,1)":  "5502845511211.25",
		// BESSELJ
		"=BESSELJ(1.9,2)": "0.329925727692387",
		// BESSELK
		"=BESSELK(0.05,0)": "3.11423403428966",
		"=BESSELK(0.05,1)": "19.9096743272486",
		"=BESSELK(0.05,2)": "799.501207124235",
		"=BESSELK(3,2)":    "0.0615104585619118",
		// BESSELY
		"=BESSELY(0.05,0)": "-1.97931100684153",
		"=BESSELY(0.05,1)": "-12.789855163794",
		"=BESSELY(0.05,2)": "-509.61489554492",
		"=BESSELY(9,2)":    "-0.229082087487741",
		// BIN2DEC
		"=BIN2DEC(\"10\")":         "2",
		"=BIN2DEC(\"11\")":         "3",
		"=BIN2DEC(\"0000000010\")": "2",
		"=BIN2DEC(\"1111111110\")": "-2",
		"=BIN2DEC(\"110\")":        "6",
		// BIN2HEX
		"=BIN2HEX(\"10\")":         "2",
		"=BIN2HEX(\"0000000001\")": "1",
		"=BIN2HEX(\"10\",10)":      "0000000002",
		"=BIN2HEX(\"1111111110\")": "FFFFFFFFFE",
		"=BIN2HEX(\"11101\")":      "1D",
		// BIN2OCT
		"=BIN2OCT(\"101\")":        "5",
		"=BIN2OCT(\"0000000001\")": "1",
		"=BIN2OCT(\"10\",10)":      "0000000002",
		"=BIN2OCT(\"1111111110\")": "7777777776",
		"=BIN2OCT(\"1110\")":       "16",
		// BITAND
		"=BITAND(13,14)": "12",
		// BITLSHIFT
		"=BITLSHIFT(5,2)": "20",
		"=BITLSHIFT(3,5)": "96",
		// BITOR
		"=BITOR(9,12)": "13",
		// BITRSHIFT
		"=BITRSHIFT(20,2)": "5",
		"=BITRSHIFT(52,4)": "3",
		// BITXOR
		"=BITXOR(5,6)":  "3",
		"=BITXOR(9,12)": "5",
		// COMPLEX
		"=COMPLEX(5,2)":         "5+2i",
		"=COMPLEX(5,-9)":        "5-9i",
		"=COMPLEX(-1,2,\"j\")":  "-1+2j",
		"=COMPLEX(10,-5,\"i\")": "10-5i",
		"=COMPLEX(0,5)":         "5i",
		"=COMPLEX(3,0)":         "3",
		"=COMPLEX(0,-2)":        "-2i",
		"=COMPLEX(0,0)":         "0",
		"=COMPLEX(0,-1,\"j\")":  "-j",
		// DEC2BIN
		"=DEC2BIN(2)":    "10",
		"=DEC2BIN(3)":    "11",
		"=DEC2BIN(2,10)": "0000000010",
		"=DEC2BIN(-2)":   "1111111110",
		"=DEC2BIN(6)":    "110",
		// DEC2HEX
		"=DEC2HEX(10)":    "A",
		"=DEC2HEX(31)":    "1F",
		"=DEC2HEX(16,10)": "0000000010",
		"=DEC2HEX(-16)":   "FFFFFFFFF0",
		"=DEC2HEX(273)":   "111",
		// DEC2OCT
		"=DEC2OCT(8)":    "10",
		"=DEC2OCT(18)":   "22",
		"=DEC2OCT(8,10)": "0000000010",
		"=DEC2OCT(-8)":   "7777777770",
		"=DEC2OCT(237)":  "355",
		// DELTA
		"=DELTA(5,4)":       "0",
		"=DELTA(1.00001,1)": "0",
		"=DELTA(1.23,1.23)": "1",
		"=DELTA(1)":         "0",
		"=DELTA(0)":         "1",
		// ERF
		"=ERF(1.5)":   "0.966105146475311",
		"=ERF(0,1.5)": "0.966105146475311",
		"=ERF(1,2)":   "0.152621472069238",
		// ERF.PRECISE
		"=ERF.PRECISE(-1)":  "-0.842700792949715",
		"=ERF.PRECISE(1.5)": "0.966105146475311",
		// ERFC
		"=ERFC(0)":   "1",
		"=ERFC(0.5)": "0.479500122186953",
		"=ERFC(-1)":  "1.84270079294971",
		// ERFC.PRECISE
		"=ERFC.PRECISE(0)":   "1",
		"=ERFC.PRECISE(0.5)": "0.479500122186953",
		"=ERFC.PRECISE(-1)":  "1.84270079294971",
		// GESTEP
		"=GESTEP(1.2,0.001)":  "1",
		"=GESTEP(0.05,0.05)":  "1",
		"=GESTEP(-0.00001,0)": "0",
		"=GESTEP(-0.00001)":   "0",
		// HEX2BIN
		"=HEX2BIN(\"2\")":          "10",
		"=HEX2BIN(\"0000000001\")": "1",
		"=HEX2BIN(\"2\",10)":       "0000000010",
		"=HEX2BIN(\"F0\")":         "11110000",
		"=HEX2BIN(\"1D\")":         "11101",
		// HEX2DEC
		"=HEX2DEC(\"A\")":          "10",
		"=HEX2DEC(\"1F\")":         "31",
		"=HEX2DEC(\"0000000010\")": "16",
		"=HEX2DEC(\"FFFFFFFFF0\")": "-16",
		"=HEX2DEC(\"111\")":        "273",
		"=HEX2DEC(\"\")":           "0",
		// HEX2OCT
		"=HEX2OCT(\"A\")":          "12",
		"=HEX2OCT(\"000000000F\")": "17",
		"=HEX2OCT(\"8\",10)":       "0000000010",
		"=HEX2OCT(\"FFFFFFFFF8\")": "7777777770",
		"=HEX2OCT(\"1F3\")":        "763",
		// IMABS
		"=IMABS(\"2j\")":              "2",
		"=IMABS(\"-1+2i\")":           "2.23606797749979",
		"=IMABS(COMPLEX(-1,2,\"j\"))": "2.23606797749979",
		// IMAGINARY
		"=IMAGINARY(\"5+2i\")": "2",
		"=IMAGINARY(\"2-i\")":  "-1",
		"=IMAGINARY(6)":        "0",
		"=IMAGINARY(\"3i\")":   "3",
		"=IMAGINARY(\"4+i\")":  "1",
		// IMARGUMENT
		"=IMARGUMENT(\"5+2i\")": "0.380506377112365",
		"=IMARGUMENT(\"2-i\")":  "-0.463647609000806",
		"=IMARGUMENT(6)":        "0",
		// IMCONJUGATE
		"=IMCONJUGATE(\"5+2i\")": "5-2i",
		"=IMCONJUGATE(\"2-i\")":  "2+i",
		"=IMCONJUGATE(6)":        "6",
		"=IMCONJUGATE(\"3i\")":   "-3i",
		"=IMCONJUGATE(\"4+i\")":  "4-i",
		// IMCOS
		"=IMCOS(0)":          "1",
		"=IMCOS(0.5)":        "0.877582561890373",
		"=IMCOS(\"3+0.5i\")": "-1.11634124452615-0.0735369737112366i",
		// IMCOSH
		"=IMCOSH(0.5)":           "1.12762596520638",
		"=IMCOSH(\"3+0.5i\")":    "8.83520460650099+4.80282508274303i",
		"=IMCOSH(\"2-i\")":       "2.03272300701967-3.0518977991518i",
		"=IMCOSH(COMPLEX(1,-1))": "0.833730025131149-0.988897705762865i",
		// IMCOT
		"=IMCOT(0.5)":           "1.83048772171245",
		"=IMCOT(\"3+0.5i\")":    "-0.479345578747373-2.01609252150623i",
		"=IMCOT(\"2-i\")":       "-0.171383612909185+0.821329797493852i",
		"=IMCOT(COMPLEX(1,-1))": "0.217621561854403+0.868014142895925i",
		// IMCSC
		"=IMCSC(\"j\")": "-0.850918128239322j",
		// IMCSCH
		"=IMCSCH(COMPLEX(1,-1))": "0.303931001628426+0.621518017170428i",
		// IMDIV
		"=IMDIV(\"5+2i\",\"1+i\")":          "3.5-1.5i",
		"=IMDIV(\"2+2i\",\"2+i\")":          "1.2+0.4i",
		"=IMDIV(COMPLEX(5,2),COMPLEX(0,1))": "2-5i",
		// IMEXP
		"=IMEXP(0)":             "1",
		"=IMEXP(0.5)":           "1.64872127070013",
		"=IMEXP(\"1-2i\")":      "-1.13120438375681-2.47172667200482i",
		"=IMEXP(COMPLEX(1,-1))": "1.46869393991589-2.28735528717884i",
		// IMLN
		"=IMLN(0.5)":           "-0.693147180559945",
		"=IMLN(\"3+0.5i\")":    "1.11231177576217+0.165148677414627i",
		"=IMLN(\"2-i\")":       "0.80471895621705-0.463647609000806i",
		"=IMLN(COMPLEX(1,-1))": "0.346573590279973-0.785398163397448i",
		// IMLOG10
		"=IMLOG10(0.5)":           "-0.301029995663981",
		"=IMLOG10(\"3+0.5i\")":    "0.483070866369516+0.0717231592947926i",
		"=IMLOG10(\"2-i\")":       "0.349485002168009-0.201359598136687i",
		"=IMLOG10(COMPLEX(1,-1))": "0.150514997831991-0.34109408846046i",
		// IMREAL
		"=IMREAL(\"5+2i\")":     "5",
		"=IMREAL(\"2+2i\")":     "2",
		"=IMREAL(6)":            "6",
		"=IMREAL(\"3i\")":       "0",
		"=IMREAL(COMPLEX(4,1))": "4",
		// IMSEC
		"=IMSEC(0.5)":           "1.13949392732455",
		"=IMSEC(\"3+0.5i\")":    "-0.89191317974033+0.0587531781817398i",
		"=IMSEC(\"2-i\")":       "-0.41314934426694-0.687527438655479i",
		"=IMSEC(COMPLEX(1,-1))": "0.498337030555187-0.591083841721045i",
		// IMSECH
		"=IMSECH(0.5)":           "0.886818883970074",
		"=IMSECH(\"3+0.5i\")":    "0.0873665779621303-0.0474925494901607i",
		"=IMSECH(\"2-i\")":       "0.151176298265577+0.226973675393722i",
		"=IMSECH(COMPLEX(1,-1))": "0.498337030555187+0.591083841721045i",
		// IMSIN
		"=IMSIN(0.5)":           "0.479425538604203",
		"=IMSIN(\"3+0.5i\")":    "0.15913058529844-0.515880442452527i",
		"=IMSIN(\"2-i\")":       "1.40311925062204+0.489056259041294i",
		"=IMSIN(COMPLEX(1,-1))": "1.29845758141598-0.634963914784736i",
		// IMSINH
		"=IMSINH(-0)":            "0",
		"=IMSINH(0.5)":           "0.521095305493747",
		"=IMSINH(\"3+0.5i\")":    "8.79151234349371+4.82669427481082i",
		"=IMSINH(\"2-i\")":       "1.95960104142161-3.16577851321617i",
		"=IMSINH(COMPLEX(1,-1))": "0.634963914784736-1.29845758141598i",
		// IMSQRT
		"=IMSQRT(\"i\")":     "0.707106781186548+0.707106781186548i",
		"=IMSQRT(\"2-i\")":   "1.45534669022535-0.343560749722512i",
		"=IMSQRT(\"5+2i\")":  "2.27872385417085+0.438842116902254i",
		"=IMSQRT(6)":         "2.44948974278318",
		"=IMSQRT(\"-2-4i\")": "1.11178594050284-1.79890743994787i",
		// IMSUB
		"=IMSUB(\"5+i\",\"1+4i\")":          "4-3i",
		"=IMSUB(\"9+2i\",6)":                "3+2i",
		"=IMSUB(COMPLEX(5,2),COMPLEX(0,1))": "5+i",
		// IMSUM
		"=IMSUM(\"1-i\",\"5+10i\",2)":       "8+9i",
		"=IMSUM(COMPLEX(5,2),COMPLEX(0,1))": "5+3i",
		// IMTAN
		"=IMTAN(-0)":            "0",
		"=IMTAN(0.5)":           "0.54630248984379",
		"=IMTAN(\"3+0.5i\")":    "-0.111621050771583+0.469469993425885i",
		"=IMTAN(\"2-i\")":       "-0.243458201185725-1.16673625724092i",
		"=IMTAN(COMPLEX(1,-1))": "0.271752585319512-1.08392332733869i",
		// OCT2BIN
		"=OCT2BIN(\"5\")":          "101",
		"=OCT2BIN(\"0000000001\")": "1",
		"=OCT2BIN(\"2\",10)":       "0000000010",
		"=OCT2BIN(\"7777777770\")": "1111111000",
		"=OCT2BIN(\"16\")":         "1110",
		// OCT2DEC
		"=OCT2DEC(\"10\")":         "8",
		"=OCT2DEC(\"22\")":         "18",
		"=OCT2DEC(\"0000000010\")": "8",
		"=OCT2DEC(\"7777777770\")": "-8",
		"=OCT2DEC(\"355\")":        "237",
		// OCT2HEX
		"=OCT2HEX(\"10\")":         "8",
		"=OCT2HEX(\"0000000007\")": "7",
		"=OCT2HEX(\"10\",10)":      "0000000008",
		"=OCT2HEX(\"7777777770\")": "FFFFFFFFF8",
		"=OCT2HEX(\"763\")":        "1F3",
		// Math and Trigonometric Functions
		// ABS
		"=ABS(-1)":      "1",
		"=ABS(-6.5)":    "6.5",
		"=ABS(6.5)":     "6.5",
		"=ABS(0)":       "0",
		"=ABS(2-4.5)":   "2.5",
		"=ABS(ABS(-1))": "1",
		// ACOS
		"=ACOS(-1)":     "3.14159265358979",
		"=ACOS(0)":      "1.5707963267949",
		"=ACOS(ABS(0))": "1.5707963267949",
		// ACOSH
		"=ACOSH(1)":        "0",
		"=ACOSH(2.5)":      "1.56679923697241",
		"=ACOSH(5)":        "2.29243166956118",
		"=ACOSH(ACOSH(5))": "1.47138332153668",
		// ACOT
		"=_xlfn.ACOT(1)":             "0.785398163397448",
		"=_xlfn.ACOT(-2)":            "2.67794504458899",
		"=_xlfn.ACOT(0)":             "1.5707963267949",
		"=_xlfn.ACOT(_xlfn.ACOT(0))": "0.566911504941009",
		// ACOTH
		"=_xlfn.ACOTH(-5)":      "-0.202732554054082",
		"=_xlfn.ACOTH(1.1)":     "1.52226121886171",
		"=_xlfn.ACOTH(2)":       "0.549306144334055",
		"=_xlfn.ACOTH(ABS(-2))": "0.549306144334055",
		// ARABIC
		"=_xlfn.ARABIC(\"IV\")":       "4",
		"=_xlfn.ARABIC(\"-IV\")":      "-4",
		"=_xlfn.ARABIC(\"MCXX\")":     "1120",
		"=_xlfn.ARABIC(\"\")":         "0",
		"=_xlfn.ARABIC(\" ll  lc \")": "-50",
		// ASIN
		"=ASIN(-1)":      "-1.5707963267949",
		"=ASIN(0)":       "0",
		"=ASIN(ASIN(0))": "0",
		// ASINH
		"=ASINH(0)":        "0",
		"=ASINH(-0.5)":     "-0.481211825059603",
		"=ASINH(2)":        "1.44363547517881",
		"=ASINH(ASINH(0))": "0",
		// ATAN
		"=ATAN(-1)":      "-0.785398163397448",
		"=ATAN(0)":       "0",
		"=ATAN(1)":       "0.785398163397448",
		"=ATAN(ATAN(0))": "0",
		// ATANH
		"=ATANH(-0.8)":     "-1.09861228866811",
		"=ATANH(0)":        "0",
		"=ATANH(0.5)":      "0.549306144334055",
		"=ATANH(ATANH(0))": "0",
		// ATAN2
		"=ATAN2(1,1)":          "0.785398163397448",
		"=ATAN2(1,-1)":         "-0.785398163397448",
		"=ATAN2(4,0)":          "0",
		"=ATAN2(4,ATAN2(4,0))": "0",
		// BASE
		"=BASE(12,2)":          "1100",
		"=BASE(12,2,8)":        "00001100",
		"=BASE(100000,16)":     "186A0",
		"=BASE(BASE(12,2),16)": "44C",
		// CEILING
		"=CEILING(22.25,0.1)":              "22.3",
		"=CEILING(22.25,0.5)":              "22.5",
		"=CEILING(22.25,1)":                "23",
		"=CEILING(22.25,10)":               "30",
		"=CEILING(22.25,20)":               "40",
		"=CEILING(-22.25,-0.1)":            "-22.3",
		"=CEILING(-22.25,-1)":              "-23",
		"=CEILING(-22.25,-5)":              "-25",
		"=CEILING(22.25)":                  "23",
		"=CEILING(CEILING(22.25,0.1),0.1)": "22.3",
		// _xlfn.CEILING.MATH
		"=_xlfn.CEILING.MATH(15.25,1)":                       "16",
		"=_xlfn.CEILING.MATH(15.25,0.1)":                     "15.3",
		"=_xlfn.CEILING.MATH(15.25,5)":                       "20",
		"=_xlfn.CEILING.MATH(-15.25,1)":                      "-15",
		"=_xlfn.CEILING.MATH(-15.25,1,1)":                    "-15", // should be 16
		"=_xlfn.CEILING.MATH(-15.25,10)":                     "-10",
		"=_xlfn.CEILING.MATH(-15.25)":                        "-15",
		"=_xlfn.CEILING.MATH(-15.25,-5,-1)":                  "-10",
		"=_xlfn.CEILING.MATH(_xlfn.CEILING.MATH(15.25,1),1)": "16",
		// _xlfn.CEILING.PRECISE
		"=_xlfn.CEILING.PRECISE(22.25,0.1)":                          "22.3",
		"=_xlfn.CEILING.PRECISE(22.25,0.5)":                          "22.5",
		"=_xlfn.CEILING.PRECISE(22.25,1)":                            "23",
		"=_xlfn.CEILING.PRECISE(22.25)":                              "23",
		"=_xlfn.CEILING.PRECISE(22.25,10)":                           "30",
		"=_xlfn.CEILING.PRECISE(22.25,0)":                            "0",
		"=_xlfn.CEILING.PRECISE(-22.25,1)":                           "-22",
		"=_xlfn.CEILING.PRECISE(-22.25,-1)":                          "-22",
		"=_xlfn.CEILING.PRECISE(-22.25,5)":                           "-20",
		"=_xlfn.CEILING.PRECISE(_xlfn.CEILING.PRECISE(22.25,0.1),5)": "25",
		// COMBIN
		"=COMBIN(6,1)":           "6",
		"=COMBIN(6,2)":           "15",
		"=COMBIN(6,3)":           "20",
		"=COMBIN(6,4)":           "15",
		"=COMBIN(6,5)":           "6",
		"=COMBIN(6,6)":           "1",
		"=COMBIN(0,0)":           "1",
		"=COMBIN(6,COMBIN(0,0))": "6",
		// _xlfn.COMBINA
		"=_xlfn.COMBINA(6,1)":                  "6",
		"=_xlfn.COMBINA(6,2)":                  "21",
		"=_xlfn.COMBINA(6,3)":                  "56",
		"=_xlfn.COMBINA(6,4)":                  "126",
		"=_xlfn.COMBINA(6,5)":                  "252",
		"=_xlfn.COMBINA(6,6)":                  "462",
		"=_xlfn.COMBINA(0,0)":                  "0",
		"=_xlfn.COMBINA(0,_xlfn.COMBINA(0,0))": "0",
		// COS
		"=COS(0.785398163)": "0.707106781467586",
		"=COS(0)":           "1",
		"=COS(COS(0))":      "0.54030230586814",
		// COSH
		"=COSH(0)":       "1",
		"=COSH(0.5)":     "1.12762596520638",
		"=COSH(-2)":      "3.76219569108363",
		"=COSH(COSH(0))": "1.54308063481524",
		// _xlfn.COT
		"=_xlfn.COT(0.785398163397448)": "1",
		"=_xlfn.COT(_xlfn.COT(0.45))":   "-0.545473116787229",
		// _xlfn.COTH
		"=_xlfn.COTH(-3.14159265358979)": "-1.00374187319732",
		"=_xlfn.COTH(_xlfn.COTH(1))":     "1.15601401811395",
		// _xlfn.CSC
		"=_xlfn.CSC(-6)":              "3.57889954725441",
		"=_xlfn.CSC(1.5707963267949)": "1",
		"=_xlfn.CSC(_xlfn.CSC(1))":    "1.07785184031088",
		// _xlfn.CSCH
		"=_xlfn.CSCH(-3.14159265358979)": "-0.0865895375300472",
		"=_xlfn.CSCH(_xlfn.CSCH(1))":     "1.04451010395518",
		// _xlfn.DECIMAL
		`=_xlfn.DECIMAL("1100",2)`:    "12",
		`=_xlfn.DECIMAL("186A0",16)`:  "100000",
		`=_xlfn.DECIMAL("31L0",32)`:   "100000",
		`=_xlfn.DECIMAL("70122",8)`:   "28754",
		`=_xlfn.DECIMAL("0x70122",8)`: "28754",
		// DEGREES
		"=DEGREES(1)":          "57.2957795130823",
		"=DEGREES(2.5)":        "143.239448782706",
		"=DEGREES(DEGREES(1))": "3282.80635001174",
		// EVEN
		"=EVEN(23)":   "24",
		"=EVEN(2.22)": "4",
		"=EVEN(0)":    "0",
		"=EVEN(-0.3)": "-2",
		"=EVEN(-11)":  "-12",
		"=EVEN(-4)":   "-4",
		"=EVEN((0))":  "0",
		// EXP
		"=EXP(100)":    "2.68811714181614E+43",
		"=EXP(0.1)":    "1.10517091807565",
		"=EXP(0)":      "1",
		"=EXP(-5)":     "0.00673794699908547",
		"=EXP(EXP(0))": "2.71828182845905",
		// FACT
		"=FACT(3)":       "6",
		"=FACT(6)":       "720",
		"=FACT(10)":      "3628800",
		"=FACT(FACT(3))": "720",
		// FACTDOUBLE
		"=FACTDOUBLE(5)":             "15",
		"=FACTDOUBLE(8)":             "384",
		"=FACTDOUBLE(13)":            "135135",
		"=FACTDOUBLE(FACTDOUBLE(1))": "1",
		// FLOOR
		"=FLOOR(26.75,0.1)":        "26.7",
		"=FLOOR(26.75,0.5)":        "26.5",
		"=FLOOR(26.75,1)":          "26",
		"=FLOOR(26.75,10)":         "20",
		"=FLOOR(26.75,20)":         "20",
		"=FLOOR(-26.75,-0.1)":      "-26.7",
		"=FLOOR(-26.75,-1)":        "-26",
		"=FLOOR(-26.75,-5)":        "-25",
		"=FLOOR(FLOOR(26.75,1),1)": "26",
		// _xlfn.FLOOR.MATH
		"=_xlfn.FLOOR.MATH(58.55)":                  "58",
		"=_xlfn.FLOOR.MATH(58.55,0.1)":              "58.5",
		"=_xlfn.FLOOR.MATH(58.55,5)":                "55",
		"=_xlfn.FLOOR.MATH(58.55,1,1)":              "58",
		"=_xlfn.FLOOR.MATH(-58.55,1)":               "-59",
		"=_xlfn.FLOOR.MATH(-58.55,1,-1)":            "-58",
		"=_xlfn.FLOOR.MATH(-58.55,1,1)":             "-59", // should be -58
		"=_xlfn.FLOOR.MATH(-58.55,10)":              "-60",
		"=_xlfn.FLOOR.MATH(_xlfn.FLOOR.MATH(1),10)": "0",
		// _xlfn.FLOOR.PRECISE
		"=_xlfn.FLOOR.PRECISE(26.75,0.1)":                     "26.7",
		"=_xlfn.FLOOR.PRECISE(26.75,0.5)":                     "26.5",
		"=_xlfn.FLOOR.PRECISE(26.75,1)":                       "26",
		"=_xlfn.FLOOR.PRECISE(26.75)":                         "26",
		"=_xlfn.FLOOR.PRECISE(26.75,10)":                      "20",
		"=_xlfn.FLOOR.PRECISE(26.75,0)":                       "0",
		"=_xlfn.FLOOR.PRECISE(-26.75,1)":                      "-27",
		"=_xlfn.FLOOR.PRECISE(-26.75,-1)":                     "-27",
		"=_xlfn.FLOOR.PRECISE(-26.75,-5)":                     "-30",
		"=_xlfn.FLOOR.PRECISE(_xlfn.FLOOR.PRECISE(26.75),-5)": "25",
		// GCD
		"=GCD(0)":        "0",
		"=GCD(1,0)":      "1",
		"=GCD(1,5)":      "1",
		"=GCD(15,10,25)": "5",
		"=GCD(0,8,12)":   "4",
		"=GCD(7,2)":      "1",
		"=GCD(1,GCD(1))": "1",
		// INT
		"=INT(100.9)":  "100",
		"=INT(5.22)":   "5",
		"=INT(5.99)":   "5",
		"=INT(-6.1)":   "-7",
		"=INT(-100.9)": "-101",
		"=INT(INT(0))": "0",
		// ISO.CEILING
		"=ISO.CEILING(22.25)":              "23",
		"=ISO.CEILING(22.25,1)":            "23",
		"=ISO.CEILING(22.25,0.1)":          "22.3",
		"=ISO.CEILING(22.25,10)":           "30",
		"=ISO.CEILING(-22.25,1)":           "-22",
		"=ISO.CEILING(-22.25,0.1)":         "-22.2",
		"=ISO.CEILING(-22.25,5)":           "-20",
		"=ISO.CEILING(-22.25,0)":           "0",
		"=ISO.CEILING(1,ISO.CEILING(1,0))": "0",
		// LCM
		"=LCM(1,5)":        "5",
		"=LCM(15,10,25)":   "150",
		"=LCM(1,8,12)":     "24",
		"=LCM(7,2)":        "14",
		"=LCM(7)":          "7",
		`=LCM("",1)`:       "1",
		`=LCM(0,0)`:        "0",
		`=LCM(0,LCM(0,0))`: "0",
		// LN
		"=LN(1)":       "0",
		"=LN(100)":     "4.60517018598809",
		"=LN(0.5)":     "-0.693147180559945",
		"=LN(LN(100))": "1.5271796258079",
		// LOG
		"=LOG(64,2)":     "6",
		"=LOG(100)":      "2",
		"=LOG(4,0.5)":    "-2",
		"=LOG(500)":      "2.69897000433602",
		"=LOG(LOG(100))": "0.301029995663981",
		// LOG10
		"=LOG10(100)":        "2",
		"=LOG10(1000)":       "3",
		"=LOG10(0.001)":      "-3",
		"=LOG10(25)":         "1.39794000867204",
		"=LOG10(LOG10(100))": "0.301029995663981",
		// IMLOG2
		"=IMLOG2(\"5+2i\")": "2.42899049756379+0.548954663286635i",
		"=IMLOG2(\"2-i\")":  "1.16096404744368-0.668902106225488i",
		"=IMLOG2(6)":        "2.58496250072116",
		"=IMLOG2(\"3i\")":   "1.58496250072116+2.2661800709136i",
		"=IMLOG2(\"4+i\")":  "2.04373142062517+0.353429502416735i",
		// IMPOWER
		"=IMPOWER(\"2-i\",2)":   "3-4i",
		"=IMPOWER(\"2-i\",3)":   "2-11i",
		"=IMPOWER(9,0.5)":       "3",
		"=IMPOWER(\"2+4i\",-2)": "-0.03-0.04i",
		// IMPRODUCT
		"=IMPRODUCT(3,6)":                       "18",
		`=IMPRODUCT("",3,SUM(6))`:               "18",
		"=IMPRODUCT(\"1-i\",\"5+10i\",2)":       "30+10i",
		"=IMPRODUCT(COMPLEX(5,2),COMPLEX(0,1))": "-2+5i",
		"=IMPRODUCT(A1:C1)":                     "4",
		// MOD
		"=MOD(6,4)":        "2",
		"=MOD(6,3)":        "0",
		"=MOD(6,2.5)":      "1",
		"=MOD(6,1.333)":    "0.668",
		"=MOD(-10.23,1)":   "0.77",
		"=MOD(MOD(1,1),1)": "0",
		// MROUND
		"=MROUND(333.7,0.5)":     "333.5",
		"=MROUND(333.8,1)":       "334",
		"=MROUND(333.3,2)":       "334",
		"=MROUND(555.3,400)":     "400",
		"=MROUND(555,1000)":      "1000",
		"=MROUND(-555.7,-1)":     "-556",
		"=MROUND(-555.4,-1)":     "-555",
		"=MROUND(-1555,-1000)":   "-2000",
		"=MROUND(MROUND(1,1),1)": "1",
		// MULTINOMIAL
		"=MULTINOMIAL(3,1,2,5)":        "27720",
		`=MULTINOMIAL("",3,1,2,5)`:     "27720",
		"=MULTINOMIAL(MULTINOMIAL(1))": "1",
		// _xlfn.MUNIT
		"=_xlfn.MUNIT(4)": "",
		// ODD
		"=ODD(22)":     "23",
		"=ODD(1.22)":   "3",
		"=ODD(1.22+4)": "7",
		"=ODD(0)":      "1",
		"=ODD(-1.3)":   "-3",
		"=ODD(-10)":    "-11",
		"=ODD(-3)":     "-3",
		"=ODD(ODD(1))": "1",
		// PI
		"=PI()": "3.14159265358979",
		// POWER
		"=POWER(4,2)":          "16",
		"=POWER(4,POWER(1,1))": "4",
		// PRODUCT
		"=PRODUCT(3,6)":            "18",
		`=PRODUCT("",3,6)`:         "18",
		`=PRODUCT(PRODUCT(1),3,6)`: "18",
		// QUOTIENT
		"=QUOTIENT(5,2)":             "2",
		"=QUOTIENT(4.5,3.1)":         "1",
		"=QUOTIENT(-10,3)":           "-3",
		"=QUOTIENT(QUOTIENT(1,2),3)": "0",
		// RADIANS
		"=RADIANS(50)":           "0.872664625997165",
		"=RADIANS(-180)":         "-3.14159265358979",
		"=RADIANS(180)":          "3.14159265358979",
		"=RADIANS(360)":          "6.28318530717959",
		"=RADIANS(RADIANS(360))": "0.109662271123215",
		// ROMAN
		"=ROMAN(499,0)":       "CDXCIX",
		"=ROMAN(1999,0)":      "MCMXCIX",
		"=ROMAN(1999,1)":      "MLMVLIV",
		"=ROMAN(1999,2)":      "MXMIX",
		"=ROMAN(1999,3)":      "MVMIV",
		"=ROMAN(1999,4)":      "MIM",
		"=ROMAN(1999,-1)":     "MCMXCIX",
		"=ROMAN(1999,5)":      "MIM",
		"=ROMAN(1999,ODD(1))": "MLMVLIV",
		// ROUND
		"=ROUND(100.319,1)":       "100.3",
		"=ROUND(5.28,1)":          "5.3",
		"=ROUND(5.9999,3)":        "6",
		"=ROUND(99.5,0)":          "100",
		"=ROUND(-6.3,0)":          "-6",
		"=ROUND(-100.5,0)":        "-101",
		"=ROUND(-22.45,1)":        "-22.5",
		"=ROUND(999,-1)":          "1000",
		"=ROUND(991,-1)":          "990",
		"=ROUND(ROUND(100,1),-1)": "100",
		// ROUNDDOWN
		"=ROUNDDOWN(99.999,1)":            "99.9",
		"=ROUNDDOWN(99.999,2)":            "99.99",
		"=ROUNDDOWN(99.999,0)":            "99",
		"=ROUNDDOWN(99.999,-1)":           "90",
		"=ROUNDDOWN(-99.999,2)":           "-99.99",
		"=ROUNDDOWN(-99.999,-1)":          "-90",
		"=ROUNDDOWN(ROUNDDOWN(100,1),-1)": "100",
		// ROUNDUP`
		"=ROUNDUP(11.111,1)":          "11.2",
		"=ROUNDUP(11.111,2)":          "11.12",
		"=ROUNDUP(11.111,0)":          "12",
		"=ROUNDUP(11.111,-1)":         "20",
		"=ROUNDUP(-11.111,2)":         "-11.12",
		"=ROUNDUP(-11.111,-1)":        "-20",
		"=ROUNDUP(ROUNDUP(100,1),-1)": "100",
		// SEC
		"=_xlfn.SEC(-3.14159265358979)": "-1",
		"=_xlfn.SEC(0)":                 "1",
		"=_xlfn.SEC(_xlfn.SEC(0))":      "0.54030230586814",
		// SECH
		"=_xlfn.SECH(-3.14159265358979)": "0.0862667383340547",
		"=_xlfn.SECH(0)":                 "1",
		"=_xlfn.SECH(_xlfn.SECH(0))":     "0.648054273663885",
		// SERIESSUM
		"=SERIESSUM(1,2,3,A1:A4)": "6",
		"=SERIESSUM(1,2,3,A1:B5)": "15",
		// SIGN
		"=SIGN(9.5)":        "1",
		"=SIGN(-9.5)":       "-1",
		"=SIGN(0)":          "0",
		"=SIGN(0.00000001)": "1",
		"=SIGN(6-7)":        "-1",
		"=SIGN(SIGN(-1))":   "-1",
		// SIN
		"=SIN(0.785398163)": "0.707106780905509",
		"=SIN(SIN(1))":      "0.745624141665558",
		// SINH
		"=SINH(0)":       "0",
		"=SINH(0.5)":     "0.521095305493747",
		"=SINH(-2)":      "-3.62686040784702",
		"=SINH(SINH(0))": "0",
		// SQRT
		"=SQRT(4)":        "2",
		"=SQRT(SQRT(16))": "2",
		// SQRTPI
		"=SQRTPI(5)":         "3.96332729760601",
		"=SQRTPI(0.2)":       "0.792665459521202",
		"=SQRTPI(100)":       "17.7245385090552",
		"=SQRTPI(0)":         "0",
		"=SQRTPI(SQRTPI(0))": "0",
		// STDEV
		"=STDEV(F2:F9)":         "10724.9782875238",
		"=STDEV(MUNIT(2))":      "0.577350269189626",
		"=STDEV(0,INT(0))":      "0",
		"=STDEV(INT(1),INT(1))": "0",
		// STDEV.S
		"=STDEV.S(F2:F9)": "10724.9782875238",
		// STDEVA
		"=STDEVA(F2:F9)":    "10724.9782875238",
		"=STDEVA(MUNIT(2))": "0.577350269189626",
		"=STDEVA(0,INT(0))": "0",
		// POISSON.DIST
		"=POISSON.DIST(20,25,FALSE)": "0.0519174686084913",
		"=POISSON.DIST(35,40,TRUE)":  "0.242414197690103",
		// POISSON
		"=POISSON(20,25,FALSE)": "0.0519174686084913",
		"=POISSON(35,40,TRUE)":  "0.242414197690103",
		// SUM
		"=SUM(1,2)":                           "3",
		`=SUM("",1,2)`:                        "3",
		"=SUM(1,2+3)":                         "6",
		"=SUM(SUM(1,2),2)":                    "5",
		"=(-2-SUM(-4+7))*5":                   "-25",
		"SUM(1,2,3,4,5,6,7)":                  "28",
		"=SUM(1,2)+SUM(1,2)":                  "6",
		"=1+SUM(SUM(1,2*3),4)":                "12",
		"=1+SUM(SUM(1,-2*3),4)":               "0",
		"=(-2-SUM(-4*(7+7)))*5":               "270",
		"=SUM(SUM(1+2/1)*2-3/2,2)":            "6.5",
		"=((3+5*2)+3)/5+(-6)/4*2+3":           "3.2",
		"=1+SUM(SUM(1,2*3),4)*-4/2+5+(4+2)*3": "2",
		"=1+SUM(SUM(1,2*3),4)*4/3+5+(4+2)*3":  "38.6666666666667",
		"=SUM(1+ROW())":                       "2",
		"=SUM((SUM(2))+1)":                    "3",
		// SUMIF
		`=SUMIF(F1:F5, "")`:             "0",
		`=SUMIF(A1:A5, "3")`:            "3",
		`=SUMIF(F1:F5, "=36693")`:       "36693",
		`=SUMIF(F1:F5, "<100")`:         "0",
		`=SUMIF(F1:F5, "<=36693")`:      "93233",
		`=SUMIF(F1:F5, ">100")`:         "146554",
		`=SUMIF(F1:F5, ">=100")`:        "146554",
		`=SUMIF(F1:F5, ">=text")`:       "0",
		`=SUMIF(F1:F5, "*Jan",F2:F5)`:   "0",
		`=SUMIF(D3:D7,"Jan",F2:F5)`:     "112114",
		`=SUMIF(D2:D9,"Feb",F2:F9)`:     "157559",
		`=SUMIF(E2:E9,"North 1",F2:F9)`: "66582",
		`=SUMIF(E2:E9,"North*",F2:F9)`:  "138772",
		"=SUMIF(D1:D3,\"Month\",D1:D3)": "0",
		// SUMPRODUCT
		"=SUMPRODUCT(A1,B1)":             "4",
		"=SUMPRODUCT(A1:A2,B1:B2)":       "14",
		"=SUMPRODUCT(A1:A3,B1:B3)":       "14",
		"=SUMPRODUCT(A1:B3)":             "15",
		"=SUMPRODUCT(A1:A3,B1:B3,B2:B4)": "20",
		// SUMSQ
		"=SUMSQ(A1:A4)":            "14",
		"=SUMSQ(A1,B1,A2,B2,6)":    "82",
		`=SUMSQ("",A1,B1,A2,B2,6)`: "82",
		`=SUMSQ(1,SUMSQ(1))`:       "2",
		"=SUMSQ(MUNIT(3))":         "0",
		// SUMX2MY2
		"=SUMX2MY2(A1:A4,B1:B4)": "-36",
		// SUMX2PY2
		"=SUMX2PY2(A1:A4,B1:B4)": "46",
		// SUMXMY2
		"=SUMXMY2(A1:A4,B1:B4)": "18",
		// TAN
		"=TAN(1.047197551)": "1.73205080678249",
		"=TAN(0)":           "0",
		"=TAN(TAN(0))":      "0",
		// TANH
		"=TANH(0)":       "0",
		"=TANH(0.5)":     "0.46211715726001",
		"=TANH(-2)":      "-0.964027580075817",
		"=TANH(TANH(0))": "0",
		// TRUNC
		"=TRUNC(99.999,1)":    "99.9",
		"=TRUNC(99.999,2)":    "99.99",
		"=TRUNC(99.999)":      "99",
		"=TRUNC(99.999,-1)":   "90",
		"=TRUNC(-99.999,2)":   "-99.99",
		"=TRUNC(-99.999,-1)":  "-90",
		"=TRUNC(TRUNC(1),-1)": "0",
		// Statistical Functions
		// AVEDEV
		"=AVEDEV(1,2)":          "0.5",
		"=AVERAGE(A1:A4,B1:B4)": "2.5",
		// AVERAGE
		"=AVERAGE(INT(1))": "1",
		"=AVERAGE(A1)":     "1",
		"=AVERAGE(A1:A2)":  "1.5",
		"=AVERAGE(D2:F9)":  "38014.125",
		// AVERAGEA
		"=AVERAGEA(INT(1))": "1",
		"=AVERAGEA(A1)":     "1",
		"=AVERAGEA(A1:A2)":  "1.5",
		"=AVERAGEA(D2:F9)":  "12671.375",
		// BETA.DIST
		"=BETA.DIST(0.4,4,5,TRUE,0,1)":  "0.4059136",
		"=BETA.DIST(0.6,4,5,FALSE,0,1)": "1.548288",
		// BETADIST
		"=BETADIST(0.4,4,5)":         "0.4059136",
		"=BETADIST(0.4,4,5,0,1)":     "0.4059136",
		"=BETADIST(0.4,4,5,0.4,1)":   "0",
		"=BETADIST(1,2,2,1,3)":       "0",
		"=BETADIST(0.4,4,5,0.2,0.4)": "1",
		"=BETADIST(0.4,4,1)":         "0.0256",
		"=BETADIST(0.4,1,5)":         "0.92224",
		"=BETADIST(3,4,6,2,4)":       "0.74609375",
		"=BETADIST(0.4,2,100)":       "1",
		"=BETADIST(0.75,3,4)":        "0.96240234375",
		"=BETADIST(0.2,0.7,4)":       "0.71794309318323",
		"=BETADIST(0.01,3,4)":        "1.955359E-05",
		"=BETADIST(0.75,130,140)":    "1",
		// BETAINV
		"=BETAINV(0.2,4,5,0,1)": "0.303225844664082",
		// BETA.INV
		"=BETA.INV(0.2,4,5,0,1)": "0.303225844664082",
		// BINOMDIST
		"=BINOMDIST(10,100,0.5,FALSE)": "1.36554263874631E-17",
		"=BINOMDIST(50,100,0.5,FALSE)": "0.0795892373871787",
		"=BINOMDIST(65,100,0.5,FALSE)": "0.000863855665741652",
		"=BINOMDIST(10,100,0.5,TRUE)":  "1.53164508771899E-17",
		"=BINOMDIST(50,100,0.5,TRUE)":  "0.539794618693589",
		"=BINOMDIST(65,100,0.5,TRUE)":  "0.999105034804256",
		// BINOM.DIST
		"=BINOM.DIST(10,100,0.5,FALSE)": "1.36554263874631E-17",
		"=BINOM.DIST(50,100,0.5,FALSE)": "0.0795892373871787",
		"=BINOM.DIST(65,100,0.5,FALSE)": "0.000863855665741652",
		"=BINOM.DIST(10,100,0.5,TRUE)":  "1.53164508771899E-17",
		"=BINOM.DIST(50,100,0.5,TRUE)":  "0.539794618693589",
		"=BINOM.DIST(65,100,0.5,TRUE)":  "0.999105034804256",
		// BINOM.DIST.RANGE
		"=BINOM.DIST.RANGE(100,0.5,0,40)":   "0.0284439668204904",
		"=BINOM.DIST.RANGE(100,0.5,45,55)":  "0.728746975926165",
		"=BINOM.DIST.RANGE(100,0.5,50,100)": "0.539794618693589",
		"=BINOM.DIST.RANGE(100,0.5,50)":     "0.0795892373871787",
		// BINOM.INV
		"=BINOM.INV(0,0.5,0.75)":   "0",
		"=BINOM.INV(0.1,0.1,0.75)": "0",
		"=BINOM.INV(0.6,0.4,0.75)": "0",
		"=BINOM.INV(2,0.4,0.75)":   "1",
		"=BINOM.INV(100,0.5,20%)":  "46",
		"=BINOM.INV(100,0.5,50%)":  "50",
		"=BINOM.INV(100,0.5,90%)":  "56",
		// CHIDIST
		"=CHIDIST(0.5,3)": "0.918891411654676",
		"=CHIDIST(8,3)":   "0.0460117056892314",
		"=CHIDIST(40,4)":  "4.32842260712097E-08",
		"=CHIDIST(42,4)":  "1.66816329414062E-08",
		// CHIINV
		"=CHIINV(0.5,1)":  "0.454936423119572",
		"=CHIINV(0.75,1)": "0.101531044267622",
		"=CHIINV(0.1,2)":  "4.60517018598809",
		"=CHIINV(0.8,2)":  "0.446287102628419",
		// CHISQ.DIST
		"=CHISQ.DIST(0,2,TRUE)":        "0",
		"=CHISQ.DIST(4,1,TRUE)":        "0.954499736103642",
		"=CHISQ.DIST(1180,1180,FALSE)": "0.00821093706387967",
		"=CHISQ.DIST(2,1,FALSE)":       "0.103776874355149",
		"=CHISQ.DIST(3,2,FALSE)":       "0.111565080074215",
		"=CHISQ.DIST(2,3,FALSE)":       "0.207553748710297",
		"=CHISQ.DIST(1425,1,FALSE)":    "3.88315098887099E-312",
		"=CHISQ.DIST(3,2,TRUE)":        "0.77686983985157",
		// CHISQ.DIST.RT
		"=CHISQ.DIST.RT(0.5,3)": "0.918891411654676",
		"=CHISQ.DIST.RT(8,3)":   "0.0460117056892314",
		"=CHISQ.DIST.RT(40,4)":  "4.32842260712097E-08",
		"=CHISQ.DIST.RT(42,4)":  "1.66816329414062E-08",
		// CHISQ.INV
		"=CHISQ.INV(0,2)":    "0",
		"=CHISQ.INV(0.75,1)": "1.32330369693147",
		"=CHISQ.INV(0.1,2)":  "0.210721031315653",
		"=CHISQ.INV(0.8,2)":  "3.2188758248682",
		"=CHISQ.INV(0.25,3)": "1.21253290304567",
		// CHISQ.INV.RT
		"=CHISQ.INV.RT(0.75,1)": "0.101531044267622",
		"=CHISQ.INV.RT(0.1,2)":  "4.60517018598809",
		"=CHISQ.INV.RT(0.8,2)":  "0.446287102628419",
		// CONFIDENCE
		"=CONFIDENCE(0.05,0.07,100)": "0.0137197479028414",
		// CONFIDENCE.NORM
		"=CONFIDENCE.NORM(0.05,0.07,100)": "0.0137197479028414",
		// CORREL
		"=CORREL(A1:A5,B1:B5)": "1",
		// COUNT
		"=COUNT()":                        "0",
		"=COUNT(E1:F2,\"text\",1,INT(2))": "3",
		// COUNTA
		"=COUNTA()":                              "0",
		"=COUNTA(A1:A5,B2:B5,\"text\",1,INT(2))": "8",
		"=COUNTA(COUNTA(1),MUNIT(1))":            "2",
		// COUNTBLANK
		"=COUNTBLANK(MUNIT(1))": "0",
		"=COUNTBLANK(1)":        "0",
		"=COUNTBLANK(B1:C1)":    "1",
		"=COUNTBLANK(C1)":       "1",
		// COUNTIF
		"=COUNTIF(D1:D9,\"Jan\")":     "4",
		"=COUNTIF(D1:D9,\"<>Jan\")":   "5",
		"=COUNTIF(A1:F9,\">=50000\")": "2",
		"=COUNTIF(A1:F9,TRUE)":        "0",
		// COUNTIFS
		"=COUNTIFS(A1:A9,2,D1:D9,\"Jan\")":          "1",
		"=COUNTIFS(F1:F9,\">20000\",D1:D9,\"Jan\")": "4",
		"=COUNTIFS(F1:F9,\">60000\",D1:D9,\"Jan\")": "0",
		// CRITBINOM
		"=CRITBINOM(0,0.5,0.75)":   "0",
		"=CRITBINOM(0.1,0.1,0.75)": "0",
		"=CRITBINOM(0.6,0.4,0.75)": "0",
		"=CRITBINOM(2,0.4,0.75)":   "1",
		"=CRITBINOM(100,0.5,20%)":  "46",
		"=CRITBINOM(100,0.5,50%)":  "50",
		"=CRITBINOM(100,0.5,90%)":  "56",
		// DEVSQ
		"=DEVSQ(1,3,5,2,9,7)": "47.5",
		"=DEVSQ(A1:D2)":       "10",
		// FISHER
		"=FISHER(-0.9)":   "-1.47221948958322",
		"=FISHER(-0.25)":  "-0.255412811882995",
		"=FISHER(0.8)":    "1.09861228866811",
		"=FISHER(INT(0))": "0",
		// FISHERINV
		"=FISHERINV(-0.2)":   "-0.197375320224904",
		"=FISHERINV(INT(0))": "0",
		"=FISHERINV(2.8)":    "0.992631520201128",
		// GAMMA
		"=GAMMA(0.1)":    "9.51350769866873",
		"=GAMMA(INT(1))": "1",
		"=GAMMA(1.5)":    "0.886226925452758",
		"=GAMMA(5.5)":    "52.3427777845535",
		// GAMMA.DIST
		"=GAMMA.DIST(6,3,2,FALSE)": "0.112020903827694",
		"=GAMMA.DIST(6,3,2,TRUE)":  "0.576809918873156",
		// GAMMADIST
		"=GAMMADIST(6,3,2,FALSE)": "0.112020903827694",
		"=GAMMADIST(6,3,2,TRUE)":  "0.576809918873156",
		// GAMMA.INV
		"=GAMMA.INV(0.5,3,2)":   "5.34812062744712",
		"=GAMMA.INV(0.5,0.5,1)": "0.227468211559786",
		// GAMMAINV
		"=GAMMAINV(0.5,3,2)":   "5.34812062744712",
		"=GAMMAINV(0.5,0.5,1)": "0.227468211559786",
		// GAMMALN
		"=GAMMALN(4.5)":    "2.45373657084244",
		"=GAMMALN(INT(1))": "0",
		// GAMMALN.PRECISE
		"=GAMMALN.PRECISE(0.4)": "0.796677817701784",
		"=GAMMALN.PRECISE(4.5)": "2.45373657084244",
		// GAUSS
		"=GAUSS(-5)":  "-0.499999713348428",
		"=GAUSS(0)":   "0",
		"=GAUSS(0.1)": "0.039827837277029",
		"=GAUSS(2.5)": "0.493790334674224",
		// GEOMEAN
		"=GEOMEAN(2.5,3,0.5,1,3)": "1.6226711115996",
		// HARMEAN
		"=HARMEAN(2.5,3,0.5,1,3)":               "1.22950819672131",
		"=HARMEAN(\"2.5\",3,0.5,1,INT(3),\"\")": "1.22950819672131",
		// HYPGEOM.DIST
		"=HYPGEOM.DIST(0,3,3,9,TRUE)":   "0.238095238095238",
		"=HYPGEOM.DIST(1,3,3,9,TRUE)":   "0.773809523809524",
		"=HYPGEOM.DIST(2,3,3,9,TRUE)":   "0.988095238095238",
		"=HYPGEOM.DIST(3,3,3,9,TRUE)":   "1",
		"=HYPGEOM.DIST(1,4,4,12,FALSE)": "0.452525252525253",
		"=HYPGEOM.DIST(2,4,4,12,FALSE)": "0.339393939393939",
		"=HYPGEOM.DIST(3,4,4,12,FALSE)": "0.0646464646464646",
		"=HYPGEOM.DIST(4,4,4,12,FALSE)": "0.00202020202020202",
		// HYPGEOMDIST
		"=HYPGEOMDIST(1,4,4,12)": "0.452525252525253",
		"=HYPGEOMDIST(2,4,4,12)": "0.339393939393939",
		"=HYPGEOMDIST(3,4,4,12)": "0.0646464646464646",
		"=HYPGEOMDIST(4,4,4,12)": "0.00202020202020202",
		// KURT
		"=KURT(F1:F9)":           "-1.03350350255137",
		"=KURT(F1,F2:F9)":        "-1.03350350255137",
		"=KURT(INT(1),MUNIT(2))": "-3.33333333333334",
		// EXPON.DIST
		"=EXPON.DIST(0.5,1,TRUE)":  "0.393469340287367",
		"=EXPON.DIST(0.5,1,FALSE)": "0.606530659712633",
		"=EXPON.DIST(2,1,TRUE)":    "0.864664716763387",
		// EXPONDIST
		"=EXPONDIST(0.5,1,TRUE)":  "0.393469340287367",
		"=EXPONDIST(0.5,1,FALSE)": "0.606530659712633",
		"=EXPONDIST(2,1,TRUE)":    "0.864664716763387",
		// FDIST
		"=FDIST(5,1,2)": "0.154845745271483",
		// F.DIST
		"=F.DIST(1,2,5,TRUE)":  "0.568798849628308",
		"=F.DIST(1,2,5,FALSE)": "0.308000821694066",
		// F.DIST.RT
		"=F.DIST.RT(5,1,2)": "0.154845745271483",
		// F.INV
		"=F.INV(0.9,2,5)": "3.77971607877395",
		// FINV
		"=FINV(0.2,1,2)":   "3.55555555555555",
		"=FINV(0.6,1,2)":   "0.380952380952381",
		"=FINV(0.6,2,2)":   "0.666666666666667",
		"=FINV(0.6,4,4)":   "0.763454070045235",
		"=FINV(0.5,4,8)":   "0.914645355977072",
		"=FINV(0.1,79,86)": "1.32646097270444",
		"=FINV(1,40,5)":    "0",
		// F.INV.RT
		"=F.INV.RT(0.2,1,2)":   "3.55555555555555",
		"=F.INV.RT(0.6,1,2)":   "0.380952380952381",
		"=F.INV.RT(0.6,2,2)":   "0.666666666666667",
		"=F.INV.RT(0.6,4,4)":   "0.763454070045235",
		"=F.INV.RT(0.5,4,8)":   "0.914645355977072",
		"=F.INV.RT(0.1,79,86)": "1.32646097270444",
		"=F.INV.RT(1,40,5)":    "0",
		// LOGINV
		"=LOGINV(0.3,2,0.2)": "6.6533460753367",
		// LOGINV
		"=LOGNORM.INV(0.3,2,0.2)": "6.6533460753367",
		// LOGNORM.DIST
		"=LOGNORM.DIST(0.5,10,5,FALSE)": "0.0162104821842127",
		"=LOGNORM.DIST(12,10,5,TRUE)":   "0.0664171147992078",
		// LOGNORMDIST
		"=LOGNORMDIST(12,10,5)": "0.0664171147992078",
		// NEGBINOM.DIST
		"=NEGBINOM.DIST(6,12,0.5,FALSE)":  "0.047210693359375",
		"=NEGBINOM.DIST(12,12,0.5,FALSE)": "0.0805901288986206",
		"=NEGBINOM.DIST(15,12,0.5,FALSE)": "0.057564377784729",
		"=NEGBINOM.DIST(12,12,0.5,TRUE)":  "0.580590128898621",
		"=NEGBINOM.DIST(15,12,0.5,TRUE)":  "0.778965830802917",
		// NEGBINOMDIST
		"=NEGBINOMDIST(6,12,0.5)":  "0.047210693359375",
		"=NEGBINOMDIST(12,12,0.5)": "0.0805901288986206",
		"=NEGBINOMDIST(15,12,0.5)": "0.057564377784729",
		// NORM.DIST
		"=NORM.DIST(0.8,1,0.3,TRUE)": "0.252492537546923",
		"=NORM.DIST(50,40,20,FALSE)": "0.017603266338215",
		// NORMDIST
		"=NORMDIST(0.8,1,0.3,TRUE)": "0.252492537546923",
		"=NORMDIST(50,40,20,FALSE)": "0.017603266338215",
		// NORM.INV
		"=NORM.INV(0.6,5,2)": "5.50669420572",
		// NORMINV
		"=NORMINV(0.6,5,2)":     "5.50669420572",
		"=NORMINV(0.99,40,1.5)": "43.489521811582",
		"=NORMINV(0.02,40,1.5)": "36.9193766364954",
		// NORM.S.DIST
		"=NORM.S.DIST(0.8,TRUE)": "0.788144601416603",
		// NORMSDIST
		"=NORMSDIST(1.333333)": "0.908788725604095",
		"=NORMSDIST(0)":        "0.5",
		// NORM.S.INV
		"=NORM.S.INV(0.25)": "-0.674489750223423",
		// NORMSINV
		"=NORMSINV(0.25)": "-0.674489750223423",
		// LARGE
		"=LARGE(A1:A5,1)": "3",
		"=LARGE(A1:B5,2)": "4",
		"=LARGE(A1,1)":    "1",
		"=LARGE(A1:F2,1)": "36693",
		// MAX
		"=MAX(1)":          "1",
		"=MAX(TRUE())":     "1",
		"=MAX(0.5,TRUE())": "1",
		"=MAX(FALSE())":    "0",
		"=MAX(MUNIT(2))":   "1",
		"=MAX(INT(1))":     "1",
		// MAXA
		"=MAXA(1)":          "1",
		"=MAXA(TRUE())":     "1",
		"=MAXA(0.5,TRUE())": "1",
		"=MAXA(FALSE())":    "0",
		"=MAXA(MUNIT(2))":   "1",
		"=MAXA(INT(1))":     "1",
		"=MAXA(A1:B4,MUNIT(1),INT(0),1,E1:F2,\"\")": "36693",
		// MAXIFS
		"=MAXIFS(F2:F4,A2:A4,\">0\")": "36693",
		// MEDIAN
		"=MEDIAN(A1:A5,12)":               "2",
		"=MEDIAN(A1:A5)":                  "1.5",
		"=MEDIAN(A1:A5,MEDIAN(A1:A5,12))": "2",
		// MIN
		"=MIN(1)":           "1",
		"=MIN(TRUE())":      "1",
		"=MIN(0.5,FALSE())": "0",
		"=MIN(FALSE())":     "0",
		"=MIN(MUNIT(2))":    "0",
		"=MIN(INT(1))":      "1",
		// MINA
		"=MINA(1)":           "1",
		"=MINA(TRUE())":      "1",
		"=MINA(0.5,FALSE())": "0",
		"=MINA(FALSE())":     "0",
		"=MINA(MUNIT(2))":    "0",
		"=MINA(INT(1))":      "1",
		"=MINA(A1:B4,MUNIT(1),INT(0),1,E1:F2,\"\")": "0",
		// MINIFS
		"=MINIFS(F2:F4,A2:A4,\">0\")": "22100",
		// PERCENTILE.EXC
		"=PERCENTILE.EXC(A1:A4,0.2)": "0",
		"=PERCENTILE.EXC(A1:A4,0.6)": "2",
		// PERCENTILE.INC
		"=PERCENTILE.INC(A1:A4,0.2)": "0.6",
		// PERCENTILE
		"=PERCENTILE(A1:A4,0.2)": "0.6",
		"=PERCENTILE(0,0)":       "0",
		// PERCENTRANK.EXC
		"=PERCENTRANK.EXC(A1:B4,0)":     "0.142",
		"=PERCENTRANK.EXC(A1:B4,2)":     "0.428",
		"=PERCENTRANK.EXC(A1:B4,2.5)":   "0.5",
		"=PERCENTRANK.EXC(A1:B4,2.6,1)": "0.5",
		"=PERCENTRANK.EXC(A1:B4,5)":     "0.857",
		// PERCENTRANK.INC
		"=PERCENTRANK.INC(A1:B4,0)":     "0",
		"=PERCENTRANK.INC(A1:B4,2)":     "0.4",
		"=PERCENTRANK.INC(A1:B4,2.5)":   "0.5",
		"=PERCENTRANK.INC(A1:B4,2.6,1)": "0.5",
		"=PERCENTRANK.INC(A1:B4,5)":     "1",
		// PERCENTRANK
		"=PERCENTRANK(A1:B4,0)":     "0",
		"=PERCENTRANK(A1:B4,2)":     "0.4",
		"=PERCENTRANK(A1:B4,2.5)":   "0.5",
		"=PERCENTRANK(A1:B4,2.6,1)": "0.5",
		"=PERCENTRANK(A1:B4,5)":     "1",
		// PERMUT
		"=PERMUT(6,6)":  "720",
		"=PERMUT(7,6)":  "5040",
		"=PERMUT(10,6)": "151200",
		// PERMUTATIONA
		"=PERMUTATIONA(6,6)": "46656",
		"=PERMUTATIONA(7,6)": "117649",
		// QUARTILE
		"=QUARTILE(A1:A4,2)": "1.5",
		// QUARTILE.EXC
		"=QUARTILE.EXC(A1:A4,1)": "0.25",
		"=QUARTILE.EXC(A1:A4,2)": "1.5",
		"=QUARTILE.EXC(A1:A4,3)": "2.75",
		// QUARTILE.INC
		"=QUARTILE.INC(A1:A4,0)": "0",
		// RANK
		"=RANK(1,A1:B5)":   "5",
		"=RANK(1,A1:B5,0)": "5",
		"=RANK(1,A1:B5,1)": "2",
		// RANK.EQ
		"=RANK.EQ(1,A1:B5)":   "5",
		"=RANK.EQ(1,A1:B5,0)": "5",
		"=RANK.EQ(1,A1:B5,1)": "2",
		// SKEW
		"=SKEW(1,2,3,4,3)": "-0.404796008910937",
		"=SKEW(A1:B2)":     "0",
		"=SKEW(A1:D3)":     "0",
		// SMALL
		"=SMALL(A1:A5,1)": "0",
		"=SMALL(A1:B5,2)": "1",
		"=SMALL(A1,1)":    "1",
		"=SMALL(A1:F2,1)": "1",
		// STANDARDIZE
		"=STANDARDIZE(5.5,5,2)":   "0.25",
		"=STANDARDIZE(12,15,1.5)": "-2",
		"=STANDARDIZE(-2,0,5)":    "-0.4",
		// STDEVP
		"=STDEVP(A1:B2,6,-1)": "2.40947204913349",
		// STDEV.P
		"=STDEV.P(A1:B2,6,-1)": "2.40947204913349",
		// T.DIST
		"=T.DIST(1,10,TRUE)":   "0.82955343384897",
		"=T.DIST(-1,10,TRUE)":  "0.17044656615103",
		"=T.DIST(-1,10,FALSE)": "0.230361989229139",
		// TDIST
		"=TDIST(1,10,1)": "0.17044656615103",
		"=TDIST(1,10,2)": "0.34089313230206",
		// TRIMMEAN
		"=TRIMMEAN(A1:B4,10%)": "2.5",
		"=TRIMMEAN(A1:B4,70%)": "2.5",
		// VAR
		"=VAR(1,3,5,0,C1)":      "4.91666666666667",
		"=VAR(1,3,5,0,C1,TRUE)": "4",
		// VARA
		"=VARA(1,3,5,0,C1)":      "4.7",
		"=VARA(1,3,5,0,C1,TRUE)": "3.86666666666667",
		// VARP
		"=VARP(A1:A5)":           "1.25",
		"=VARP(1,3,5,0,C1,TRUE)": "3.2",
		// VAR.P
		"=VAR.P(A1:A5)": "1.25",
		// VAR.S
		"=VAR.S(1,3,5,0,C1)":      "4.91666666666667",
		"=VAR.S(1,3,5,0,C1,TRUE)": "4",
		// VARPA
		"=VARPA(1,3,5,0,C1)":      "3.76",
		"=VARPA(1,3,5,0,C1,TRUE)": "3.22222222222222",
		// WEIBULL
		"=WEIBULL(1,3,1,FALSE)":  "1.10363832351433",
		"=WEIBULL(2,5,1.5,TRUE)": "0.985212776817482",
		// WEIBULL.DIST
		"=WEIBULL.DIST(1,3,1,FALSE)":  "1.10363832351433",
		"=WEIBULL.DIST(2,5,1.5,TRUE)": "0.985212776817482",
		// Information Functions
		// ERROR.TYPE
		"=ERROR.TYPE(1/0)":           "2",
		"=ERROR.TYPE(COT(0))":        "2",
		"=ERROR.TYPE(XOR(\"text\"))": "3",
		"=ERROR.TYPE(HEX2BIN(2,1))":  "6",
		"=ERROR.TYPE(NA())":          "7",
		// ISBLANK
		"=ISBLANK(A1)": "FALSE",
		"=ISBLANK(A5)": "TRUE",
		// ISERR
		"=ISERR(A1)":           "FALSE",
		"=ISERR(NA())":         "FALSE",
		"=ISERR(POWER(0,-1)))": "TRUE",
		// ISERROR
		"=ISERROR(A1)":          "FALSE",
		"=ISERROR(NA())":        "TRUE",
		"=ISERROR(\"#VALUE!\")": "FALSE",
		// ISEVEN
		"=ISEVEN(A1)": "FALSE",
		"=ISEVEN(A2)": "TRUE",
		// ISFORMULA
		"=ISFORMULA(A1)":    "FALSE",
		"=ISFORMULA(\"A\")": "FALSE",
		// ISLOGICAL
		"=ISLOGICAL(TRUE)":      "TRUE",
		"=ISLOGICAL(FALSE)":     "TRUE",
		"=ISLOGICAL(A1=A2)":     "TRUE",
		"=ISLOGICAL(\"true\")":  "TRUE",
		"=ISLOGICAL(\"false\")": "TRUE",
		"=ISLOGICAL(A1)":        "FALSE",
		"=ISLOGICAL(20/5)":      "FALSE",
		// ISNA
		"=ISNA(A1)":   "FALSE",
		"=ISNA(NA())": "TRUE",
		// ISNONTEXT
		"=ISNONTEXT(A1)":         "FALSE",
		"=ISNONTEXT(A5)":         "TRUE",
		`=ISNONTEXT("Excelize")`: "FALSE",
		"=ISNONTEXT(NA())":       "TRUE",
		// ISNUMBER
		"=ISNUMBER(A1)": "TRUE",
		"=ISNUMBER(D1)": "FALSE",
		// ISODD
		"=ISODD(A1)": "TRUE",
		"=ISODD(A2)": "FALSE",
		// ISREF
		"=ISREF(B1)":       "TRUE",
		"=ISREF(B1:B2)":    "TRUE",
		"=ISREF(\"text\")": "FALSE",
		"=ISREF(B1*B2)":    "FALSE",
		// ISTEXT
		"=ISTEXT(D1)": "TRUE",
		"=ISTEXT(A1)": "FALSE",
		// N
		"=N(10)":     "10",
		"=N(\"10\")": "10",
		"=N(\"x\")":  "0",
		"=N(TRUE)":   "1",
		"=N(FALSE)":  "0",
		// SHEET
		"=SHEET()":           "1",
		"=SHEET(\"Sheet1\")": "1",
		// SHEETS
		"=SHEETS()":   "1",
		"=SHEETS(A1)": "1",
		// TYPE
		"=TYPE(2)":        "1",
		"=TYPE(10/2)":     "1",
		"=TYPE(C1)":       "1",
		"=TYPE(\"text\")": "2",
		"=TYPE(TRUE)":     "4",
		"=TYPE(NA())":     "16",
		"=TYPE(MUNIT(2))": "64",
		// T
		"=T(\"text\")": "text",
		"=T(N(10))":    "",
		// Logical Functions
		// AND
		"=AND(0)":               "FALSE",
		"=AND(1)":               "TRUE",
		"=AND(1,0)":             "FALSE",
		"=AND(0,1)":             "FALSE",
		"=AND(1=1)":             "TRUE",
		"=AND(1<2)":             "TRUE",
		"=AND(1>2,2<3,2>0,3>1)": "FALSE",
		"=AND(1=1),1=1":         "TRUE",
		// FALSE
		"=FALSE()": "FALSE",
		// IFERROR
		"=IFERROR(1/2,0)":       "0.5",
		"=IFERROR(ISERROR(),0)": "0",
		"=IFERROR(1/0,0)":       "0",
		// IFNA
		"=IFNA(1,\"not found\")":    "1",
		"=IFNA(NA(),\"not found\")": "not found",
		// IFS
		"=IFS(4>1,5/4,4<-1,-5/4,TRUE,0)":     "1.25",
		"=IFS(-2>1,5/-2,-2<-1,-5/-2,TRUE,0)": "2.5",
		"=IFS(0>1,5/0,0<-1,-5/0,TRUE,0)":     "0",
		// NOT
		"=NOT(FALSE())":     "TRUE",
		"=NOT(\"false\")":   "TRUE",
		"=NOT(\"true\")":    "FALSE",
		"=NOT(ISBLANK(B1))": "TRUE",
		// OR
		"=OR(1)":       "TRUE",
		"=OR(0)":       "FALSE",
		"=OR(1=2,2=2)": "TRUE",
		"=OR(1=2,2=3)": "FALSE",
		// SWITCH
		"=SWITCH(1,1,\"A\",2,\"B\",3,\"C\",\"N\")": "A",
		"=SWITCH(3,1,\"A\",2,\"B\",3,\"C\",\"N\")": "C",
		"=SWITCH(4,1,\"A\",2,\"B\",3,\"C\",\"N\")": "N",
		// TRUE
		"=TRUE()": "TRUE",
		// XOR
		"=XOR(1>0,2>0)":                       "FALSE",
		"=XOR(1>0,0>1)":                       "TRUE",
		"=XOR(1>0,0>1,INT(0),INT(1),A1:A4,2)": "FALSE",
		// Date and Time Functions
		// DATE
		"=DATE(2020,10,21)": "2020-10-21 00:00:00 +0000 UTC",
		"=DATE(1900,1,1)":   "1899-12-31 00:00:00 +0000 UTC",
		// DATEDIF
		"=DATEDIF(43101,43101,\"D\")":  "0",
		"=DATEDIF(43101,43891,\"d\")":  "790",
		"=DATEDIF(43101,43891,\"Y\")":  "2",
		"=DATEDIF(42156,44242,\"y\")":  "5",
		"=DATEDIF(43101,43891,\"M\")":  "26",
		"=DATEDIF(42171,44242,\"m\")":  "67",
		"=DATEDIF(42156,44454,\"MD\")": "14",
		"=DATEDIF(42171,44242,\"md\")": "30",
		"=DATEDIF(43101,43891,\"YM\")": "2",
		"=DATEDIF(42171,44242,\"ym\")": "7",
		"=DATEDIF(43101,43891,\"YD\")": "59",
		"=DATEDIF(36526,73110,\"YD\")": "60",
		"=DATEDIF(42171,44242,\"yd\")": "244",
		// DATEVALUE
		"=DATEVALUE(\"01/01/16\")":   "42370",
		"=DATEVALUE(\"01/01/2016\")": "42370",
		"=DATEVALUE(\"01/01/29\")":   "47119",
		"=DATEVALUE(\"01/01/30\")":   "10959",
		// DAY
		"=DAY(0)":                                "0",
		"=DAY(INT(7))":                           "7",
		"=DAY(\"35\")":                           "4",
		"=DAY(42171)":                            "16",
		"=DAY(\"2-28-1900\")":                    "28",
		"=DAY(\"31-May-2015\")":                  "31",
		"=DAY(\"01/03/2019 12:14:16\")":          "3",
		"=DAY(\"January 25, 2020 01 AM\")":       "25",
		"=DAY(\"January 25, 2020 01:03 AM\")":    "25",
		"=DAY(\"January 25, 2020 12:00:00 AM\")": "25",
		"=DAY(\"1900-1-1\")":                     "1",
		"=DAY(\"12-1-1900\")":                    "1",
		"=DAY(\"3-January-1900\")":               "3",
		"=DAY(\"3-February-2000\")":              "3",
		"=DAY(\"3-February-2008\")":              "3",
		"=DAY(\"01/25/20\")":                     "25",
		"=DAY(\"01/25/31\")":                     "25",
		// DAYS
		"=DAYS(2,1)":                           "1",
		"=DAYS(INT(2),INT(1))":                 "1",
		"=DAYS(\"02/02/2015\",\"01/01/2015\")": "32",
		// HOUR
		"=HOUR(1)":                    "0",
		"=HOUR(43543.5032060185)":     "12",
		"=HOUR(\"43543.5032060185\")": "12",
		"=HOUR(\"13:00:55\")":         "13",
		"=HOUR(\"1:00 PM\")":          "13",
		"=HOUR(\"12/09/2015 08:55\")": "8",
		// ISOWEEKNUM
		"=ISOWEEKNUM(42370)":          "53",
		"=ISOWEEKNUM(\"42370\")":      "53",
		"=ISOWEEKNUM(\"01/01/2005\")": "53",
		"=ISOWEEKNUM(\"02/02/2005\")": "5",
		// MINUTE
		"=MINUTE(1)":                    "0",
		"=MINUTE(0.04)":                 "57",
		"=MINUTE(\"0.04\")":             "57",
		"=MINUTE(\"13:35:55\")":         "35",
		"=MINUTE(\"12/09/2015 08:55\")": "55",
		// MONTH
		"=MONTH(42171)":           "6",
		"=MONTH(\"31-May-2015\")": "5",
		// YEAR
		"=YEAR(15)":              "1900",
		"=YEAR(\"15\")":          "1900",
		"=YEAR(2048)":            "1905",
		"=YEAR(42171)":           "2015",
		"=YEAR(\"29-May-2015\")": "2015",
		"=YEAR(\"05/03/1984\")":  "1984",
		// YEARFRAC
		"=YEARFRAC(42005,42005)":                      "0",
		"=YEARFRAC(42005,42094)":                      "0.25",
		"=YEARFRAC(42005,42094,0)":                    "0.25",
		"=YEARFRAC(42005,42094,1)":                    "0.243835616438356",
		"=YEARFRAC(42005,42094,2)":                    "0.247222222222222",
		"=YEARFRAC(42005,42094,3)":                    "0.243835616438356",
		"=YEARFRAC(42005,42094,4)":                    "0.247222222222222",
		"=YEARFRAC(\"01/01/2015\",\"03/31/2015\")":    "0.25",
		"=YEARFRAC(\"01/01/2015\",\"03/31/2015\",0)":  "0.25",
		"=YEARFRAC(\"01/01/2015\",\"03/31/2015\",1)":  "0.243835616438356",
		"=YEARFRAC(\"01/01/2015\",\"03/31/2015\",2)":  "0.247222222222222",
		"=YEARFRAC(\"01/01/2015\",\"03/31/2015\",3)":  "0.243835616438356",
		"=YEARFRAC(\"01/01/2015\",\"03/31/2015\",4)":  "0.247222222222222",
		"=YEARFRAC(\"01/01/2015\",42094)":             "0.25",
		"=YEARFRAC(42005,\"03/31/2015\",0)":           "0.25",
		"=YEARFRAC(\"01/31/2015\",\"03/31/2015\")":    "0.166666666666667",
		"=YEARFRAC(\"01/30/2015\",\"03/31/2015\")":    "0.166666666666667",
		"=YEARFRAC(\"02/29/2000\", \"02/29/2008\")":   "8",
		"=YEARFRAC(\"02/29/2000\", \"02/29/2008\",1)": "7.99817518248175",
		"=YEARFRAC(\"02/29/2000\", \"01/29/2001\",1)": "0.915300546448087",
		"=YEARFRAC(\"02/29/2000\", \"03/29/2000\",1)": "0.0792349726775956",
		"=YEARFRAC(\"01/31/2000\", \"03/29/2000\",4)": "0.163888888888889",
		// SECOND
		"=SECOND(\"13:35:55\")":            "55",
		"=SECOND(\"13:10:60\")":            "0",
		"=SECOND(\"13:10:61\")":            "1",
		"=SECOND(\"08:17:00\")":            "0",
		"=SECOND(\"12/09/2015 08:55\")":    "0",
		"=SECOND(\"12/09/2011 08:17:23\")": "23",
		"=SECOND(\"43543.5032060185\")":    "37",
		"=SECOND(43543.5032060185)":        "37",
		// TIME
		"=TIME(5,44,32)":             "0.239259259259259",
		"=TIME(\"5\",\"44\",\"32\")": "0.239259259259259",
		"=TIME(0,0,73)":              "0.000844907407407407",
		// TIMEVALUE
		"=TIMEVALUE(\"2:23\")":             "0.0993055555555555",
		"=TIMEVALUE(\"2:23 am\")":          "0.0993055555555555",
		"=TIMEVALUE(\"2:23 PM\")":          "0.599305555555556",
		"=TIMEVALUE(\"14:23:00\")":         "0.599305555555556",
		"=TIMEVALUE(\"00:02:23\")":         "0.00165509259259259",
		"=TIMEVALUE(\"01/01/2011 02:23\")": "0.0993055555555555",
		// WEEKDAY
		"=WEEKDAY(0)":                 "7",
		"=WEEKDAY(47119)":             "2",
		"=WEEKDAY(\"12/25/2012\")":    "3",
		"=WEEKDAY(\"12/25/2012\",1)":  "3",
		"=WEEKDAY(\"12/25/2012\",2)":  "2",
		"=WEEKDAY(\"12/25/2012\",3)":  "1",
		"=WEEKDAY(\"12/25/2012\",11)": "2",
		"=WEEKDAY(\"12/25/2012\",12)": "1",
		"=WEEKDAY(\"12/25/2012\",13)": "7",
		"=WEEKDAY(\"12/25/2012\",14)": "6",
		"=WEEKDAY(\"12/25/2012\",15)": "5",
		"=WEEKDAY(\"12/25/2012\",16)": "4",
		"=WEEKDAY(\"12/25/2012\",17)": "3",
		// Text Functions
		// CHAR
		"=CHAR(65)": "A",
		"=CHAR(97)": "a",
		"=CHAR(63)": "?",
		"=CHAR(51)": "3",
		// CLEAN
		"=CLEAN(\"\u0009clean text\")": "clean text",
		"=CLEAN(0)":                    "0",
		// CODE
		"=CODE(\"Alpha\")": "65",
		"=CODE(\"alpha\")": "97",
		"=CODE(\"?\")":     "63",
		"=CODE(\"3\")":     "51",
		"=CODE(\"\")":      "0",
		// CONCAT
		"=CONCAT(TRUE(),1,FALSE(),\"0\",INT(2))": "TRUE1FALSE02",
		// CONCATENATE
		"=CONCATENATE(TRUE(),1,FALSE(),\"0\",INT(2))": "TRUE1FALSE02",
		// EXACT
		"=EXACT(1,\"1\")":     "TRUE",
		"=EXACT(1,1)":         "TRUE",
		"=EXACT(\"A\",\"a\")": "FALSE",
		// FIXED
		"=FIXED(5123.591)":         "5,123.591",
		"=FIXED(5123.591,1)":       "5,123.6",
		"=FIXED(5123.591,0)":       "5,124",
		"=FIXED(5123.591,-1)":      "5,120",
		"=FIXED(5123.591,-2)":      "5,100",
		"=FIXED(5123.591,-3,TRUE)": "5000",
		"=FIXED(5123.591,-5)":      "0",
		"=FIXED(-77262.23973,-5)":  "-100,000",
		// FIND
		"=FIND(\"T\",\"Original Text\")":   "10",
		"=FIND(\"t\",\"Original Text\")":   "13",
		"=FIND(\"i\",\"Original Text\")":   "3",
		"=FIND(\"i\",\"Original Text\",4)": "5",
		"=FIND(\"\",\"Original Text\")":    "1",
		"=FIND(\"\",\"Original Text\",2)":  "2",
		// FINDB
		"=FINDB(\"T\",\"Original Text\")":   "10",
		"=FINDB(\"t\",\"Original Text\")":   "13",
		"=FINDB(\"i\",\"Original Text\")":   "3",
		"=FINDB(\"i\",\"Original Text\",4)": "5",
		"=FINDB(\"\",\"Original Text\")":    "1",
		"=FINDB(\"\",\"Original Text\",2)":  "2",
		// LEFT
		"=LEFT(\"Original Text\")":    "O",
		"=LEFT(\"Original Text\",4)":  "Orig",
		"=LEFT(\"Original Text\",0)":  "",
		"=LEFT(\"Original Text\",13)": "Original Text",
		"=LEFT(\"Original Text\",20)": "Original Text",
		// LEFTB
		"=LEFTB(\"Original Text\")":    "O",
		"=LEFTB(\"Original Text\",4)":  "Orig",
		"=LEFTB(\"Original Text\",0)":  "",
		"=LEFTB(\"Original Text\",13)": "Original Text",
		"=LEFTB(\"Original Text\",20)": "Original Text",
		// LEN
		"=LEN(\"\")": "0",
		"=LEN(D1)":   "5",
		// LENB
		"=LENB(\"\")": "0",
		"=LENB(D1)":   "5",
		// LOWER
		"=LOWER(\"test\")":     "test",
		"=LOWER(\"TEST\")":     "test",
		"=LOWER(\"Test\")":     "test",
		"=LOWER(\"TEST 123\")": "test 123",
		// MID
		"=MID(\"Original Text\",7,1)": "a",
		"=MID(\"Original Text\",4,7)": "ginal T",
		"=MID(\"255 years\",3,1)":     "5",
		"=MID(\"text\",3,6)":          "xt",
		"=MID(\"text\",6,0)":          "",
		// MIDB
		"=MIDB(\"Original Text\",7,1)": "a",
		"=MIDB(\"Original Text\",4,7)": "ginal T",
		"=MIDB(\"255 years\",3,1)":     "5",
		"=MIDB(\"text\",3,6)":          "xt",
		"=MIDB(\"text\",6,0)":          "",
		// PROPER
		"=PROPER(\"this is a test sentence\")": "This Is A Test Sentence",
		"=PROPER(\"THIS IS A TEST SENTENCE\")": "This Is A Test Sentence",
		"=PROPER(\"123tEST teXT\")":            "123Test Text",
		"=PROPER(\"Mr. SMITH's address\")":     "Mr. Smith'S Address",
		// REPLACE
		"=REPLACE(\"test string\",7,3,\"X\")":          "test sXng",
		"=REPLACE(\"second test string\",8,4,\"XXX\")": "second XXX string",
		"=REPLACE(\"text\",5,0,\" and char\")":         "text and char",
		"=REPLACE(\"text\",1,20,\"char and \")":        "char and ",
		// REPLACEB
		"=REPLACEB(\"test string\",7,3,\"X\")":          "test sXng",
		"=REPLACEB(\"second test string\",8,4,\"XXX\")": "second XXX string",
		"=REPLACEB(\"text\",5,0,\" and char\")":         "text and char",
		"=REPLACEB(\"text\",1,20,\"char and \")":        "char and ",
		// REPT
		"=REPT(\"*\",0)":  "",
		"=REPT(\"*\",1)":  "*",
		"=REPT(\"**\",2)": "****",
		// RIGHT
		"=RIGHT(\"Original Text\")":    "t",
		"=RIGHT(\"Original Text\",4)":  "Text",
		"=RIGHT(\"Original Text\",0)":  "",
		"=RIGHT(\"Original Text\",13)": "Original Text",
		"=RIGHT(\"Original Text\",20)": "Original Text",
		// RIGHTB
		"=RIGHTB(\"Original Text\")":    "t",
		"=RIGHTB(\"Original Text\",4)":  "Text",
		"=RIGHTB(\"Original Text\",0)":  "",
		"=RIGHTB(\"Original Text\",13)": "Original Text",
		"=RIGHTB(\"Original Text\",20)": "Original Text",
		// SUBSTITUTE
		"=SUBSTITUTE(\"abab\",\"a\",\"X\")":                      "XbXb",
		"=SUBSTITUTE(\"abab\",\"a\",\"X\",2)":                    "abXb",
		"=SUBSTITUTE(\"abab\",\"x\",\"X\",2)":                    "abab",
		"=SUBSTITUTE(\"John is 5 years old\",\"John\",\"Jack\")": "Jack is 5 years old",
		"=SUBSTITUTE(\"John is 5 years old\",\"5\",\"6\")":       "John is 6 years old",
		// TEXTJOIN
		"=TEXTJOIN(\"-\",TRUE,1,2,3,4)":  "1-2-3-4",
		"=TEXTJOIN(A4,TRUE,A1:B2)":       "1040205",
		"=TEXTJOIN(\",\",FALSE,A1:C2)":   "1,4,,2,5,",
		"=TEXTJOIN(\",\",TRUE,A1:C2)":    "1,4,2,5",
		"=TEXTJOIN(\",\",TRUE,MUNIT(2))": "1,0,0,1",
		// TRIM
		"=TRIM(\" trim text \")": "trim text",
		"=TRIM(0)":               "0",
		// UNICHAR
		"=UNICHAR(65)": "A",
		"=UNICHAR(97)": "a",
		"=UNICHAR(63)": "?",
		"=UNICHAR(51)": "3",
		// UNICODE
		"=UNICODE(\"Alpha\")": "65",
		"=UNICODE(\"alpha\")": "97",
		"=UNICODE(\"?\")":     "63",
		"=UNICODE(\"3\")":     "51",
		// UPPER
		"=UPPER(\"test\")":     "TEST",
		"=UPPER(\"TEST\")":     "TEST",
		"=UPPER(\"Test\")":     "TEST",
		"=UPPER(\"TEST 123\")": "TEST 123",
		// VALUE
		"=VALUE(\"50\")":                  "50",
		"=VALUE(\"1.0E-07\")":             "1E-07",
		"=VALUE(\"5,000\")":               "5000",
		"=VALUE(\"20%\")":                 "0.2",
		"=VALUE(\"12:00:00\")":            "0.5",
		"=VALUE(\"01/02/2006 15:04:05\")": "38719.6278356481",
		// Conditional Functions
		// IF
		"=IF(1=1)":                              "TRUE",
		"=IF(1<>1)":                             "FALSE",
		"=IF(5<0, \"negative\", \"positive\")":  "positive",
		"=IF(-2<0, \"negative\", \"positive\")": "negative",
		`=IF(1=1, "equal", "notequal")`:         "equal",
		`=IF(1<>1, "equal", "notequal")`:        "notequal",
		`=IF("A"="A", "equal", "notequal")`:     "equal",
		`=IF("A"<>"A", "equal", "notequal")`:    "notequal",
		`=IF(FALSE,0,ROUND(4/2,0))`:             "2",
		`=IF(TRUE,ROUND(4/2,0),0)`:              "2",
		// Excel Lookup and Reference Functions
		// ADDRESS
		"=ADDRESS(1,1,1,TRUE)":            "$A$1",
		"=ADDRESS(1,1,1,FALSE)":           "R1C1",
		"=ADDRESS(1,1,2,TRUE)":            "A$1",
		"=ADDRESS(1,1,2,FALSE)":           "R1C[1]",
		"=ADDRESS(1,1,3,TRUE)":            "$A1",
		"=ADDRESS(1,1,3,FALSE)":           "R[1]C1",
		"=ADDRESS(1,1,4,TRUE)":            "A1",
		"=ADDRESS(1,1,4,FALSE)":           "R[1]C[1]",
		"=ADDRESS(1,1,4,TRUE,\"\")":       "A1",
		"=ADDRESS(1,1,4,TRUE,\"Sheet1\")": "Sheet1!A1",
		// CHOOSE
		"=CHOOSE(4,\"red\",\"blue\",\"green\",\"brown\")": "brown",
		"=CHOOSE(1,\"red\",\"blue\",\"green\",\"brown\")": "red",
		"=SUM(CHOOSE(A2,A1,B1:B2,A1:A3,A1:A4))":           "9",
		// COLUMN
		"=COLUMN()":                "3",
		"=COLUMN(Sheet1!A1)":       "1",
		"=COLUMN(Sheet1!A1:B1:C1)": "1",
		"=COLUMN(Sheet1!F1:G1)":    "6",
		"=COLUMN(H1)":              "8",
		// COLUMNS
		"=COLUMNS(B1)":                   "1",
		"=COLUMNS(1:1)":                  "16384",
		"=COLUMNS(Sheet1!1:1)":           "16384",
		"=COLUMNS(B1:E5)":                "4",
		"=COLUMNS(Sheet1!E5:H7:B1)":      "7",
		"=COLUMNS(E5:H7:B1:C1:Z1:C1:B1)": "25",
		"=COLUMNS(E5:B1)":                "4",
		"=COLUMNS(EM38:HZ81)":            "92",
		// HLOOKUP
		"=HLOOKUP(D2,D2:D8,1,FALSE)":          "Jan",
		"=HLOOKUP(F3,F3:F8,3,FALSE)":          "34440",
		"=HLOOKUP(INT(F3),F3:F8,3,FALSE)":     "34440",
		"=HLOOKUP(MUNIT(1),MUNIT(1),1,FALSE)": "1",
		// VLOOKUP
		"=VLOOKUP(D2,D:D,1,FALSE)":            "Jan",
		"=VLOOKUP(D2,D1:D10,1)":               "Jan",
		"=VLOOKUP(D2,D1:D11,1)":               "Feb",
		"=VLOOKUP(D2,D1:D10,1,FALSE)":         "Jan",
		"=VLOOKUP(INT(36693),F2:F2,1,FALSE)":  "36693",
		"=VLOOKUP(INT(F2),F3:F9,1)":           "32080",
		"=VLOOKUP(INT(F2),F3:F9,1,TRUE)":      "32080",
		"=VLOOKUP(MUNIT(3),MUNIT(3),1)":       "0",
		"=VLOOKUP(A1,A3:B5,1)":                "0",
		"=VLOOKUP(MUNIT(1),MUNIT(1),1,FALSE)": "1",
		// INDEX
		"=INDEX(0,0,0)":          "0",
		"=INDEX(A1,0,0)":         "1",
		"=INDEX(A1:A1,0,0)":      "1",
		"=SUM(INDEX(A1:B1,1))":   "5",
		"=SUM(INDEX(A1:B1,1,0))": "5",
		"=SUM(INDEX(A1:B2,2,0))": "7",
		"=SUM(INDEX(A1:B4,0,2))": "9",
		"=SUM(INDEX(E1:F5,5,2))": "34440",
		// INDIRECT
		"=INDIRECT(\"E1\")":                   "Team",
		"=INDIRECT(\"E\"&1)":                  "Team",
		"=INDIRECT(\"E\"&ROW())":              "Team",
		"=INDIRECT(\"E\"&ROW(),TRUE)":         "Team",
		"=INDIRECT(\"R1C5\",FALSE)":           "Team",
		"=INDIRECT(\"R\"&1&\"C\"&5,FALSE)":    "Team",
		"=SUM(INDIRECT(\"A1:B2\"))":           "12",
		"=SUM(INDIRECT(\"A1:B2\",TRUE))":      "12",
		"=SUM(INDIRECT(\"R1C1:R2C2\",FALSE))": "12",
		// LOOKUP
		"=LOOKUP(F8,F8:F9,F8:F9)":      "32080",
		"=LOOKUP(F8,F8:F9,D8:D9)":      "Feb",
		"=LOOKUP(E3,E2:E5,F2:F5)":      "22100",
		"=LOOKUP(E3,E2:F5)":            "22100",
		"=LOOKUP(F3+1,F3:F4,F3:F4)":    "22100",
		"=LOOKUP(F4+1,F3:F4,F3:F4)":    "53321",
		"=LOOKUP(1,MUNIT(1))":          "1",
		"=LOOKUP(1,MUNIT(1),MUNIT(1))": "1",
		// ROW
		"=ROW()":                "1",
		"=ROW(Sheet1!A1)":       "1",
		"=ROW(Sheet1!A1:B2:C3)": "1",
		"=ROW(Sheet1!F5:G6)":    "5",
		"=ROW(A8)":              "8",
		// ROWS
		"=ROWS(B1)":                    "1",
		"=ROWS(B:B)":                   "1048576",
		"=ROWS(Sheet1!B:B)":            "1048576",
		"=ROWS(B1:E5)":                 "5",
		"=ROWS(Sheet1!E5:H7:B1)":       "7",
		"=ROWS(E5:H8:B2:C3:Z26:C3:B2)": "25",
		"=ROWS(E5:B1)":                 "5",
		"=ROWS(EM38:HZ81)":             "44",
		// Web Functions
		// ENCODEURL
		"=ENCODEURL(\"https://xuri.me/excelize/en/?q=Save As\")": "https%3A%2F%2Fxuri.me%2Fexcelize%2Fen%2F%3Fq%3DSave%20As",
		// Financial Functions
		// ACCRINT
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,4,0,TRUE)":  "1600",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,4,0,FALSE)": "1600",
		// ACCRINTM
		"=ACCRINTM(\"01/01/2012\",\"12/31/2012\",8%,10000)":   "800",
		"=ACCRINTM(\"01/01/2012\",\"12/31/2012\",8%,10000,3)": "800",
		// AMORDEGRC
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,20%)":    "42",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,20%,4)":  "42",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,40%,4)":  "42",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,25%,4)":  "41",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",109,1,25%,4)": "54",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",110,2,25%,4)": "0",
		// AMORLINC
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,1,20%,4)":  "30",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,1,0%,4)":   "0",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,20,15%,4)": "0",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,6,15%,4)":  "0.6875",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,0,15%,4)":  "16.8125",
		// COUPDAYBS
		"=COUPDAYBS(\"02/24/2000\",\"11/24/2000\",4,4)": "0",
		"=COUPDAYBS(\"03/27/2000\",\"11/29/2000\",4,4)": "28",
		"=COUPDAYBS(\"02/29/2000\",\"04/01/2000\",4,4)": "58",
		"=COUPDAYBS(\"01/01/2011\",\"10/25/2012\",4)":   "66",
		"=COUPDAYBS(\"01/01/2011\",\"10/25/2012\",4,1)": "68",
		"=COUPDAYBS(\"10/31/2011\",\"02/26/2012\",4,0)": "65",
		// COUPDAYS
		"=COUPDAYS(\"01/01/2011\",\"10/25/2012\",4)":   "90",
		"=COUPDAYS(\"01/01/2011\",\"10/25/2012\",4,1)": "92",
		// COUPDAYSNC
		"=COUPDAYSNC(\"01/01/2011\",\"10/25/2012\",4)": "24",
		"=COUPDAYSNC(\"04/01/2012\",\"03/31/2020\",2)": "179",
		// COUPNCD
		"=COUPNCD(\"01/01/2011\",\"10/25/2012\",4)":   "40568",
		"=COUPNCD(\"01/01/2011\",\"10/25/2012\",4,0)": "40568",
		"=COUPNCD(\"10/25/2011\",\"01/01/2012\",4)":   "40909",
		"=COUPNCD(\"04/01/2012\",\"03/31/2020\",2)":   "41182",
		"=COUPNCD(\"01/01/2000\",\"08/30/2001\",2)":   "36585",
		// COUPNUM
		"=COUPNUM(\"01/01/2011\",\"10/25/2012\",4)":   "8",
		"=COUPNUM(\"01/01/2011\",\"10/25/2012\",4,0)": "8",
		"=COUPNUM(\"09/30/2017\",\"03/31/2021\",4,0)": "14",
		// COUPPCD
		"=COUPPCD(\"01/01/2011\",\"10/25/2012\",4)":   "40476",
		"=COUPPCD(\"01/01/2011\",\"10/25/2012\",4,0)": "40476",
		"=COUPPCD(\"10/25/2011\",\"01/01/2012\",4)":   "40817",
		// CUMIPMT
		"=CUMIPMT(0.05/12,60,50000,1,12,0)":  "-2294.97753732664",
		"=CUMIPMT(0.05/12,60,50000,13,24,0)": "-1833.10006657389",
		// CUMPRINC
		"=CUMPRINC(0.05/12,60,50000,1,12,0)":  "-9027.76264907988",
		"=CUMPRINC(0.05/12,60,50000,13,24,0)": "-9489.64011983263",
		// DB
		"=DB(0,1000,5,1)":       "0",
		"=DB(10000,1000,5,1)":   "3690",
		"=DB(10000,1000,5,2)":   "2328.39",
		"=DB(10000,1000,5,1,6)": "1845",
		"=DB(10000,1000,5,6,6)": "238.527124587882",
		// DDB
		"=DDB(0,1000,5,1)":     "0",
		"=DDB(10000,1000,5,1)": "4000",
		"=DDB(10000,1000,5,2)": "2400",
		"=DDB(10000,1000,5,3)": "1440",
		"=DDB(10000,1000,5,4)": "864",
		"=DDB(10000,1000,5,5)": "296",
		// DISC
		"=DISC(\"04/01/2016\",\"03/31/2021\",95,100)": "0.01",
		// DOLLARDE
		"=DOLLARDE(1.01,16)": "1.0625",
		// DOLLARFR
		"=DOLLARFR(1.0625,16)": "1.01",
		// DURATION
		"=DURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,4)": "6.67442279848313",
		// EFFECT
		"=EFFECT(0.1,4)":   "0.103812890625",
		"=EFFECT(0.025,2)": "0.02515625",
		// FV
		"=FV(0.05/12,60,-1000)":   "68006.0828408434",
		"=FV(0.1/4,16,-2000,0,1)": "39729.4608941662",
		"=FV(0,16,-2000)":         "32000",
		// FVSCHEDULE
		"=FVSCHEDULE(10000,A1:A5)": "240000",
		"=FVSCHEDULE(10000,0.5)":   "15000",
		// INTRATE
		"=INTRATE(\"04/01/2005\",\"03/31/2010\",1000,2125)": "0.225",
		// IPMT
		"=IPMT(0.05/12,2,60,50000)":   "-205.26988187972",
		"=IPMT(0.035/4,2,8,0,5000,1)": "5.25745523782908",
		// ISPMT
		"=ISPMT(0.05/12,1,60,50000)": "-204.861111111111",
		"=ISPMT(0.05/12,2,60,50000)": "-201.388888888889",
		"=ISPMT(0.05/12,2,1,50000)":  "208.333333333333",
		// MDURATION
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,4)": "6.54355176321876",
		// NOMINAL
		"=NOMINAL(0.025,12)": "0.0247180352381129",
		// NPER
		"=NPER(0.04,-6000,50000)":           "10.3380350715077",
		"=NPER(0,-6000,50000)":              "8.33333333333333",
		"=NPER(0.06/4,-2000,60000,30000,1)": "52.7947737092748",
		// NPV
		"=NPV(0.02,-5000,\"\",800)": "-4133.02575932334",
		// ODDFPRICE
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,2)":              "107.691830256629",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,4,1)":            "106.766915010929",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,4,3)":            "106.7819138147",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,4,4)":            "106.771913772467",
		"=ODDFPRICE(\"11/11/2008\",\"03/01/2021\",\"10/15/2008\",\"03/01/2009\",7.85%,6.25%,100,2,1)":          "113.597717474079",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"09/30/2017\",5.5%,3.5%,100,4,0)":            "106.72930611878",
		"=ODDFPRICE(\"11/11/2008\",\"03/29/2021\", \"08/15/2008\", \"03/29/2009\", 0.0785, 0.0625, 100, 2, 1)": "113.61826640814",
		// PDURATION
		"=PDURATION(0.04,10000,15000)": "10.3380350715076",
		// PMT
		"=PMT(0,8,0,5000,1)":       "-625",
		"=PMT(0.035/4,8,0,5000,1)": "-600.852027180466",
		// PRICE
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,100,2)":   "110.655105178443",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,100,2,4)": "110.655105178443",
		"=PRICE(\"04/01/2012\",\"03/31/2020\",12%,10%,100,2)":   "110.834483593216",
		"=PRICE(\"01/01/2010\",\"06/30/2010\",0.5,1,1,1,4)":     "8.92419088847661",
		// PPMT
		"=PPMT(0.05/12,2,60,50000)":   "-738.291800320824",
		"=PPMT(0.035/4,2,8,0,5000,1)": "-606.109482418295",
		// PRICEDISC
		"=PRICEDISC(\"04/01/2017\",\"03/31/2021\",2.5%,100)":   "90",
		"=PRICEDISC(\"04/01/2017\",\"03/31/2021\",2.5%,100,3)": "90",
		// PRICEMAT
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",4.5%,2.5%)":   "107.170454545455",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",4.5%,2.5%,0)": "107.170454545455",
		// PV
		"=PV(0,60,1000)":         "-60000",
		"=PV(5%/12,60,1000)":     "-52990.7063239275",
		"=PV(10%/4,16,2000,0,1)": "-26762.7554528811",
		// RATE
		"=RATE(60,-1000,50000)":       "0.0061834131621292",
		"=RATE(24,-800,0,20000,1)":    "0.00325084350160374",
		"=RATE(48,-200,8000,3,1,0.5)": "0.0080412665831637",
		// RECEIVED
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",1000,4.5%)":   "1290.32258064516",
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",1000,4.5%,0)": "1290.32258064516",
		// RRI
		"=RRI(10,10000,15000)": "0.0413797439924106",
		// SLN
		"=SLN(10000,1000,5)": "1800",
		// SYD
		"=SYD(10000,1000,5,1)": "3000",
		"=SYD(10000,1000,5,2)": "2400",
		// TBILLEQ
		"=TBILLEQ(\"01/01/2017\",\"06/30/2017\",2.5%)": "0.0256680731364276",
		// TBILLPRICE
		"=TBILLPRICE(\"02/01/2017\",\"06/30/2017\",2.75%)": "98.8618055555556",
		// TBILLYIELD
		"=TBILLYIELD(\"02/01/2017\",\"06/30/2017\",99)": "0.024405125076266",
		// VDB
		"=VDB(10000,1000,5,0,1)":           "4000",
		"=VDB(10000,1000,5,1,3)":           "3840",
		"=VDB(10000,1000,5,3,5)":           "1160",
		"=VDB(10000,1000,5,3,5,0.2,FALSE)": "3600",
		"=VDB(10000,1000,5,3,5,0.2,TRUE)":  "693.633024",
		"=VDB(24000,3000,10,0,0.875,2)":    "4200",
		"=VDB(24000,3000,10,0.1,1)":        "4233.6",
		"=VDB(24000,3000,10,0.1,1,1)":      "2138.4",
		"=VDB(24000,3000,100,50,100,1)":    "10377.2944184652",
		"=VDB(24000,3000,100,50,100,2)":    "5740.0723220908",
		// YIELD
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,100,4)":               "0.0975631546829798",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,100,4,4)":             "0.0976269355643988",
		"=YIELD(\"01/01/2010\",\"06/30/2010\",0.5,1,1,1,4)":                 "1.91285866099894",
		"=YIELD(\"01/01/2010\",\"06/30/2010\",0,1,1,1,4)":                   "0",
		"=YIELD(\"01/01/2010\",\"01/02/2020\",100,68.15518653988686,1,1,1)": "64",
		// YIELDDISC
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",97,100)":   "0.0622012325059031",
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",97,100,0)": "0.0622012325059031",
		// YIELDMAT
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",5.5%,101)":   "0.0419422478838651",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",5.5%,101,0)": "0.0419422478838651",
	}
	for formula, expected := range mathCalc {
		f := prepareCalcData(cellData)
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	mathCalcError := map[string]string{
		"=1/0":       "#DIV/0!",
		"1^\"text\"": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"\"text\"^1": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"1+\"text\"": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"\"text\"+1": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"1-\"text\"": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"\"text\"-1": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"1*\"text\"": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"\"text\"*1": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"1/\"text\"": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		"\"text\"/1": "strconv.ParseFloat: parsing \"text\": invalid syntax",
		// Engineering Functions
		// BESSELI
		"=BESSELI()":       "BESSELI requires 2 numeric arguments",
		"=BESSELI(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BESSELI(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// BESSELJ
		"=BESSELJ()":       "BESSELJ requires 2 numeric arguments",
		"=BESSELJ(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BESSELJ(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// BESSELK
		"=BESSELK()":       "BESSELK requires 2 numeric arguments",
		"=BESSELK(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BESSELK(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BESSELK(-1,0)":   "#NUM!",
		"=BESSELK(1,-1)":   "#NUM!",
		// BESSELY
		"=BESSELY()":       "BESSELY requires 2 numeric arguments",
		"=BESSELY(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BESSELY(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BESSELY(-1,0)":   "#NUM!",
		"=BESSELY(1,-1)":   "#NUM!",
		// BIN2DEC
		"=BIN2DEC()":     "BIN2DEC requires 1 numeric argument",
		"=BIN2DEC(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// BIN2HEX
		"=BIN2HEX()":               "BIN2HEX requires at least 1 argument",
		"=BIN2HEX(1,1,1)":          "BIN2HEX allows at most 2 arguments",
		"=BIN2HEX(\"\",1)":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BIN2HEX(1,\"\")":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BIN2HEX(12345678901,10)": "#NUM!",
		"=BIN2HEX(1,-1)":           "#NUM!",
		"=BIN2HEX(31,1)":           "#NUM!",
		// BIN2OCT
		"=BIN2OCT()":                 "BIN2OCT requires at least 1 argument",
		"=BIN2OCT(1,1,1)":            "BIN2OCT allows at most 2 arguments",
		"=BIN2OCT(\"\",1)":           "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BIN2OCT(1,\"\")":           "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BIN2OCT(-12345678901 ,10)": "#NUM!",
		"=BIN2OCT(1,-1)":             "#NUM!",
		"=BIN2OCT(8,1)":              "#NUM!",
		// BITAND
		"=BITAND()":        "BITAND requires 2 numeric arguments",
		"=BITAND(-1,2)":    "#NUM!",
		"=BITAND(2^48,2)":  "#NUM!",
		"=BITAND(1,-1)":    "#NUM!",
		"=BITAND(\"\",-1)": "#NUM!",
		"=BITAND(1,\"\")":  "#NUM!",
		"=BITAND(1,2^48)":  "#NUM!",
		// BITLSHIFT
		"=BITLSHIFT()":        "BITLSHIFT requires 2 numeric arguments",
		"=BITLSHIFT(-1,2)":    "#NUM!",
		"=BITLSHIFT(2^48,2)":  "#NUM!",
		"=BITLSHIFT(1,-1)":    "#NUM!",
		"=BITLSHIFT(\"\",-1)": "#NUM!",
		"=BITLSHIFT(1,\"\")":  "#NUM!",
		"=BITLSHIFT(1,2^48)":  "#NUM!",
		// BITOR
		"=BITOR()":        "BITOR requires 2 numeric arguments",
		"=BITOR(-1,2)":    "#NUM!",
		"=BITOR(2^48,2)":  "#NUM!",
		"=BITOR(1,-1)":    "#NUM!",
		"=BITOR(\"\",-1)": "#NUM!",
		"=BITOR(1,\"\")":  "#NUM!",
		"=BITOR(1,2^48)":  "#NUM!",
		// BITRSHIFT
		"=BITRSHIFT()":        "BITRSHIFT requires 2 numeric arguments",
		"=BITRSHIFT(-1,2)":    "#NUM!",
		"=BITRSHIFT(2^48,2)":  "#NUM!",
		"=BITRSHIFT(1,-1)":    "#NUM!",
		"=BITRSHIFT(\"\",-1)": "#NUM!",
		"=BITRSHIFT(1,\"\")":  "#NUM!",
		"=BITRSHIFT(1,2^48)":  "#NUM!",
		// BITXOR
		"=BITXOR()":        "BITXOR requires 2 numeric arguments",
		"=BITXOR(-1,2)":    "#NUM!",
		"=BITXOR(2^48,2)":  "#NUM!",
		"=BITXOR(1,-1)":    "#NUM!",
		"=BITXOR(\"\",-1)": "#NUM!",
		"=BITXOR(1,\"\")":  "#NUM!",
		"=BITXOR(1,2^48)":  "#NUM!",
		// COMPLEX
		"=COMPLEX()":              "COMPLEX requires at least 2 arguments",
		"=COMPLEX(10,-5,\"\")":    "#VALUE!",
		"=COMPLEX(\"\",0)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=COMPLEX(0,\"\")":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=COMPLEX(10,-5,\"i\",0)": "COMPLEX allows at most 3 arguments",
		// DEC2BIN
		"=DEC2BIN()":        "DEC2BIN requires at least 1 argument",
		"=DEC2BIN(1,1,1)":   "DEC2BIN allows at most 2 arguments",
		"=DEC2BIN(\"\",1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DEC2BIN(1,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DEC2BIN(-513,10)": "#NUM!",
		"=DEC2BIN(1,-1)":    "#NUM!",
		"=DEC2BIN(2,1)":     "#NUM!",
		// DEC2HEX
		"=DEC2HEX()":                 "DEC2HEX requires at least 1 argument",
		"=DEC2HEX(1,1,1)":            "DEC2HEX allows at most 2 arguments",
		"=DEC2HEX(\"\",1)":           "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DEC2HEX(1,\"\")":           "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DEC2HEX(-549755813888,10)": "#NUM!",
		"=DEC2HEX(1,-1)":             "#NUM!",
		"=DEC2HEX(31,1)":             "#NUM!",
		// DEC2OCT
		"=DEC2OCT()":               "DEC2OCT requires at least 1 argument",
		"=DEC2OCT(1,1,1)":          "DEC2OCT allows at most 2 arguments",
		"=DEC2OCT(\"\",1)":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DEC2OCT(1,\"\")":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DEC2OCT(-536870912 ,10)": "#NUM!",
		"=DEC2OCT(1,-1)":           "#NUM!",
		"=DEC2OCT(8,1)":            "#NUM!",
		// DELTA
		"=DELTA()":       "DELTA requires at least 1 argument",
		"=DELTA(0,0,0)":  "DELTA allows at most 2 arguments",
		"=DELTA(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DELTA(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// ERF
		"=ERF()":       "ERF requires at least 1 argument",
		"=ERF(0,0,0)":  "ERF allows at most 2 arguments",
		"=ERF(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ERF(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// ERF.PRECISE
		"=ERF.PRECISE()":     "ERF.PRECISE requires 1 argument",
		"=ERF.PRECISE(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// ERFC
		"=ERFC()":     "ERFC requires 1 argument",
		"=ERFC(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// ERFC.PRECISE
		"=ERFC.PRECISE()":     "ERFC.PRECISE requires 1 argument",
		"=ERFC.PRECISE(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// GESTEP
		"=GESTEP()":       "GESTEP requires at least 1 argument",
		"=GESTEP(0,0,0)":  "GESTEP allows at most 2 arguments",
		"=GESTEP(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GESTEP(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// HEX2BIN
		"=HEX2BIN()":        "HEX2BIN requires at least 1 argument",
		"=HEX2BIN(1,1,1)":   "HEX2BIN allows at most 2 arguments",
		"=HEX2BIN(\"X\",1)": "strconv.ParseInt: parsing \"X\": invalid syntax",
		"=HEX2BIN(1,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HEX2BIN(-513,10)": "strconv.ParseInt: parsing \"-\": invalid syntax",
		"=HEX2BIN(1,-1)":    "#NUM!",
		"=HEX2BIN(2,1)":     "#NUM!",
		// HEX2DEC
		"=HEX2DEC()":      "HEX2DEC requires 1 numeric argument",
		"=HEX2DEC(\"X\")": "strconv.ParseInt: parsing \"X\": invalid syntax",
		// HEX2OCT
		"=HEX2OCT()":        "HEX2OCT requires at least 1 argument",
		"=HEX2OCT(1,1,1)":   "HEX2OCT allows at most 2 arguments",
		"=HEX2OCT(\"X\",1)": "strconv.ParseInt: parsing \"X\": invalid syntax",
		"=HEX2OCT(1,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HEX2OCT(-513,10)": "strconv.ParseInt: parsing \"-\": invalid syntax",
		"=HEX2OCT(1,-1)":    "#NUM!",
		// IMABS
		"=IMABS()":     "IMABS requires 1 argument",
		"=IMABS(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMAGINARY
		"=IMAGINARY()":     "IMAGINARY requires 1 argument",
		"=IMAGINARY(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMARGUMENT
		"=IMARGUMENT()":     "IMARGUMENT requires 1 argument",
		"=IMARGUMENT(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMCONJUGATE
		"=IMCONJUGATE()":     "IMCONJUGATE requires 1 argument",
		"=IMCONJUGATE(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMCOS
		"=IMCOS()":     "IMCOS requires 1 argument",
		"=IMCOS(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMCOSH
		"=IMCOSH()":     "IMCOSH requires 1 argument",
		"=IMCOSH(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMCOT
		"=IMCOT()":     "IMCOT requires 1 argument",
		"=IMCOT(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMCSC
		"=IMCSC()":     "IMCSC requires 1 argument",
		"=IMCSC(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMCSC(0)":    "#NUM!",
		// IMCSCH
		"=IMCSCH()":     "IMCSCH requires 1 argument",
		"=IMCSCH(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMCSCH(0)":    "#NUM!",
		// IMDIV
		"=IMDIV()":       "IMDIV requires 2 arguments",
		"=IMDIV(0,\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMDIV(\"\",0)": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMDIV(1,0)":    "#NUM!",
		// IMEXP
		"=IMEXP()":     "IMEXP requires 1 argument",
		"=IMEXP(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMLN
		"=IMLN()":     "IMLN requires 1 argument",
		"=IMLN(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMLN(0)":    "#NUM!",
		// IMLOG10
		"=IMLOG10()":     "IMLOG10 requires 1 argument",
		"=IMLOG10(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMLOG10(0)":    "#NUM!",
		// IMLOG2
		"=IMLOG2()":     "IMLOG2 requires 1 argument",
		"=IMLOG2(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMLOG2(0)":    "#NUM!",
		// IMPOWER
		"=IMPOWER()":       "IMPOWER requires 2 arguments",
		"=IMPOWER(0,\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMPOWER(\"\",0)": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMPOWER(0,0)":    "#NUM!",
		"=IMPOWER(0,-1)":   "#NUM!",
		// IMPRODUCT
		"=IMPRODUCT(\"x\")": "strconv.ParseComplex: parsing \"x\": invalid syntax",
		"=IMPRODUCT(A1:D1)": "strconv.ParseComplex: parsing \"Month\": invalid syntax",
		// IMREAL
		"=IMREAL()":     "IMREAL requires 1 argument",
		"=IMREAL(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMSEC
		"=IMSEC()":     "IMSEC requires 1 argument",
		"=IMSEC(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMSECH
		"=IMSECH()":     "IMSECH requires 1 argument",
		"=IMSECH(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMSIN
		"=IMSIN()":     "IMSIN requires 1 argument",
		"=IMSIN(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMSINH
		"=IMSINH()":     "IMSINH requires 1 argument",
		"=IMSINH(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMSQRT
		"=IMSQRT()":     "IMSQRT requires 1 argument",
		"=IMSQRT(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMSUB
		"=IMSUB()":       "IMSUB requires 2 arguments",
		"=IMSUB(0,\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		"=IMSUB(\"\",0)": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMSUM
		"=IMSUM()":     "IMSUM requires at least 1 argument",
		"=IMSUM(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// IMTAN
		"=IMTAN()":     "IMTAN requires 1 argument",
		"=IMTAN(\"\")": "strconv.ParseComplex: parsing \"\": invalid syntax",
		// OCT2BIN
		"=OCT2BIN()":               "OCT2BIN requires at least 1 argument",
		"=OCT2BIN(1,1,1)":          "OCT2BIN allows at most 2 arguments",
		"=OCT2BIN(\"\",1)":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=OCT2BIN(1,\"\")":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=OCT2BIN(-536870912 ,10)": "#NUM!",
		"=OCT2BIN(1,-1)":           "#NUM!",
		// OCT2DEC
		"=OCT2DEC()":     "OCT2DEC requires 1 numeric argument",
		"=OCT2DEC(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// OCT2HEX
		"=OCT2HEX()":               "OCT2HEX requires at least 1 argument",
		"=OCT2HEX(1,1,1)":          "OCT2HEX allows at most 2 arguments",
		"=OCT2HEX(\"\",1)":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=OCT2HEX(1,\"\")":         "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=OCT2HEX(-536870912 ,10)": "#NUM!",
		"=OCT2HEX(1,-1)":           "#NUM!",
		// Math and Trigonometric Functions
		// ABS
		"=ABS()":    "ABS requires 1 numeric argument",
		`=ABS("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=ABS(~)":   newInvalidColumnNameError("~").Error(),
		// ACOS
		"=ACOS()":        "ACOS requires 1 numeric argument",
		`=ACOS("X")`:     "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=ACOS(ACOS(0))": "#NUM!",
		// ACOSH
		"=ACOSH()":    "ACOSH requires 1 numeric argument",
		`=ACOSH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.ACOT
		"=_xlfn.ACOT()":    "ACOT requires 1 numeric argument",
		`=_xlfn.ACOT("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.ACOTH
		"=_xlfn.ACOTH()":               "ACOTH requires 1 numeric argument",
		`=_xlfn.ACOTH("X")`:            "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=_xlfn.ACOTH(_xlfn.ACOTH(2))": "#NUM!",
		// _xlfn.ARABIC
		"=_xlfn.ARABIC()": "ARABIC requires 1 numeric argument",
		"=_xlfn.ARABIC(\"" + strings.Repeat("I", 256) + "\")": "#VALUE!",
		// ASIN
		"=ASIN()":    "ASIN requires 1 numeric argument",
		`=ASIN("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ASINH
		"=ASINH()":    "ASINH requires 1 numeric argument",
		`=ASINH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ATAN
		"=ATAN()":    "ATAN requires 1 numeric argument",
		`=ATAN("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ATANH
		"=ATANH()":    "ATANH requires 1 numeric argument",
		`=ATANH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ATAN2
		"=ATAN2()":      "ATAN2 requires 2 numeric arguments",
		`=ATAN2("X",0)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=ATAN2(0,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// BASE
		"=BASE()":        "BASE requires at least 2 arguments",
		"=BASE(1,2,3,4)": "BASE allows at most 3 arguments",
		"=BASE(1,1)":     "radix must be an integer >= 2 and <= 36",
		`=BASE("X",2)`:   "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=BASE(1,"X")`:   "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=BASE(1,2,"X")`: "strconv.Atoi: parsing \"X\": invalid syntax",
		// CEILING
		"=CEILING()":      "CEILING requires at least 1 argument",
		"=CEILING(1,2,3)": "CEILING allows at most 2 arguments",
		"=CEILING(1,-1)":  "negative sig to CEILING invalid",
		`=CEILING("X",0)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=CEILING(0,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.CEILING.MATH
		"=_xlfn.CEILING.MATH()":        "CEILING.MATH requires at least 1 argument",
		"=_xlfn.CEILING.MATH(1,2,3,4)": "CEILING.MATH allows at most 3 arguments",
		`=_xlfn.CEILING.MATH("X")`:     "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=_xlfn.CEILING.MATH(1,"X")`:   "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=_xlfn.CEILING.MATH(1,2,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.CEILING.PRECISE
		"=_xlfn.CEILING.PRECISE()":      "CEILING.PRECISE requires at least 1 argument",
		"=_xlfn.CEILING.PRECISE(1,2,3)": "CEILING.PRECISE allows at most 2 arguments",
		`=_xlfn.CEILING.PRECISE("X",2)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=_xlfn.CEILING.PRECISE(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// COMBIN
		"=COMBIN()":       "COMBIN requires 2 argument",
		"=COMBIN(-1,1)":   "COMBIN requires number >= number_chosen",
		`=COMBIN("X",1)`:  "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=COMBIN(-1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.COMBINA
		"=_xlfn.COMBINA()":       "COMBINA requires 2 argument",
		"=_xlfn.COMBINA(-1,1)":   "COMBINA requires number > number_chosen",
		"=_xlfn.COMBINA(-1,-1)":  "COMBIN requires number >= number_chosen",
		`=_xlfn.COMBINA("X",1)`:  "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=_xlfn.COMBINA(-1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// COS
		"=COS()":    "COS requires 1 numeric argument",
		`=COS("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// COSH
		"=COSH()":    "COSH requires 1 numeric argument",
		`=COSH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.COT
		"=COT()":    "COT requires 1 numeric argument",
		`=COT("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=COT(0)":   "#DIV/0!",
		// _xlfn.COTH
		"=COTH()":    "COTH requires 1 numeric argument",
		`=COTH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=COTH(0)":   "#DIV/0!",
		// _xlfn.CSC
		"=_xlfn.CSC()":    "CSC requires 1 numeric argument",
		`=_xlfn.CSC("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=_xlfn.CSC(0)":   "#DIV/0!",
		// _xlfn.CSCH
		"=_xlfn.CSCH()":    "CSCH requires 1 numeric argument",
		`=_xlfn.CSCH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=_xlfn.CSCH(0)":   "#DIV/0!",
		// _xlfn.DECIMAL
		"=_xlfn.DECIMAL()":          "DECIMAL requires 2 numeric arguments",
		`=_xlfn.DECIMAL("X", 2)`:    "strconv.ParseInt: parsing \"X\": invalid syntax",
		`=_xlfn.DECIMAL(2000, "X")`: "strconv.Atoi: parsing \"X\": invalid syntax",
		// DEGREES
		"=DEGREES()":    "DEGREES requires 1 numeric argument",
		`=DEGREES("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=DEGREES(0)":   "#DIV/0!",
		// EVEN
		"=EVEN()":    "EVEN requires 1 numeric argument",
		`=EVEN("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// EXP
		"=EXP()":    "EXP requires 1 numeric argument",
		`=EXP("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// FACT
		"=FACT()":    "FACT requires 1 numeric argument",
		`=FACT("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=FACT(-1)":  "#NUM!",
		// FACTDOUBLE
		"=FACTDOUBLE()":    "FACTDOUBLE requires 1 numeric argument",
		`=FACTDOUBLE("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=FACTDOUBLE(-1)":  "#NUM!",
		// FLOOR
		"=FLOOR()":       "FLOOR requires 2 numeric arguments",
		`=FLOOR("X",-1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=FLOOR(1,"X")`:  "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=FLOOR(1,-1)":   "invalid arguments to FLOOR",
		// _xlfn.FLOOR.MATH
		"=_xlfn.FLOOR.MATH()":        "FLOOR.MATH requires at least 1 argument",
		"=_xlfn.FLOOR.MATH(1,2,3,4)": "FLOOR.MATH allows at most 3 arguments",
		`=_xlfn.FLOOR.MATH("X",2,3)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=_xlfn.FLOOR.MATH(1,"X",3)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=_xlfn.FLOOR.MATH(1,2,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.FLOOR.PRECISE
		"=_xlfn.FLOOR.PRECISE()":      "FLOOR.PRECISE requires at least 1 argument",
		"=_xlfn.FLOOR.PRECISE(1,2,3)": "FLOOR.PRECISE allows at most 2 arguments",
		`=_xlfn.FLOOR.PRECISE("X",2)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=_xlfn.FLOOR.PRECISE(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// GCD
		"=GCD()":     "GCD requires at least 1 argument",
		"=GCD(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GCD(-1)":   "GCD only accepts positive arguments",
		"=GCD(1,-1)": "GCD only accepts positive arguments",
		`=GCD("X")`:  "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// INT
		"=INT()":    "INT requires 1 numeric argument",
		`=INT("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ISO.CEILING
		"=ISO.CEILING()":      "ISO.CEILING requires at least 1 argument",
		"=ISO.CEILING(1,2,3)": "ISO.CEILING allows at most 2 arguments",
		`=ISO.CEILING("X",2)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=ISO.CEILING(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// LCM
		"=LCM()":     "LCM requires at least 1 argument",
		"=LCM(-1)":   "LCM only accepts positive arguments",
		"=LCM(1,-1)": "LCM only accepts positive arguments",
		`=LCM("X")`:  "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// LN
		"=LN()":    "LN requires 1 numeric argument",
		`=LN("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// LOG
		"=LOG()":      "LOG requires at least 1 argument",
		"=LOG(1,2,3)": "LOG allows at most 2 arguments",
		`=LOG("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=LOG(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=LOG(0,0)":   "#DIV/0!",
		"=LOG(1,0)":   "#DIV/0!",
		"=LOG(1,1)":   "#DIV/0!",
		// LOG10
		"=LOG10()":    "LOG10 requires 1 numeric argument",
		`=LOG10("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// MDETERM
		"MDETERM()": "MDETERM requires at least 1 argument",
		// MOD
		"=MOD()":      "MOD requires 2 numeric arguments",
		"=MOD(6,0)":   "MOD divide by zero",
		`=MOD("X",0)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=MOD(6,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// MROUND
		"=MROUND()":      "MROUND requires 2 numeric arguments",
		"=MROUND(1,0)":   "#NUM!",
		"=MROUND(1,-1)":  "#NUM!",
		`=MROUND("X",0)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=MROUND(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// MULTINOMIAL
		`=MULTINOMIAL("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.MUNIT
		"=_xlfn.MUNIT()":    "MUNIT requires 1 numeric argument",
		`=_xlfn.MUNIT("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=_xlfn.MUNIT(-1)":  "",
		// ODD
		"=ODD()":    "ODD requires 1 numeric argument",
		`=ODD("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// PI
		"=PI(1)": "PI accepts no arguments",
		// POWER
		`=POWER("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=POWER(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=POWER(0,0)":   "#NUM!",
		"=POWER(0,-1)":  "#DIV/0!",
		"=POWER(1)":     "POWER requires 2 numeric arguments",
		// PRODUCT
		`=PRODUCT("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// QUOTIENT
		`=QUOTIENT("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=QUOTIENT(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=QUOTIENT(1,0)":   "#DIV/0!",
		"=QUOTIENT(1)":     "QUOTIENT requires 2 numeric arguments",
		// RADIANS
		`=RADIANS("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=RADIANS()":    "RADIANS requires 1 numeric argument",
		// RAND
		"=RAND(1)": "RAND accepts no arguments",
		// RANDBETWEEN
		`=RANDBETWEEN("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=RANDBETWEEN(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=RANDBETWEEN()":      "RANDBETWEEN requires 2 numeric arguments",
		"=RANDBETWEEN(2,1)":   "#NUM!",
		// ROMAN
		"=ROMAN()":      "ROMAN requires at least 1 argument",
		"=ROMAN(1,2,3)": "ROMAN allows at most 2 arguments",
		`=ROMAN("X")`:   "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=ROMAN("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ROUND
		"=ROUND()":      "ROUND requires 2 numeric arguments",
		`=ROUND("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=ROUND(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ROUNDDOWN
		"=ROUNDDOWN()":      "ROUNDDOWN requires 2 numeric arguments",
		`=ROUNDDOWN("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=ROUNDDOWN(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// ROUNDUP
		"=ROUNDUP()":      "ROUNDUP requires 2 numeric arguments",
		`=ROUNDUP("X",1)`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=ROUNDUP(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// SEC
		"=_xlfn.SEC()":    "SEC requires 1 numeric argument",
		`=_xlfn.SEC("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// _xlfn.SECH
		"=_xlfn.SECH()":    "SECH requires 1 numeric argument",
		`=_xlfn.SECH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// SERIESSUM
		"=SERIESSUM()":               "SERIESSUM requires 4 arguments",
		"=SERIESSUM(\"\",2,3,A1:A4)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=SERIESSUM(1,\"\",3,A1:A4)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=SERIESSUM(1,2,\"\",A1:A4)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=SERIESSUM(1,2,3,A1:D1)":    "strconv.ParseFloat: parsing \"Month\": invalid syntax",
		// SIGN
		"=SIGN()":    "SIGN requires 1 numeric argument",
		`=SIGN("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// SIN
		"=SIN()":    "SIN requires 1 numeric argument",
		`=SIN("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// SINH
		"=SINH()":    "SINH requires 1 numeric argument",
		`=SINH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// SQRT
		"=SQRT()":    "SQRT requires 1 numeric argument",
		`=SQRT("")`:  "strconv.ParseFloat: parsing \"\": invalid syntax",
		`=SQRT("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=SQRT(-1)":  "#NUM!",
		// SQRTPI
		"=SQRTPI()":    "SQRTPI requires 1 numeric argument",
		`=SQRTPI("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// STDEV
		"=STDEV()":      "STDEV requires at least 1 argument",
		"=STDEV(E2:E9)": "#DIV/0!",
		// STDEV.S
		"=STDEV.S()": "STDEV.S requires at least 1 argument",
		// STDEVA
		"=STDEVA()":      "STDEVA requires at least 1 argument",
		"=STDEVA(E2:E9)": "#DIV/0!",
		// POISSON.DIST
		"=POISSON.DIST()": "POISSON.DIST requires 3 arguments",
		// POISSON
		"=POISSON()":             "POISSON requires 3 arguments",
		"=POISSON(\"\",0,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=POISSON(0,\"\",FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=POISSON(0,0,\"\")":     "strconv.ParseBool: parsing \"\": invalid syntax",
		"=POISSON(0,-1,TRUE)":    "#N/A",
		// SUM
		"=SUM((":             ErrInvalidFormula.Error(),
		"=SUM(-)":            ErrInvalidFormula.Error(),
		"=SUM(1+)":           ErrInvalidFormula.Error(),
		"=SUM(1-)":           ErrInvalidFormula.Error(),
		"=SUM(1*)":           ErrInvalidFormula.Error(),
		"=SUM(1/)":           ErrInvalidFormula.Error(),
		"=SUM(1*SUM(1/0))":   "#DIV/0!",
		"=SUM(1*SUM(1/0)*1)": "#DIV/0!",
		// SUMIF
		"=SUMIF()": "SUMIF requires at least 2 arguments",
		// SUMSQ
		`=SUMSQ("X")`:   "strconv.ParseFloat: parsing \"X\": invalid syntax",
		"=SUMSQ(C1:D2)": "strconv.ParseFloat: parsing \"Month\": invalid syntax",
		// SUMPRODUCT
		"=SUMPRODUCT()":            "SUMPRODUCT requires at least 1 argument",
		"=SUMPRODUCT(A1,B1:B2)":    "#VALUE!",
		"=SUMPRODUCT(A1,D1)":       "#VALUE!",
		"=SUMPRODUCT(A1:A3,D1:D3)": "#VALUE!",
		"=SUMPRODUCT(A1:A2,B1:B3)": "#VALUE!",
		"=SUMPRODUCT(A1,NA())":     "#N/A",
		// SUMX2MY2
		"=SUMX2MY2()":         "SUMX2MY2 requires 2 arguments",
		"=SUMX2MY2(A1,B1:B2)": "#N/A",
		// SUMX2PY2
		"=SUMX2PY2()":         "SUMX2PY2 requires 2 arguments",
		"=SUMX2PY2(A1,B1:B2)": "#N/A",
		// SUMXMY2
		"=SUMXMY2()":         "SUMXMY2 requires 2 arguments",
		"=SUMXMY2(A1,B1:B2)": "#N/A",
		// TAN
		"=TAN()":    "TAN requires 1 numeric argument",
		`=TAN("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// TANH
		"=TANH()":    "TANH requires 1 numeric argument",
		`=TANH("X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// TRUNC
		"=TRUNC()":      "TRUNC requires at least 1 argument",
		`=TRUNC("X")`:   "strconv.ParseFloat: parsing \"X\": invalid syntax",
		`=TRUNC(1,"X")`: "strconv.ParseFloat: parsing \"X\": invalid syntax",
		// Statistical Functions
		// AVEDEV
		"=AVEDEV()":       "AVEDEV requires at least 1 argument",
		"=AVEDEV(\"\")":   "#VALUE!",
		"=AVEDEV(1,\"\")": "#VALUE!",
		// AVERAGE
		"=AVERAGE(H1)": "AVERAGE divide by zero",
		// AVERAGEA
		"=AVERAGEA(H1)": "AVERAGEA divide by zero",
		// AVERAGEIF
		"=AVERAGEIF()":                      "AVERAGEIF requires at least 2 arguments",
		"=AVERAGEIF(H1,\"\")":               "AVERAGEIF divide by zero",
		"=AVERAGEIF(D1:D3,\"Month\",D1:D3)": "AVERAGEIF divide by zero",
		"=AVERAGEIF(C1:C3,\"Month\",D1:D3)": "AVERAGEIF divide by zero",
		// BETA.DIST
		"=BETA.DIST()":                     "BETA.DIST requires at least 4 arguments",
		"=BETA.DIST(0.4,4,5,TRUE,0,1,0)":   "BETA.DIST requires at most 6 arguments",
		"=BETA.DIST(\"\",4,5,TRUE,0,1)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.DIST(0.4,\"\",5,TRUE,0,1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.DIST(0.4,4,\"\",TRUE,0,1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.DIST(0.4,4,5,\"\",0,1)":     "strconv.ParseBool: parsing \"\": invalid syntax",
		"=BETA.DIST(0.4,4,5,TRUE,\"\",1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.DIST(0.4,4,5,TRUE,0,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.DIST(0.4,0,5,TRUE,0,1)":     "#NUM!",
		"=BETA.DIST(0.4,4,0,TRUE,0,0)":     "#NUM!",
		"=BETA.DIST(0.4,4,5,TRUE,0.5,1)":   "#NUM!",
		"=BETA.DIST(0.4,4,5,TRUE,0,0.3)":   "#NUM!",
		"=BETA.DIST(0.4,4,5,TRUE,0.4,0.4)": "#NUM!",
		// BETADIST
		"=BETADIST()":                "BETADIST requires at least 3 arguments",
		"=BETADIST(0.4,4,5,0,1,0)":   "BETADIST requires at most 5 arguments",
		"=BETADIST(\"\",4,5,0,1)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETADIST(0.4,\"\",5,0,1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETADIST(0.4,4,\"\",0,1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETADIST(0.4,4,5,\"\",1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETADIST(0.4,4,5,0,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETADIST(2,4,5,3,1)":       "#NUM!",
		"=BETADIST(2,4,5,0,1)":       "#NUM!",
		"=BETADIST(0.4,0,5,0,1)":     "#NUM!",
		"=BETADIST(0.4,4,0,0,1)":     "#NUM!",
		"=BETADIST(0.4,4,5,0.4,0.4)": "#NUM!",
		// BETAINV
		"=BETAINV()":               "BETAINV requires at least 3 arguments",
		"=BETAINV(0.2,4,5,0,1,0)":  "BETAINV requires at most 5 arguments",
		"=BETAINV(\"\",4,5,0,1)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETAINV(0.2,\"\",5,0,1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETAINV(0.2,4,\"\",0,1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETAINV(0.2,4,5,\"\",1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETAINV(0.2,4,5,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETAINV(0,4,5,0,1)":      "#NUM!",
		"=BETAINV(1,4,5,0,1)":      "#NUM!",
		"=BETAINV(0.2,0,5,0,1)":    "#NUM!",
		"=BETAINV(0.2,4,0,0,1)":    "#NUM!",
		"=BETAINV(0.2,4,5,2,2)":    "#NUM!",
		// BETA.INV
		"=BETA.INV()":               "BETA.INV requires at least 3 arguments",
		"=BETA.INV(0.2,4,5,0,1,0)":  "BETA.INV requires at most 5 arguments",
		"=BETA.INV(\"\",4,5,0,1)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.INV(0.2,\"\",5,0,1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.INV(0.2,4,\"\",0,1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.INV(0.2,4,5,\"\",1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.INV(0.2,4,5,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BETA.INV(0,4,5,0,1)":      "#NUM!",
		"=BETA.INV(1,4,5,0,1)":      "#NUM!",
		"=BETA.INV(0.2,0,5,0,1)":    "#NUM!",
		"=BETA.INV(0.2,4,0,0,1)":    "#NUM!",
		"=BETA.INV(0.2,4,5,2,2)":    "#NUM!",
		// BINOMDIST
		"=BINOMDIST()":                   "BINOMDIST requires 4 arguments",
		"=BINOMDIST(\"\",100,0.5,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOMDIST(10,\"\",0.5,FALSE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOMDIST(10,100,\"\",FALSE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOMDIST(10,100,0.5,\"\")":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=BINOMDIST(-1,100,0.5,FALSE)":   "#NUM!",
		"=BINOMDIST(110,100,0.5,FALSE)":  "#NUM!",
		"=BINOMDIST(10,100,-1,FALSE)":    "#NUM!",
		"=BINOMDIST(10,100,2,FALSE)":     "#NUM!",
		// BINOM.DIST
		"=BINOM.DIST()":                   "BINOM.DIST requires 4 arguments",
		"=BINOM.DIST(\"\",100,0.5,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.DIST(10,\"\",0.5,FALSE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.DIST(10,100,\"\",FALSE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.DIST(10,100,0.5,\"\")":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=BINOM.DIST(-1,100,0.5,FALSE)":   "#NUM!",
		"=BINOM.DIST(110,100,0.5,FALSE)":  "#NUM!",
		"=BINOM.DIST(10,100,-1,FALSE)":    "#NUM!",
		"=BINOM.DIST(10,100,2,FALSE)":     "#NUM!",
		// BINOM.DIST.RANGE
		"=BINOM.DIST.RANGE()":                "BINOM.DIST.RANGE requires at least 3 arguments",
		"=BINOM.DIST.RANGE(100,0.5,0,40,0)":  "BINOM.DIST.RANGE requires at most 4 arguments",
		"=BINOM.DIST.RANGE(\"\",0.5,0,40)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.DIST.RANGE(100,\"\",0,40)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.DIST.RANGE(100,0.5,\"\",40)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.DIST.RANGE(100,0.5,0,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.DIST.RANGE(100,-1,0,40)":     "#NUM!",
		"=BINOM.DIST.RANGE(100,2,0,40)":      "#NUM!",
		"=BINOM.DIST.RANGE(100,0.5,-1,40)":   "#NUM!",
		"=BINOM.DIST.RANGE(100,0.5,110,40)":  "#NUM!",
		"=BINOM.DIST.RANGE(100,0.5,0,-1)":    "#NUM!",
		"=BINOM.DIST.RANGE(100,0.5,0,110)":   "#NUM!",
		// BINOM.INV
		"=BINOM.INV()":             "BINOM.INV requires 3 numeric arguments",
		"=BINOM.INV(\"\",0.5,20%)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.INV(100,\"\",20%)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.INV(100,0.5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=BINOM.INV(-1,0.5,20%)":   "#NUM!",
		"=BINOM.INV(100,-1,20%)":   "#NUM!",
		"=BINOM.INV(100,2,20%)":    "#NUM!",
		"=BINOM.INV(100,0.5,-1)":   "#NUM!",
		"=BINOM.INV(100,0.5,2)":    "#NUM!",
		"=BINOM.INV(1,1,20%)":      "#NUM!",
		// CHIDIST
		"=CHIDIST()":         "CHIDIST requires 2 numeric arguments",
		"=CHIDIST(\"\",3)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHIDIST(0.5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// CHIINV
		"=CHIINV()":         "CHIINV requires 2 numeric arguments",
		"=CHIINV(\"\",1)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHIINV(0.5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHIINV(0,1)":      "#NUM!",
		"=CHIINV(2,1)":      "#NUM!",
		"=CHIINV(0.5,0.5)":  "#NUM!",
		// CHISQ.DIST
		"=CHISQ.DIST()":            "CHISQ.DIST requires 3 arguments",
		"=CHISQ.DIST(\"\",2,TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHISQ.DIST(3,\"\",TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHISQ.DIST(3,2,\"\")":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=CHISQ.DIST(-1,2,TRUE)":   "#NUM!",
		"=CHISQ.DIST(3,0,TRUE)":    "#NUM!",
		// CHISQ.DIST.RT
		"=CHISQ.DIST.RT()":         "CHISQ.DIST.RT requires 2 numeric arguments",
		"=CHISQ.DIST.RT(\"\",3)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHISQ.DIST.RT(0.5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// CHISQ.INV
		"=CHISQ.INV()":                "CHISQ.INV requires 2 numeric arguments",
		"=CHISQ.INV(\"\",1)":          "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHISQ.INV(0.5,\"\")":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHISQ.INV(-1,1)":            "#NUM!",
		"=CHISQ.INV(1,1)":             "#NUM!",
		"=CHISQ.INV(0.5,0.5)":         "#NUM!",
		"=CHISQ.INV(0.5,10000000001)": "#NUM!",
		// CHISQ.INV.RT
		"=CHISQ.INV.RT()":         "CHISQ.INV.RT requires 2 numeric arguments",
		"=CHISQ.INV.RT(\"\",1)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHISQ.INV.RT(0.5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CHISQ.INV.RT(0,1)":      "#NUM!",
		"=CHISQ.INV.RT(2,1)":      "#NUM!",
		"=CHISQ.INV.RT(0.5,0.5)":  "#NUM!",
		// CONFIDENCE
		"=CONFIDENCE()":               "CONFIDENCE requires 3 numeric arguments",
		"=CONFIDENCE(\"\",0.07,100)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CONFIDENCE(0.05,\"\",100)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CONFIDENCE(0.05,0.07,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CONFIDENCE(0,0.07,100)":     "#NUM!",
		"=CONFIDENCE(1,0.07,100)":     "#NUM!",
		"=CONFIDENCE(0.05,0,100)":     "#NUM!",
		"=CONFIDENCE(0.05,0.07,0.5)":  "#NUM!",
		// CONFIDENCE.NORM
		"=CONFIDENCE.NORM()":               "CONFIDENCE.NORM requires 3 numeric arguments",
		"=CONFIDENCE.NORM(\"\",0.07,100)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CONFIDENCE.NORM(0.05,\"\",100)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CONFIDENCE.NORM(0.05,0.07,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CONFIDENCE.NORM(0,0.07,100)":     "#NUM!",
		"=CONFIDENCE.NORM(1,0.07,100)":     "#NUM!",
		"=CONFIDENCE.NORM(0.05,0,100)":     "#NUM!",
		"=CONFIDENCE.NORM(0.05,0.07,0.5)":  "#NUM!",
		// CORREL
		"=CORREL()":            "CORREL requires 2 arguments",
		"=CORREL(A1:A3,B1:B5)": "#N/A",
		"=CORREL(A1:A1,B1:B1)": "#DIV/0!",
		// COUNTBLANK
		"=COUNTBLANK()":    "COUNTBLANK requires 1 argument",
		"=COUNTBLANK(1,2)": "COUNTBLANK requires 1 argument",
		// COUNTIF
		"=COUNTIF()": "COUNTIF requires 2 arguments",
		// COUNTIFS
		"=COUNTIFS()":              "COUNTIFS requires at least 2 arguments",
		"=COUNTIFS(A1:A9,2,D1:D9)": "#N/A",
		// CRITBINOM
		"=CRITBINOM()":             "CRITBINOM requires 3 numeric arguments",
		"=CRITBINOM(\"\",0.5,20%)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CRITBINOM(100,\"\",20%)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CRITBINOM(100,0.5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CRITBINOM(-1,0.5,20%)":   "#NUM!",
		"=CRITBINOM(100,-1,20%)":   "#NUM!",
		"=CRITBINOM(100,2,20%)":    "#NUM!",
		"=CRITBINOM(100,0.5,-1)":   "#NUM!",
		"=CRITBINOM(100,0.5,2)":    "#NUM!",
		"=CRITBINOM(1,1,20%)":      "#NUM!",
		// DEVSQ
		"=DEVSQ()":      "DEVSQ requires at least 1 numeric argument",
		"=DEVSQ(D1:D2)": "#N/A",
		// FISHER
		"=FISHER()":         "FISHER requires 1 numeric argument",
		"=FISHER(2)":        "#N/A",
		"=FISHER(INT(-2)))": "#N/A",
		"=FISHER(F1)":       "FISHER requires 1 numeric argument",
		// FISHERINV
		"=FISHERINV()":   "FISHERINV requires 1 numeric argument",
		"=FISHERINV(F1)": "FISHERINV requires 1 numeric argument",
		// GAMMA
		"=GAMMA()":       "GAMMA requires 1 numeric argument",
		"=GAMMA(F1)":     "GAMMA requires 1 numeric argument",
		"=GAMMA(0)":      "#N/A",
		"=GAMMA(INT(0))": "#N/A",
		// GAMMA.DIST
		"=GAMMA.DIST()":               "GAMMA.DIST requires 4 arguments",
		"=GAMMA.DIST(\"\",3,2,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMA.DIST(6,\"\",2,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMA.DIST(6,3,\"\",FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMA.DIST(6,3,2,\"\")":     "strconv.ParseBool: parsing \"\": invalid syntax",
		"=GAMMA.DIST(-1,3,2,FALSE)":   "#NUM!",
		"=GAMMA.DIST(6,0,2,FALSE)":    "#NUM!",
		"=GAMMA.DIST(6,3,0,FALSE)":    "#NUM!",
		// GAMMADIST
		"=GAMMADIST()":               "GAMMADIST requires 4 arguments",
		"=GAMMADIST(\"\",3,2,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMADIST(6,\"\",2,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMADIST(6,3,\"\",FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMADIST(6,3,2,\"\")":     "strconv.ParseBool: parsing \"\": invalid syntax",
		"=GAMMADIST(-1,3,2,FALSE)":   "#NUM!",
		"=GAMMADIST(6,0,2,FALSE)":    "#NUM!",
		"=GAMMADIST(6,3,0,FALSE)":    "#NUM!",
		// GAMMA.INV
		"=GAMMA.INV()":           "GAMMA.INV requires 3 arguments",
		"=GAMMA.INV(\"\",3,2)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMA.INV(0.5,\"\",2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMA.INV(0.5,3,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMA.INV(-1,3,2)":     "#NUM!",
		"=GAMMA.INV(2,3,2)":      "#NUM!",
		"=GAMMA.INV(0.5,0,2)":    "#NUM!",
		"=GAMMA.INV(0.5,3,0)":    "#NUM!",
		// GAMMAINV
		"=GAMMAINV()":           "GAMMAINV requires 3 arguments",
		"=GAMMAINV(\"\",3,2)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMAINV(0.5,\"\",2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMAINV(0.5,3,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMAINV(-1,3,2)":     "#NUM!",
		"=GAMMAINV(2,3,2)":      "#NUM!",
		"=GAMMAINV(0.5,0,2)":    "#NUM!",
		"=GAMMAINV(0.5,3,0)":    "#NUM!",
		// GAMMALN
		"=GAMMALN()":       "GAMMALN requires 1 numeric argument",
		"=GAMMALN(F1)":     "GAMMALN requires 1 numeric argument",
		"=GAMMALN(0)":      "#N/A",
		"=GAMMALN(INT(0))": "#N/A",
		// GAMMALN.PRECISE
		"=GAMMALN.PRECISE()":     "GAMMALN.PRECISE requires 1 numeric argument",
		"=GAMMALN.PRECISE(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=GAMMALN.PRECISE(0)":    "#NUM!",
		// GAUSS
		"=GAUSS()":     "GAUSS requires 1 numeric argument",
		"=GAUSS(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// GEOMEAN
		"=GEOMEAN()":      "GEOMEAN requires at least 1 numeric argument",
		"=GEOMEAN(0)":     "#NUM!",
		"=GEOMEAN(D1:D2)": "strconv.ParseFloat: parsing \"Month\": invalid syntax",
		// HARMEAN
		"=HARMEAN()":   "HARMEAN requires at least 1 argument",
		"=HARMEAN(-1)": "#N/A",
		"=HARMEAN(0)":  "#N/A",
		// HYPGEOM.DIST
		"=HYPGEOM.DIST()":                  "HYPGEOM.DIST requires 5 arguments",
		"=HYPGEOM.DIST(\"\",4,4,12,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOM.DIST(1,\"\",4,12,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOM.DIST(1,4,\"\",12,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOM.DIST(1,4,4,\"\",FALSE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOM.DIST(1,4,4,12,\"\")":     "strconv.ParseBool: parsing \"\": invalid syntax",
		"=HYPGEOM.DIST(-1,4,4,12,FALSE)":   "#NUM!",
		"=HYPGEOM.DIST(2,1,4,12,FALSE)":    "#NUM!",
		"=HYPGEOM.DIST(2,4,1,12,FALSE)":    "#NUM!",
		"=HYPGEOM.DIST(2,2,2,1,FALSE)":     "#NUM!",
		"=HYPGEOM.DIST(1,0,4,12,FALSE)":    "#NUM!",
		"=HYPGEOM.DIST(1,4,4,2,FALSE)":     "#NUM!",
		"=HYPGEOM.DIST(1,4,0,12,FALSE)":    "#NUM!",
		"=HYPGEOM.DIST(1,4,4,0,FALSE)":     "#NUM!",
		// HYPGEOMDIST
		"=HYPGEOMDIST()":            "HYPGEOMDIST requires 4 numeric arguments",
		"=HYPGEOMDIST(\"\",4,4,12)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOMDIST(1,\"\",4,12)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOMDIST(1,4,\"\",12)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOMDIST(1,4,4,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=HYPGEOMDIST(-1,4,4,12)":   "#NUM!",
		"=HYPGEOMDIST(2,1,4,12)":    "#NUM!",
		"=HYPGEOMDIST(2,4,1,12)":    "#NUM!",
		"=HYPGEOMDIST(2,2,2,1)":     "#NUM!",
		"=HYPGEOMDIST(1,0,4,12)":    "#NUM!",
		"=HYPGEOMDIST(1,4,4,2)":     "#NUM!",
		"=HYPGEOMDIST(1,4,0,12)":    "#NUM!",
		"=HYPGEOMDIST(1,4,4,0)":     "#NUM!",
		// KURT
		"=KURT()":          "KURT requires at least 1 argument",
		"=KURT(F1,INT(1))": "#DIV/0!",
		// EXPON.DIST
		"=EXPON.DIST()":            "EXPON.DIST requires 3 arguments",
		"=EXPON.DIST(\"\",1,TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=EXPON.DIST(0,\"\",TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=EXPON.DIST(0,1,\"\")":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=EXPON.DIST(-1,1,TRUE)":   "#NUM!",
		"=EXPON.DIST(1,0,TRUE)":    "#NUM!",
		// EXPONDIST
		"=EXPONDIST()":            "EXPONDIST requires 3 arguments",
		"=EXPONDIST(\"\",1,TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=EXPONDIST(0,\"\",TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=EXPONDIST(0,1,\"\")":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=EXPONDIST(-1,1,TRUE)":   "#NUM!",
		"=EXPONDIST(1,0,TRUE)":    "#NUM!",
		// FDIST
		"=FDIST()":                "FDIST requires 3 arguments",
		"=FDIST(\"\",1,2)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FDIST(5,\"\",2)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FDIST(5,1,\"\")":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FDIST(-1,1,2)":          "#NUM!",
		"=FDIST(5,0,2)":           "#NUM!",
		"=FDIST(5,10000000000,2)": "#NUM!",
		"=FDIST(5,1,0)":           "#NUM!",
		"=FDIST(5,1,10000000000)": "#NUM!",
		// F.DIST
		"=F.DIST()":                     "F.DIST requires 4 arguments",
		"=F.DIST(\"\",2,5,TRUE)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.DIST(1,\"\",5,TRUE)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.DIST(1,2,\"\",TRUE)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.DIST(1,2,5,\"\")":           "strconv.ParseBool: parsing \"\": invalid syntax",
		"=F.DIST(-1,1,2,TRUE)":          "#NUM!",
		"=F.DIST(5,0,2,TRUE)":           "#NUM!",
		"=F.DIST(5,10000000000,2,TRUE)": "#NUM!",
		"=F.DIST(5,1,0,TRUE)":           "#NUM!",
		"=F.DIST(5,1,10000000000,TRUE)": "#NUM!",
		// F.DIST.RT
		"=F.DIST.RT()":                "F.DIST.RT requires 3 arguments",
		"=F.DIST.RT(\"\",1,2)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.DIST.RT(5,\"\",2)":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.DIST.RT(5,1,\"\")":        "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.DIST.RT(-1,1,2)":          "#NUM!",
		"=F.DIST.RT(5,0,2)":           "#NUM!",
		"=F.DIST.RT(5,10000000000,2)": "#NUM!",
		"=F.DIST.RT(5,1,0)":           "#NUM!",
		"=F.DIST.RT(5,1,10000000000)": "#NUM!",
		// F.INV
		"=F.INV()":           "F.INV requires 3 arguments",
		"=F.INV(\"\",1,2)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.INV(0.2,\"\",2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.INV(0.2,1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.INV(0,1,2)":      "#NUM!",
		"=F.INV(0.2,0.5,2)":  "#NUM!",
		"=F.INV(0.2,1,0.5)":  "#NUM!",
		// FINV
		"=FINV()":           "FINV requires 3 arguments",
		"=FINV(\"\",1,2)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FINV(0.2,\"\",2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FINV(0.2,1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FINV(0,1,2)":      "#NUM!",
		"=FINV(0.2,0.5,2)":  "#NUM!",
		"=FINV(0.2,1,0.5)":  "#NUM!",
		// F.INV.RT
		"=F.INV.RT()":           "F.INV.RT requires 3 arguments",
		"=F.INV.RT(\"\",1,2)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.INV.RT(0.2,\"\",2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.INV.RT(0.2,1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=F.INV.RT(0,1,2)":      "#NUM!",
		"=F.INV.RT(0.2,0.5,2)":  "#NUM!",
		"=F.INV.RT(0.2,1,0.5)":  "#NUM!",
		// LOGINV
		"=LOGINV()":             "LOGINV requires 3 arguments",
		"=LOGINV(\"\",2,0.2)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGINV(0.3,\"\",0.2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGINV(0.3,2,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGINV(0,2,0.2)":      "#NUM!",
		"=LOGINV(1,2,0.2)":      "#NUM!",
		"=LOGINV(0.3,2,0)":      "#NUM!",
		// LOGNORM.INV
		"=LOGNORM.INV()":             "LOGNORM.INV requires 3 arguments",
		"=LOGNORM.INV(\"\",2,0.2)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORM.INV(0.3,\"\",0.2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORM.INV(0.3,2,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORM.INV(0,2,0.2)":      "#NUM!",
		"=LOGNORM.INV(1,2,0.2)":      "#NUM!",
		"=LOGNORM.INV(0.3,2,0)":      "#NUM!",
		// LOGNORM.DIST
		"=LOGNORM.DIST()":                  "LOGNORM.DIST requires 4 arguments",
		"=LOGNORM.DIST(\"\",10,5,FALSE)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORM.DIST(0.5,\"\",5,FALSE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORM.DIST(0.5,10,\"\",FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORM.DIST(0.5,10,5,\"\")":     "strconv.ParseBool: parsing \"\": invalid syntax",
		"=LOGNORM.DIST(0,10,5,FALSE)":      "#NUM!",
		"=LOGNORM.DIST(0.5,10,0,FALSE)":    "#NUM!",
		// LOGNORMDIST
		"=LOGNORMDIST()":           "LOGNORMDIST requires 3 arguments",
		"=LOGNORMDIST(\"\",10,5)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORMDIST(12,\"\",5)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORMDIST(12,10,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LOGNORMDIST(0,2,5)":      "#NUM!",
		"=LOGNORMDIST(12,10,0)":    "#NUM!",
		// NEGBINOM.DIST
		"=NEGBINOM.DIST()":                 "NEGBINOM.DIST requires 4 arguments",
		"=NEGBINOM.DIST(\"\",12,0.5,TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NEGBINOM.DIST(6,\"\",0.5,TRUE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NEGBINOM.DIST(6,12,\"\",TRUE)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NEGBINOM.DIST(6,12,0.5,\"\")":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=NEGBINOM.DIST(-1,12,0.5,TRUE)":   "#NUM!",
		"=NEGBINOM.DIST(6,0,0.5,TRUE)":     "#NUM!",
		"=NEGBINOM.DIST(6,12,-1,TRUE)":     "#NUM!",
		"=NEGBINOM.DIST(6,12,2,TRUE)":      "#NUM!",
		// NEGBINOMDIST
		"=NEGBINOMDIST()":            "NEGBINOMDIST requires 3 arguments",
		"=NEGBINOMDIST(\"\",12,0.5)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NEGBINOMDIST(6,\"\",0.5)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NEGBINOMDIST(6,12,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NEGBINOMDIST(-1,12,0.5)":   "#NUM!",
		"=NEGBINOMDIST(6,0,0.5)":     "#NUM!",
		"=NEGBINOMDIST(6,12,-1)":     "#NUM!",
		"=NEGBINOMDIST(6,12,2)":      "#NUM!",
		// NORM.DIST
		"=NORM.DIST()": "NORM.DIST requires 4 arguments",
		// NORMDIST
		"=NORMDIST()":               "NORMDIST requires 4 arguments",
		"=NORMDIST(\"\",0,0,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NORMDIST(0,\"\",0,FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NORMDIST(0,0,\"\",FALSE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NORMDIST(0,0,0,\"\")":     "strconv.ParseBool: parsing \"\": invalid syntax",
		"=NORMDIST(0,0,-1,TRUE)":    "#N/A",
		// NORM.INV
		"=NORM.INV()": "NORM.INV requires 3 arguments",
		// NORMINV
		"=NORMINV()":         "NORMINV requires 3 arguments",
		"=NORMINV(\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NORMINV(0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NORMINV(0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NORMINV(0,0,-1)":   "#N/A",
		"=NORMINV(-1,0,0)":   "#N/A",
		"=NORMINV(0,0,0)":    "#NUM!",
		// NORM.S.DIST
		"=NORM.S.DIST()": "NORM.S.DIST requires 2 numeric arguments",
		// NORMSDIST
		"=NORMSDIST()": "NORMSDIST requires 1 numeric argument",
		// NORM.S.INV
		"=NORM.S.INV()": "NORM.S.INV requires 1 numeric argument",
		// NORMSINV
		"=NORMSINV()": "NORMSINV requires 1 numeric argument",
		// LARGE
		"=LARGE()":           "LARGE requires 2 arguments",
		"=LARGE(A1:A5,0)":    "k should be > 0",
		"=LARGE(A1:A5,6)":    "k should be <= length of array",
		"=LARGE(A1:A5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// MAX
		"=MAX()":     "MAX requires at least 1 argument",
		"=MAX(NA())": "#N/A",
		// MAXA
		"=MAXA()":     "MAXA requires at least 1 argument",
		"=MAXA(NA())": "#N/A",
		// MAXIFS
		"=MAXIFS()":                         "MAXIFS requires at least 3 arguments",
		"=MAXIFS(F2:F4,A2:A4,\">0\",D2:D9)": "#N/A",
		// MEDIAN
		"=MEDIAN()":      "MEDIAN requires at least 1 argument",
		"=MEDIAN(\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MEDIAN(D1:D2)": "strconv.ParseFloat: parsing \"Month\": invalid syntax",
		// MIN
		"=MIN()":     "MIN requires at least 1 argument",
		"=MIN(NA())": "#N/A",
		// MINA
		"=MINA()":     "MINA requires at least 1 argument",
		"=MINA(NA())": "#N/A",
		// MINIFS
		"=MINIFS()":                         "MINIFS requires at least 3 arguments",
		"=MINIFS(F2:F4,A2:A4,\"<0\",D2:D9)": "#N/A",
		// PERCENTILE.EXC
		"=PERCENTILE.EXC()":           "PERCENTILE.EXC requires 2 arguments",
		"=PERCENTILE.EXC(A1:A4,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTILE.EXC(A1:A4,-1)":   "#NUM!",
		"=PERCENTILE.EXC(A1:A4,0)":    "#NUM!",
		"=PERCENTILE.EXC(A1:A4,1)":    "#NUM!",
		"=PERCENTILE.EXC(NA(),0.5)":   "#NUM!",
		// PERCENTILE.INC
		"=PERCENTILE.INC()": "PERCENTILE.INC requires 2 arguments",
		// PERCENTILE
		"=PERCENTILE()":       "PERCENTILE requires 2 arguments",
		"=PERCENTILE(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTILE(0,-1)":   "#N/A",
		"=PERCENTILE(NA(),1)": "#N/A",
		// PERCENTRANK.EXC
		"=PERCENTRANK.EXC()":             "PERCENTRANK.EXC requires 2 or 3 arguments",
		"=PERCENTRANK.EXC(A1:B4,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTRANK.EXC(A1:B4,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTRANK.EXC(A1:B4,0,0)":    "PERCENTRANK.EXC arguments significance should be > 1",
		"=PERCENTRANK.EXC(A1:B4,6)":      "#N/A",
		"=PERCENTRANK.EXC(NA(),1)":       "#N/A",
		// PERCENTRANK.INC
		"=PERCENTRANK.INC()":             "PERCENTRANK.INC requires 2 or 3 arguments",
		"=PERCENTRANK.INC(A1:B4,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTRANK.INC(A1:B4,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTRANK.INC(A1:B4,0,0)":    "PERCENTRANK.INC arguments significance should be > 1",
		"=PERCENTRANK.INC(A1:B4,6)":      "#N/A",
		"=PERCENTRANK.INC(NA(),1)":       "#N/A",
		// PERCENTRANK
		"=PERCENTRANK()":             "PERCENTRANK requires 2 or 3 arguments",
		"=PERCENTRANK(A1:B4,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTRANK(A1:B4,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERCENTRANK(A1:B4,0,0)":    "PERCENTRANK arguments significance should be > 1",
		"=PERCENTRANK(A1:B4,6)":      "#N/A",
		"=PERCENTRANK(NA(),1)":       "#N/A",
		// PERMUT
		"=PERMUT()":       "PERMUT requires 2 numeric arguments",
		"=PERMUT(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERMUT(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERMUT(6,8)":    "#N/A",
		// PERMUTATIONA
		"=PERMUTATIONA()":       "PERMUTATIONA requires 2 numeric arguments",
		"=PERMUTATIONA(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERMUTATIONA(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PERMUTATIONA(-1,0)":   "#N/A",
		"=PERMUTATIONA(0,-1)":   "#N/A",
		// QUARTILE
		"=QUARTILE()":           "QUARTILE requires 2 arguments",
		"=QUARTILE(A1:A4,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=QUARTILE(A1:A4,-1)":   "#NUM!",
		"=QUARTILE(A1:A4,5)":    "#NUM!",
		// QUARTILE.EXC
		"=QUARTILE.EXC()":           "QUARTILE.EXC requires 2 arguments",
		"=QUARTILE.EXC(A1:A4,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=QUARTILE.EXC(A1:A4,0)":    "#NUM!",
		"=QUARTILE.EXC(A1:A4,4)":    "#NUM!",
		// QUARTILE.INC
		"=QUARTILE.INC()": "QUARTILE.INC requires 2 arguments",
		// RANK
		"=RANK()":             "RANK requires at least 2 arguments",
		"=RANK(1,A1:B5,0,0)":  "RANK requires at most 3 arguments",
		"=RANK(-1,A1:B5)":     "#N/A",
		"=RANK(\"\",A1:B5)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RANK(1,A1:B5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// RANK.EQ
		"=RANK.EQ()":             "RANK.EQ requires at least 2 arguments",
		"=RANK.EQ(1,A1:B5,0,0)":  "RANK.EQ requires at most 3 arguments",
		"=RANK.EQ(-1,A1:B5)":     "#N/A",
		"=RANK.EQ(\"\",A1:B5)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RANK.EQ(1,A1:B5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// SKEW
		"=SKEW()":     "SKEW requires at least 1 argument",
		"=SKEW(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=SKEW(0)":    "#DIV/0!",
		// SMALL
		"=SMALL()":           "SMALL requires 2 arguments",
		"=SMALL(A1:A5,0)":    "k should be > 0",
		"=SMALL(A1:A5,6)":    "k should be <= length of array",
		"=SMALL(A1:A5,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// STANDARDIZE
		"=STANDARDIZE()":         "STANDARDIZE requires 3 arguments",
		"=STANDARDIZE(\"\",0,5)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=STANDARDIZE(0,\"\",5)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=STANDARDIZE(0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=STANDARDIZE(0,0,0)":    "#N/A",
		// STDEVP
		"=STDEVP()":     "STDEVP requires at least 1 argument",
		"=STDEVP(\"\")": "#DIV/0!",
		// STDEV.P
		"=STDEV.P()":     "STDEV.P requires at least 1 argument",
		"=STDEV.P(\"\")": "#DIV/0!",
		// T.DIST
		"=T.DIST()":             "T.DIST requires 3 arguments",
		"=T.DIST(\"\",10,TRUE)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=T.DIST(1,\"\",TRUE)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=T.DIST(1,10,\"\")":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=T.DIST(1,0,TRUE)":     "#NUM!",
		"=T.DIST(1,-1,FALSE)":   "#NUM!",
		"=T.DIST(1,0,FALSE)":    "#DIV/0!",
		// TDIST
		"=TDIST()":          "TDIST requires 3 arguments",
		"=TDIST(\"\",10,1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=TDIST(1,\"\",1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=TDIST(1,10,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=TDIST(-1,10,1)":   "#NUM!",
		"=TDIST(1,0,1)":     "#NUM!",
		"=TDIST(1,10,0)":    "#NUM!",
		// TRIMMEAN
		"=TRIMMEAN()":        "TRIMMEAN requires 2 arguments",
		"=TRIMMEAN(A1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=TRIMMEAN(A1,1)":    "#NUM!",
		"=TRIMMEAN(A1,-1)":   "#NUM!",
		// VAR
		"=VAR()": "VAR requires at least 1 argument",
		// VARA
		"=VARA()": "VARA requires at least 1 argument",
		// VARP
		"=VARP()":     "VARP requires at least 1 argument",
		"=VARP(\"\")": "#DIV/0!",
		// VAR.P
		"=VAR.P()":     "VAR.P requires at least 1 argument",
		"=VAR.P(\"\")": "#DIV/0!",
		// VAR.S
		"=VAR.S()": "VAR.S requires at least 1 argument",
		// VARPA
		"=VARPA()": "VARPA requires at least 1 argument",
		// WEIBULL
		"=WEIBULL()":               "WEIBULL requires 4 arguments",
		"=WEIBULL(\"\",1,1,FALSE)": "#VALUE!",
		"=WEIBULL(1,0,1,FALSE)":    "#N/A",
		"=WEIBULL(1,1,-1,FALSE)":   "#N/A",
		// WEIBULL.DIST
		"=WEIBULL.DIST()":               "WEIBULL.DIST requires 4 arguments",
		"=WEIBULL.DIST(\"\",1,1,FALSE)": "#VALUE!",
		"=WEIBULL.DIST(1,0,1,FALSE)":    "#N/A",
		"=WEIBULL.DIST(1,1,-1,FALSE)":   "#N/A",
		// Z.TEST
		"=Z.TEST(A1)":        "Z.TEST requires at least 2 arguments",
		"=Z.TEST(A1,0,0,0)":  "Z.TEST accepts at most 3 arguments",
		"=Z.TEST(H1,0)":      "#N/A",
		"=Z.TEST(A1,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=Z.TEST(A1,1)":      "#DIV/0!",
		"=Z.TEST(A1,1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// ZTEST
		"=ZTEST(A1)":        "ZTEST requires at least 2 arguments",
		"=ZTEST(A1,0,0,0)":  "ZTEST accepts at most 3 arguments",
		"=ZTEST(H1,0)":      "#N/A",
		"=ZTEST(A1,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ZTEST(A1,1)":      "#DIV/0!",
		"=ZTEST(A1,1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// Information Functions
		// ERROR.TYPE
		"=ERROR.TYPE()":  "ERROR.TYPE requires 1 argument",
		"=ERROR.TYPE(1)": "#N/A",
		// ISBLANK
		"=ISBLANK(A1,A2)": "ISBLANK requires 1 argument",
		// ISERR
		"=ISERR()": "ISERR requires 1 argument",
		// ISERROR
		"=ISERROR()": "ISERROR requires 1 argument",
		// ISEVEN
		"=ISEVEN()":       "ISEVEN requires 1 argument",
		`=ISEVEN("text")`: "strconv.Atoi: parsing \"text\": invalid syntax",
		// ISFORMULA
		"=ISFORMULA()": "ISFORMULA requires 1 argument",
		// ISLOGICAL
		"=ISLOGICAL()": "ISLOGICAL requires 1 argument",
		// ISNA
		"=ISNA()": "ISNA requires 1 argument",
		// ISNONTEXT
		"=ISNONTEXT()": "ISNONTEXT requires 1 argument",
		// ISNUMBER
		"=ISNUMBER()": "ISNUMBER requires 1 argument",
		// ISODD
		"=ISODD()":       "ISODD requires 1 argument",
		`=ISODD("text")`: "strconv.Atoi: parsing \"text\": invalid syntax",
		// ISREF
		"=ISREF()": "ISREF requires 1 argument",
		// ISTEXT
		"=ISTEXT()": "ISTEXT requires 1 argument",
		// N
		"=N()":     "N requires 1 argument",
		"=N(NA())": "#N/A",
		// NA
		"=NA()":  "#N/A",
		"=NA(1)": "NA accepts no arguments",
		// SHEET
		"=SHEET(\"\",\"\")":  "SHEET accepts at most 1 argument",
		"=SHEET(\"Sheet2\")": "#N/A",
		// SHEETS
		"=SHEETS(\"\",\"\")":  "SHEETS accepts at most 1 argument",
		"=SHEETS(\"Sheet1\")": "#N/A",
		// TYPE
		"=TYPE()": "TYPE requires 1 argument",
		// T
		"=T()":     "T requires 1 argument",
		"=T(NA())": "#N/A",
		// Logical Functions
		// AND
		`=AND("text")`: "strconv.ParseFloat: parsing \"text\": invalid syntax",
		`=AND(A1:B1)`:  "#VALUE!",
		"=AND()":       "AND requires at least 1 argument",
		"=AND(1" + strings.Repeat(",1", 30) + ")": "AND accepts at most 30 arguments",
		// FALSE
		"=FALSE(A1)": "FALSE takes no arguments",
		// IFERROR
		"=IFERROR()": "IFERROR requires 2 arguments",
		// IFNA
		"=IFNA()": "IFNA requires 2 arguments",
		// IFS
		"=IFS()":            "IFS requires at least 2 arguments",
		"=IFS(FALSE,FALSE)": "#N/A",
		// NOT
		"=NOT()":      "NOT requires 1 argument",
		"=NOT(NOT())": "NOT requires 1 argument",
		"=NOT(\"\")":  "NOT expects 1 boolean or numeric argument",
		// OR
		`=OR("text")`:                            "strconv.ParseFloat: parsing \"text\": invalid syntax",
		`=OR(A1:B1)`:                             "#VALUE!",
		"=OR()":                                  "OR requires at least 1 argument",
		"=OR(1" + strings.Repeat(",1", 30) + ")": "OR accepts at most 30 arguments",
		// SWITCH
		"=SWITCH()":      "SWITCH requires at least 3 arguments",
		"=SWITCH(0,1,2)": "#N/A",
		// TRUE
		"=TRUE(A1)": "TRUE takes no arguments",
		// XOR
		"=XOR()":              "XOR requires at least 1 argument",
		"=XOR(\"text\")":      "#VALUE!",
		"=XOR(XOR(\"text\"))": "#VALUE!",
		// Date and Time Functions
		// DATE
		"=DATE()":               "DATE requires 3 number arguments",
		`=DATE("text",10,21)`:   "DATE requires 3 number arguments",
		`=DATE(2020,"text",21)`: "DATE requires 3 number arguments",
		`=DATE(2020,10,"text")`: "DATE requires 3 number arguments",
		// DATEDIF
		"=DATEDIF()":                  "DATEDIF requires 3 number arguments",
		"=DATEDIF(\"\",\"\",\"\")":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DATEDIF(43891,43101,\"Y\")": "start_date > end_date",
		"=DATEDIF(43101,43891,\"x\")": "DATEDIF has invalid unit",
		// DATEVALUE
		"=DATEVALUE()":             "DATEVALUE requires 1 argument",
		"=DATEVALUE(\"01/01\")":    "#VALUE!", // valid in Excel, which uses years by the system date
		"=DATEVALUE(\"1900-0-0\")": "#VALUE!",
		// DAY
		"=DAY()":         "DAY requires exactly 1 argument",
		"=DAY(-1)":       "DAY only accepts positive argument",
		"=DAY(0,0)":      "DAY requires exactly 1 argument",
		"=DAY(\"text\")": "#VALUE!",
		"=DAY(\"January 25, 2020 9223372036854775808 AM\")":                   "#VALUE!",
		"=DAY(\"January 25, 2020 9223372036854775808:00 AM\")":                "#VALUE!",
		"=DAY(\"January 25, 2020 00:9223372036854775808 AM\")":                "#VALUE!",
		"=DAY(\"January 25, 2020 9223372036854775808:00.0 AM\")":              "#VALUE!",
		"=DAY(\"January 25, 2020 0:1" + strings.Repeat("0", 309) + ".0 AM\")": "#VALUE!",
		"=DAY(\"January 25, 2020 9223372036854775808:00:00 AM\")":             "#VALUE!",
		"=DAY(\"January 25, 2020 0:9223372036854775808:0 AM\")":               "#VALUE!",
		"=DAY(\"January 25, 2020 0:0:1" + strings.Repeat("0", 309) + " AM\")": "#VALUE!",
		"=DAY(\"January 25, 2020 0:61:0 AM\")":                                "#VALUE!",
		"=DAY(\"January 25, 2020 0:00:60 AM\")":                               "#VALUE!",
		"=DAY(\"January 25, 2020 24:00:00\")":                                 "#VALUE!",
		"=DAY(\"January 25, 2020 00:00:10001\")":                              "#VALUE!",
		"=DAY(\"9223372036854775808/25/2020\")":                               "#VALUE!",
		"=DAY(\"01/9223372036854775808/2020\")":                               "#VALUE!",
		"=DAY(\"01/25/9223372036854775808\")":                                 "#VALUE!",
		"=DAY(\"01/25/10000\")":                                               "#VALUE!",
		"=DAY(\"01/25/100\")":                                                 "#VALUE!",
		"=DAY(\"January 9223372036854775808, 2020\")":                         "#VALUE!",
		"=DAY(\"January 25, 9223372036854775808\")":                           "#VALUE!",
		"=DAY(\"January 25, 10000\")":                                         "#VALUE!",
		"=DAY(\"January 25, 100\")":                                           "#VALUE!",
		"=DAY(\"9223372036854775808-25-2020\")":                               "#VALUE!",
		"=DAY(\"01-9223372036854775808-2020\")":                               "#VALUE!",
		"=DAY(\"01-25-9223372036854775808\")":                                 "#VALUE!",
		"=DAY(\"1900-0-0\")":                                                  "#VALUE!",
		"=DAY(\"14-25-1900\")":                                                "#VALUE!",
		"=DAY(\"3-January-9223372036854775808\")":                             "#VALUE!",
		"=DAY(\"9223372036854775808-January-1900\")":                          "#VALUE!",
		"=DAY(\"0-January-1900\")":                                            "#VALUE!",
		// DAYS
		"=DAYS()":       "DAYS requires 2 arguments",
		"=DAYS(\"\",0)": "#VALUE!",
		"=DAYS(0,\"\")": "#VALUE!",
		"=DAYS(NA(),0)": "#VALUE!",
		"=DAYS(0,NA())": "#VALUE!",
		// HOUR
		"=HOUR()":             "HOUR requires exactly 1 argument",
		"=HOUR(-1)":           "HOUR only accepts positive argument",
		"=HOUR(\"\")":         "#VALUE!",
		"=HOUR(\"25:10:55\")": "#VALUE!",
		// ISOWEEKNUM
		"=ISOWEEKNUM()":                    "ISOWEEKNUM requires 1 argument",
		"=ISOWEEKNUM(\"\")":                "#VALUE!",
		"=ISOWEEKNUM(\"January 25, 100\")": "#VALUE!",
		"=ISOWEEKNUM(-1)":                  "#NUM!",
		// MINUTE
		"=MINUTE()":             "MINUTE requires exactly 1 argument",
		"=MINUTE(-1)":           "MINUTE only accepts positive argument",
		"=MINUTE(\"\")":         "#VALUE!",
		"=MINUTE(\"13:60:55\")": "#VALUE!",
		// MONTH
		"=MONTH()":                    "MONTH requires exactly 1 argument",
		"=MONTH(0,0)":                 "MONTH requires exactly 1 argument",
		"=MONTH(-1)":                  "MONTH only accepts positive argument",
		"=MONTH(\"text\")":            "#VALUE!",
		"=MONTH(\"January 25, 100\")": "#VALUE!",
		// YEAR
		"=YEAR()":                    "YEAR requires exactly 1 argument",
		"=YEAR(0,0)":                 "YEAR requires exactly 1 argument",
		"=YEAR(-1)":                  "YEAR only accepts positive argument",
		"=YEAR(\"text\")":            "#VALUE!",
		"=YEAR(\"January 25, 100\")": "#VALUE!",
		// YEARFRAC
		"=YEARFRAC()":                 "YEARFRAC requires 3 or 4 arguments",
		"=YEARFRAC(42005,42094,5)":    "invalid basis",
		"=YEARFRAC(\"\",42094,5)":     "#VALUE!",
		"=YEARFRAC(42005,\"\",5)":     "#VALUE!",
		"=YEARFRAC(42005,42094,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// NOW
		"=NOW(A1)": "NOW accepts no arguments",
		// SECOND
		"=SECOND()":          "SECOND requires exactly 1 argument",
		"=SECOND(-1)":        "SECOND only accepts positive argument",
		"=SECOND(\"\")":      "#VALUE!",
		"=SECOND(\"25:55\")": "#VALUE!",
		// TIME
		"=TIME()":         "TIME requires 3 number arguments",
		"=TIME(\"\",0,0)": "TIME requires 3 number arguments",
		"=TIME(0,0,-1)":   "#NUM!",
		// TIMEVALUE
		"=TIMEVALUE()":          "TIMEVALUE requires exactly 1 argument",
		"=TIMEVALUE(1)":         "#VALUE!",
		"=TIMEVALUE(-1)":        "#VALUE!",
		"=TIMEVALUE(\"25:55\")": "#VALUE!",
		// TODAY
		"=TODAY(A1)": "TODAY accepts no arguments",
		// WEEKDAY
		"=WEEKDAY()":                    "WEEKDAY requires at least 1 argument",
		"=WEEKDAY(0,1,0)":               "WEEKDAY allows at most 2 arguments",
		"=WEEKDAY(0,\"\")":              "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=WEEKDAY(\"\",1)":              "#VALUE!",
		"=WEEKDAY(0,0)":                 "#VALUE!",
		"=WEEKDAY(\"January 25, 100\")": "#VALUE!",
		"=WEEKDAY(-1,1)":                "#NUM!",
		// Text Functions
		// CHAR
		"=CHAR()":     "CHAR requires 1 argument",
		"=CHAR(-1)":   "#VALUE!",
		"=CHAR(256)":  "#VALUE!",
		"=CHAR(\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// CLEAN
		"=CLEAN()":    "CLEAN requires 1 argument",
		"=CLEAN(1,2)": "CLEAN requires 1 argument",
		// CODE
		"=CODE()":    "CODE requires 1 argument",
		"=CODE(1,2)": "CODE requires 1 argument",
		// CONCAT
		"=CONCAT(MUNIT(2))": "CONCAT requires arguments to be strings",
		// CONCATENATE
		"=CONCATENATE(MUNIT(2))": "CONCATENATE requires arguments to be strings",
		// EXACT
		"=EXACT()":      "EXACT requires 2 arguments",
		"=EXACT(1,2,3)": "EXACT requires 2 arguments",
		// FIXED
		"=FIXED()":         "FIXED requires at least 1 argument",
		"=FIXED(0,1,2,3)":  "FIXED allows at most 3 arguments",
		"=FIXED(\"\")":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FIXED(0,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FIXED(0,0,\"\")": "strconv.ParseBool: parsing \"\": invalid syntax",
		// FIND
		"=FIND()":                 "FIND requires at least 2 arguments",
		"=FIND(1,2,3,4)":          "FIND allows at most 3 arguments",
		"=FIND(\"x\",\"\")":       "#VALUE!",
		"=FIND(\"x\",\"x\",-1)":   "#VALUE!",
		"=FIND(\"x\",\"x\",\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// FINDB
		"=FINDB()":                 "FINDB requires at least 2 arguments",
		"=FINDB(1,2,3,4)":          "FINDB allows at most 3 arguments",
		"=FINDB(\"x\",\"\")":       "#VALUE!",
		"=FINDB(\"x\",\"x\",-1)":   "#VALUE!",
		"=FINDB(\"x\",\"x\",\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// LEFT
		"=LEFT()":          "LEFT requires at least 1 argument",
		"=LEFT(\"\",2,3)":  "LEFT allows at most 2 arguments",
		"=LEFT(\"\",\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LEFT(\"\",-1)":   "#VALUE!",
		// LEFTB
		"=LEFTB()":          "LEFTB requires at least 1 argument",
		"=LEFTB(\"\",2,3)":  "LEFTB allows at most 2 arguments",
		"=LEFTB(\"\",\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=LEFTB(\"\",-1)":   "#VALUE!",
		// LEN
		"=LEN()": "LEN requires 1 string argument",
		// LENB
		"=LENB()": "LENB requires 1 string argument",
		// LOWER
		"=LOWER()":    "LOWER requires 1 argument",
		"=LOWER(1,2)": "LOWER requires 1 argument",
		// MID
		"=MID()":            "MID requires 3 arguments",
		"=MID(\"\",-1,1)":   "#VALUE!",
		"=MID(\"\",\"\",1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MID(\"\",1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// MIDB
		"=MIDB()":            "MIDB requires 3 arguments",
		"=MIDB(\"\",-1,1)":   "#VALUE!",
		"=MIDB(\"\",\"\",1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MIDB(\"\",1,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// PROPER
		"=PROPER()":    "PROPER requires 1 argument",
		"=PROPER(1,2)": "PROPER requires 1 argument",
		// REPLACE
		"=REPLACE()":                           "REPLACE requires 4 arguments",
		"=REPLACE(\"text\",0,4,\"string\")":    "#VALUE!",
		"=REPLACE(\"text\",\"\",0,\"string\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=REPLACE(\"text\",1,\"\",\"string\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// REPLACEB
		"=REPLACEB()":                           "REPLACEB requires 4 arguments",
		"=REPLACEB(\"text\",0,4,\"string\")":    "#VALUE!",
		"=REPLACEB(\"text\",\"\",0,\"string\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=REPLACEB(\"text\",1,\"\",\"string\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// REPT
		"=REPT()":            "REPT requires 2 arguments",
		"=REPT(INT(0),2)":    "REPT requires first argument to be a string",
		"=REPT(\"*\",\"*\")": "REPT requires second argument to be a number",
		"=REPT(\"*\",-1)":    "REPT requires second argument to be >= 0",
		// RIGHT
		"=RIGHT()":          "RIGHT requires at least 1 argument",
		"=RIGHT(\"\",2,3)":  "RIGHT allows at most 2 arguments",
		"=RIGHT(\"\",\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RIGHT(\"\",-1)":   "#VALUE!",
		// RIGHTB
		"=RIGHTB()":          "RIGHTB requires at least 1 argument",
		"=RIGHTB(\"\",2,3)":  "RIGHTB allows at most 2 arguments",
		"=RIGHTB(\"\",\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RIGHTB(\"\",-1)":   "#VALUE!",
		// SUBSTITUTE
		"=SUBSTITUTE()":                    "SUBSTITUTE requires 3 or 4 arguments",
		"=SUBSTITUTE(\"\",\"\",\"\",\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=SUBSTITUTE(\"\",\"\",\"\",0)":    "instance_num should be > 0",
		// TEXTJOIN
		"=TEXTJOIN()":               "TEXTJOIN requires at least 3 arguments",
		"=TEXTJOIN(\"\",\"\",1)":    "strconv.ParseBool: parsing \"\": invalid syntax",
		"=TEXTJOIN(\"\",TRUE,NA())": "#N/A",
		"=TEXTJOIN(\"\",TRUE," + strings.Repeat("0,", 250) + ",0)": "TEXTJOIN accepts at most 252 arguments",
		"=TEXTJOIN(\",\",FALSE,REPT(\"*\",32768))":                 "TEXTJOIN function exceeds 32767 characters",
		// TRIM
		"=TRIM()":    "TRIM requires 1 argument",
		"=TRIM(1,2)": "TRIM requires 1 argument",
		// UNICHAR
		"=UNICHAR()":      "UNICHAR requires 1 argument",
		"=UNICHAR(\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=UNICHAR(55296)": "#VALUE!",
		"=UNICHAR(0)":     "#VALUE!",
		// UNICODE
		"=UNICODE()":     "UNICODE requires 1 argument",
		"=UNICODE(\"\")": "#VALUE!",
		// VALUE
		"=VALUE()":     "VALUE requires 1 argument",
		"=VALUE(\"\")": "#VALUE!",
		// UPPER
		"=UPPER()":    "UPPER requires 1 argument",
		"=UPPER(1,2)": "UPPER requires 1 argument",
		// Conditional Functions
		// IF
		"=IF()":        "IF requires at least 1 argument",
		"=IF(0,1,2,3)": "IF accepts at most 3 arguments",
		"=IF(D1,1,2)":  "strconv.ParseBool: parsing \"Month\": invalid syntax",
		// Excel Lookup and Reference Functions
		// ADDRESS
		"=ADDRESS()":                        "ADDRESS requires at least 2 arguments",
		"=ADDRESS(1,1,1,TRUE,\"Sheet1\",0)": "ADDRESS requires at most 5 arguments",
		"=ADDRESS(\"\",1,1,TRUE)":           "#VALUE!",
		"=ADDRESS(1,\"\",1,TRUE)":           "#VALUE!",
		"=ADDRESS(1,1,\"\",TRUE)":           "#VALUE!",
		"=ADDRESS(1,1,1,\"\")":              "#VALUE!",
		"=ADDRESS(1,1,0,TRUE)":              "#NUM!",
		"=ADDRESS(1,16385,2,TRUE)":          "#VALUE!",
		"=ADDRESS(1,16385,3,TRUE)":          "#VALUE!",
		"=ADDRESS(1048576,1,1,TRUE)":        "#VALUE!",
		// CHOOSE
		"=CHOOSE()":                "CHOOSE requires 2 arguments",
		"=CHOOSE(\"index_num\",0)": "CHOOSE requires first argument of type number",
		"=CHOOSE(2,0)":             "index_num should be <= to the number of values",
		"=CHOOSE(1,NA())":          "#N/A",
		// COLUMN
		"=COLUMN(1,2)":          "COLUMN requires at most 1 argument",
		"=COLUMN(\"\")":         "invalid reference",
		"=COLUMN(Sheet1)":       newInvalidColumnNameError("Sheet1").Error(),
		"=COLUMN(Sheet1!A1!B1)": newInvalidColumnNameError("Sheet1").Error(),
		// COLUMNS
		"=COLUMNS()":              "COLUMNS requires 1 argument",
		"=COLUMNS(1)":             "invalid reference",
		"=COLUMNS(\"\")":          "invalid reference",
		"=COLUMNS(Sheet1)":        newInvalidColumnNameError("Sheet1").Error(),
		"=COLUMNS(Sheet1!A1!B1)":  newInvalidColumnNameError("Sheet1").Error(),
		"=COLUMNS(Sheet1!Sheet1)": newInvalidColumnNameError("Sheet1").Error(),
		// FORMULATEXT
		"=FORMULATEXT()":  "FORMULATEXT requires 1 argument",
		"=FORMULATEXT(1)": "#VALUE!",
		// HLOOKUP
		"=HLOOKUP()":                     "HLOOKUP requires at least 3 arguments",
		"=HLOOKUP(D2,D1,1,FALSE)":        "HLOOKUP requires second argument of table array",
		"=HLOOKUP(D2,D:D,FALSE,FALSE)":   "HLOOKUP requires numeric row argument",
		"=HLOOKUP(D2,D:D,1,FALSE,FALSE)": "HLOOKUP requires at most 4 arguments",
		"=HLOOKUP(D2,D:D,1,2)":           "strconv.ParseBool: parsing \"2\": invalid syntax",
		"=HLOOKUP(D2,D10:D10,1,FALSE)":   "HLOOKUP no result found",
		"=HLOOKUP(D2,D2:D3,4,FALSE)":     "HLOOKUP has invalid row index",
		"=HLOOKUP(D2,C:C,1,FALSE)":       "HLOOKUP no result found",
		"=HLOOKUP(ISNUMBER(1),F3:F9,1)":  "HLOOKUP no result found",
		"=HLOOKUP(INT(1),E2:E9,1)":       "HLOOKUP no result found",
		"=HLOOKUP(MUNIT(2),MUNIT(3),1)":  "HLOOKUP no result found",
		"=HLOOKUP(A1:B2,B2:B3,1)":        "HLOOKUP no result found",
		// MATCH
		"=MATCH()":              "MATCH requires 1 or 2 arguments",
		"=MATCH(0,A1:A1,0,0)":   "MATCH requires 1 or 2 arguments",
		"=MATCH(0,A1:A1,\"x\")": "MATCH requires numeric match_type argument",
		"=MATCH(0,A1)":          "MATCH arguments lookup_array should be one-dimensional array",
		"=MATCH(0,A1:B1)":       "MATCH arguments lookup_array should be one-dimensional array",
		// TRANSPOSE
		"=TRANSPOSE()": "TRANSPOSE requires 1 argument",
		// VLOOKUP
		"=VLOOKUP()":                     "VLOOKUP requires at least 3 arguments",
		"=VLOOKUP(D2,D1,1,FALSE)":        "VLOOKUP requires second argument of table array",
		"=VLOOKUP(D2,D:D,FALSE,FALSE)":   "VLOOKUP requires numeric col argument",
		"=VLOOKUP(D2,D:D,1,FALSE,FALSE)": "VLOOKUP requires at most 4 arguments",
		"=VLOOKUP(D2,D:D,1,2)":           "strconv.ParseBool: parsing \"2\": invalid syntax",
		"=VLOOKUP(D2,D10:D10,1,FALSE)":   "VLOOKUP no result found",
		"=VLOOKUP(D2,D:D,2,FALSE)":       "VLOOKUP has invalid column index",
		"=VLOOKUP(D2,C:C,1,FALSE)":       "VLOOKUP no result found",
		"=VLOOKUP(ISNUMBER(1),F3:F9,1)":  "VLOOKUP no result found",
		"=VLOOKUP(INT(1),E2:E9,1)":       "VLOOKUP no result found",
		"=VLOOKUP(MUNIT(2),MUNIT(3),1)":  "VLOOKUP no result found",
		"=VLOOKUP(1,G1:H2,1,FALSE)":      "VLOOKUP no result found",
		// INDEX
		"=INDEX()":          "INDEX requires 2 or 3 arguments",
		"=INDEX(A1,2)":      "INDEX row_num out of range",
		"=INDEX(A1,0,2)":    "INDEX col_num out of range",
		"=INDEX(A1:A1,2)":   "INDEX row_num out of range",
		"=INDEX(A1:A1,0,2)": "INDEX col_num out of range",
		"=INDEX(A1:B2,2,3)": "INDEX col_num out of range",
		"=INDEX(A1:A2,0,0)": "#VALUE!",
		"=INDEX(0,\"\")":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=INDEX(0,0,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		// INDIRECT
		"=INDIRECT()":                     "INDIRECT requires 1 or 2 arguments",
		"=INDIRECT(\"E\"&1,TRUE,1)":       "INDIRECT requires 1 or 2 arguments",
		"=INDIRECT(\"R1048577C1\",\"\")":  "#VALUE!",
		"=INDIRECT(\"E1048577\")":         "#REF!",
		"=INDIRECT(\"R1048577C1\",FALSE)": "#REF!",
		"=INDIRECT(\"R1C16385\",FALSE)":   "#REF!",
		"=INDIRECT(\"\",FALSE)":           "#REF!",
		"=INDIRECT(\"R C1\",FALSE)":       "#REF!",
		"=INDIRECT(\"R1C \",FALSE)":       "#REF!",
		"=INDIRECT(\"R1C1:R2C \",FALSE)":  "#REF!",
		// LOOKUP
		"=LOOKUP()":                     "LOOKUP requires at least 2 arguments",
		"=LOOKUP(D2,D1,D2)":             "LOOKUP requires second argument of table array",
		"=LOOKUP(D2,D1,D2,FALSE)":       "LOOKUP requires at most 3 arguments",
		"=LOOKUP(1,MUNIT(0))":           "LOOKUP requires not empty range as second argument",
		"=LOOKUP(D1,MUNIT(1),MUNIT(1))": "LOOKUP no result found",
		// ROW
		"=ROW(1,2)":          "ROW requires at most 1 argument",
		"=ROW(\"\")":         "invalid reference",
		"=ROW(Sheet1)":       newInvalidColumnNameError("Sheet1").Error(),
		"=ROW(Sheet1!A1!B1)": newInvalidColumnNameError("Sheet1").Error(),
		// ROWS
		"=ROWS()":              "ROWS requires 1 argument",
		"=ROWS(1)":             "invalid reference",
		"=ROWS(\"\")":          "invalid reference",
		"=ROWS(Sheet1)":        newInvalidColumnNameError("Sheet1").Error(),
		"=ROWS(Sheet1!A1!B1)":  newInvalidColumnNameError("Sheet1").Error(),
		"=ROWS(Sheet1!Sheet1)": newInvalidColumnNameError("Sheet1").Error(),
		// Web Functions
		// ENCODEURL
		"=ENCODEURL()": "ENCODEURL requires 1 argument",
		// Financial Functions
		// ACCRINT
		"=ACCRINT()": "ACCRINT requires at least 6 arguments",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,4,1,FALSE,0)":  "ACCRINT allows at most 8 arguments",
		"=ACCRINT(\"\",\"04/01/2012\",\"12/31/2013\",8%,10000,4,1,FALSE)":              "#VALUE!",
		"=ACCRINT(\"01/01/2012\",\"\",\"12/31/2013\",8%,10000,4,1,FALSE)":              "#VALUE!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"\",8%,10000,4,1,FALSE)":              "#VALUE!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",\"\",10000,4,1,FALSE)":  "#NUM!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,\"\",4,1,FALSE)":     "#NUM!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,3)":            "#NUM!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,\"\",1,FALSE)": "#NUM!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,4,\"\",FALSE)": "#NUM!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,4,1,\"\")":     "#VALUE!",
		"=ACCRINT(\"01/01/2012\",\"04/01/2012\",\"12/31/2013\",8%,10000,4,5,FALSE)":    "invalid basis",
		// ACCRINTM
		"=ACCRINTM()": "ACCRINTM requires 4 or 5 arguments",
		"=ACCRINTM(\"\",\"01/01/2012\",8%,10000)":                "#VALUE!",
		"=ACCRINTM(\"01/01/2012\",\"\",8%,10000)":                "#VALUE!",
		"=ACCRINTM(\"12/31/2012\",\"01/01/2012\",8%,10000)":      "#NUM!",
		"=ACCRINTM(\"01/01/2012\",\"12/31/2012\",\"\",10000)":    "#NUM!",
		"=ACCRINTM(\"01/01/2012\",\"12/31/2012\",8%,\"\",10000)": "#NUM!",
		"=ACCRINTM(\"01/01/2012\",\"12/31/2012\",8%,-1,10000)":   "#NUM!",
		"=ACCRINTM(\"01/01/2012\",\"12/31/2012\",8%,10000,\"\")": "#NUM!",
		"=ACCRINTM(\"01/01/2012\",\"12/31/2012\",8%,10000,5)":    "invalid basis",
		// AMORDEGRC
		"=AMORDEGRC()": "AMORDEGRC requires 6 or 7 arguments",
		"=AMORDEGRC(\"\",\"01/01/2015\",\"09/30/2015\",20,1,20%)":     "AMORDEGRC requires cost to be number argument",
		"=AMORDEGRC(-1,\"01/01/2015\",\"09/30/2015\",20,1,20%)":       "AMORDEGRC requires cost >= 0",
		"=AMORDEGRC(150,\"\",\"09/30/2015\",20,1,20%)":                "#VALUE!",
		"=AMORDEGRC(150,\"01/01/2015\",\"\",20,1,20%)":                "#VALUE!",
		"=AMORDEGRC(150,\"09/30/2015\",\"01/01/2015\",20,1,20%)":      "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",\"\",1,20%)":    "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",-1,1,20%)":      "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,\"\",20%)":   "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,-1,20%)":     "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,\"\")":     "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,-1)":       "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,20%,\"\")": "#NUM!",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,50%)":      "AMORDEGRC requires rate to be < 0.5",
		"=AMORDEGRC(150,\"01/01/2015\",\"09/30/2015\",20,1,20%,5)":    "invalid basis",
		// AMORLINC
		"=AMORLINC()": "AMORLINC requires 6 or 7 arguments",
		"=AMORLINC(\"\",\"01/01/2015\",\"09/30/2015\",20,1,20%)":     "AMORLINC requires cost to be number argument",
		"=AMORLINC(-1,\"01/01/2015\",\"09/30/2015\",20,1,20%)":       "AMORLINC requires cost >= 0",
		"=AMORLINC(150,\"\",\"09/30/2015\",20,1,20%)":                "#VALUE!",
		"=AMORLINC(150,\"01/01/2015\",\"\",20,1,20%)":                "#VALUE!",
		"=AMORLINC(150,\"09/30/2015\",\"01/01/2015\",20,1,20%)":      "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",\"\",1,20%)":    "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",-1,1,20%)":      "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,\"\",20%)":   "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,-1,20%)":     "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,1,\"\")":     "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,1,-1)":       "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,1,20%,\"\")": "#NUM!",
		"=AMORLINC(150,\"01/01/2015\",\"09/30/2015\",20,1,20%,5)":    "invalid basis",
		// COUPDAYBS
		"=COUPDAYBS()":                                     "COUPDAYBS requires 3 or 4 arguments",
		"=COUPDAYBS(\"\",\"10/25/2012\",4)":                "#VALUE!",
		"=COUPDAYBS(\"01/01/2011\",\"\",4)":                "#VALUE!",
		"=COUPDAYBS(\"01/01/2011\",\"10/25/2012\",\"\")":   "#VALUE!",
		"=COUPDAYBS(\"01/01/2011\",\"10/25/2012\",4,\"\")": "#NUM!",
		"=COUPDAYBS(\"10/25/2012\",\"01/01/2011\",4)":      "COUPDAYBS requires maturity > settlement",
		// COUPDAYS
		"=COUPDAYS()":                                     "COUPDAYS requires 3 or 4 arguments",
		"=COUPDAYS(\"\",\"10/25/2012\",4)":                "#VALUE!",
		"=COUPDAYS(\"01/01/2011\",\"\",4)":                "#VALUE!",
		"=COUPDAYS(\"01/01/2011\",\"10/25/2012\",\"\")":   "#VALUE!",
		"=COUPDAYS(\"01/01/2011\",\"10/25/2012\",4,\"\")": "#NUM!",
		"=COUPDAYS(\"10/25/2012\",\"01/01/2011\",4)":      "COUPDAYS requires maturity > settlement",
		// COUPDAYSNC
		"=COUPDAYSNC()":                                     "COUPDAYSNC requires 3 or 4 arguments",
		"=COUPDAYSNC(\"\",\"10/25/2012\",4)":                "#VALUE!",
		"=COUPDAYSNC(\"01/01/2011\",\"\",4)":                "#VALUE!",
		"=COUPDAYSNC(\"01/01/2011\",\"10/25/2012\",\"\")":   "#VALUE!",
		"=COUPDAYSNC(\"01/01/2011\",\"10/25/2012\",4,\"\")": "#NUM!",
		"=COUPDAYSNC(\"10/25/2012\",\"01/01/2011\",4)":      "COUPDAYSNC requires maturity > settlement",
		// COUPNCD
		"=COUPNCD()": "COUPNCD requires 3 or 4 arguments",
		"=COUPNCD(\"01/01/2011\",\"10/25/2012\",4,0,0)":  "COUPNCD requires 3 or 4 arguments",
		"=COUPNCD(\"\",\"10/25/2012\",4)":                "#VALUE!",
		"=COUPNCD(\"01/01/2011\",\"\",4)":                "#VALUE!",
		"=COUPNCD(\"01/01/2011\",\"10/25/2012\",\"\")":   "#VALUE!",
		"=COUPNCD(\"01/01/2011\",\"10/25/2012\",4,\"\")": "#NUM!",
		"=COUPNCD(\"01/01/2011\",\"10/25/2012\",3)":      "#NUM!",
		"=COUPNCD(\"10/25/2012\",\"01/01/2011\",4)":      "COUPNCD requires maturity > settlement",
		// COUPNUM
		"=COUPNUM()": "COUPNUM requires 3 or 4 arguments",
		"=COUPNUM(\"01/01/2011\",\"10/25/2012\",4,0,0)":  "COUPNUM requires 3 or 4 arguments",
		"=COUPNUM(\"\",\"10/25/2012\",4)":                "#VALUE!",
		"=COUPNUM(\"01/01/2011\",\"\",4)":                "#VALUE!",
		"=COUPNUM(\"01/01/2011\",\"10/25/2012\",\"\")":   "#VALUE!",
		"=COUPNUM(\"01/01/2011\",\"10/25/2012\",4,\"\")": "#NUM!",
		"=COUPNUM(\"01/01/2011\",\"10/25/2012\",3)":      "#NUM!",
		"=COUPNUM(\"10/25/2012\",\"01/01/2011\",4)":      "COUPNUM requires maturity > settlement",
		// COUPPCD
		"=COUPPCD()": "COUPPCD requires 3 or 4 arguments",
		"=COUPPCD(\"01/01/2011\",\"10/25/2012\",4,0,0)":  "COUPPCD requires 3 or 4 arguments",
		"=COUPPCD(\"\",\"10/25/2012\",4)":                "#VALUE!",
		"=COUPPCD(\"01/01/2011\",\"\",4)":                "#VALUE!",
		"=COUPPCD(\"01/01/2011\",\"10/25/2012\",\"\")":   "#VALUE!",
		"=COUPPCD(\"01/01/2011\",\"10/25/2012\",4,\"\")": "#NUM!",
		"=COUPPCD(\"01/01/2011\",\"10/25/2012\",3)":      "#NUM!",
		"=COUPPCD(\"10/25/2012\",\"01/01/2011\",4)":      "COUPPCD requires maturity > settlement",
		// CUMIPMT
		"=CUMIPMT()":               "CUMIPMT requires 6 arguments",
		"=CUMIPMT(0,0,0,0,0,2)":    "#N/A",
		"=CUMIPMT(0,0,0,-1,0,0)":   "#N/A",
		"=CUMIPMT(0,0,0,1,0,0)":    "#N/A",
		"=CUMIPMT(\"\",0,0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMIPMT(0,\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMIPMT(0,0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMIPMT(0,0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMIPMT(0,0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMIPMT(0,0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// CUMPRINC
		"=CUMPRINC()":               "CUMPRINC requires 6 arguments",
		"=CUMPRINC(0,0,0,0,0,2)":    "#N/A",
		"=CUMPRINC(0,0,0,-1,0,0)":   "#N/A",
		"=CUMPRINC(0,0,0,1,0,0)":    "#N/A",
		"=CUMPRINC(\"\",0,0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMPRINC(0,\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMPRINC(0,0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMPRINC(0,0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMPRINC(0,0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=CUMPRINC(0,0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// DB
		"=DB()":             "DB requires at least 4 arguments",
		"=DB(0,0,0,0,0,0)":  "DB allows at most 5 arguments",
		"=DB(-1,0,0,0)":     "#N/A",
		"=DB(\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DB(0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DB(0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DB(0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DB(0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// DDB
		"=DDB()":             "DDB requires at least 4 arguments",
		"=DDB(0,0,0,0,0,0)":  "DDB allows at most 5 arguments",
		"=DDB(-1,0,0,0)":     "#N/A",
		"=DDB(\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DDB(0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DDB(0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DDB(0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DDB(0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// DISC
		"=DISC()":                                          "DISC requires 4 or 5 arguments",
		"=DISC(\"\",\"03/31/2021\",95,100)":                "#VALUE!",
		"=DISC(\"04/01/2016\",\"\",95,100)":                "#VALUE!",
		"=DISC(\"04/01/2016\",\"03/31/2021\",\"\",100)":    "#VALUE!",
		"=DISC(\"04/01/2016\",\"03/31/2021\",95,\"\")":     "#VALUE!",
		"=DISC(\"04/01/2016\",\"03/31/2021\",95,100,\"\")": "#NUM!",
		"=DISC(\"03/31/2021\",\"04/01/2016\",95,100)":      "DISC requires maturity > settlement",
		"=DISC(\"04/01/2016\",\"03/31/2021\",0,100)":       "DISC requires pr > 0",
		"=DISC(\"04/01/2016\",\"03/31/2021\",95,0)":        "DISC requires redemption > 0",
		"=DISC(\"04/01/2016\",\"03/31/2021\",95,100,5)":    "invalid basis",
		// DOLLARDE
		"=DOLLARDE()":       "DOLLARDE requires 2 arguments",
		"=DOLLARDE(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DOLLARDE(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DOLLARDE(0,-1)":   "#NUM!",
		"=DOLLARDE(0,0)":    "#DIV/0!",
		// DOLLARFR
		"=DOLLARFR()":       "DOLLARFR requires 2 arguments",
		"=DOLLARFR(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DOLLARFR(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DOLLARFR(0,-1)":   "#NUM!",
		"=DOLLARFR(0,0)":    "#DIV/0!",
		// DURATION
		"=DURATION()": "DURATION requires 5 or 6 arguments",
		"=DURATION(\"\",\"03/31/2025\",10%,8%,4)":                "#VALUE!",
		"=DURATION(\"04/01/2015\",\"\",10%,8%,4)":                "#VALUE!",
		"=DURATION(\"03/31/2025\",\"04/01/2015\",10%,8%,4)":      "DURATION requires maturity > settlement",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",-1,8%,4)":       "DURATION requires coupon >= 0",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",10%,-1,4)":      "DURATION requires yld >= 0",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",\"\",8%,4)":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",10%,\"\",4)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,3)":      "#NUM!",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,4,\"\")": "#NUM!",
		"=DURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,4,5)":    "invalid basis",
		// EFFECT
		"=EFFECT()":       "EFFECT requires 2 arguments",
		"=EFFECT(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=EFFECT(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=EFFECT(0,0)":    "#NUM!",
		"=EFFECT(1,0)":    "#NUM!",
		// FV
		"=FV()":              "FV requires at least 3 arguments",
		"=FV(0,0,0,0,0,0,0)": "FV allows at most 5 arguments",
		"=FV(0,0,0,0,2)":     "#N/A",
		"=FV(\"\",0,0,0,0)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FV(0,\"\",0,0,0)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FV(0,0,\"\",0,0)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FV(0,0,0,\"\",0)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FV(0,0,0,0,\"\")":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		// FVSCHEDULE
		"=FVSCHEDULE()":        "FVSCHEDULE requires 2 arguments",
		"=FVSCHEDULE(\"\",0)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=FVSCHEDULE(0,\"x\")": "strconv.ParseFloat: parsing \"x\": invalid syntax",
		// INTRATE
		"=INTRATE()":                                          "INTRATE requires 4 or 5 arguments",
		"=INTRATE(\"\",\"03/31/2021\",95,100)":                "#VALUE!",
		"=INTRATE(\"04/01/2016\",\"\",95,100)":                "#VALUE!",
		"=INTRATE(\"04/01/2016\",\"03/31/2021\",\"\",100)":    "#VALUE!",
		"=INTRATE(\"04/01/2016\",\"03/31/2021\",95,\"\")":     "#VALUE!",
		"=INTRATE(\"04/01/2016\",\"03/31/2021\",95,100,\"\")": "#NUM!",
		"=INTRATE(\"03/31/2021\",\"04/01/2016\",95,100)":      "INTRATE requires maturity > settlement",
		"=INTRATE(\"04/01/2016\",\"03/31/2021\",0,100)":       "INTRATE requires investment > 0",
		"=INTRATE(\"04/01/2016\",\"03/31/2021\",95,0)":        "INTRATE requires redemption > 0",
		"=INTRATE(\"04/01/2016\",\"03/31/2021\",95,100,5)":    "invalid basis",
		// IPMT
		"=IPMT()":               "IPMT requires at least 4 arguments",
		"=IPMT(0,0,0,0,0,0,0)":  "IPMT allows at most 6 arguments",
		"=IPMT(0,0,0,0,0,2)":    "#N/A",
		"=IPMT(0,-1,0,0,0,0)":   "#N/A",
		"=IPMT(0,1,0,0,0,0)":    "#N/A",
		"=IPMT(\"\",0,0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=IPMT(0,\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=IPMT(0,0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=IPMT(0,0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=IPMT(0,0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=IPMT(0,0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// ISPMT
		"=ISPMT()":           "ISPMT requires 4 arguments",
		"=ISPMT(\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ISPMT(0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ISPMT(0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ISPMT(0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// MDURATION
		"=MDURATION()": "MDURATION requires 5 or 6 arguments",
		"=MDURATION(\"\",\"03/31/2025\",10%,8%,4)":                "#VALUE!",
		"=MDURATION(\"04/01/2015\",\"\",10%,8%,4)":                "#VALUE!",
		"=MDURATION(\"03/31/2025\",\"04/01/2015\",10%,8%,4)":      "MDURATION requires maturity > settlement",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",-1,8%,4)":       "MDURATION requires coupon >= 0",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",10%,-1,4)":      "MDURATION requires yld >= 0",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",\"\",8%,4)":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",10%,\"\",4)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,3)":      "#NUM!",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,4,\"\")": "#NUM!",
		"=MDURATION(\"04/01/2015\",\"03/31/2025\",10%,8%,4,5)":    "invalid basis",
		// NOMINAL
		"=NOMINAL()":       "NOMINAL requires 2 arguments",
		"=NOMINAL(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NOMINAL(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NOMINAL(0,0)":    "#NUM!",
		"=NOMINAL(1,0)":    "#NUM!",
		// NPER
		"=NPER()":             "NPER requires at least 3 arguments",
		"=NPER(0,0,0,0,0,0)":  "NPER allows at most 5 arguments",
		"=NPER(0,0,0)":        "#NUM!",
		"=NPER(0,0,0,0,2)":    "#N/A",
		"=NPER(\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NPER(0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NPER(0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NPER(0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=NPER(0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// NPV
		"=NPV()":       "NPV requires at least 2 arguments",
		"=NPV(\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// ODDFPRICE
		"=ODDFPRICE()": "ODDFPRICE requires 8 or 9 arguments",
		"=ODDFPRICE(\"\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,2)":                "#VALUE!",
		"=ODDFPRICE(\"02/01/2017\",\"\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,2)":                "#VALUE!",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"\",\"03/31/2017\",5.5%,3.5%,100,2)":                "#VALUE!",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"\",5.5%,3.5%,100,2)":                "#VALUE!",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",\"\",3.5%,100,2)":      "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,\"\",100,2)":      "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,\"\",2)":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"02/01/2017\",\"03/31/2017\",5.5%,3.5%,100,2)":      "ODDFPRICE requires settlement > issue",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"02/01/2017\",5.5%,3.5%,100,2)":      "ODDFPRICE requires first_coupon > settlement",
		"=ODDFPRICE(\"02/01/2017\",\"02/01/2017\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,2)":      "ODDFPRICE requires maturity > first_coupon",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",-1,3.5%,100,2)":        "ODDFPRICE requires rate >= 0",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,-1,100,2)":        "ODDFPRICE requires yld >= 0",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,0,2)":        "ODDFPRICE requires redemption > 0",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,2,\"\")": "#NUM!",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,3)":      "#NUM!",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/30/2017\",5.5%,3.5%,100,4)":      "#NUM!",
		"=ODDFPRICE(\"02/01/2017\",\"03/31/2021\",\"12/01/2016\",\"03/31/2017\",5.5%,3.5%,100,2,5)":    "invalid basis",
		// PDURATION
		"=PDURATION()":         "PDURATION requires 3 arguments",
		"=PDURATION(\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PDURATION(0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PDURATION(0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PDURATION(0,0,0)":    "#NUM!",
		// PMT
		"=PMT()":             "PMT requires at least 3 arguments",
		"=PMT(0,0,0,0,0,0)":  "PMT allows at most 5 arguments",
		"=PMT(0,0,0,0,2)":    "#N/A",
		"=PMT(\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PMT(0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PMT(0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PMT(0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PMT(0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// PRICE
		"=PRICE()": "PRICE requires 6 or 7 arguments",
		"=PRICE(\"\",\"02/01/2020\",12%,10%,100,2,4)":              "#VALUE!",
		"=PRICE(\"04/01/2012\",\"\",12%,10%,100,2,4)":              "#VALUE!",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",\"\",10%,100,2,4)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,\"\",100,2,4)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,\"\",2,4)":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,100,\"\",4)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",-1,10%,100,2,4)":     "PRICE requires rate >= 0",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,-1,100,2,4)":     "PRICE requires yld >= 0",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,0,2,4)":      "PRICE requires redemption > 0",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,100,2,\"\")": "#NUM!",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,100,3,4)":    "#NUM!",
		"=PRICE(\"04/01/2012\",\"02/01/2020\",12%,10%,100,2,5)":    "invalid basis",
		// PPMT
		"=PPMT()":               "PPMT requires at least 4 arguments",
		"=PPMT(0,0,0,0,0,0,0)":  "PPMT allows at most 6 arguments",
		"=PPMT(0,0,0,0,0,2)":    "#N/A",
		"=PPMT(0,-1,0,0,0,0)":   "#N/A",
		"=PPMT(0,1,0,0,0,0)":    "#N/A",
		"=PPMT(\"\",0,0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PPMT(0,\"\",0,0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PPMT(0,0,\"\",0,0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PPMT(0,0,0,\"\",0,0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PPMT(0,0,0,0,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PPMT(0,0,0,0,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// PRICEDISC
		"=PRICEDISC()":                                          "PRICEDISC requires 4 or 5 arguments",
		"=PRICEDISC(\"\",\"03/31/2021\",95,100)":                "#VALUE!",
		"=PRICEDISC(\"04/01/2016\",\"\",95,100)":                "#VALUE!",
		"=PRICEDISC(\"04/01/2016\",\"03/31/2021\",\"\",100)":    "#VALUE!",
		"=PRICEDISC(\"04/01/2016\",\"03/31/2021\",95,\"\")":     "#VALUE!",
		"=PRICEDISC(\"04/01/2016\",\"03/31/2021\",95,100,\"\")": "#NUM!",
		"=PRICEDISC(\"03/31/2021\",\"04/01/2016\",95,100)":      "PRICEDISC requires maturity > settlement",
		"=PRICEDISC(\"04/01/2016\",\"03/31/2021\",0,100)":       "PRICEDISC requires discount > 0",
		"=PRICEDISC(\"04/01/2016\",\"03/31/2021\",95,0)":        "PRICEDISC requires redemption > 0",
		"=PRICEDISC(\"04/01/2016\",\"03/31/2021\",95,100,5)":    "invalid basis",
		// PRICEMAT
		"=PRICEMAT()": "PRICEMAT requires 5 or 6 arguments",
		"=PRICEMAT(\"\",\"03/31/2021\",\"01/01/2017\",4.5%,2.5%)":                "#VALUE!",
		"=PRICEMAT(\"04/01/2017\",\"\",\"01/01/2017\",4.5%,2.5%)":                "#VALUE!",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"\",4.5%,2.5%)":                "#VALUE!",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",\"\",2.5%)":      "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",4.5%,\"\")":      "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",4.5%,2.5%,\"\")": "#NUM!",
		"=PRICEMAT(\"03/31/2021\",\"04/01/2017\",\"01/01/2017\",4.5%,2.5%)":      "PRICEMAT requires maturity > settlement",
		"=PRICEMAT(\"01/01/2017\",\"03/31/2021\",\"04/01/2017\",4.5%,2.5%)":      "PRICEMAT requires settlement > issue",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",-1,2.5%)":        "PRICEMAT requires rate >= 0",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",4.5%,-1)":        "PRICEMAT requires yld >= 0",
		"=PRICEMAT(\"04/01/2017\",\"03/31/2021\",\"01/01/2017\",4.5%,2.5%,5)":    "invalid basis",
		// PV
		"=PV()":                     "PV requires at least 3 arguments",
		"=PV(10%/4,16,2000,0,1,0)":  "PV allows at most 5 arguments",
		"=PV(\"\",16,2000,0,1)":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PV(10%/4,\"\",2000,0,1)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PV(10%/4,16,\"\",0,1)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PV(10%/4,16,2000,\"\",1)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=PV(10%/4,16,2000,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		// RATE
		"=RATE()":                        "RATE requires at least 3 arguments",
		"=RATE(48,-200,8000,3,1,0.5,0)":  "RATE allows at most 6 arguments",
		"=RATE(\"\",-200,8000,3,1,0.5)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RATE(48,\"\",8000,3,1,0.5)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RATE(48,-200,\"\",3,1,0.5)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RATE(48,-200,8000,\"\",1,0.5)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RATE(48,-200,8000,3,\"\",0.5)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RATE(48,-200,8000,3,1,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		// RECEIVED
		"=RECEIVED()": "RECEIVED requires at least 4 arguments",
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",1000,4.5%,1,0)":  "RECEIVED allows at most 5 arguments",
		"=RECEIVED(\"\",\"03/31/2016\",1000,4.5%,1)":              "#VALUE!",
		"=RECEIVED(\"04/01/2011\",\"\",1000,4.5%,1)":              "#VALUE!",
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",\"\",4.5%,1)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",1000,\"\",1)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",1000,4.5%,\"\")": "#NUM!",
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",1000,0)":         "RECEIVED requires discount > 0",
		"=RECEIVED(\"04/01/2011\",\"03/31/2016\",1000,4.5%,5)":    "invalid basis",
		// RRI
		"=RRI()":               "RRI requires 3 arguments",
		"=RRI(\"\",\"\",\"\")": "#NUM!",
		"=RRI(0,10000,15000)":  "RRI requires nper argument to be > 0",
		"=RRI(10,0,15000)":     "RRI requires pv argument to be > 0",
		"=RRI(10,10000,-1)":    "RRI requires fv argument to be >= 0",
		// SLN
		"=SLN()":               "SLN requires 3 arguments",
		"=SLN(\"\",\"\",\"\")": "#NUM!",
		"=SLN(10000,1000,0)":   "SLN requires life argument to be > 0",
		// SYD
		"=SYD()":                    "SYD requires 4 arguments",
		"=SYD(\"\",\"\",\"\",\"\")": "#NUM!",
		"=SYD(10000,1000,0,1)":      "SYD requires life argument to be > 0",
		"=SYD(10000,1000,5,0)":      "SYD requires per argument to be > 0",
		"=SYD(10000,1000,1,5)":      "#NUM!",
		// TBILLEQ
		"=TBILLEQ()":                                   "TBILLEQ requires 3 arguments",
		"=TBILLEQ(\"\",\"06/30/2017\",2.5%)":           "#VALUE!",
		"=TBILLEQ(\"01/01/2017\",\"\",2.5%)":           "#VALUE!",
		"=TBILLEQ(\"01/01/2017\",\"06/30/2017\",\"\")": "#VALUE!",
		"=TBILLEQ(\"01/01/2017\",\"06/30/2017\",0)":    "#NUM!",
		"=TBILLEQ(\"01/01/2017\",\"06/30/2018\",2.5%)": "#NUM!",
		"=TBILLEQ(\"06/30/2017\",\"01/01/2017\",2.5%)": "#NUM!",
		// TBILLPRICE
		"=TBILLPRICE()":                                   "TBILLPRICE requires 3 arguments",
		"=TBILLPRICE(\"\",\"06/30/2017\",2.5%)":           "#VALUE!",
		"=TBILLPRICE(\"01/01/2017\",\"\",2.5%)":           "#VALUE!",
		"=TBILLPRICE(\"01/01/2017\",\"06/30/2017\",\"\")": "#VALUE!",
		"=TBILLPRICE(\"01/01/2017\",\"06/30/2017\",0)":    "#NUM!",
		"=TBILLPRICE(\"01/01/2017\",\"06/30/2018\",2.5%)": "#NUM!",
		"=TBILLPRICE(\"06/30/2017\",\"01/01/2017\",2.5%)": "#NUM!",
		// TBILLYIELD
		"=TBILLYIELD()":                                   "TBILLYIELD requires 3 arguments",
		"=TBILLYIELD(\"\",\"06/30/2017\",2.5%)":           "#VALUE!",
		"=TBILLYIELD(\"01/01/2017\",\"\",2.5%)":           "#VALUE!",
		"=TBILLYIELD(\"01/01/2017\",\"06/30/2017\",\"\")": "#VALUE!",
		"=TBILLYIELD(\"01/01/2017\",\"06/30/2017\",0)":    "#NUM!",
		"=TBILLYIELD(\"01/01/2017\",\"06/30/2018\",2.5%)": "#NUM!",
		"=TBILLYIELD(\"06/30/2017\",\"01/01/2017\",2.5%)": "#NUM!",
		// VDB
		"=VDB()":                          "VDB requires 5 or 7 arguments",
		"=VDB(-1,1000,5,0,1)":             "VDB requires cost >= 0",
		"=VDB(10000,-1,5,0,1)":            "VDB requires salvage >= 0",
		"=VDB(10000,1000,0,0,1)":          "VDB requires life > 0",
		"=VDB(10000,1000,5,-1,1)":         "VDB requires start_period > 0",
		"=VDB(10000,1000,5,2,1)":          "VDB requires start_period <= end_period",
		"=VDB(10000,1000,5,0,6)":          "VDB requires end_period <= life",
		"=VDB(10000,1000,5,0,1,-0.2)":     "VDB requires factor >= 0",
		"=VDB(\"\",1000,5,0,1)":           "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=VDB(10000,\"\",5,0,1)":          "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=VDB(10000,1000,\"\",0,1)":       "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=VDB(10000,1000,5,\"\",1)":       "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=VDB(10000,1000,5,0,\"\")":       "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=VDB(10000,1000,5,0,1,\"\")":     "#NUM!",
		"=VDB(10000,1000,5,0,1,0.2,\"\")": "#NUM!",
		// YIELD
		"=YIELD()": "YIELD requires 6 or 7 arguments",
		"=YIELD(\"\",\"06/30/2015\",10%,101,100,4)":                "#VALUE!",
		"=YIELD(\"01/01/2010\",\"\",10%,101,100,4)":                "#VALUE!",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",\"\",101,100,4)":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,\"\",100,4)":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,\"\",4)":     "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,100,\"\")":   "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,100,4,\"\")": "#NUM!",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,100,3)":      "#NUM!",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,100,4,5)":    "invalid basis",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",-1,101,100,4)":       "PRICE requires rate >= 0",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,0,100,4)":        "PRICE requires pr > 0",
		"=YIELD(\"01/01/2010\",\"06/30/2015\",10%,101,-1,4)":       "PRICE requires redemption >= 0",
		// YIELDDISC
		"=YIELDDISC()": "YIELDDISC requires 4 or 5 arguments",
		"=YIELDDISC(\"\",\"06/30/2017\",97,100,0)":              "#VALUE!",
		"=YIELDDISC(\"01/01/2017\",\"\",97,100,0)":              "#VALUE!",
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",\"\",100,0)":  "#VALUE!",
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",97,\"\",0)":   "#VALUE!",
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",97,100,\"\")": "#NUM!",
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",0,100)":       "YIELDDISC requires pr > 0",
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",97,0)":        "YIELDDISC requires redemption > 0",
		"=YIELDDISC(\"01/01/2017\",\"06/30/2017\",97,100,5)":    "invalid basis",
		// YIELDMAT
		"=YIELDMAT()": "YIELDMAT requires 5 or 6 arguments",
		"=YIELDMAT(\"\",\"06/30/2018\",\"06/01/2014\",5.5%,101,0)":            "#VALUE!",
		"=YIELDMAT(\"01/01/2017\",\"\",\"06/01/2014\",5.5%,101,0)":            "#VALUE!",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"\",5.5%,101,0)":            "#VALUE!",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",\"\",101,0)":  "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",5,\"\",0)":    "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",5,5.5%,\"\")": "#NUM!",
		"=YIELDMAT(\"06/01/2014\",\"06/30/2018\",\"01/01/2017\",5.5%,101,0)":  "YIELDMAT requires settlement > issue",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",-1,101,0)":    "YIELDMAT requires rate >= 0",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",1,0,0)":       "YIELDMAT requires pr > 0",
		"=YIELDMAT(\"01/01/2017\",\"06/30/2018\",\"06/01/2014\",5.5%,101,5)":  "invalid basis",
	}
	for formula, expected := range mathCalcError {
		f := prepareCalcData(cellData)
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}

	referenceCalc := map[string]string{
		// MDETERM
		"=MDETERM(A1:B2)": "-3",
		// PRODUCT
		"=PRODUCT(Sheet1!A1:Sheet1!A1:A2,A2)": "4",
		// IMPRODUCT
		"=IMPRODUCT(Sheet1!A1:Sheet1!A1:A2,A2)": "4",
		// SUM
		"=A1/A3":                          "0.333333333333333",
		"=SUM(A1:A2)":                     "3",
		"=SUM(Sheet1!A1,A2)":              "3",
		"=(-2-SUM(-4+A2))*5":              "0",
		"=SUM(Sheet1!A1:Sheet1!A1:A2,A2)": "5",
		"=SUM(A1,A2,A3)*SUM(2,3)":         "30",
		"=1+SUM(SUM(A1+A2/A3)*(2-3),2)":   "1.33333333333333",
		"=A1/A2/SUM(A1:A2:B1)":            "0.0416666666666667",
		"=A1/A2/SUM(A1:A2:B1)*A3":         "0.125",
		"=SUM(B1:D1)":                     "4",
		"=SUM(\"X\")":                     "0",
	}
	for formula, expected := range referenceCalc {
		f := prepareCalcData(cellData)
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err)
		assert.Equal(t, expected, result, formula)
	}

	referenceCalcError := map[string]string{
		// MDETERM
		"=MDETERM(A1:B3)": "#VALUE!",
		// SUM
		"=1+SUM(SUM(A1+A2/A4)*(2-3),2)": "#DIV/0!",
	}
	for formula, expected := range referenceCalcError {
		f := prepareCalcData(cellData)
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}

	volatileFuncs := []string{
		"=NOW()",
		"=RAND()",
		"=RANDBETWEEN(1,2)",
		"=TODAY()",
	}
	for _, formula := range volatileFuncs {
		f := prepareCalcData(cellData)
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		_, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err)
	}

	// Test get calculated cell value on not formula cell.
	f := prepareCalcData(cellData)
	result, err := f.CalcCellValue("Sheet1", "A1")
	assert.NoError(t, err)
	assert.Equal(t, "", result)
	// Test get calculated cell value on not exists worksheet.
	f = prepareCalcData(cellData)
	_, err = f.CalcCellValue("SheetN", "A1")
	assert.EqualError(t, err, "sheet SheetN is not exist")
	// Test get calculated cell value with not support formula.
	f = prepareCalcData(cellData)
	assert.NoError(t, f.SetCellFormula("Sheet1", "A1", "=UNSUPPORT(A1)"))
	_, err = f.CalcCellValue("Sheet1", "A1")
	assert.EqualError(t, err, "not support UNSUPPORT function")
	assert.NoError(t, f.SaveAs(filepath.Join("test", "TestCalcCellValue.xlsx")))
}

func TestCalculate(t *testing.T) {
	err := `strconv.ParseFloat: parsing "string": invalid syntax`
	opd := NewStack()
	opd.Push(efp.Token{TValue: "string"})
	opt := efp.Token{TValue: "-", TType: efp.TokenTypeOperatorPrefix}
	assert.EqualError(t, calculate(opd, opt), err)
	opd.Push(efp.Token{TValue: "string"})
	opd.Push(efp.Token{TValue: "string"})
	opt = efp.Token{TValue: "-", TType: efp.TokenTypeOperatorInfix}
	assert.EqualError(t, calculate(opd, opt), err)
}

func TestCalcWithDefinedName(t *testing.T) {
	cellData := [][]interface{}{
		{"A1_as_string", "B1_as_string", 123, nil},
	}
	f := prepareCalcData(cellData)
	assert.NoError(t, f.SetDefinedName(&DefinedName{Name: "defined_name1", RefersTo: "Sheet1!A1", Scope: "Workbook"}))
	assert.NoError(t, f.SetDefinedName(&DefinedName{Name: "defined_name1", RefersTo: "Sheet1!B1", Scope: "Sheet1"}))
	assert.NoError(t, f.SetDefinedName(&DefinedName{Name: "defined_name2", RefersTo: "Sheet1!C1", Scope: "Workbook"}))

	assert.NoError(t, f.SetCellFormula("Sheet1", "D1", "=defined_name1"))
	result, err := f.CalcCellValue("Sheet1", "D1")
	assert.NoError(t, err)
	// DefinedName with scope WorkSheet takes precedence over DefinedName with scope Workbook, so we should get B1 value
	assert.Equal(t, "B1_as_string", result, "=defined_name1")

	assert.NoError(t, f.SetCellFormula("Sheet1", "D1", `=CONCATENATE("<",defined_name1,">")`))
	result, err = f.CalcCellValue("Sheet1", "D1")
	assert.NoError(t, err)
	assert.Equal(t, "<B1_as_string>", result, "=defined_name1")

	// comparing numeric values
	assert.NoError(t, f.SetCellFormula("Sheet1", "D1", `=123=defined_name2`))
	result, err = f.CalcCellValue("Sheet1", "D1")
	assert.NoError(t, err)
	assert.Equal(t, "TRUE", result, "=123=defined_name2")

	// comparing text values
	assert.NoError(t, f.SetCellFormula("Sheet1", "D1", `="B1_as_string"=defined_name1`))
	result, err = f.CalcCellValue("Sheet1", "D1")
	assert.NoError(t, err)
	assert.Equal(t, "TRUE", result, `="B1_as_string"=defined_name1`)

	// comparing text values
	assert.NoError(t, f.SetCellFormula("Sheet1", "D1", `=IF("B1_as_string"=defined_name1,"YES","NO")`))
	result, err = f.CalcCellValue("Sheet1", "D1")
	assert.NoError(t, err)
	assert.Equal(t, "YES", result, `=IF("B1_as_string"=defined_name1,"YES","NO")`)
}

func TestCalcISBLANK(t *testing.T) {
	argsList := list.New()
	argsList.PushBack(formulaArg{
		Type: ArgUnknown,
	})
	fn := formulaFuncs{}
	result := fn.ISBLANK(argsList)
	assert.Equal(t, result.String, "TRUE")
	assert.Empty(t, result.Error)
}

func TestCalcAND(t *testing.T) {
	argsList := list.New()
	argsList.PushBack(formulaArg{
		Type: ArgUnknown,
	})
	fn := formulaFuncs{}
	result := fn.AND(argsList)
	assert.Equal(t, result.String, "")
	assert.Empty(t, result.Error)
}

func TestCalcOR(t *testing.T) {
	argsList := list.New()
	argsList.PushBack(formulaArg{
		Type: ArgUnknown,
	})
	fn := formulaFuncs{}
	result := fn.OR(argsList)
	assert.Equal(t, result.String, "FALSE")
	assert.Empty(t, result.Error)
}

func TestCalcDet(t *testing.T) {
	assert.Equal(t, det([][]float64{
		{1, 2, 3, 4},
		{2, 3, 4, 5},
		{3, 4, 5, 6},
		{4, 5, 6, 7},
	}), float64(0))
}

func TestCalcToBool(t *testing.T) {
	b := newBoolFormulaArg(true).ToBool()
	assert.Equal(t, b.Boolean, true)
	assert.Equal(t, b.Number, 1.0)
}

func TestCalcToList(t *testing.T) {
	assert.Equal(t, []formulaArg(nil), newEmptyFormulaArg().ToList())
	formulaList := []formulaArg{newEmptyFormulaArg()}
	assert.Equal(t, formulaList, newListFormulaArg(formulaList).ToList())
}

func TestCalcCompareFormulaArg(t *testing.T) {
	assert.Equal(t, compareFormulaArg(newEmptyFormulaArg(), newEmptyFormulaArg(), newNumberFormulaArg(matchModeMaxLess), false), criteriaEq)
	lhs := newListFormulaArg([]formulaArg{newEmptyFormulaArg()})
	rhs := newListFormulaArg([]formulaArg{newEmptyFormulaArg(), newEmptyFormulaArg()})
	assert.Equal(t, compareFormulaArg(lhs, rhs, newNumberFormulaArg(matchModeMaxLess), false), criteriaL)
	assert.Equal(t, compareFormulaArg(rhs, lhs, newNumberFormulaArg(matchModeMaxLess), false), criteriaG)

	lhs = newListFormulaArg([]formulaArg{newBoolFormulaArg(true)})
	rhs = newListFormulaArg([]formulaArg{newBoolFormulaArg(true)})
	assert.Equal(t, compareFormulaArg(lhs, rhs, newNumberFormulaArg(matchModeMaxLess), false), criteriaEq)

	assert.Equal(t, compareFormulaArg(formulaArg{Type: ArgUnknown}, formulaArg{Type: ArgUnknown}, newNumberFormulaArg(matchModeMaxLess), false), criteriaErr)
}

func TestCalcMatchPattern(t *testing.T) {
	assert.True(t, matchPattern("", ""))
	assert.True(t, matchPattern("file/*", "file/abc/bcd/def"))
	assert.True(t, matchPattern("*", ""))
	assert.False(t, matchPattern("file/?", "file/abc/bcd/def"))
}

func TestCalcTRANSPOSE(t *testing.T) {
	cellData := [][]interface{}{
		{"a", "d"},
		{"b", "e"},
		{"c", "f"},
	}
	formula := "=TRANSPOSE(A1:A3)"
	f := prepareCalcData(cellData)
	formulaType, ref := STCellFormulaTypeArray, "D1:F2"
	assert.NoError(t, f.SetCellFormula("Sheet1", "D1", formula,
		FormulaOpts{Ref: &ref, Type: &formulaType}))
	_, err := f.CalcCellValue("Sheet1", "D1")
	assert.NoError(t, err, formula)
}

func TestCalcVLOOKUP(t *testing.T) {
	cellData := [][]interface{}{
		{nil, nil, nil, nil, nil, nil},
		{nil, "Score", "Grade", nil, nil, nil},
		{nil, 0, "F", nil, "Score", 85},
		{nil, 60, "D", nil, "Grade"},
		{nil, 70, "C", nil, nil, nil},
		{nil, 80, "b", nil, nil, nil},
		{nil, 90, "A", nil, nil, nil},
		{nil, 85, "B", nil, nil, nil},
		{nil, nil, nil, nil, nil, nil},
	}
	f := prepareCalcData(cellData)
	calc := map[string]string{
		"=VLOOKUP(F3,B3:C8,2)":       "b",
		"=VLOOKUP(F3,B3:C8,2,TRUE)":  "b",
		"=VLOOKUP(F3,B3:C8,2,FALSE)": "B",
	}
	for formula, expected := range calc {
		assert.NoError(t, f.SetCellFormula("Sheet1", "F4", formula))
		result, err := f.CalcCellValue("Sheet1", "F4")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=VLOOKUP(INT(1),C3:C3,1,FALSE)": "VLOOKUP no result found",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "F4", formula))
		result, err := f.CalcCellValue("Sheet1", "F4")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcBoolean(t *testing.T) {
	cellData := [][]interface{}{
		{0.5, "TRUE", -0.5, "FALSE"},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=AVERAGEA(A1:C1)":  "0.333333333333333",
		"=MAX(0.5,B1)":      "0.5",
		"=MAX(A1:B1)":       "0.5",
		"=MAXA(A1:B1)":      "1",
		"=MAXA(0.5,B1)":     "1",
		"=MIN(-0.5,D1)":     "-0.5",
		"=MIN(C1:D1)":       "-0.5",
		"=MINA(C1:D1)":      "-0.5",
		"=MINA(-0.5,D1)":    "-0.5",
		"=STDEV(A1:C1)":     "0.707106781186548",
		"=STDEV(A1,B1,C1)":  "0.707106781186548",
		"=STDEVA(A1:C1,B1)": "0.707106781186548",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "B10", formula))
		result, err := f.CalcCellValue("Sheet1", "B10")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
}

func TestCalcAVERAGEIF(t *testing.T) {
	f := prepareCalcData([][]interface{}{
		{"Monday", 500},
		{"Tuesday", 50},
		{"Thursday", 100},
		{"Friday", 100},
		{"Thursday", 200},
		{5, 300},
		{2, 200},
		{3, 100},
		{4, 50},
		{5, 100},
		{1, 50},
		{"TRUE", 200},
		{"TRUE", 250},
		{"FALSE", 50},
	})
	for formula, expected := range map[string]string{
		"=AVERAGEIF(A1:A14,\"Thursday\",B1:B14)": "150",
		"=AVERAGEIF(A1:A14,5,B1:B14)":            "200",
		"=AVERAGEIF(A1:A14,\">2\",B1:B14)":       "137.5",
		"=AVERAGEIF(A1:A14,TRUE,B1:B14)":         "225",
		"=AVERAGEIF(A1:A14,\"<>TRUE\",B1:B14)":   "150",
	} {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
}

func TestCalcCOVAR(t *testing.T) {
	cellData := [][]interface{}{
		{"array1", "array2"},
		{2, 22.9},
		{7, 33.49},
		{8, 34.5},
		{3, 27.61},
		{4, 19.5},
		{1, 10.11},
		{6, 37.9},
		{5, 31.08},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=COVAR(A1:A9,B1:B9)":        "16.633125",
		"=COVAR(A2:A9,B2:B9)":        "16.633125",
		"=COVARIANCE.P(A1:A9,B1:B9)": "16.633125",
		"=COVARIANCE.P(A2:A9,B2:B9)": "16.633125",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=COVAR()":                   "COVAR requires 2 arguments",
		"=COVAR(A2:A9,B3:B3)":        "#N/A",
		"=COVARIANCE.P()":            "COVARIANCE.P requires 2 arguments",
		"=COVARIANCE.P(A2:A9,B3:B3)": "#N/A",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcFORMULATEXT(t *testing.T) {
	f, formulaText := NewFile(), "=SUM(B1:C1)"
	assert.NoError(t, f.SetCellFormula("Sheet1", "A1", formulaText))
	for _, formula := range []string{"=FORMULATEXT(A1)", "=FORMULATEXT(A:A)", "=FORMULATEXT(A1:B1)"} {
		assert.NoError(t, f.SetCellFormula("Sheet1", "D1", formula), formula)
		result, err := f.CalcCellValue("Sheet1", "D1")
		assert.NoError(t, err, formula)
		assert.Equal(t, formulaText, result, formula)
	}
}

func TestCalcHLOOKUP(t *testing.T) {
	cellData := [][]interface{}{
		{"Example Result Table"},
		{nil, "A", "B", "C", "E", "F"},
		{"Math", .58, .9, .67, .76, .8},
		{"French", .61, .71, .59, .59, .76},
		{"Physics", .75, .45, .39, .52, .69},
		{"Biology", .39, .55, .77, .61, .45},
		{},
		{"Individual Student Score"},
		{"Student:", "Biology Score:"},
		{"E"},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=HLOOKUP(A10,A2:F6,5,FALSE)":  "0.61",
		"=HLOOKUP(D3,D3:D3,1,TRUE)":    "0.67",
		"=HLOOKUP(F3,D3:F3,1,TRUE)":    "0.8",
		"=HLOOKUP(A5,A2:F2,1,TRUE)":    "F",
		"=HLOOKUP(\"D\",A2:F2,1,TRUE)": "C",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "B10", formula))
		result, err := f.CalcCellValue("Sheet1", "B10")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=HLOOKUP(INT(1),A3:A3,1,FALSE)": "HLOOKUP no result found",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "B10", formula))
		result, err := f.CalcCellValue("Sheet1", "B10")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcCHITESTandCHISQdotTEST(t *testing.T) {
	cellData := [][]interface{}{
		{nil, "Observed Frequencies", nil, nil, "Expected Frequencies"},
		{nil, "men", "women", nil, nil, "men", "women"},
		{"answer a", 33, 39, nil, "answer a", 26.25, 31.5},
		{"answer b", 62, 62, nil, "answer b", 57.75, 61.95},
		{"answer c", 10, 4, nil, "answer c", 21, 11.55},
		{nil, -1, 0},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=CHITEST(B3:C5,F3:G5)":    "0.000699102758787672",
		"=CHITEST(B3:C3,F3:G3)":    "0.0605802098655177",
		"=CHITEST(B3:B4,F3:F4)":    "0.152357748933542",
		"=CHITEST(B4:B6,F3:F5)":    "7.07076951440726E-25",
		"=CHISQ.TEST(B3:C5,F3:G5)": "0.000699102758787672",
		"=CHISQ.TEST(B3:C3,F3:G3)": "0.0605802098655177",
		"=CHISQ.TEST(B3:B4,F3:F4)": "0.152357748933542",
		"=CHISQ.TEST(B4:B6,F3:F5)": "7.07076951440726E-25",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "I1", formula))
		result, err := f.CalcCellValue("Sheet1", "I1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=CHITEST()":               "CHITEST requires 2 arguments",
		"=CHITEST(B3:C5,F3:F4)":    "#N/A",
		"=CHITEST(B3:B3,F3:F3)":    "#N/A",
		"=CHITEST(F3:F5,B4:B6)":    "#NUM!",
		"=CHITEST(F3:F5,C4:C6)":    "#DIV/0!",
		"=CHISQ.TEST()":            "CHISQ.TEST requires 2 arguments",
		"=CHISQ.TEST(B3:C5,F3:F4)": "#N/A",
		"=CHISQ.TEST(B3:B3,F3:F3)": "#N/A",
		"=CHISQ.TEST(F3:F5,B4:B6)": "#NUM!",
		"=CHISQ.TEST(F3:F5,C4:C6)": "#DIV/0!",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "I1", formula))
		result, err := f.CalcCellValue("Sheet1", "I1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcFTEST(t *testing.T) {
	cellData := [][]interface{}{
		{"Group 1", "Group 2"},
		{3.5, 9.2},
		{4.7, 8.2},
		{6.2, 7.3},
		{4.9, 6.1},
		{3.8, 5.4},
		{5.5, 7.8},
		{7.1, 5.9},
		{6.7, 8.4},
		{3.9, 7.7},
		{4.6, 6.6},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=FTEST(A2:A11,B2:B11)":  "0.95403555939413",
		"=F.TEST(A2:A11,B2:B11)": "0.95403555939413",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=FTEST()":               "FTEST requires 2 arguments",
		"=FTEST(A2:A2,B2:B2)":    "#DIV/0!",
		"=FTEST(A12:A14,B2:B4)":  "#DIV/0!",
		"=FTEST(A2:A4,B2:B2)":    "#DIV/0!",
		"=FTEST(A2:A4,B12:B14)":  "#DIV/0!",
		"=F.TEST()":              "F.TEST requires 2 arguments",
		"=F.TEST(A2:A2,B2:B2)":   "#DIV/0!",
		"=F.TEST(A12:A14,B2:B4)": "#DIV/0!",
		"=F.TEST(A2:A4,B2:B2)":   "#DIV/0!",
		"=F.TEST(A2:A4,B12:B14)": "#DIV/0!",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcIRR(t *testing.T) {
	cellData := [][]interface{}{{-1}, {0.2}, {0.24}, {0.288}, {0.3456}, {0.4147}}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=IRR(A1:A4)":      "-0.136189509034157",
		"=IRR(A1:A6)":      "0.130575760006905",
		"=IRR(A1:A4,-0.1)": "-0.136189514994621",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "B1", formula))
		result, err := f.CalcCellValue("Sheet1", "B1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=IRR()":       "IRR requires at least 1 argument",
		"=IRR(0,0,0)":  "IRR allows at most 2 arguments",
		"=IRR(0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=IRR(A2:A3)":  "#NUM!",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "B1", formula))
		result, err := f.CalcCellValue("Sheet1", "B1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcMAXMINIFS(t *testing.T) {
	f := NewFile()
	for cell, row := range map[string][]interface{}{
		"A1": {1, -math.MaxFloat64 - 1},
		"A2": {2, -math.MaxFloat64 - 2},
		"A3": {3, math.MaxFloat64 + 1},
		"A4": {4, math.MaxFloat64 + 2},
	} {
		assert.NoError(t, f.SetSheetRow("Sheet1", cell, &row))
	}
	formulaList := map[string]string{
		"=MAX(B1:B2)":                 "0",
		"=MAXIFS(B1:B2,A1:A2,\">0\")": "0",
		"=MIN(B3:B4)":                 "0",
		"=MINIFS(B3:B4,A3:A4,\"<0\")": "0",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
}

func TestCalcMIRR(t *testing.T) {
	cellData := [][]interface{}{{-100}, {18}, {22.5}, {28}, {35.5}, {45}}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=MIRR(A1:A6,0.055,0.05)": "0.1000268752662",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "B1", formula))
		result, err := f.CalcCellValue("Sheet1", "B1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=MIRR()":             "MIRR requires 3 arguments",
		"=MIRR(A1:A5,\"\",0)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MIRR(A1:A5,0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=MIRR(B1:B5,0,0)":    "#DIV/0!",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "B1", formula))
		result, err := f.CalcCellValue("Sheet1", "B1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcSUMIFSAndAVERAGEIFS(t *testing.T) {
	cellData := [][]interface{}{
		{"Quarter", "Area", "Sales Rep.", "Sales"},
		{1, "North", "Jeff", 223000},
		{1, "North", "Chris", 125000},
		{1, "South", "Carol", 456000},
		{2, "North", "Jeff", 322000},
		{2, "North", "Chris", 340000},
		{2, "South", "Carol", 198000},
		{3, "North", "Jeff", 310000},
		{3, "North", "Chris", 250000},
		{3, "South", "Carol", 460000},
		{4, "North", "Jeff", 261000},
		{4, "North", "Chris", 389000},
		{4, "South", "Carol", 305000},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=AVERAGEIFS(D2:D13,A2:A13,1,B2:B13,\"North\")":     "174000",
		"=AVERAGEIFS(D2:D13,A2:A13,\">2\",C2:C13,\"Jeff\")": "285500",
		"=SUMIFS(D2:D13,A2:A13,1,B2:B13,\"North\")":         "348000",
		"=SUMIFS(D2:D13,A2:A13,\">2\",C2:C13,\"Jeff\")":     "571000",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "E1", formula))
		result, err := f.CalcCellValue("Sheet1", "E1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=AVERAGEIFS()":                   "AVERAGEIFS requires at least 3 arguments",
		"=AVERAGEIFS(H1,\"\")":            "AVERAGEIFS requires at least 3 arguments",
		"=AVERAGEIFS(H1,\"\",TRUE,1)":     "#N/A",
		"=AVERAGEIFS(H1,\"\",TRUE)":       "AVERAGEIF divide by zero",
		"=SUMIFS()":                       "SUMIFS requires at least 3 arguments",
		"=SUMIFS(D2:D13,A2:A13,1,B2:B13)": "#N/A",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "E1", formula))
		result, err := f.CalcCellValue("Sheet1", "E1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcXIRR(t *testing.T) {
	cellData := [][]interface{}{
		{-100.00, "01/01/2016"},
		{20.00, "04/01/2016"},
		{40.00, "10/01/2016"},
		{25.00, "02/01/2017"},
		{8.00, "03/01/2017"},
		{15.00, "06/01/2017"},
		{-1e-10, "09/01/2017"},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=XIRR(A1:A4,B1:B4)":     "-0.196743861298328",
		"=XIRR(A1:A6,B1:B6,0.5)": "0.0944390744445204",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=XIRR()":                 "XIRR requires 2 or 3 arguments",
		"=XIRR(A1:A4,B1:B4,-1)":   "XIRR requires guess > -1",
		"=XIRR(\"\",B1:B4)":       "#NUM!",
		"=XIRR(A1:A4,\"\")":       "#NUM!",
		"=XIRR(A1:A4,B1:B4,\"\")": "#NUM!",
		"=XIRR(A2:A6,B2:B6)":      "#NUM!",
		"=XIRR(A2:A7,B2:B7)":      "#NUM!",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcXLOOKUP(t *testing.T) {
	cellData := [][]interface{}{
		{},
		{nil, nil, "Quarter", "Gross Profit", "Net profit", "Profit %"},
		{nil, nil, "Qtr1", nil, 19342, 29.30},
		{},
		{nil, "Income Statement", "Qtr1", "Qtr2", "Qtr3", "Qtr4", "Total"},
		{nil, "Total sales", 50000, 78200, 89500, 91250, 308.95},
		{nil, "Cost of sales", -25000, -42050, -59450, -60450, -186950},
		{nil, "Gross Profit", 25000, 36150, 30050, 30800, 122000},
		{},
		{nil, "Depreciation", -899, -791, -202, -412, -2304},
		{nil, "Interest", -513, -853, -150, -956, -2472},
		{nil, "Earnings before Tax", 23588, 34506, 29698, 29432, 117224},
		{},
		{nil, "Tax", -4246, -6211, -5346, -5298, 21100},
		{},
		{nil, "Net profit", 19342, 28295, 24352, 24134, 96124},
		{nil, "Profit %", 0.293, 0.278, 0.234, 0.276, 0.269},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=SUM(XLOOKUP($C3,$C5:$C5,$C6:$C17,NA(),0,2))":        "87272.293",
		"=SUM(XLOOKUP($C3,$C5:$C5,$C6:$G6,NA(),0,-2))":        "309258.95",
		"=SUM(XLOOKUP($C3,$C5:$C5,$C6:$C17,NA(),0,-2))":       "87272.293",
		"=SUM(XLOOKUP($C3,$C5:$G5,$C6:$G17,NA(),0,2))":        "87272.293",
		"=SUM(XLOOKUP(D2,$B6:$B17,$C6:$G17,NA(),0,2))":        "244000",
		"=XLOOKUP(D2,$B6:$B17,C6:C17)":                        "25000",
		"=XLOOKUP(D2,$B6:$B17,XLOOKUP($C3,$C5:$G5,$C6:$G17))": "25000",
		"=XLOOKUP(\"*p*\",B2:B9,C2:C9,NA(),2)":                "25000",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "D3", formula))
		result, err := f.CalcCellValue("Sheet1", "D3")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=XLOOKUP()": "XLOOKUP requires at least 3 arguments",
		"=XLOOKUP($C3,$C5:$C5,$C6:$C17,NA(),0,2,1)":  "XLOOKUP allows at most 6 arguments",
		"=XLOOKUP($C3,$C5,$C6,NA(),0,2)":             "#N/A",
		"=XLOOKUP($C3,$C4:$D5,$C6:$C17,NA(),0,2)":    "#VALUE!",
		"=XLOOKUP($C3,$C5:$C5,$C6:$G17,NA(),0,-2)":   "#VALUE!",
		"=XLOOKUP($C3,$C5:$G5,$C6:$F7,NA(),0,2)":     "#VALUE!",
		"=XLOOKUP(D2,$B6:$B17,$C6:$G16,NA(),0,2)":    "#VALUE!",
		"=XLOOKUP(D2,$B6:$B17,$C6:$G17,NA(),3,2)":    "#VALUE!",
		"=XLOOKUP(D2,$B6:$B17,$C6:$G17,NA(),0,0)":    "#VALUE!",
		"=XLOOKUP(D2,$B6:$B17,$C6:$G17,NA(),\"\",2)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=XLOOKUP(D2,$B6:$B17,$C6:$G17,NA(),0,\"\")": "strconv.ParseFloat: parsing \"\": invalid syntax",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "D3", formula))
		result, err := f.CalcCellValue("Sheet1", "D3")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}

	cellData = [][]interface{}{
		{"Salesperson", "Item", "Amont"},
		{"B", "Apples", 30, 25, 15, 50, 45, 18},
		{"L", "Oranges", 25, "D3", "E3"},
		{"C", "Grapes", 15},
		{"L", "Lemons", 50},
		{"L", "Oranges", 45},
		{"C", "Peaches", 18},
		{"B", "Pears", 40},
		{"B", "Apples", 55},
	}
	f = prepareCalcData(cellData)
	formulaList = map[string]string{
		// Test match mode with partial match (wildcards)
		"=XLOOKUP(\"*p*\",B2:B9,C2:C9,NA(),2)": "30",
		// Test match mode with approximate match in vertical (next larger item)
		"=XLOOKUP(32,B2:B9,C2:C9,NA(),1)": "30",
		// Test match mode with approximate match in horizontal (next larger item)
		"=XLOOKUP(30,C2:F2,C3:F3,NA(),1)": "25",
		// Test match mode with approximate match in vertical (next smaller item)
		"=XLOOKUP(40,C2:C9,B2:B9,NA(),-1)": "Pears",
		// Test match mode with approximate match in horizontal (next smaller item)
		"=XLOOKUP(29,C2:F2,C3:F3,NA(),-1)": "D3",
		// Test search mode
		"=XLOOKUP(\"L\",A2:A9,C2:C9,NA(),0,1)":  "25",
		"=XLOOKUP(\"L\",A2:A9,C2:C9,NA(),0,-1)": "45",
		"=XLOOKUP(\"L\",A2:A9,C2:C9,NA(),0,2)":  "50",
		"=XLOOKUP(\"L\",A2:A9,C2:C9,NA(),0,-2)": "45",
		// Test match mode and search mode
		"=XLOOKUP(29,C2:H2,C3:H3,NA(),-1,-1)": "D3",
		"=XLOOKUP(29,C2:H2,C3:H3,NA(),-1,1)":  "D3",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "D3", formula))
		result, err := f.CalcCellValue("Sheet1", "D3")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError = map[string]string{
		// Test match mode with exact match
		"=XLOOKUP(\"*p*\",B2:B9,C2:C9,NA(),0)": "#N/A",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "D3", formula))
		result, err := f.CalcCellValue("Sheet1", "D3")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcXNPV(t *testing.T) {
	cellData := [][]interface{}{
		{nil, 0.05},
		{"01/01/2016", -10000, nil},
		{"02/01/2016", 2000},
		{"05/01/2016", 2400},
		{"07/01/2016", 2900},
		{"11/01/2016", 3500},
		{"01/01/2017", 4100},
		{},
		{"02/01/2016"},
		{"01/01/2016"},
	}
	f := prepareCalcData(cellData)
	formulaList := map[string]string{
		"=XNPV(B1,B2:B7,A2:A7)": "4447.93800944052",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=XNPV()":                 "XNPV requires 3 arguments",
		"=XNPV(\"\",B2:B7,A2:A7)": "strconv.ParseFloat: parsing \"\": invalid syntax",
		"=XNPV(0,B2:B7,A2:A7)":    "XNPV requires rate > 0",
		"=XNPV(B1,\"\",A2:A7)":    "#NUM!",
		"=XNPV(B1,B2:B7,\"\")":    "#NUM!",
		"=XNPV(B1,B2:B7,C2:C7)":   "#NUM!",
		"=XNPV(B1,B2,A2)":         "#NUM!",
		"=XNPV(B1,B2:B3,A2:A5)":   "#NUM!",
		"=XNPV(B1,B2:B3,A9:A10)":  "#VALUE!",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "C1", formula))
		result, err := f.CalcCellValue("Sheet1", "C1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
}

func TestCalcMATCH(t *testing.T) {
	f := NewFile()
	for cell, row := range map[string][]interface{}{
		"A1": {"cccc", 7, 4, 16},
		"A2": {"dddd", 2, 6, 11},
		"A3": {"aaaa", 4, 7, 10},
		"A4": {"bbbb", 1, 10, 7},
		"A5": {"eeee", 8, 11, 6},
		"A6": {nil, 11, 16, 4},
	} {
		assert.NoError(t, f.SetSheetRow("Sheet1", cell, &row))
	}
	formulaList := map[string]string{
		"=MATCH(\"aaaa\",A1:A6,0)": "3",
		"=MATCH(\"*b\",A1:A5,0)":   "4",
		"=MATCH(\"?eee\",A1:A5,0)": "5",
		"=MATCH(\"?*?e\",A1:A5,0)": "5",
		"=MATCH(\"aaaa\",A1:A6,1)": "3",
		"=MATCH(10,B1:B6)":         "5",
		"=MATCH(8,C1:C6,1)":        "3",
		"=MATCH(6,B1:B6,-1)":       "1",
		"=MATCH(10,D1:D6,-1)":      "3",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "E1", formula))
		result, err := f.CalcCellValue("Sheet1", "E1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
	calcError := map[string]string{
		"=MATCH(3,C1:C6,1)":  "#N/A",
		"=MATCH(5,C1:C6,-1)": "#N/A",
	}
	for formula, expected := range calcError {
		assert.NoError(t, f.SetCellFormula("Sheet1", "E1", formula))
		result, err := f.CalcCellValue("Sheet1", "E1")
		assert.EqualError(t, err, expected, formula)
		assert.Equal(t, "", result, formula)
	}
	assert.Equal(t, newErrorFormulaArg(formulaErrorNA, formulaErrorNA), calcMatch(2, nil, []formulaArg{}))
}

func TestCalcISFORMULA(t *testing.T) {
	f := NewFile()
	assert.NoError(t, f.SetCellFormula("Sheet1", "B1", "=ISFORMULA(A1)"))
	for _, formula := range []string{"=NA()", "=SUM(A1:A3)"} {
		assert.NoError(t, f.SetCellFormula("Sheet1", "A1", formula))
		result, err := f.CalcCellValue("Sheet1", "B1")
		assert.NoError(t, err, formula)
		assert.Equal(t, "TRUE", result, formula)
	}
}

func TestCalcSHEET(t *testing.T) {
	f := NewFile()
	f.NewSheet("Sheet2")
	formulaList := map[string]string{
		"=SHEET(\"Sheet2\")":   "2",
		"=SHEET(Sheet2!A1)":    "2",
		"=SHEET(Sheet2!A1:A2)": "2",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "A1", formula))
		result, err := f.CalcCellValue("Sheet1", "A1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
}

func TestCalcSHEETS(t *testing.T) {
	f := NewFile()
	f.NewSheet("Sheet2")
	formulaList := map[string]string{
		"=SHEETS(Sheet1!A1:B1)":        "1",
		"=SHEETS(Sheet1!A1:Sheet1!A1)": "1",
		"=SHEETS(Sheet1!A1:Sheet2!A1)": "2",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "A1", formula))
		result, err := f.CalcCellValue("Sheet1", "A1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
}

func TestCalcZTEST(t *testing.T) {
	f := NewFile()
	assert.NoError(t, f.SetSheetRow("Sheet1", "A1", &[]int{4, 5, 2, 5, 8, 9, 3, 2, 3, 8, 9, 5}))
	formulaList := map[string]string{
		"=Z.TEST(A1:L1,5)":   "0.371103278558538",
		"=Z.TEST(A1:L1,6)":   "0.838129187019751",
		"=Z.TEST(A1:L1,5,1)": "0.193238115385616",
		"=ZTEST(A1:L1,5)":    "0.371103278558538",
		"=ZTEST(A1:L1,6)":    "0.838129187019751",
		"=ZTEST(A1:L1,5,1)":  "0.193238115385616",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "M1", formula))
		result, err := f.CalcCellValue("Sheet1", "M1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
}

func TestStrToDate(t *testing.T) {
	_, _, _, _, err := strToDate("")
	assert.Equal(t, formulaErrorVALUE, err.Error)
}

func TestGetYearDays(t *testing.T) {
	for _, data := range [][]int{{2021, 0, 360}, {2000, 1, 366}, {2021, 1, 365}, {2000, 3, 365}} {
		assert.Equal(t, data[2], getYearDays(data[0], data[1]))
	}
}

func TestCalcGetBetaHelperContFrac(t *testing.T) {
	assert.Equal(t, 1.0, getBetaHelperContFrac(1, 0, 1))
}

func TestCalcGetBetaDistPDF(t *testing.T) {
	assert.Equal(t, 0.0, getBetaDistPDF(0.5, 2000, 3))
	assert.Equal(t, 0.0, getBetaDistPDF(0, 1, 0))
}

func TestCalcD1mach(t *testing.T) {
	assert.Equal(t, 0.0, d1mach(6))
}

func TestCalcChebyshevInit(t *testing.T) {
	assert.Equal(t, 0, chebyshevInit(0, 0, nil))
	assert.Equal(t, 0, chebyshevInit(1, 0, []float64{0}))
}

func TestCalcChebyshevEval(t *testing.T) {
	assert.True(t, math.IsNaN(chebyshevEval(0, 0, nil)))
}

func TestCalcLgammacor(t *testing.T) {
	assert.True(t, math.IsNaN(lgammacor(9)))
	assert.Equal(t, 4.930380657631324e-32, lgammacor(3.7451940309632633e+306))
	assert.Equal(t, 8.333333333333334e-10, lgammacor(10e+07))
}

func TestCalcLgammaerr(t *testing.T) {
	assert.True(t, math.IsNaN(logrelerr(-2)))
}

func TestCalcLogBeta(t *testing.T) {
	assert.True(t, math.IsNaN(logBeta(-1, -1)))
	assert.Equal(t, math.MaxFloat64, logBeta(0, 0))
}

func TestCalcBetainvProbIterator(t *testing.T) {
	assert.Equal(t, 1.0, betainvProbIterator(1, 1, 1, 1, 1, 1, 1, 1, 1))
}

func TestNestedFunctionsWithOperators(t *testing.T) {
	f := NewFile()
	formulaList := map[string]string{
		`=LEN("KEEP")`:                                               "4",
		`=LEN("REMOVEKEEP") - LEN("REMOVE")`:                         "4",
		`=RIGHT("REMOVEKEEP", 4)`:                                    "KEEP",
		`=RIGHT("REMOVEKEEP", 10 - 6))`:                              "KEEP",
		`=RIGHT("REMOVEKEEP", LEN("REMOVEKEEP") - 6)`:                "KEEP",
		`=RIGHT("REMOVEKEEP", LEN("REMOVEKEEP") - LEN("REMOV") - 1)`: "KEEP",
		`=RIGHT("REMOVEKEEP", 10 - LEN("REMOVE"))`:                   "KEEP",
		`=RIGHT("REMOVEKEEP", LEN("REMOVEKEEP") - LEN("REMOVE"))`:    "KEEP",
	}
	for formula, expected := range formulaList {
		assert.NoError(t, f.SetCellFormula("Sheet1", "E1", formula))
		result, err := f.CalcCellValue("Sheet1", "E1")
		assert.NoError(t, err, formula)
		assert.Equal(t, expected, result, formula)
	}
}
