package service

import (
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

	h := (b - a) / float64(n)
	res := 0.0
	for i := 1; i < n+1; i++ {
		y_i := f(h * (float64(i-1) + 0.5))
		res += h * y_i
	}

	return dto.CalcIntegralResponseDto{
		Value:    res,
		N:        n,
		Messages: messages,
	}, nil
}
