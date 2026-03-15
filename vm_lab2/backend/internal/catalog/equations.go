package catalog

import (
	"math"
	"vm_lab2/internal/equations"
)

// x^2 + e^x - 10x sinx - 5x
var eq1 equations.Equation = func(x float64) float64 {
	return x*x + math.Exp(x) - 10*x*math.Sin(x) - 5*x
}

// ln(x) + (sinx)^3
var eq2 equations.Equation = func(x float64) float64 {
	return math.Log(x) + math.Pow(math.Sin(x), 3)
}

// sqrt(x) - 0.07x^3 + 0.5x^2 - 5
var eq3 equations.Equation = func(x float64) float64 {
	return math.Sqrt(x) - 0.07*x*x*x + 0.5*x*x - 5
}

// x^3 + 4x^2 + 3x - 1
var eq4 equations.Equation = func(x float64) float64 {
	return x*x*x + 4*x*x + 3*x - 1
}

// -0.4x^3 + 2x^2 - 2
var eq5 equations.Equation = func(x float64) float64 {
	return -0.4*x*x*x + 2*x*x - 2
}

// иммутабельная мапа с доступными уравнениями
var equationsMap = map[int]equations.Equation{
	1: eq1,
	2: eq2,
	3: eq3,
	4: eq4,
	5: eq5,
}

func GetEquations() map[int]equations.Equation {
	return equationsMap
}
