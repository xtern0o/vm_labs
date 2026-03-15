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

type EquationResponse struct {
	Result solver.Result `json:"result"`
}
