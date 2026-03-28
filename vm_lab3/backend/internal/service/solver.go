package service

import (
	"fmt"
	"math"
	"vm_lab3/internal/dto"
)

const (
	divCheck = 100
)

type IntegralSolver interface {
	solve(f func(x float64) float64, a, b, eps float64, n int) (dto.CalcIntegralResponseDto, error)
}

func BasicCheck(f func(x float64) float64, a, b, eps float64, n int) error {
	if f == nil {
		return fmt.Errorf("func cant be nil")
	}
	if eps <= 0 {
		return fmt.Errorf("eps must be more than 0")
	}
	if a == b {
		return fmt.Errorf("a must not be equal to b")
	}
	if n <= 0 {
		return fmt.Errorf("n must be more than 0")
	}

	h := (b - a) / float64(divCheck)
	for i := 0; i <= divCheck; i++ {
		x := a + h*float64(i)
		y := f(x)
		if math.IsInf(y, 0) || math.IsNaN(y) {
			return fmt.Errorf("f(x) не определена корректно при x=%g (бесконечность или nan)", x)
		}
	}

	return nil
}

// расчет коэффициента R по правилу Рунге
func CalcR(I1, I2, k float64) float64 {
	return math.Abs(I1-I2) / (math.Pow(2.0, k) - 1)
}
