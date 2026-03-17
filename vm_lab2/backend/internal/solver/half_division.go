package solver

import (
	"errors"
	"fmt"
	"math"
	"vm_lab2/internal/equations"
)

type HalfDivisionSolver struct{}

func NewHalfDivisionSolver() *HalfDivisionSolver {
	return &HalfDivisionSolver{}
}

func (s *HalfDivisionSolver) Solve(
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

	steps := []Point{}

	x := (a + b) / 2
	steps = append(steps, Point{x, f(x)})
	iter := 1
	for {
		if math.Abs(a-b) <= eps || math.Abs(f(x)) < eps {
			break
		}
		if f(a)*f(x) > 0 {
			a = x
		} else {
			b = x
		}
		x = (a + b) / 2

		fx := f(x)
		if math.IsNaN(fx) || math.IsInf(fx, 0) {
			return Result{}, fmt.Errorf("func is not defined at x=%f", x)
		}
		steps = append(steps, Point{x, fx})

		if CheckAns(f, x, eps) {
			messages = append(messages, "достигли нужной точности")
			break
		}

		iter++
		if iter > maxIter {
			return Result{}, errors.New("iter limit")
		}
	}

	return Result{
		Solution:   x,
		Value:      f(x),
		Iterations: iter,
		ArgError:   math.Abs(a - b),
		Messages:   messages,
		Steps:      steps,
	}, nil

}
