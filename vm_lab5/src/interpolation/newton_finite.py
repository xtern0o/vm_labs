import math

def calculate_finite_differences(y: list[float]) -> list[list[float]]:
    """
    Строит таблицу конечных разностей (\delta^i y) для интерполяционного многочлена Ньютона
    для равноотстоящих узлов.
    """
    n = len(y)
    table = [[0.0] * n for _ in range(n)]
    
    for i in range(n):
        table[i][0] = y[i]
        
    for j in range(1, n):
        for i in range(n - j):
            table[i][j] = table[i + 1][j - 1] - table[i][j - 1]
            
    return table

def interpolate_newton_finite(x: list[float], y: list[float], target_x: float) -> tuple[float, str]:
    """
    Вычисляет интерполированное значение функции в точке target_x,
    используя многочлен Ньютона с конечными разностями (для равноотстящих узлов)
    Автоматически выбирает первую или вторую формулу в зависимости от target_x
    Возвращает кортеж: (результат, строка с названием метода)
    """
    table = calculate_finite_differences(y)
    n = len(x)
    h = x[1] - x[0]
    
    if target_x <= x[n // 2]:
        used = set()

        method_name = "первая интерполяционная формула Ньютона (вперед)"
        result = table[0][0]
        t = (target_x - x[0]) / h
        t_term = 1.0
        
        for i in range(1, n):
            t_term *= (t - i + 1)
            result += (t_term * table[0][i]) / math.factorial(i)
            used.add(table[0][i])

        # print(f"Использованы конечные разонсти: {used}")
    else:
        method_name = "вторая интерполяционная формула Ньютона (назад)"
        result = table[n - 1][0]
        t = (target_x - x[n - 1]) / h
        t_term = 1.0
        
        for i in range(1, n):
            t_term *= (t + i - 1)
            result += (t_term * table[n - 1 - i][i]) / math.factorial(i)
            
    return result, method_name
