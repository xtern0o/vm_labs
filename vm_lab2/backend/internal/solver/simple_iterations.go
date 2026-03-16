package solver

import (
	"errors"
	"fmt"
	"math"
	"vm_lab2/internal/equations"
	"vm_lab2/internal/numeric"
)

type SimpleIterationsSolver struct{}

func NewSimpleIterationsSolver() *SimpleIterationsSolver {
	return &SimpleIterationsSolver{}
}

func (s *SimpleIterationsSolver) Solve(
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

	derivativeFA := numeric.Derivative(f, a)
	derivativeFB := numeric.Derivative(f, b)

	lambda := 1 / math.Max(derivativeFA, derivativeFB)
	if derivativeFA > 0 && derivativeFB > 0 {
		lambda *= -1
	}

	phi := func(x float64) float64 {
		return x + lambda*f(x)
	}
	if numeric.Derivative(phi, a) >= 1 || numeric.Derivative(phi, b) >= 1 {
		messages = append(messages, fmt.Sprintf("достаточное условие сходимости не выполнено - сходимость не гарантирована, q = %f > 1. Пробуем вычислить корень", math.Max(derivativeFA, derivativeFB)))
	}

	steps := []Point{}

	xPrev := a
	xCurr := phi(xPrev)

	iter := 0
	for ; iter < maxIter && math.Abs(xCurr-xPrev) > eps; iter++ {
		xPrev = xCurr
		xCurr = phi(xCurr)

		steps = append(steps, Point{xCurr, f(xCurr)})
	}
	if iter >= maxIter && math.Abs(xCurr-xPrev) > eps {
		return Result{}, errors.New("iter limit")
	}
	if math.IsNaN(xCurr) || math.IsInf(xPrev, 0) {
		return Result{}, fmt.Errorf("функция phi не определена на %f. Метод не сходится", xPrev)
	}
	if math.IsNaN(f(xCurr)) {
		return Result{}, fmt.Errorf("функция не определена на %f. Метод не сходится", xCurr)
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
