package solver

import (
	"fmt"
	"math"
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

	messages := make([]string, 0)

	if !matrix.TryToMakeDiagonallyDominant(ACopy, bCopy) {
		messages = append(messages, "Диагонального преобладания нет. Сходимость НЕ гарантирована")
	}

	fmt.Println("- Матрица А с преобладанием: ")
	for i := 0; i < len(ACopy); i++ {
		for j := 0; j < len(ACopy); j++ {
			fmt.Printf("%f ", ACopy[i][j])
		}
		fmt.Println()
	}

	C, d, err := matrix.BuildCanonicalForm(ACopy, bCopy)
	if err != nil {
		return Result{}, err
	}

	normC := matrix.NormOfMatrix(C)
	if normC >= 1 {
		fmt.Println("Норма больше 1")
		messages = append(messages, "||C|| >= 1. Сходимость не гаранитирована")
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
			fmt.Println("- невязка с конечным вектором:")
			final := make([]float64, len(x))
			for i := 0; i < len(x); i++ {
				sum := 0.0
				for j := 0; j < len(x); j++ {
					sum += ACopy[i][j] * x[j]
				}
				final[i] = math.Abs(sum - bCopy[i])
				fmt.Println(final[i])
			}

			return Result{
				Solution:     x,
				Iterations:   iter + 1,
				Errors:       errors,
				NormOfMatrix: normC,
				Messages:     messages,
			}, nil
		}
	}

	return Result{}, fmt.Errorf("не сошлось за %d итераций. Последняя погрешность: %.6f", maxIter, errors[len(errors)-1])

}
