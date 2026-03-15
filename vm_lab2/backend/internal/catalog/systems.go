package catalog

import (
	"math"
	"vm_lab2/internal/equations"
)

var sys1func1 = func(x float64, y float64) float64 {
	return math.Tan(x*y+0.3) - x*x
}
var sys1func2 = func(x float64, y float64) float64 {
	return 0.9*x*x + 2*y*y - 1
}

var systemsMap = map[int]equations.System2{
	1: *equations.NewSystem2(sys1func1, sys1func2),
}

func GetSystems() map[int]equations.System2 {
	return systemsMap
}
