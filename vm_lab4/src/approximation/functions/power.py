"""Степенная аппроксимация: phi(x) = a * x**b"""

import numpy as np
from ..base import ApproximationFunction


class PowerApproximation(ApproximationFunction):
    """Степенная аппроксимация: phi(x) = a * x**b"""
    
    def fit(self) -> np.ndarray:
        """Найти коэффициенты a и b путем линеаризации"""
        # Линеаризация: ln(phi) = ln(a) + b * ln(x)
        # Y = ln(y), X = ln(x), A = ln(a), B = b
        
        if np.any(self.x <= 0) or np.any(self.y <= 0):
            raise ValueError("Для степенной функции все x и y должны быть > 0")
        
        X = np.log(self.x)
        Y = np.log(self.y)
        
        # Применяем линейную регрессию к (X, Y)
        sX = np.sum(X)
        sXX = np.sum(X ** 2)
        sY = np.sum(Y)
        sXY = np.sum(X * Y)
        
        delta = sXX * self.n - sX * sX
        
        if delta == 0:
            raise ValueError("Система вырождена, невозможно найти решение")
        
        B = (sXY * self.n - sX * sY) / delta
        A = (sXX * sY - sX * sXY) / delta
        
        a = np.exp(A)
        b = B
        
        self.coefficients = np.array([a, b])
        self.approximated_y = None
        self.errors = None
        
        return self.coefficients
    
    def evaluate(self, x_values: np.ndarray) -> np.ndarray:
        """Вычислить phi(x) = a * x**b"""
        if self.coefficients is None:
            self.fit()
        
        a, b = self.coefficients
        return a * np.power(x_values, b)
