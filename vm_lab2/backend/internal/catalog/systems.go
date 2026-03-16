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

var sys2func1 = func(x float64, y float64) float64 {
	return x*x + y*y - 1
}
var sys2func2 = func(x float64, y float64) float64 {
	return x*x - 3*y - 0.5
}

// приводим к виду x = phi(x)
var sys1phi1 = func(x float64, y float64) float64 {
	return math.Sqrt(math.Tan(x*y + 0.3))
}
var sys1phi2 = func(x float64, y float64) float64 {
	if y > 0 {
		return math.Sqrt((1 - 0.9*x*x) / 2)
	} else {
		return -math.Sqrt((1 - 0.9*x*x) / 2)
	}
}

var sys2phi1 = func(x, y float64) float64 {
	if x >= 0 {
		return math.Sqrt(1 - y*y)
	} else {
		return -math.Sqrt(1 - y*y)
	}
}
var sys2phi2 = func(x, y float64) float64 {
	return (x*x - 0.5) / 3
}

var systemsMap = map[int]*equations.System2{
	1: equations.NewSystem2(sys1func1, sys1func2),
	2: equations.NewSystem2(sys2func1, sys2func2),
}

var systemsPhiMap = map[int]*equations.System2{
	1: equations.NewSystem2(sys1phi1, sys1phi2),
	2: equations.NewSystem2(sys2phi1, sys2phi2),
}

func GetSystem(id int) *equations.System2 {
	return systemsMap[id]
}

func GetPhiSystem(id int) *equations.System2 {
	return systemsPhiMap[id]
}
