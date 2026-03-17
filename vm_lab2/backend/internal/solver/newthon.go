package solver

import (
	"errors"
	"fmt"
	"math"
	"vm_lab2/internal/equations"
	"vm_lab2/internal/numeric"
)

const minDerivative = 1e-12

type NewthonSolver struct{}

func NewNewthonSolver() *NewthonSolver {
	return &NewthonSolver{}
}

func (s *NewthonSolver) Solve(
	f equations.Equation,
	a float64,
	b float64,
	eps float64,
	maxIter int,
) (Result, error) {
	messages := make([]string, 0)

	if a > b {
		a, b = b, a
		messages = append(messages, fmt.Sprintf("a and b were swapped (a <= b): a=%f, b=%f", a, b))
	}

	err := BasicCheck(f, a, b)
	if err != nil {
		return Result{}, err
	}

	var xPrev float64

	secondDerA := numeric.SecondDerivative(f, a)
	secondDerB := numeric.SecondDerivative(f, b)

	if f(a)*secondDerA > 0 {
		xPrev = a
		messages = append(messages, "x_0 = a, так как f(a)f''(a)>0")
	} else if f(b)*secondDerB > 0 {
		xPrev = b
		messages = append(messages, "x_0=b, так как f(b)f''(b)>0")
	} else {
		xPrev = (a + b) / 2
		messages = append(messages, fmt.Sprintf("быстрая сходимость не гарантирована, за начальное приближение выбрана точка x_0=a. f''(a)=%f, f''(b)=%f", secondDerA, secondDerB))
	}

	xCurr, err := getNextX(f, xPrev)
	if err != nil {
		return Result{}, err
	}

	steps := []Point{}
	steps = append(steps, Point{xCurr, f(xCurr)})

	iter := 0
	for ; iter < maxIter && math.Abs(xCurr-xPrev) > eps && math.Abs(f(xCurr)) > eps; iter++ {
		xPrev = xCurr
		xCurr, err = getNextX(f, xCurr)
		if err != nil {
			return Result{}, err
		}
		steps = append(steps, Point{xCurr, f(xCurr)})

		if CheckAns(f, xCurr, eps) {
			messages = append(messages, "достигли нужной точности")
			break
		}

	}

	if iter >= maxIter && math.Abs(xCurr-xPrev) > eps && math.Abs(f(xCurr)) > eps {
		return Result{}, errors.New("iter limit")
	}

	return Result{
		Solution:   xCurr,
		Value:      f(xCurr),
		Iterations: iter,
		ArgError:   math.Abs(xCurr - xPrev),
		Messages:   messages,
		Steps:      steps,
	}, nil
}

func getNextX(f equations.Equation, x float64) (float64, error) {
	deriv := numeric.Derivative(f, x)
	if math.Abs(deriv) <= minDerivative {
		return 0, fmt.Errorf("производная слишком близка к 0 в точке x=%f", x)
	}

	return x - (f(x) / deriv), nil
}
