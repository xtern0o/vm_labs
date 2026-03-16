package solver

import (
	"errors"
	"fmt"
	"math"
	"vm_lab2/internal/equations"
	"vm_lab2/internal/numeric"
)

type SystemSimpleIterSolver struct{}

func NewSystemSimpleIterSolver() *SystemSimpleIterSolver {
	return &SystemSimpleIterSolver{}
}

func (s *SystemSimpleIterSolver) Solve(
	system equations.System2,
	phiSystem equations.System2,
	x0 float64,
	y0 float64,
	eps float64,
	maxIter int,
) (System2Result, error) {
	f1 := system[0]
	f2 := system[1]
	if f1 == nil || f2 == nil {
		return System2Result{}, errors.New("system equations should not be nil")
	}

	phi1 := phiSystem[0]
	phi2 := phiSystem[1]

	if phi1 == nil || phi2 == nil {
		return System2Result{}, errors.New("phi functions should not be nil")
	}

	messages := []string{}

	// проверка условной сходимости в начальной точке
	dPhi1Dx := numeric.Derivative(
		func(v float64) float64 {
			return phi1(v, y0)
		}, x0,
	)
	dPhi1Dy := numeric.Derivative(
		func(v float64) float64 {
			return phi1(x0, v)
		}, y0,
	)
	dPhi2Dx := numeric.Derivative(
		func(v float64) float64 {
			return phi2(v, y0)
		}, x0,
	)
	dPhi2Dy := numeric.Derivative(
		func(v float64) float64 {
			return phi2(x0, v)
		}, y0,
	)
	q1 := math.Abs(dPhi1Dx) + math.Abs(dPhi1Dy)
	q2 := math.Abs(dPhi2Dx) + math.Abs(dPhi2Dy)
	q := math.Max(q1, q2)
	if math.IsNaN(q) || math.IsInf(q, 0) {
		messages = append(messages, "не удалось оценить достаточное условие сходимости (q нечисловое)")
	} else if q >= 1 {
		messages = append(messages, fmt.Sprintf("достаточное условие сходимости не выполнено: q=%.6f >= 1; пробуем решить!", q))
	}

	steps := []Point{}
	errPoints := []Point{}

	xPrev, yPrev := x0, y0
	xCurr, yCurr := phi1(x0, y0), phi2(x0, y0)
	if math.IsNaN(xCurr) || math.IsInf(xCurr, 0) || math.IsNaN(yCurr) || math.IsInf(yCurr, 0) {
		return System2Result{}, fmt.Errorf("итерации вышли за зону определения функции phi на первом шаге: x=%f, y=%f", xCurr, yCurr)
	}

	steps = append(steps, Point{xCurr, yCurr})
	errPoints = append(errPoints, Point{math.Abs(xCurr - xPrev), math.Abs(yCurr - yPrev)})

	iter := 1
	for ; iter < maxIter && math.Max(math.Abs(xPrev-xCurr), math.Abs(yPrev-yCurr)) > eps && !checkSystemAns(system, xCurr, yCurr, eps); iter++ {
		xPrev, yPrev = xCurr, yCurr
		xCurr, yCurr = phi1(xPrev, yPrev), phi2(xPrev, yPrev)
		if math.IsNaN(xCurr) || math.IsInf(xCurr, 0) || math.IsNaN(yCurr) || math.IsInf(yCurr, 0) {
			return System2Result{}, fmt.Errorf("итерации вышли за зону определения функции phi на шаге %d: xCurr=%f, yCurr=%f, xPrev=%f, yPrev=%f", iter, xCurr, yCurr, xPrev, yPrev)
		}

		steps = append(steps, Point{xCurr, yCurr})
		errPoints = append(errPoints, Point{math.Abs(xCurr - xPrev), math.Abs(yCurr - yPrev)})

	}

	if iter >= maxIter && math.Max(math.Abs(xPrev-xCurr), math.Abs(yPrev-yCurr)) > eps {
		return System2Result{}, errors.New("iter limit")
	}

	return System2Result{
		Solution:   Point{xCurr, yCurr},
		Iterations: iter,
		Messages:   messages,
		Steps:      steps,
		Errors:     errPoints,
	}, nil

}

func checkSystemAns(s equations.System2, x float64, y float64, eps float64) bool {
	f1Val := s[0](x, y)
	f2Val := s[1](x, y)
	return math.Abs(f1Val) < eps && math.Abs(f2Val) < eps
}
