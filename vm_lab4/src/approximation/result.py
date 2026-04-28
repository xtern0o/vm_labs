"""Результаты аппроксимации"""

from dataclasses import dataclass
import numpy as np


@dataclass
class ApproximationResult:
    """Результаты аппроксимации функции"""
    coefficients: np.ndarray      # коэффициенты функции
    deviation_sum: float          # сумма квадратов отклонений (S)
    sigma: float                  # ско (sigma = sqrt(S/n))
    pearson_r: float              # коэффициент корреляции Пирсона (для линейной)
    r_squared: float              # коэффициент детерминации (R_squared)
    approximated_y: np.ndarray    # вычисленные значения y
    errors: np.ndarray            # отклонения для каждой точки
