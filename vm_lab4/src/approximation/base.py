"""Базовый класс для аппроксимирующих функций"""

import numpy as np
from abc import ABC, abstractmethod


class ApproximationFunction(ABC):
    """Абстрактный базовый класс для аппроксимирующих функций"""
    
    def __init__(self, x: np.ndarray, y: np.ndarray):
        """
        Args:
            x: Массив значений x
            y: Массив значений y
        """
        self.x = np.array(x, dtype=float)
        self.y = np.array(y, dtype=float)
        self.n = len(x)
        self.coefficients = None
        self.approximated_y = None
        self.errors = None
    
    @abstractmethod
    def fit(self) -> np.ndarray:
        """Найти коэффициенты функции методом наименьших квадратов"""
        pass
    
    @abstractmethod
    def evaluate(self, x_values: np.ndarray) -> np.ndarray:
        """Вычислить значения аппроксимирующей функции"""
        pass
    
    def get_approximation_values(self) -> np.ndarray:
        """Получить аппроксимированные значения y"""
        if self.approximated_y is None:
            self.approximated_y = self.evaluate(self.x)
        return self.approximated_y
    
    def get_errors(self) -> np.ndarray:
        """Получить отклонения для каждой точки"""
        if self.errors is None:
            self.errors = self.get_approximation_values() - self.y
        return self.errors
    
    def calculate_deviation_sum(self) -> float:
        """Вычислить сумму квадратов отклонений S"""
        errors = self.get_errors()
        return np.sum(errors ** 2)
    
    def calculate_sigma(self) -> float:
        """Вычислить сигму (sigma = sqrt(S/n))"""
        deviation_sum = self.calculate_deviation_sum()
        return np.sqrt(deviation_sum / self.n)
    
    def calculate_r_squared(self) -> float:
        """Вычислить коэффициент детерминации (R^2)"""
        approx_y = self.get_approximation_values()
        mean_approx = np.mean(approx_y)
        
        ss_res = np.sum((self.y - approx_y) ** 2)
        ss_tot = np.sum((self.y - mean_approx) ** 2)
        
        if ss_tot == 0:
            return 1.0 if ss_res == 0 else 0.0
        return 1 - (ss_res / ss_tot)
    
    def calculate_pearson_r(self) -> float:
        """Вычислить коэффициент корреляции Пирсона"""
        x_mean = np.mean(self.x)
        y_mean = np.mean(self.y)
        
        numerator = np.sum((self.x - x_mean) * (self.y - y_mean))
        denominator = np.sqrt(
            np.sum((self.x - x_mean) ** 2) * np.sum((self.y - y_mean) ** 2)
        )
        
        if denominator == 0:
            return 0.0
        return numerator / denominator
