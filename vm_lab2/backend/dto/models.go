package dto

import (
	"vm_lab2/internal/solver"
)

// тут еще есть дто, не хочу плодить
// кучу мапперов для учебного проекта

type EquationSolveRequest struct {
	EquationId int      `json:"equation_id"`
	A          float64  `json:"a"`
	B          float64  `json:"b"`
	Eps        *float64 `json:"eps,omitempty"`
	MaxIter    *int     `json:"max_iter,omitempty"`
}

type System2SolveRequest struct {
	SystemId int      `json:"system_id"`
	x0       *float64 `json:"x0,omitempty"`
	y0       *float64 `json:"y0,omitempty"`
	Eps      *float64 `json:"eps,omitempty"`
	MaxIter  *int     `json:"max_iter,omitempty"`
}

type EquationResponse struct {
	Result solver.Result `json:"result"`
}
