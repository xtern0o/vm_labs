package numeric

import "vm_lab2/internal/equations"

type Interval struct {
	Left  float64
	Right float64
}

// разбивает функцию на сегменты и возвращает интервалы, на которых функция меняет знак
func RootBrackets(f equations.Equation, a float64, b float64, segments int) []Interval {
	if f == nil || segments < 1 {
		return nil
	}
	if a > b {
		a, b = b, a
	}

	step := (b - a) / float64(segments)
	if step == 0 {
		return nil
	}

	out := make([]Interval, 0)
	xLeft := a
	yLeft := f(xLeft)

	for i := 1; i <= segments; i++ {
		xRight := a + float64(i)*step
		yRight := f(xRight)

		if yLeft == 0 {
			out = append(out, Interval{Left: xLeft, Right: xLeft})
		}
		if yLeft*yRight < 0 {
			out = append(out, Interval{Left: xLeft, Right: xRight})
		}

		xLeft = xRight
		yLeft = yRight
	}

	if yLeft == 0 {
		out = append(out, Interval{Left: xLeft, Right: xLeft})
	}

	return out
}

// возвращает количество корней на исследуемом промежутке
func HasMultipleRootCandidates(f equations.Equation, a float64, b float64, segments int) bool {
	return len(RootBrackets(f, a, b, segments)) > 1
}

// проверка на то, что функция имеет 1 коренб на промежутке
func HasOnlyOneRootCandidate(f equations.Equation, a float64, b float64, segments int) bool {
	return len(RootBrackets(f, a, b, segments)) == 1
}
