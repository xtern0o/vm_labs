package catalog

import "vm_lab2/internal/solver"

var equationSolversMap = map[int]solver.EquationSolver{
	1: solver.NewHalfDivisionSolver(),
	2: solver.NewNewthonSolver(),
	3: solver.NewSimpleIterationsSolver(),
}

func GetEquationSolversMap() map[int]solver.EquationSolver {
	return equationSolversMap
}
