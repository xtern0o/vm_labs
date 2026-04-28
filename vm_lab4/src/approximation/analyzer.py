"""Анализатор для работы со всеми типами аппроксимаций"""

import numpy as np
from typing import Dict, Tuple
from .result import ApproximationResult
from .functions import (
    LinearApproximation,
    PolynomialApproximation,
    ExponentialApproximation,
    LogarithmicApproximation,
    PowerApproximation
)


class ApproximationAnalyzer:
    """Анализатор для работы со всеми типами аппроксимаций"""
    
    def __init__(self, x: np.ndarray, y: np.ndarray):
        """
        Args:
            x: Массив значений x
            y: Массив значений y
        """
        self.x = np.array(x, dtype=float)
        self.y = np.array(y, dtype=float)
        self.results = {}
    
    def fit_all(self) -> Dict[str, ApproximationResult]:
        """Аппроксимировать данные всеми доступными функциями"""
        
        approximations = {}
        
        try:
            approx = LinearApproximation(self.x, self.y)
            approx.fit()
            approximations['linear'] = approx
        except Exception as e:
            print(f"Ошибка линейной аппроксимации: {e}")
        
        for degree in [2, 3]:
            try:
                approx = PolynomialApproximation(self.x, self.y, degree)
                approx.fit()
                approximations[f'polynomial_{degree}'] = approx
            except Exception as e:
                print(f"Ошибка полиномиальной аппроксимации (степень {degree}): {e}")
        
        try:
            approx = ExponentialApproximation(self.x, self.y)
            approx.fit()
            approximations['exponential'] = approx
        except Exception as e:
            print(f"Ошибка экспоненциальной аппроксимации: {e}")
        
        try:
            approx = LogarithmicApproximation(self.x, self.y)
            approx.fit()
            approximations['logarithmic'] = approx
        except Exception as e:
            print(f"Ошибка логарифмической аппроксимации: {e}")
        
        try:
            approx = PowerApproximation(self.x, self.y)
            approx.fit()
            approximations['power'] = approx
        except Exception as e:
            print(f"Ошибка степенной аппроксимации: {e}")
        
        self.results = {}
        for name, approx in approximations.items():
            self.results[name] = ApproximationResult(
                coefficients=approx.coefficients.copy(),
                deviation_sum=approx.calculate_deviation_sum(),
                sigma=approx.calculate_sigma(),
                pearson_r=approx.calculate_pearson_r() if name == 'linear' else None,
                r_squared=approx.calculate_r_squared(),
                approximated_y=approx.get_approximation_values().copy(),
                errors=approx.get_errors().copy()
            )
        
        return self.results
    
    def get_best_approximation(self) -> Tuple[str, ApproximationResult]:
        """Получить наилучшую аппроксимацию по минимальной сигме (sigma)"""
        if not self.results:
            self.fit_all()
        
        best_name = min(self.results.keys(), 
                       key=lambda k: self.results[k].sigma)
        return best_name, self.results[best_name]
    
    def get_results_summary(self) -> Dict[str, Dict]:
        """Получить сводку результатов для всех функций"""
        if not self.results:
            self.fit_all()
        
        summary = {}
        for name, result in self.results.items():
            summary[name] = {
                'coefficients': result.coefficients,
                'deviation_sum': result.deviation_sum,
                'sigma': result.sigma,
                'pearson_r': result.pearson_r,
                'r_squared': result.r_squared
            }
        
        return summary
