package solver

import (
	"errors"
	"fmt"
	"math"
	"vm_lab2/internal/equations"
	"vm_lab2/internal/numeric"
)

const segments = 200

type EquationSolver interface {
	Solve(
		f equations.Equation,
		a float64,
		b float64,
		eps float64,
		maxIter int,
	) (Result, error)
}

// проверка на область определения функции на отрезке, знакопостоянство и наличие только 1 корня
func BasicCheck(f equations.Equation, a float64, b float64) error {
	if f == nil {
		return errors.New("equation should not be nil")
	}
	if math.IsNaN(f(a)) || math.IsNaN(f(b)) || math.IsInf(f(a), 0) || math.IsInf(f(b), 0) {
		return fmt.Errorf("функция не определена на границах промежутка: f(%f)=%f, f(%f)=%f", a, f(a), b, f(b))
	}
	if f(a)*f(b) > 0 {
		return fmt.Errorf("нет смены знака на [%f, %f]: f(a)=%f, f(b)=%f. Значения на границах должны быть разных знаков", a, b, f(a), f(b))
	}
	numOfRoots := len(numeric.RootBrackets(f, a, b, segments))
	if numOfRoots > 1 {
		return fmt.Errorf("функция имеет более 1 корня на интервале [%f, %f]. На этом промежутке %d смен(-ы) знака", a, b, numOfRoots)
	}

	return nil
}

// проверка нужной точности по ответу
func CheckAns(f equations.Equation, x, eps float64) bool {
	return math.Abs(f(x)) < eps
}
