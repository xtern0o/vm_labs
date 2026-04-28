"""Линейная аппроксимация: phi(x) = ax + b"""

import numpy as np
from ..base import ApproximationFunction


class LinearApproximation(ApproximationFunction):
    """Линейная аппроксимация: phi(x) = ax + b"""
    
    def fit(self) -> np.ndarray:
        """Найти коэффициенты a и b"""
        sx = np.sum(self.x)
        sxx = np.sum(self.x ** 2)
        sy = np.sum(self.y)
        sxy = np.sum(self.x * self.y)
        
        delta = sxx * self.n - sx * sx
        
        if delta == 0:
            raise ValueError("Система вырождена, невозможно найти решение")
        
        delta1 = sxy * self.n - sx * sy
        delta2 = sxx * sy - sx * sxy
        
        a = delta1 / delta
        b = delta2 / delta
        
        self.coefficients = np.array([a, b])
        
        self.approximated_y = None
        self.errors = None
        
        return self.coefficients
    
    def evaluate(self, x_values: np.ndarray) -> np.ndarray:
        """Вычислить линейную аппроксимацию"""
        if self.coefficients is None:
            self.fit()
        
        a, b = self.coefficients
        return a * x_values + b
