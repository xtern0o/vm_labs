package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"vm_lab3/internal/dto"
	"vm_lab3/internal/registry"
	"vm_lab3/internal/service"

	"github.com/go-chi/chi/v5"
)

const (
	defaultN = 4.0
)

func IntegralHandler(w http.ResponseWriter, r *http.Request) {
	solverId, err := strconv.Atoi(chi.URLParam(r, "solver_id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "incorrect method id")
		return
	}
	solver, err := registry.GetMethod(solverId)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}

	var req dto.CalcIntegralRequestDto
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
	}

	fun, err := registry.GetFunc(req.FuncID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("error: %v", err))
	}

	n := req.N
	if n == nil {
		*n = defaultN
	}

	res, err := service.SolveIntegral(*solver, fun, req.A, req.B, req.Eps, *n)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
	}
	writeJSON(w, http.StatusOK, res)

}
