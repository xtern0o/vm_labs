package service

import "vm_lab3/internal/dto"

// SolveIntegral - обобщенный метод для вычисления интеграла.
// Автоматически проверяет наличие точек разрыва и, если они есть, вызывает алгоритм несобственного интегрирования.
func SolveIntegral(
	method IntegralSolver,
	f func(x float64) float64,
	a, b, eps float64,
	n int,
) (dto.CalcIntegralResponseDto, error) {
	breakpoints := FindBreakpoints(f, a, b)
	if len(breakpoints) == 0 {
		return method.solve(f, a, b, eps, n)
	}

	return solveImproper(method, f, breakpoints, a, b, eps, n)
}
