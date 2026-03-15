package numeric

import (
	"math"
	"vm_lab2/internal/equations"
)

const defaultStep = 1e-6

// функция для расчета значения производной в точке с точностью из конфигурации
func Derivative(f equations.Equation, a float64) float64 {
	if f == nil {
		return math.NaN()
	}

	h := defaultStep * (1 + math.Abs(a))
	return (f(a+h) - f(a-h)) / (2 * h)
}

// функция для расчета значения второй производной в точке
func SecondDerivative(f equations.Equation, a float64) float64 {
	if f == nil {
		return math.NaN()
	}

	h := defaultStep * (1 + math.Abs(a))
	return (f(a+h) - 2*f(a) + f(a-h)) / (h * h)
}
