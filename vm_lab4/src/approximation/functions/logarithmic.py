"""Логарифмическая аппроксимация: phi(x) = a * ln(x) + b"""

import numpy as np
from ..base import ApproximationFunction


class LogarithmicApproximation(ApproximationFunction):
    """Логарифмическая аппроксимация: phi(x) = a * ln(x) + b"""
    
    def fit(self) -> np.ndarray:
        """Найти коэффициенты a и b путем линеаризации"""
        # линеаризация: X = ln(x), phi = aX + b
        
        if np.any(self.x <= 0):
            raise ValueError("Для логарифмической функции все x должны быть > 0")
        
        X = np.log(self.x)
        
        sX = np.sum(X)
        sXX = np.sum(X ** 2)
        sy = np.sum(self.y)
        sXy = np.sum(X * self.y)
        
        delta = sXX * self.n - sX * sX
        
        if delta == 0:
            raise ValueError("Система вырождена, невозможно найти решение")
        
        a = (sXy * self.n - sX * sy) / delta
        b = (sXX * sy - sX * sXy) / delta
        
        self.coefficients = np.array([a, b])
        self.approximated_y = None
        self.errors = None
        
        return self.coefficients
    
    def evaluate(self, x_values: np.ndarray) -> np.ndarray:
        """Вычислить phi(x) = a * ln(x) + b"""
        if self.coefficients is None:
            self.fit()
        
        a, b = self.coefficients
        return a * np.log(x_values) + b
