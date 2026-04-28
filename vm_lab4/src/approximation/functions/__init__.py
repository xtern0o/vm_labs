"""Модуль с реализациями функций аппроксимации"""

from .linear import LinearApproximation
from .polynomial import PolynomialApproximation
from .exponential import ExponentialApproximation
from .logarithmic import LogarithmicApproximation
from .power import PowerApproximation

__all__ = [
    'LinearApproximation',
    'PolynomialApproximation',
    'ExponentialApproximation',
    'LogarithmicApproximation',
    'PowerApproximation'
]
