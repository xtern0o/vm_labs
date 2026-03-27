package service

import "vm_lab3/internal/dto"

func SolveIntegral(
	method IntegralSolver,
	f func(x float64) float64,
	a, b, eps float64,
	n int,
) (dto.CalcIntegralResponseDto, error) {
	return method.solve(f, a, b, eps, n)
}
