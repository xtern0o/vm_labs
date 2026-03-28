package service

import (
	"fmt"
	"math"
	"vm_lab3/internal/dto"
)

type SimpsonSolver struct{}

func NewSimpsonSolver() *SimpsonSolver {
	return &SimpsonSolver{}
}

func (s SimpsonSolver) solve(
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

	k := 4.0
	R := math.Inf(0)
	var IPrev, ICurr float64
	for R > eps {
		IPrev = solveSimpsonMethod(f, a, b, n, &messages)
		n *= 2
		ICurr = solveSimpsonMethod(f, a, b, n, &messages)
		R = CalcR(IPrev, ICurr, k)
	}

	return dto.CalcIntegralResponseDto{
		Value:    ICurr,
		N:        n,
		Messages: messages,
		RungeR:   R,
	}, nil

}

func solveSimpsonMethod(
	f func(float64) float64,
	a, b float64,
	n int,
	messages *[]string,
) float64 {
	if n%2 != 0 {
		*messages = append(*messages, fmt.Sprintf("n = n+1 = %d - n should be even", n))
		n++
	}
	h := (b - a) / float64(n)
	sum := f(a) + f(b)
	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		coef := 4.0
		if i%2 == 0 {
			coef = 2.0
		}
		sum += coef * f(x)
	}
	return h * sum / 3.0
}
