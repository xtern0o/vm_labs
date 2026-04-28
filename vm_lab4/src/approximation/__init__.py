"""
Модуль для аппроксимации функций методом наименьших квадратов.
"""

from .result import ApproximationResult
from .base import ApproximationFunction
from .functions import (
    LinearApproximation,
    PolynomialApproximation,
    ExponentialApproximation,
    LogarithmicApproximation,
    PowerApproximation
)
from .analyzer import ApproximationAnalyzer
from .io_handler import InputHandler, OutputHandler
from .plotter import ApproximationPlotter

__all__ = [
    'ApproximationResult',
    'ApproximationFunction',
    'LinearApproximation',
    'PolynomialApproximation',
    'ExponentialApproximation',
    'LogarithmicApproximation',
    'PowerApproximation',
    'ApproximationAnalyzer',
    'InputHandler',
    'OutputHandler',
    'ApproximationPlotter'
]
