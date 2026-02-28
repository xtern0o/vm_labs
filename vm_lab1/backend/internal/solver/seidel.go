package solver

import (
	"fmt"
	"vm_lab1/internal/matrix"
)

type SeidelSolver struct{}

func NewSeidelSolver() *SeidelSolver {
	return &SeidelSolver{}
}

func (s *SeidelSolver) Solve(
	A matrix.Matrix,
	b matrix.Vector,
	x0 matrix.Vector,
	eps float64,
	maxIter int,
) (Result, error) {
	n := len(A)

	if len(b) != n {
		return Result{}, fmt.Errorf("неподходящие измерения: %d для A != %d для B", len(A), len(b))
	}
	if len(x0) != n {
		return Result{}, fmt.Errorf("начальное приближение имеет размерность %d != %d у матрицы A", len(x0), len(A))
	}

	ACopy := matrix.CopyMatrix(A)
	bCopy := matrix.CopyVector(b)

	if !matrix.TryToMakeDiagonallyDominant(ACopy, bCopy) {
		return Result{}, fmt.Errorf("матрицу A не удалось привести к диагональному преобладанию")
	}

	C, d, err := matrix.BuildCanonicalForm(ACopy, bCopy)
	if err != nil {
		return Result{}, err
	}

	normC := matrix.NormOfMatrix(C)
	if normC >= 1 {
		return Result{}, fmt.Errorf("||C|| = %.6f >= 1, сходимость не гарантирована", normC)
	}

	x := matrix.CopyVector(x0)   // текущее решение
	xPrev := matrix.NewVector(n) // прошлое решение
	errors := []float64{}        // вектор погрешностей

	for iter := 0; iter < maxIter; iter++ {
		copy(xPrev, x)

		for i := 0; i < n; i++ {
			sum := d[i]

			// подстановка новых приближений (j < i)
			for j := 0; j < i; j++ {
				sum += C[i][j] * x[j]
			}

			// подстановка старых приближений (j > i)
			for j := i + 1; j < n; j++ {
				sum += C[i][j] * x[j]
			}

			x[i] = sum
		}

		diff, _ := matrix.SubVectors(x, xPrev)
		normVal := matrix.MaxNormOfVector(diff)
		errors = append(errors, normVal)

		if normVal < eps {
			return Result{
				Solution:   x,
				Iterations: iter + 1,
				Errors:     errors,
			}, nil
		}
	}

	return Result{}, fmt.Errorf("не сошлось за %d итераций. Последняя погрешность: %.6f", maxIter, errors[len(errors)-1])

}
