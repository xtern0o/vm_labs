package service

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"vm_lab3/internal/dto"
)

const minDelta = 1e-8
const deltaStep = 0.1

const numOfEqual = 3

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

	// if !spec.Convergent { // TODO: определять разрывы
	// 	return dto.CalcIntegralResponseDto{}, fmt.Errorf("интеграл не существует - расходится")
	// }

	res := dto.CalcIntegralResponseDto{
		Messages: []string{"обнаружены точки разрыва - интеграл вычисляется как несобственный"},
	}
	res.Messages = append(res.Messages, fmt.Sprintf("точки разрыва: %v", activeBreakpoints))

	delta := 1e-2
	var prevValue float64
	firstIter := true
	var lastIntervals []interval
	var lastParts []dto.CalcIntegralResponseDto

	for {
		intervals := buildSubIntervals(a, b, activeBreakpoints, delta)
		if len(intervals) == 0 {
			return dto.CalcIntegralResponseDto{}, fmt.Errorf("не удалось построить интервалы интегрирования")
		}

		value := 0.0
		nSum := 0
		rungeR := 0.0
		parts := make([]dto.CalcIntegralResponseDto, 0, len(intervals))

		for _, in := range intervals {
			partRes, err := SolveIntegral(method, spec.Fn, in.A, in.B, eps, n)
			if err != nil {
				return dto.CalcIntegralResponseDto{}, err
			}
			value += partRes.Value
			nSum += partRes.N
			rungeR = math.Max(rungeR, partRes.RungeR)
			parts = append(parts, partRes)
		}

		dif := math.Abs(value - prevValue)

		if !firstIter {
			res.Messages = append(res.Messages, fmt.Sprintf("currDelta=%.9f ; |I_prev - I_curr| = %f", delta, dif))
		}

		if !firstIter && dif < eps {
			res.Value = value
			res.N = nSum
			res.RungeR = rungeR
			res.Messages = append(res.Messages, fmt.Sprintf("интервалы интегрирования: %s", intervalsToString(lastIntervals)))
			for i, in := range lastIntervals {
				res.Messages = append(res.Messages, lastParts[i].Messages...)
				res.Messages = append(res.Messages, fmt.Sprintf("вычислен интеграл на интервале (%.6f, %.6f)", in.A, in.B))
			}

			res.Messages = append(res.Messages, fmt.Sprintf("достигнута сходимость по eps=%.2e", eps))
			return res, nil
		}

		prevValue = value
		lastIntervals = intervals
		lastParts = parts
		firstIter = false
		delta *= deltaStep

		if delta <= minDelta {
			res.Messages = append(res.Messages, "достигнут минимальный delta, но сходимость не достигнута")
			res.Messages = append(res.Messages, "[!] возможно интеграл расходится")
			res.Value = value
			res.N = nSum
			res.RungeR = rungeR
			return res, nil
		}

	}
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
