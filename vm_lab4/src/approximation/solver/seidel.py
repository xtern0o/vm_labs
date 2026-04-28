"""
Метод Гаусса-Зейделя для решения систем линейных уравнений
Переведено с Golang на Python с помощью Claude
"""

import numpy as np
from dataclasses import dataclass
from typing import List


@dataclass
class SeidelResult:
    """Результаты решения методом Гаусса-Зейделя"""
    solution: np.ndarray       # Найденное решение
    iterations: int            # Количество итераций
    errors: List[float]        # История погрешностей
    norm_of_matrix: float      # Норма матрицы C
    messages: List[str]        # Предупреждения и сообщения


class SeidelSolver:
    """Решатель систем линейных уравнений методом Гаусса-Зейделя"""
    
    def __init__(self, eps: float = 1e-14, max_iter: int = 1000):
        """
        Args:
            eps: Точность сходимости
            max_iter: Максимум итераций
        """
        self.eps = eps
        self.max_iter = max_iter
    
    @staticmethod
    def make_diagonally_dominant(A: np.ndarray, b: np.ndarray) -> bool:
        """
        Попытка переставить строки матрицы для получения диагонального преобладания.
        Используется для гарантии сходимости метода.
        
        Args:
            A: Матрица системы
            b: Вектор свободных членов
            
        Returns:
            True если удалось сделать матрицу диагонально доминантной, иначе False
        """
        n = len(A)
        used = [False] * n
        perm = []
        
        # Жадный алгоритм: для каждой строки найдем строку с максимальным диагональным элементом
        for i in range(n):
            best_row = -1
            max_diag = -1
            
            for j in range(n):
                if not used[j]:
                    # Проверяем диагональное преобладание
                    diag = abs(A[j, i])
                    sum_off_diag = sum(abs(A[j, k]) for k in range(n) if k != i)
                    
                    if diag > sum_off_diag and diag > max_diag:
                        best_row = j
                        max_diag = diag
            
            if best_row == -1:
                # Не удалось найти подходящую строку
                return False
            
            perm.append(best_row)
            used[best_row] = True
        
        # Переставляем строки согласно найденной перестановке
        A_new = np.zeros_like(A)
        b_new = np.zeros_like(b)
        for i, j in enumerate(perm):
            A_new[i] = A[j]
            b_new[i] = b[j]
        
        A[:] = A_new
        b[:] = b_new
        
        # Проверяем диагональное преобладание
        for i in range(n):
            if abs(A[i, i]) <= sum(abs(A[i, j]) for j in range(n) if j != i):
                return False
        
        return True
    
    @staticmethod
    def build_canonical_form(A: np.ndarray, b: np.ndarray) -> tuple:
        """
        Приведение системы Ax = b к каноническому виду x = Cx + d
        
        Args:
            A: Матрица системы (должна иметь диагональное преобладание)
            b: Вектор свободных членов
            
        Returns:
            (C, d) где x = Cx + d
        """
        n = len(A)
        C = np.zeros_like(A, dtype=float)
        d = np.zeros(n, dtype=float)
        
        for i in range(n):
            if abs(A[i, i]) < 1e-15:
                raise ValueError(f"Диагональный элемент A[{i},{i}] = 0, не могу разделить")
            
            d[i] = b[i] / A[i, i]
            
            for j in range(n):
                if i != j:
                    C[i, j] = -A[i, j] / A[i, i]
        
        return C, d
    
    @staticmethod
    def norm_of_matrix(C: np.ndarray) -> float:
        """Вычислить норму матрицы (max норма строк)"""
        n = len(C)
        max_norm = 0
        for i in range(n):
            row_sum = sum(abs(C[i, j]) for j in range(n))
            max_norm = max(max_norm, row_sum)
        return max_norm
    
    @staticmethod
    def max_norm_of_vector(v: np.ndarray) -> float:
        """Вычислить максимальную норму вектора"""
        return np.max(np.abs(v))
    
    def solve(self, A: np.ndarray, b: np.ndarray, x0: np.ndarray = None) -> SeidelResult:
        """
        Решить систему Ax = b методом Гаусса-Зейделя
        
        Args:
            A: Матрица системы (n x n)
            b: Вектор свободных членов (n,)
            x0: Начальное приближение (по умолчанию нулевое)
            
        Returns:
            SeidelResult с решением и информацией о сходимости
        """
        n = len(A)
        
        # Проверка размерности
        if len(b) != n:
            raise ValueError(f"Несоответствие размерности: A={n}x{n}, b={len(b)}")
        
        # Инициализация
        A_copy = A.copy().astype(float)
        b_copy = b.copy().astype(float)
        messages = []
        
        # Попытка переставить строки для диагонального преобладания
        if not self.make_diagonally_dominant(A_copy, b_copy):
            messages.append("Диагонального преобладания нет. Сходимость НЕ гарантирована")
        
        # Построение канонической формы
        try:
            C, d = self.build_canonical_form(A_copy, b_copy)
        except ValueError as e:
            raise ValueError(f"Ошибка построения канонической формы: {e}")
        
        # Проверка нормы матрицы C
        norm_C = self.norm_of_matrix(C)
        if norm_C >= 1:
            messages.append(f"||C|| = {norm_C:.6f} >= 1. Сходимость не гарантирована")
        
        # Инициализация переменных решения
        if x0 is None:
            x = np.zeros(n, dtype=float)
        else:
            x = x0.copy().astype(float)
        
        x_prev = np.zeros(n, dtype=float)
        errors = []
        
        # Итерационный процесс Гаусса-Зейделя
        for iteration in range(self.max_iter):
            x_prev = x.copy()
            
            # Обновление каждой переменной
            for i in range(n):
                sum_val = d[i]
                
                # Подстановка новых приближений (j < i)
                for j in range(i):
                    sum_val += C[i, j] * x[j]
                
                # Подстановка старых приближений (j > i)
                for j in range(i + 1, n):
                    sum_val += C[i, j] * x_prev[j]
                
                x[i] = sum_val
            
            # Вычисление погрешности
            diff = x - x_prev
            error_norm = self.max_norm_of_vector(diff)
            errors.append(error_norm)
            
            # Проверка сходимости
            if error_norm < self.eps:
                # Вычислим невязку
                residuals = A @ x - b
                max_residual = self.max_norm_of_vector(residuals)
                messages.append(f"  Сходимость достигнута за {iteration + 1} итераций")
                messages.append(f"  Максимальная невязка: {max_residual:.2e}")
                
                return SeidelResult(
                    solution=x,
                    iterations=iteration + 1,
                    errors=errors,
                    norm_of_matrix=norm_C,
                    messages=messages
                )
        
        # Не сошлось
        raise RuntimeError(
            f"Метод Гаусса-Зейделя не сошелся за {self.max_iter} итераций. "
            f"Последняя погрешность: {errors[-1]:.2e}"
        )
