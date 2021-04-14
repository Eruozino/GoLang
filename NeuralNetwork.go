package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	a10 = 0.1
	a11 = 0.3
	a12 = 0.4
	a20 = 0.5
	a21 = 0.8
	a22 = 0.3
	a30 = 0.7
	a31 = 0.6
	a32 = 0.6
	b10 = 0.5
	b11 = 0.3
	b12 = 0.7
	b13 = 0.1
)

// func hiddenNode( )

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func directOutput(X1, X2 float64) float64 {
	Z1 := sigmoid(a10 + a11*X1 + a12*X2)
	Z2 := sigmoid(a20 + a21*X1 + a22*X2)
	Z3 := sigmoid(a30 + a31*X1 + a32*X2)
	T1 := sigmoid(b10 + b11*Z1 + b12*Z2 + b13*Z3)
	return T1
}

func main() {
	N, err := strconv.Atoi(os.Args[1]) // 16, false
	if err != nil {
		fmt.Print("Not a number: ")
		fmt.Println(os.Args[1])
		fmt.Println(err)
		return
	}
	for k := 0; k < N; k++ {
		X1 := math.Sin(2.0 * math.Pi * float64(k) / float64(N))
		X2 := math.Cos(2.0 * math.Pi * float64(k) / float64(N))
		fmt.Printf("%f %f -> %f\n", X1, X2, directOutput(X1, X2))
	}
}
