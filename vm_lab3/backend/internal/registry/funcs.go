package registry

import "fmt"

var func1 func(x float64) float64 = func(x float64) float64 {
	return -2*x*x*x - 3*x*x + x + 5
}

var funcs []func(x float64) float64 = []func(x float64) float64{
	func1,
}

func GetFunc(id int) (func(x float64) float64, error) {
	if 0 <= id && id < len(funcs) {
		return funcs[id], nil
	}
	return nil, fmt.Errorf("none of func with id=%d", id)
}
