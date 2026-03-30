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

func ImproperIntegralHandler(w http.ResponseWriter, r *http.Request) {
	solverId, err := strconv.Atoi(chi.URLParam(r, "solver_id"))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "incorrect method id")
		return
	}
	solver, err := registry.GetMethod(solverId)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("error: %v", err))
		return
	}

	var req dto.CalcIntegralRequestDto
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	spec, err := service.GetImproperFuncSpec(req.FuncID)
	if err != nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("error: %v", err))
		return
	}

	n := defaultN
	if req.N != nil {
		n = *req.N
	}

	res, err := service.SolveImproperIntegral(*solver, spec, req.A, req.B, req.Eps, n)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)

}
