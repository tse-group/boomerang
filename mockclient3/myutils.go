package main

import (
	"log"
	"strconv"
	"math/rand"
)


// MY UTILITIES

// abstract float type ...

type float float64


// parsing numeric types

func myparseint(s string) int {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return int(v)
	} else {
		log.Panicln("Error parsing to int:", s)
	}
	return 0
}

func myparsefloat32(s string) float32 {
	if v, err := strconv.ParseFloat(s, 32); err == nil {
		return float32(v)
	} else {
		log.Panicln("Error parsing to float32:", s)
	}
	return 0.0
}

func myparsefloat(s string) float {
	if v, err := strconv.ParseFloat(s, 64); err == nil {
		return float(v)
	} else {
		log.Panicln("Error parsing to float64:", s)
	}
	return 0.0
}


// routine functions: int

func mymini(v1 int, v2 int) int {
	if v1 <= v2 {
		return v1
	} else {
		return v2
	}
}

func myreversei(numbers []int) []int {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
	return numbers
}

func myreversei_outofplace(in []int) []int {
	var out []int
	for i := len(in)-1; i >= 0; i-- {
		out = append(out, in[i])
	}
	return out
}

func myshufflei(lst []int) {
	rand.Shuffle(len(lst), func(i int, j int) { lst[i], lst[j] = lst[j], lst[i] })
}

func mypopi(slice []int) ([]int, int) {
	if len(slice) == 0 {
		panic("Cannot pop from empty slice")
	} else {
		return slice[0:len(slice)-1], slice[len(slice)-1]
	}
}

func myfirsti(vs []int) int {
	return vs[0]
}

func mylasti(vs []int) int {
	return vs[len(vs)-1]
}

func myfindi(lst []int, val int) (int, bool) {
	for i, v := range lst {
		if v == val {
			return i, true
		}
	}
	return -1, false
}


// routine functions: float

func myminf32(v1 float32, v2 float32) float32 {
	if v1 <= v2 {
		return v1
	} else {
		return v2
	}
}

func myminf(v1 float, v2 float) float {
	if v1 <= v2 {
		return v1
	} else {
		return v2
	}
}

func mypopf32(slice []float32) ([]float32, float32) {
	if len(slice) == 0 {
		panic("Cannot pop from empty slice")
	} else {
		return slice[0:len(slice)-1], slice[len(slice)-1]
	}
}

func myfirstf32(vs []float32) float32 {
	return vs[0]
}

func mylastf32(vs []float32) float32 {
	return vs[len(vs)-1]
}

func mypopf(slice []float) ([]float, float) {
	if len(slice) == 0 {
		panic("Cannot pop from empty slice")
	} else {
		return slice[0:len(slice)-1], slice[len(slice)-1]
	}
}

func myfirstf(vs []float) float {
	return vs[0]
}

func mylastf(vs []float) float {
	return vs[len(vs)-1]
}


// random number in range

func myrandrangei(frm int, to int) int {
	return frm + rand.Intn(to - frm + 1)
}


// assert is always good to have

func myassert(exp bool, str string) {
	if !exp {
		log.Panicln("Assertion error:", str)
	}
}

