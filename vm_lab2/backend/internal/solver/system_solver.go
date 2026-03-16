package solver

import "vm_lab2/internal/equations"

type Systen2Solver interface {
	Solve(
		s equations.System2,
		phiSystem equations.System2,
		x0 float64,
		y0 float64,
		eps float64,
		maxIter int,
	) (System2Result, error)
}
