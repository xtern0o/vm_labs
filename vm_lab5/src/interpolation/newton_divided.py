def calculate_divided_differences(x: list[float], y: list[float]) -> list[list[float]]:
    """
    Строит таблицу разделенных разностей для интерполяционного многочлена Ньютона
    для неравноотстоящих узлов
    """
    n = len(y)
    table = [[0.0] * n for _ in range(n)]
    
    for i in range(n):
        table[i][0] = y[i]
        
    for j in range(1, n):
        for i in range(n - j):
            table[i][j] = (table[i + 1][j - 1] - table[i][j - 1]) / (x[i + j] - x[i])
            
    return table

def interpolate_newton_divided(x: list[float], y: list[float], target_x: float) -> float:
    """
    Вычисляет интерполированное значение функции в точке target_x,
    используя многочлен Ньютона с разд. разностями
    """
    table = calculate_divided_differences(x, y)
    result = table[0][0]
    term = 1.0
    
    for i in range(1, len(x)):
        term *= (target_x - x[i - 1])
        result += table[0][i] * term
        
    return result

def interpolate_newton_divided2(x: list[float], y: list[float], target_x: float) -> float:
    """
    Вычисляет интерполированное значение функции в точке target_x,
    используя многочлен Ньютона с разд. разностями
    """
    table = calculate_divided_differences(x, y)
    n = len(x)
    result = table[n - 1][0]
    term = 1.0
    
    for i in range(1, len(x)):
        term *= (target_x - x[n - i])
        result += table[n - i - 1][i] * term
        # print(f"2: -> {table[n - i - 1][i]}")
    return result
