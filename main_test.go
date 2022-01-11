package main

import (
	"fmt"
	"github.com/gonum/floats"
	"testing"
)

func TestConvertList(t *testing.T) {
	ans := ConvertList("0.028501 1.59E-01 0.026755")
	for i, tt := range []float64{0.028501, 1.59E-01, 0.026755} {
		testname := fmt.Sprintf("%f", tt)
		t.Run(testname, func(t *testing.T) {
			tmpAns := ans[i]
			if tmpAns != tt {
				t.Errorf("got %f, want %f", tmpAns, tt)
			}
		})
	}
}

func TestNormalizeVector(t *testing.T) {
	vec := []float64{1., 1.}
	expected := []float64{0.5, 0.5}
	res := NormalizeVector(vec)
	if !floats.Equal(res, expected) {
		t.Errorf("got %v, want %v", res, expected)
	}
}

//[ 48.26724806 -24.28220614   7.43963471   0.79681006   7.61531982   -1.15476944  -1.78742913]
