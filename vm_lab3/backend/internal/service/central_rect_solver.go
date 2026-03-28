package service

import (
	"math"
	"vm_lab3/internal/dto"
)

type CentralRectSolver struct{}

func NewCentralRectSolver() *CentralRectSolver {
	return &CentralRectSolver{}
}

func (s CentralRectSolver) solve(
	f func(x float64) float64,
	a, b, eps float64,
	n int,
) (dto.CalcIntegralResponseDto, error) {
	if err := BasicCheck(f, a, b, eps, n); err != nil {
		return dto.CalcIntegralResponseDto{}, err
	}

	messages := make([]string, 0)
	if a > b {
		messages = append(messages, "a and b were swapped")
		a, b = b, a
	}

	k := 2.0
	R := math.Inf(0)
	var IPrev, ICurr float64
	for R > eps {
		IPrev = solveCentralRectMethod(f, a, b, n)
		n *= 2
		ICurr = solveCentralRectMethod(f, a, b, n)
		R = CalcR(IPrev, ICurr, k)
	}

	return dto.CalcIntegralResponseDto{
		Value:    ICurr,
		N:        n,
		Messages: messages,
		RungeR:   R,
	}, nil
}

func solveCentralRectMethod(
	f func(x float64) float64,
	a, b float64,
	n int,
) float64 {
	h := (b - a) / float64(n)
	I := 0.0
	for i := 1; i < n+1; i++ {
		x := a + h*(float64(i-1)+0.5)
		y_i := f(x)
		I += h * y_i
	}
	return I
}
