package registry

import "fmt"

var (
	func1 = func(x float64) float64 {
		return -2*x*x*x - 3*x*x + x + 5
	}
)

var id2func = map[int]func(x float64) float64{
	1: func1,
}

func GetFunc(id int) (func(x float64) float64, error) {
	fun := id2func[id]
	if fun == nil {
		return nil, fmt.Errorf("none of func with id=%d", id)
	}
	return fun, nil
}
