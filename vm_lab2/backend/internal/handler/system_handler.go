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
	defaultX0             = 0.0
	defaultY0             = 0.0
	defaultEpsSystem2     = 0.001
	defaultMaxIterSystem2 = 1000
)

func SystemHandler(w http.ResponseWriter, r *http.Request) {
	solverIDStr := chi.URLParam(r, "solver_id")
	solverID, err := strconv.Atoi(solverIDStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid solver id")
		return
	}

	solvers := catalog.GetSystem2SolversMap()
	s := solvers[solverID]
	if s == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("system2 solver with id=%d not found", solverID))
		return
	}

	var req dto.System2SolveRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, fmt.Sprintf("invalid request body: %v", err))
		return
	}

	if req.SystemId <= 0 {
		writeJSONError(w, http.StatusBadRequest, "equation_id must be provided and > 0")
		return
	}
	system := catalog.GetSystem(req.SystemId)
	if system == nil {
		writeJSONError(w, http.StatusNotFound, fmt.Sprintf("system with id=%d not found", req.SystemId))
		return
	}

	eps := defaultEpsSystem2
	if req.Eps != nil {
		eps = *req.Eps
	}
	if eps <= 0 {
		writeJSONError(w, http.StatusBadRequest, "eps must be > 0")
		return
	}

	maxIter := defaultMaxIterSystem2
	if req.MaxIter != nil {
		maxIter = *req.MaxIter
	}
	if maxIter <= 0 {
		writeJSONError(w, http.StatusBadRequest, "maxIter must be > 0")
		return
	}

	x0 := defaultX0
	if req.X0 != nil {
		x0 = *req.X0
	}
	y0 := defaultX0
	if req.Y0 != nil {
		y0 = *req.Y0
	}

	res, err := solver.SolveSystem2(*system, *catalog.GetPhiSystem(req.SystemId), s, x0, y0, eps, maxIter)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, res)

}
