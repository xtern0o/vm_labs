def interpolate_lagrange(x: list[float], y: list[float], target_x: float) -> float:
    """
    Вычисляет интерпол. значение ф. в т. target_x,
    используя интерп м-н Лагранжа
    """
    result = 0.0
    n = len(x)
    
    for i in range(n):
        term = y[i]
        for j in range(n):
            if i != j:
                term *= (target_x - x[j]) / (x[i] - x[j])
        result += term
        
    return result
