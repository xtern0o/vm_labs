package matrix

import (
	"errors"
	"fmt"
	"math"
)

type Vector []float64
type Matrix []Vector

// создание вектора
func NewVector(n int) Vector {
	return make(Vector, n)
}

// создание матрицы
func NewMatrix(n, m int) Matrix {
	mat := make(Matrix, n)
	for i := range mat {
		mat[i] = NewVector(m)
	}
	return mat
}

// для парсинга матрицы из двумерного слайса
func MatrixFromFloat64(m [][]float64) Matrix {
	mat := make(Matrix, len(m))
	for i := range mat {
		mat[i] = Vector(m[i])
	}
	return mat
}

// скопировать вектор
func CopyVector(v Vector) Vector {
	res := make(Vector, len(v))
	copy(res, v)
	return res
}

// скопировать матрицу
func CopyMatrix(mat Matrix) Matrix {
	res := make(Matrix, len(mat))
	for i := range mat {
		res[i] = CopyVector(mat[i])
	}
	return res
}

// добавить число к вектору
func AddNumToVector(v Vector, n float64) {
	for i := range v {
		v[i] += n
	}
}

// добавить число к матрице
func AddNumToMatrix(mat Matrix, n float64) {
	for i := range mat {
		AddNumToVector(mat[i], n)
	}
}

// A - B
func SubVectors(a, b Vector) (Vector, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("невозможно сложить: dim(a) = %d != dim(b) = %d", len(a), len(b))
	}
	c := NewVector(len(a))
	for i := range c {
		c[i] = a[i] - b[i]
	}
	return c, nil
}

// вычисление нормы вектора
func NormVec(v Vector) float64 {
	res := 0.0
	for _, val := range v {
		res += val * val
	}
	return math.Sqrt(res)
}

// норма матрицы
func NormOfMatrix(A Matrix) float64 {
	maxSum := 0.0
	for i := range A {
		rowSum := 0.0
		for j := range A[i] {
			rowSum += math.Abs(A[i][j])
		}
		if rowSum > maxSum {
			maxSum = rowSum
		}
	}
	return maxSum
}

// максимальный элемент вектора по модулю
func MaxNormOfVector(v Vector) float64 {
	maxVal := 0.0
	for _, val := range v {
		if abs := math.Abs(val); abs > maxVal {
			maxVal = abs
		}
	}
	return maxVal
}

// проверка, есть ли 0 на диагонали
func HasZeroOnDiagonal(A Matrix) bool {
	for i := range A {
		if A[i][i] == 0 {
			return true
		}
	}
	return false
}

// приводит систему Ax = B к системе вида x = Cx + d
func BuildCanonicalForm(A Matrix, b Vector) (C Matrix, d Vector, e error) {
	n := len(A)

	if HasZeroOnDiagonal(A) {
		return nil, nil, errors.New("0 на диагонали")
	}

	C = NewMatrix(n, n)
	d = NewVector(n)

	for i := 0; i < n; i++ {
		d[i] = b[i] / A[i][i]

		for j := 0; j < n; j++ {
			if i == j {
				C[i][j] = 0
			} else {
				C[i][j] = -A[i][j] / A[i][i]
			}
		}
	}

	return C, d, nil
}

// проверка диагонального преобладания
func IsDiagonallyDominant(A Matrix) bool {
	for i := range A {
		diag := math.Abs(A[i][i])
		sumExceptDiag := 0.0
		for j, val := range A[i] {
			if i != j {
				sumExceptDiag += math.Abs(val)
			}
		}
		if diag <= sumExceptDiag {
			return false
		}
	}
	return true
}

// применяем перестановку
func applyPermutation(A Matrix, b Vector, index []int) {
	n := len(A)

	ACopy := CopyMatrix(A)
	bCopy := CopyVector(b)

	for i := 0; i < n; i++ {
		copy(A[i], ACopy[index[i]])
		b[i] = bCopy[index[i]]
	}
}

// модификация IsDiagonallyDominant на перестановки
func checkPermutation(A Matrix, b Vector, index []int) bool {
	ACopy := CopyMatrix(A)
	bCopy := CopyVector(b)

	applyPermutation(ACopy, bCopy, index)
	return IsDiagonallyDominant(ACopy)
}

// рекурсивная проверка всех перестановок
func tryPermutations(A Matrix, b Vector, index []int, start int) bool {
	n := len(A)

	if start == n {
		if checkPermutation(A, b, index) {
			applyPermutation(A, b, index)
			return true
		}
		return false
	}

	for i := start; i < n; i++ {
		index[start], index[i] = index[i], index[start]
		if tryPermutations(A, b, index, start+1) {
			return true
		}
		index[start], index[i] = index[i], index[start]
	}

	return false

}

// если возможно диагональное преобладание - возвращаем true, иначе false
func TryToMakeDiagonallyDominant(A Matrix, b Vector) bool {
	if IsDiagonallyDominant(A) {
		return true
	}

	n := len(A)
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}

	return tryPermutations(A, b, indices, 0)
}
