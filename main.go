package main

import (
	"flag"
	"fmt"
	"github.com/gonum/floats"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func ConvertList(str string) []float64 {
	// get a string and try to convert it to array of float64
	res := []float64{}
	sp := []string{}
	if strings.Contains(str, ",") {
		sp = strings.Split(str, ",")
	} else {
		sp = strings.Split(str, " ")
	}
	for _, num := range sp {
		convNum, err := strconv.ParseFloat(num, 64)
		if err != nil {
			print(err)
		}
		res = append(res, convNum)
	}
	return res
}

func ComputeLDA(v1 []float64, weights [][]float64) []float64 {
	// preforms the dot product for LDA
	res := []float64{}
	for _, row := range weights {
		val := floats.Dot(v1, row)
		res = append(res, val)
	}
	return res
}

func NormalizeVector(v []float64) []float64 {
	res := []float64{}
	sum := floats.Sum(v)
	for _, val := range v {
		res = append(res, val/sum)
	}
	return res
}

func FormatFloat(num float64, prc int) string {
	var (
		zero, dot = "0", "."
	)
	str := ""
	if prc == -1 {
		str = fmt.Sprintf("%f", num)
	} else {
		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	}

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

func RoundVector(v []float64, prc int) []string {
	res := []string{}
	for _, val := range v {
		res = append(res, FormatFloat(val, prc))
	}
	return res
}

func ComputeDelta(v []float64, means [][]float64) []float64 {
	res := []float64{}
	for i, val := range v {
		res = append(res, val-means[i][0])
	}
	return res
}

func EncodeParams(v []string, ApiUrl string) string {
	params := url.Values{}
	for i, val := range v {
		params.Add("v"+strconv.Itoa(i), val)
	}
	return ApiUrl + params.Encode()
}

func ReadWeights(path string) [][]float64 {
	fileBytes, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sliceData := strings.Split(string(fileBytes), "\n")

	res := [][]float64{}
	for _, tmp := range sliceData {
		convLst := ConvertList(tmp)
		if len(tmp) > 0 {
			res = append(res, convLst)
		}
	}
	return res
}

func main() {
	DataPath := flag.String("DataPath", "data.csv", "csv file with samples, each row represents one sample")
	ScalingPath := flag.String("scaling-path", "scalings.csv", "path for the scalings (weights of LDA)")
	XbarPath := flag.String("xbar-path", "xbar.csv", "path for the xbars (weights of LDA)")
	NumOfFreqs := flag.Int("freqs", 680, "number of frequencies")
	cli := flag.Bool("cli", false, "print QR code in cli")
	DisableImg := flag.Bool("disable-image", false, "don't save image")
	verbose := flag.Bool("v", false, "verbose")
	heroku := flag.Bool("web", false, "heroku url")
	Round := flag.Int("round", -1, "The number of decimals to use when rounding the number")
	OutputPath := flag.String("out-path", "./output/", "path for output, should end with /")
	ApiUrl := flag.String("url", "", "url for the API")
	// image path
	flag.Parse()

	ScalingMat := ReadWeights(*ScalingPath)
	XbarMat := ReadWeights(*XbarPath)
	sampleText := ReadWeights(*DataPath)
	if *verbose {
		fmt.Println("\n@#=====! configuration !=====#@")
		fmt.Println("[data] samples:", len(sampleText), "frequencies:", len(sampleText[0]))
		fmt.Println("[scaling] components:", len(ScalingMat)+1, "frequencies:", len(ScalingMat[0]))
		fmt.Println("[xbar] components:", len(XbarMat))
		fmt.Println("@#===========================#@")
		fmt.Println("@#=====! results !=====#@")
	}
	for i, sample := range sampleText {
		vector := NormalizeVector(sample[:*NumOfFreqs])
		vector = ComputeDelta(vector, XbarMat)
		LdaRes := ComputeLDA(vector, ScalingMat)
		StrRes := RoundVector(LdaRes, *Round)
		content := fmt.Sprintf("%v", StrRes)
		if *verbose {
			fmt.Println(i, " ==> ", content)
		}
		if *heroku {
			content = EncodeParams(StrRes, "https://p435.herokuapp.com/predict?")
			if *verbose {
				fmt.Println("url ", i, " --> ", content)
			}
		}
		if *ApiUrl != "" {
			content = EncodeParams(StrRes, *ApiUrl)
			if *verbose {
				fmt.Println("url ", i, " --> ", content)
			}
		}
		if *cli {
			var q *qrcode.QRCode
			q, _ = qrcode.New(content, qrcode.Highest)
			fmt.Println()
			fmt.Println(q.ToString(false))
		}
		if *DisableImg == false {
			err := os.MkdirAll(*OutputPath, os.ModePerm)
			if err != nil {
				fmt.Println(err)
			}
			if err := qrcode.WriteFile(content, qrcode.Medium, 256, *OutputPath+strconv.Itoa(i)+".png"); err != nil {
				fmt.Print("problem!!!")
				panic(err)
			}
		}
	}
}
