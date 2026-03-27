package service

import (
	"vm_lab3/internal/dto"
)

type TrapezoidRectSolver struct{}

func NewTrapezoidRectSolver() *TrapezoidRectSolver {
	return &TrapezoidRectSolver{}
}

func (s TrapezoidRectSolver) solve(
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
	y0 := f(a)
	yn := f(b)
	sum := 0.0
	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		sum += f(x)
	}
	res := (h / 2) * (y0 + yn + 2*sum)

	return dto.CalcIntegralResponseDto{
		Value:    res,
		N:        n,
		Messages: messages,
	}, nil
}
