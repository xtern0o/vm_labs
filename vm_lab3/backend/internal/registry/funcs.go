package registry

import (
	"fmt"
	"math"
)

var (
	func1 = func(x float64) float64 {
		return -2*x*x*x - 3*x*x + x + 5
	}
	func2 = func(x float64) float64 {
		return math.Exp(0.1 * x)
	}
	func3 = func(x float64) float64 {
		return math.Sin(math.Cos(x)) + 2
	}
	func4 = func(x float64) float64 {
		return 2 * math.Log(x)
	}
)

var id2func = map[int]func(x float64) float64{
	1: func1,
	2: func2,
	3: func3,
	4: func4,
}

func GetFunc(id int) (func(x float64) float64, error) {
	fun := id2func[id]
	if fun == nil {
		return nil, fmt.Errorf("none of func with id=%d", id)
	}
	return fun, nil
}
