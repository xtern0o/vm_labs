package service

import (
	"math"
	"vm_lab3/internal/dto"
)

type LeftRectSolver struct{}

func NewLeftRectSolver() *LeftRectSolver {
	return &LeftRectSolver{}
}

func (s LeftRectSolver) solve(
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
		prevCh := make(chan float64)
		currCh := make(chan float64)

		go func() {
			prevCh <- solveLeftRectMethod(f, a, b, n)
		}()

		n2 := n * 2
		go func() {
			currCh <- solveLeftRectMethod(f, a, b, n)
		}()

		IPrev = <-prevCh
		ICurr = <-currCh

		n = n2
		R = CalcR(IPrev, ICurr, k)
	}

	return dto.CalcIntegralResponseDto{
		Value:    ICurr,
		N:        n,
		Messages: messages,
		RungeR:   R,
	}, nil
}

func solveLeftRectMethod(
	f func(x float64) float64,
	a, b float64,
	n int,
) float64 {
	h := (b - a) / float64(n)
	I := 0.0
	for i := 0; i < n; i++ {
		x := a + float64(i)*h
		y_i := f(x)
		I += h * y_i
	}
	return I
}
