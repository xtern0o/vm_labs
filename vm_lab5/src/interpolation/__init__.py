"""
Interpolation modules for Lab 5.
"""
from .lagrange import interpolate_lagrange
from .newton_divided import interpolate_newton_divided, calculate_divided_differences, interpolate_newton_divided2
from .newton_finite import interpolate_newton_finite, calculate_finite_differences

__all__ = [
    "interpolate_lagrange",
    "interpolate_newton_divided",
    "calculate_divided_differences",
    "interpolate_newton_finite",
    "calculate_finite_differences",
    "interpolate_newton_divided2"
]
