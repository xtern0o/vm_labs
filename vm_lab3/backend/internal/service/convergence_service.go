package service

import (
	"fmt"
	"math"
)

// checkConvergence проверяет сходимость интеграла в точках разрыва.
// Идея основана на численном предельном признаке сравнения.
// Если в окрестности точки разрыва c функция ведет себя как 1/|x-c|^p,
// то при p < 1 интеграл сходится, а при p >= 1 - расходится.
// Мы можем найти p, взяв отношения значений функции при приближении к точке разрыва:
// p ~= log2( f(c +- \delta/2)| / |f(c +- \delta)| )
func checkConvergence(fn func(float64) float64, breakpoints []float64, a, b float64) (bool, string) {
	delta := 1e-4 // начальный небольшой отступ для оценки предела

	for _, c := range breakpoints {
		// проверяем справа, если точка разрыва левее верхней границы
		if c < b {
			pRight := estimateP(fn, c, delta, true)
			if pRight >= 0.99 {
				return false, fmt.Sprintf("расходится справа от точки %.2f (порядок особенности p = %.3f >= 1)", c, pRight)
			}
		}
		// проверяем слева, если правее
		if c > a {
			pLeft := estimateP(fn, c, delta, false)
			if pLeft >= 0.99 {
				return false, fmt.Sprintf("расходится слева от точки %.2f (порядок особенности p = %.3f >= 1)", c, pLeft)
			}
		}
	}
	return true, "численная проверка сходимости пройдена: p < 1 во всех точках разрыва"
}

// estimateP вычисляет порядок особенности (роста функции) вблизи точки разрыва
func estimateP(fn func(float64) float64, c float64, delta float64, right bool) float64 {
	var x1, x2 float64
	if right {
		x1 = c + delta
		x2 = c + delta/2.0
	} else {
		x1 = c - delta
		x2 = c - delta/2.0
	}

	f1 := math.Abs(fn(x1))
	f2 := math.Abs(fn(x2))

	// eсли значения функции близки к нулю или NaN (не бесконечность),
	// то разрыва второго рода здесь может и не быть. Возвращаем 0 (сходится).
	if f1 < 1e-5 || math.IsNaN(f1) || math.IsNaN(f2) {
		return 0
	}

	// вычисляем порядок p
	p := math.Log2(f2 / f1)
	return p
}

// метод ищет точки разрыва второго рода на отрезке [a, b]
// метод основан на поиске нулей обратной функции g(x) = 1 / f(x)
// f(x) -> \infty => g(x) -> 0
func FindBreakpoints(fn func(float64) float64, a, b float64) []float64 {
	breakpoints := make([]float64, 0)

	// g(x) = 1 / f(x) OR 0 if f(x) == NaN|Inf
	g := func(x float64) float64 {
		val := fn(x)
		if math.IsInf(val, 0) || math.IsNaN(val) {
			return 0
		}
		return 1.0 / val
	}

	// разбиваем интервал
	const intervals = 1000
	step := (b - a) / float64(intervals)

	// ищем нули ИЛИ экстремумы близкие к нулю
	for i := 0; i < intervals; i++ {
		x0 := a + float64(i)*step
		x1 := a + float64(i+1)*step

		v0 := g(x0)
		v1 := g(x1)

		// точное попадание в ноль
		if v0 == 0 {
			breakpoints = append(breakpoints, x0)
		}
		if v1 == 0 {
			breakpoints = append(breakpoints, x1)
		}
		if v0 == 0 || v1 == 0 {
			continue
		}

		if v0*v1 < 0 {
			root := bisectionRoot(g, x0, x1, 1e-7)

			val := fn(root)
			// разрыв второго рода (f(x) улетает в бесконечность) будет, если
			// модуль функции в этой точке NaN или Inf
			if math.IsNaN(val) || math.IsInf(val, 0) || math.Abs(val) > 1e6 {
				breakpoints = append(breakpoints, root)
			}
			continue
		}

		// обрабатываем касание Ox (разрыв без смены знака - пример: 1/x^2)
		if math.Abs(v0) < 1e-2 && math.Abs(v1) < 1e-2 {
			mid := (x0 + x1) / 2
			vmid := g(mid)
			if math.Abs(vmid) < math.Abs(v0) && math.Abs(vmid) < math.Abs(v1) {
				root := findMinimum(g, x0, x1, 1e-7)

				// действительно близко к 0 ????
				if math.Abs(g(root)) < 1e-4 {
					breakpoints = append(breakpoints, root)
				}
			}
		}
	}

	// чистим дубликаты из за погрешностей
	unique := make([]float64, 0)
	for _, p := range breakpoints {
		isNew := true
		for _, u := range unique {
			if math.Abs(p-u) < 1e-4 {
				isNew = false
				break
			}
		}
		if isNew && p >= a && p <= b {
			unique = append(unique, p)
		}
	}

	return unique
}

// bisectionRoot ищет корень функции g(x) = 0 на отрезке [a, b] методом половинного деления
func bisectionRoot(g func(float64) float64, a, b float64, eps float64) float64 {
	x := (a + b) / 2.0
	for math.Abs(a-b) > eps {
		x = (a + b) / 2.0
		if math.Abs(g(x)) < eps {
			return x
		}
		if g(a)*g(x) <= 0 {
			b = x
		} else {
			a = x
		}
	}
	return x
}

// поиск локального минимума методом золотого сечения
func findMinimum(g func(float64) float64, a, b float64, eps float64) float64 {
	phi := (1.0 + math.Sqrt(5.0)) / 2.0
	resphi := 2.0 - phi

	x1 := a + resphi*(b-a)
	x2 := b - resphi*(b-a)

	f1 := math.Abs(g(x1))
	f2 := math.Abs(g(x2))

	for math.Abs(b-a) > eps {
		if f1 < f2 {
			b = x2
			x2 = x1
			f2 = f1
			x1 = a + resphi*(b-a)
			f1 = math.Abs(g(x1))
		} else {
			a = x1
			x1 = x2
			f1 = f2
			x2 = b - resphi*(b-a)
			f2 = math.Abs(g(x2))
		}
	}

	return (a + b) / 2.0
}
