package handler

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"vm_lab1/internal/dto"
	"vm_lab1/internal/matrix"
	"vm_lab1/internal/solver"
)

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.SolveRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", 400)
		respondJSON(w, 400, dto.ErrorResponse{
			Error:   "invalid json",
			Details: err.Error(),
		})
		return
	}

	if err := validateRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, dto.ErrorResponse{
			Error:   "validation failed",
			Details: err.Error(),
		})
		return
	}

	A := matrix.MatrixFromFloat64(req.Matrix)
	b := matrix.Vector(req.Vector)
	x0 := getInitialGuess(req.InitialGuess, len(A))
	eps := getEpsilon(req.Epsilon)
	maxIter := getMaxIterations(req.MaxIter)

	s := solver.NewSeidelSolver()
	res, err := s.Solve(A, b, x0, eps, maxIter)
	if err != nil {
		respondJSON(w, http.StatusUnprocessableEntity, dto.ErrorResponse{
			Error:   "solve failed",
			Details: err.Error(),
		})
		return
	}

	resp := dto.SolveResponse{
		Solution:     res.Solution,
		Errors:       res.Errors,
		Iterations:   res.Iterations,
		NormOfMatrix: res.NormOfMatrix,
		Messages:     res.Messages,
	}

	respondJSON(w, 200, resp)

}

func validateRequest(req *dto.SolveRequest) error {
	if len(req.Matrix) == 0 {
		return errors.New("матрица пустая")
	}

	n := len(req.Matrix)

	// матрица - квадратная
	for _, row := range req.Matrix {
		if len(row) != n {
			return errors.New("матрица должна быть квадратной")
		}
	}

	for _, row := range req.Matrix {
		for _, val := range row {
			if !isValidNumber(val) {
				return errors.New("матрица содержит некорректное число (не должно быть бесконечностей или nan)")
			}
		}
	}

	if len(req.Vector) != n {
		return errors.New("длина вектора должна совпадать с размером матрицы")
	}

	for _, val := range req.Vector {
		if !isValidNumber(val) {
			return errors.New("вектор должен состоять из чисел!")
		}
	}

	if len(req.InitialGuess) > 0 && len(req.InitialGuess) != n {
		return errors.New("начальное приближение должно совпадать по длине с вектором и матрицей")
	}

	if req.Epsilon != nil && *req.Epsilon <= 0 {
		return errors.New("eps должен быть положительным")
	}

	if req.MaxIter != nil && *req.MaxIter <= 0 {
		return errors.New("max_iterations должен быть положительным")
	}

	return nil
}

func isValidNumber(val float64) bool {
	return !math.IsNaN(val) && !math.IsInf(val, 0)
}

func getInitialGuess(guess []float64, n int) matrix.Vector {
	if len(guess) == n {
		return matrix.Vector(guess)
	}
	return matrix.NewVector(n)
}

func getEpsilon(eps *float64) float64 {
	if eps != nil && *eps > 0 {
		return *eps
	}
	return 0.0001
}

func getMaxIterations(max *int) int {
	if max != nil && *max > 0 {
		return *max
	}
	return 1000
}
