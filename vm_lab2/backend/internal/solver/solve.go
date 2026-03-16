package solver

import (
	"vm_lab2/internal/equations"
)

func SolveEquation(
	f equations.Equation,
	s EquationSolver,
	a float64,
	b float64,
	eps float64,
	maxIter int,
) (Result, error) {
	return s.Solve(f, a, b, eps, maxIter)
}

func SolveSystem2(
	system equations.System2,
	phiSystem equations.System2,
	s Systen2Solver,
	x0 float64,
	y0 float64,
	eps float64,
	maxIter int,
) (System2Result, error) {
	return s.Solve(system, phiSystem, x0, y0, eps, maxIter)
}
