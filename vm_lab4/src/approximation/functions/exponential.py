"""Экспоненциальная аппроксимация: phi(x) = a * exp(bx)"""

import numpy as np
from ..base import ApproximationFunction


class ExponentialApproximation(ApproximationFunction):
    """Экспоненциальная аппроксимация: phi(x) = a * exp(bx)"""
    
    def fit(self) -> np.ndarray:
        """Найти коэффициенты a и b путем линеаризации"""
        # Линеаризация: ln(phi) = ln(a) + bx
        # Y = ln(y), A = ln(a), B = b
        
        if np.any(self.y <= 0):
            raise ValueError("Для экспоненциальной функции все y должны быть > 0")
        
        Y = np.log(self.y)
        
        sx = np.sum(self.x)
        sxx = np.sum(self.x ** 2)
        sY = np.sum(Y)
        sxY = np.sum(self.x * Y)
        
        delta = sxx * self.n - sx * sx
        
        if delta == 0:
            raise ValueError("Система вырождена, невозможно найти решение")
        
        B = (sxY * self.n - sx * sY) / delta
        A = (sxx * sY - sx * sxY) / delta
        
        a = np.exp(A)   # обратное преобр.
        b = B
        
        self.coefficients = np.array([a, b])
        self.approximated_y = None
        self.errors = None
        
        return self.coefficients
    
    def evaluate(self, x_values: np.ndarray) -> np.ndarray:
        """Вычислить phi(x) = a * exp(b*x)"""
        if self.coefficients is None:
            self.fit()
        
        a, b = self.coefficients
        return a * np.exp(b * x_values)
