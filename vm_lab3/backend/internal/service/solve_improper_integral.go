package service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"vm_lab3/internal/dto"
)

const improperDelta = 1e-6

type interval struct {
	A float64
	B float64
}

func SolveImproperIntegral(
	method IntegralSolver,
	spec ImproperFuncSpec,
	a, b, eps float64,
	n int,
) (dto.CalcIntegralResponseDto, error) {

	activeBreakpoints := filterBreakpoints(spec.Breakpoints, a, b)
	if len(activeBreakpoints) == 0 {
		res, err := SolveIntegral(method, spec.Fn, a, b, eps, n)
		if err != nil {
			return res, err
		}
		res.Messages = append(res.Messages, "точки разрыва не попадают в пределы интегрирования - интеграл вычислен как обычный")
		return res, nil
	}

	if !spec.Convergent {
		return dto.CalcIntegralResponseDto{}, fmt.Errorf("интеграл не существует - расходится")
	}

	intervals := buildSubIntervals(a, b, activeBreakpoints, improperDelta)
	if len(intervals) == 0 {
		return dto.CalcIntegralResponseDto{}, fmt.Errorf("не удалось построить интервалы интегрирования")
	}

	res := dto.CalcIntegralResponseDto{
		Messages: []string{"обнаружены точки разрыва - интеграл вычисляется как несобственный"},
	}
	res.Messages = append(res.Messages, fmt.Sprintf("точки разрыва: %v", activeBreakpoints))
	res.Messages = append(res.Messages, fmt.Sprintf("интервалы интегрирования: %s", intervalsToString(intervals)))

	for _, in := range intervals {
		partRes, err := SolveIntegral(method, spec.Fn, in.A, in.B, eps, n)
		if err != nil {
			return dto.CalcIntegralResponseDto{}, err
		}

		res.Value += partRes.Value
		res.N += partRes.N

		res.RungeR = math.Max(res.RungeR, partRes.RungeR)

		res.Messages = append(res.Messages, partRes.Messages...)
		res.Messages = append(res.Messages, fmt.Sprintf("вычислен интеграл на интервале (%.6f, %.6f)", in.A, in.B))
	}

	return res, nil
}

// проверка точек разрыва
func filterBreakpoints(points []float64, a, b float64) []float64 {
	res := make([]float64, 0)
	for _, p := range points {
		if a <= p && p <= b {
			res = append(res, p)
		}
	}
	sort.Float64s(res)
	return res
}

// разбиение интервала с учетом точек разрыва - отступаем от них на малое delta
func buildSubIntervals(a, b float64, breakpoints []float64, delta float64) []interval {
	if len(breakpoints) == 0 {
		return []interval{{A: a, B: b}}
	}

	res := make([]interval, 0)

	first := breakpoints[0]
	if a < first-delta {
		res = append(res, interval{
			A: a,
			B: first - delta,
		})
	}

	for i := 0; i < len(breakpoints)-1; i++ {
		left := breakpoints[i] + delta
		right := breakpoints[i+1] - delta

		if left < right {
			res = append(res, interval{
				A: left,
				B: right,
			})
		}
	}

	last := breakpoints[len(breakpoints)-1]
	if last+delta < b {
		res = append(res, interval{
			A: last + delta,
			B: b,
		})
	}

	return res
}

// преобразование ну короче в toString
func intervalsToString(intervals []interval) string {
	if len(intervals) == 0 {
		return ""
	}

	var sb strings.Builder

	for i, in := range intervals {
		sb.WriteString(fmt.Sprintf("(%.6f, %.6f)", in.A, in.B))
		if i != len(intervals)-1 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}
