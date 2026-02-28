package solver

import "vm_lab1/internal/matrix"

type Result struct {
	Solution     matrix.Vector
	Iterations   int
	Errors       []float64
	NormOfMatrix float64
}

type Solver interface {
	Solve(
		A matrix.Matrix, // коэффициенты при x
		b matrix.Vector, // значения справа
		x0 matrix.Vector, // начальное приближение
		eps float64, // точность
		maxIter int, // предел итераций
	) (Result, error)
}
