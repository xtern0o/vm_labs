package service

import (
	"fmt"
	"math"
)

type ImproperFuncSpec struct {
	Fn          func(x float64) float64 // функция
	Breakpoints []float64               // точки разрыва
	Convergent  bool                    // сходится ли
}

type BreakType string // классификация точек разрыва

const (
	BreakAtA    BreakType = "break_at_a"
	BreakAtB    BreakType = "break_at_b"
	BreakInside BreakType = "break_inside"
	BreakAtBoth BreakType = "break_at_both"
)

var (
	// разрыв в левой границе x = 0, интеграл сходится на [0, b], b > 0
	improperFunc1 = func(x float64) float64 {
		return 1 / math.Sqrt(x)
	}

	// разрыв в правой границе x = 1, интеграл расходится на [a, 1]
	improperFunc2 = func(x float64) float64 {
		return 1 / (1 - x)
	}

	// разрывы в x = 0 и x = 2, интеграл сходится на [0, 2]
	improperFunc3 = func(x float64) float64 {
		return 1 / math.Sqrt(2*x-x*x)
	}
)

var id2improperFunc = map[int]ImproperFuncSpec{
	1: {
		Fn:          improperFunc1,
		Breakpoints: []float64{0},
		Convergent:  true,
	},
	2: {
		Fn:          improperFunc2,
		Breakpoints: []float64{1},
		Convergent:  false,
	},
	3: {
		Fn:          improperFunc3,
		Breakpoints: []float64{0, 2},
		Convergent:  true,
	},
}

func GetImproperFuncSpec(id int) (ImproperFuncSpec, error) {
	spec, ok := id2improperFunc[id]
	if !ok {
		return ImproperFuncSpec{}, fmt.Errorf("none of improper func with id=%d", id)
	}
	return spec, nil
}

func GetImproperFunc(id int) (func(x float64) float64, error) {
	spec, err := GetImproperFuncSpec(id)
	if err != nil {
		return nil, err
	}
	return spec.Fn, nil
}

func GetImproperBreakpoints(id int) ([]float64, error) {
	spec, err := GetImproperFuncSpec(id)
	if err != nil {
		return nil, err
	}
	return spec.Breakpoints, nil
}

func IsImproperConvergent(id int) (bool, error) {
	spec, err := GetImproperFuncSpec(id)
	if err != nil {
		return false, err
	}
	return spec.Convergent, nil
}
