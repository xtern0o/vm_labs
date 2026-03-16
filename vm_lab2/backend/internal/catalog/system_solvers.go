package catalog

import "vm_lab2/internal/solver"

var system2SolversMap = map[int]solver.Systen2Solver{
	1: solver.NewSystemSimpleIterSolver(),
}

func GetSystem2SolversMap() map[int]solver.Systen2Solver {
	return system2SolversMap
}
