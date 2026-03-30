package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vm_lab2/dto"
	"vm_lab2/internal/catalog"
	"vm_lab2/internal/solver"

	"github.com/go-chi/chi/v5"
)

const (
	defaultEpsEquation     = 1e-6
	defaultMaxIterEquation = 100
)

type errorResponse struct {
	Error string `json:"error"`
}

func EquationHandler(w http.ResponseWriter, r *http.Request) {
	solverIDStr := chi.URLParam(r, "solver_id")
	solverID, err := strconv.Atoi(solverIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid solver id")
		return
	}

	solvers := catalog.GetEquationSolversMap()
	s := solvers[solverID]
	if s == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("equation solver with id=%d not found", solverID))
		return
	}

	var req dto.EquationSolveRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if req.EquationId <= 0 {
		writeJSONError(w, http.StatusBadRequest, "equation_id must be provided and > 0")
		return
	}

	equationsMap := catalog.GetEquations()
	f := equationsMap[req.EquationId]
	if f == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("equation with id=%d not found", req.EquationId))
		return
	}

	if req.A == req.B {
		writeJSONError(w, http.StatusBadRequest, "a and b must be different")
		return
	}

	eps := defaultEpsEquation
	if req.Eps != nil {
		eps = *req.Eps
	}
	if eps <= 0 {
		writeJSONError(w, http.StatusBadRequest, "eps must be > 0")
		return
	}

	maxIter := defaultMaxIterEquation
	if req.MaxIter != nil {
		maxIter = *req.MaxIter
	}
	if maxIter <= 0 {
		writeJSONError(w, http.StatusBadRequest, "maxIter must be > 0")
		return
	}

	res, err := solver.SolveEquation(f, s, req.A, req.B, eps, maxIter)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, dto.EquationResponse{Result: res})

}
