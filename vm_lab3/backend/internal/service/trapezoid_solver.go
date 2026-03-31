package service

import (
	"fmt"
	"math"
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

	k := 2.0
	R := math.Inf(0)
	var IPrev, ICurr float64

	iter := 0
	for R > eps {
		iter++
		prevCh := make(chan float64)
		currCh := make(chan float64)

		go func() {
			prevCh <- solveTrapezoidMethod(f, a, b, n)
		}()

		n2 := n * 2
		go func() {
			currCh <- solveTrapezoidMethod(f, a, b, n2)
		}()

		IPrev = <-prevCh
		ICurr = <-currCh

		messages = append(messages, fmt.Sprintf("iter: %d. I_0 = %f (n=%d)", iter, IPrev, n))
		messages = append(messages, fmt.Sprintf("iter: %d. I_1 = %f (n=%d)", iter, ICurr, n2))

		n = n2
		R = CalcR(IPrev, ICurr, k)

		messages = append(messages, fmt.Sprintf("R = %f", R))

	}

	return dto.CalcIntegralResponseDto{
		Value:    ICurr,
		N:        n,
		Messages: messages,
		RungeR:   R,
	}, nil
}

func solveTrapezoidMethod(
	f func(x float64) float64,
	a, b float64,
	n int,
) float64 {
	h := (b - a) / float64(n)
	y0 := f(a)
	yn := f(b)
	sum := 0.0
	for i := 1; i < n; i++ {
		x := a + float64(i)*h
		sum += f(x)
	}
	return (h / 2) * (y0 + yn + 2*sum)
}
