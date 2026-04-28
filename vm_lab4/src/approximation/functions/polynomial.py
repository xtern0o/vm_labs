"""Полиномиальная аппроксимация"""

import numpy as np
from ..base import ApproximationFunction
from ..solver import SeidelSolver


class PolynomialApproximation(ApproximationFunction):
    """Полиномиальная аппроксимация: phi(x) = a_0 + a_1*x + a_2*x^2 + ... + a_n*x^n"""
    
    def __init__(self, x: np.ndarray, y: np.ndarray, degree: int):
        """
        Args:
            x: Массив значений x
            y: Массив значений y
            degree: Степень полинома
        """
        super().__init__(x, y)
        self.degree = degree
        self.solver = SeidelSolver(eps=1e-9, max_iter=5000)
    
    def fit(self) -> np.ndarray:
        """Найти коэффициенты полинома методом наименьших квадратов"""
        
        A = np.column_stack([self.x ** i for i in range(self.degree + 1)])
        AtA = A.T @ A
        Aty = A.T @ self.y
        
        # гаусс-зейдель
        try:
            result = self.solver.solve(AtA, Aty)
            self.coefficients = result.solution
        except RuntimeError as e:
            print(f"[Error]: {e}")
            print("  Используется прямой метод решения (numpy.linalg.solve)...")
            self.coefficients = np.linalg.solve(AtA, Aty)
        
        self.approximated_y = None  # ресет для перерасчета
        self.errors = None
        
        return self.coefficients
    
    def evaluate(self, x_values: np.ndarray) -> np.ndarray:
        """Вычислить полином"""
        if self.coefficients is None:
            self.fit()
        
        result = np.zeros_like(x_values, dtype=float)
        for i, coef in enumerate(self.coefficients):
            result += coef * (x_values ** i)
        
        return result
