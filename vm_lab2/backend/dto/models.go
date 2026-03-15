package dto

import (
	"vm_lab2/internal/plots"
	"vm_lab2/internal/solver"
)

type EquationResponse struct {
	Result    solver.Result           `json:"result"`
	CurveData plots.EquationCurveData `json:"curve_data"`
}
