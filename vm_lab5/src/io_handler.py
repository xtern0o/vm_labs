import math
from typing import Callable, Tuple

def read_from_file(filepath: str) -> Tuple[list[float], list[float]]:
    """Чтение узлов из файла с обработкой ошибок."""
    x_nodes, y_nodes = [], []
    try:
        with open(filepath, 'r') as f:
            for line_num, line in enumerate(f, 1):
                if not line.strip():
                    continue
                parts = line.strip().split()
                if len(parts) != 2:
                    raise ValueError(f"Строка {line_num} должна содержать ровно два числа.")
                x_nodes.append(float(parts[0]))
                y_nodes.append(float(parts[1]))
        
        if len(x_nodes) < 2:
            raise ValueError("Файл должен содержать как минимум 2 точки.")
    except FileNotFoundError:
        raise FileNotFoundError(f"Файл '{filepath}' не найден.")
    except ValueError as e:
        raise ValueError(f"Ошибка парсинга данных: {e}")
        
    # Сортируем узлы по возрастанию координаты X
    zipped = sorted(zip(x_nodes, y_nodes), key=lambda p: p[0])
    return [p[0] for p in zipped], [p[1] for p in zipped]

def read_from_console() -> Tuple[list[float], list[float]]:
    """Чтение узлов с клавиатуры с обработкой ошибок."""
    while True:
        try:
            n = int(input("Введите количество точек (n >= 2): "))
            if n < 2:
                print("Ошибка: количество точек должно быть 2 или больше.")
                continue
            break
        except ValueError:
            print("Ошибка: введите целое число.")

    x_nodes, y_nodes = [], []
    print("Введите координаты точек (x y) парами через пробел:")
    i = 0
    while i < n:
        try:
            user_input = input(f"Точка {i+1}: ").replace(',', '.')
            parts = user_input.split()
            if len(parts) != 2:
                raise ValueError()
            
            x, y = float(parts[0]), float(parts[1])
            x_nodes.append(x)
            y_nodes.append(y)
            i += 1
        except ValueError:
            print("Ошибка: введите ровно два числа через пробел (например, '1.5 2.3').")
            
    # Сортируем узлы по возрастанию координаты X
    zipped = sorted(zip(x_nodes, y_nodes), key=lambda p: p[0])
    return [p[0] for p in zipped], [p[1] for p in zipped]

def generate_from_function(func: Callable[[float], float], start: float, end: float, n: int) -> Tuple[list[float], list[float]]:
    """Генерация узлов на основе выбранной функции."""
    if n < 2:
        raise ValueError("Нужно как минимум 2 точки.")
    if start >= end:
        raise ValueError("Начало интервала должно быть строго меньше конца.")
        
    step = (end - start) / (n - 1)
    x_nodes = [start + i * step for i in range(n)]
    y_nodes = [func(x) for x in x_nodes]
    return x_nodes, y_nodes

def print_difference_table(x_nodes: list[float], table: list[list[float]], is_finite: bool = False) -> None:
    """
    Выводит таблицу разделенных или конечных разностей.
    """
    n = len(table)
    title = "конечных" if is_finite else "разделенных"
    print(f"\nТаблица {title} разностей:")
    
    # header
    header = f"{'x':>8} | {'y':>8} | "
    for i in range(1, n):
        if is_finite:
            header += f"{'d^'+str(i)+'y':>8} | "
        else:
            header += f"{'f'+str(i):>8} | "
    print(header)
    print("-" * len(header))
    
    # rows
    for i in range(n):
        row_str = f"{x_nodes[i]:8.4f} | "
        for j in range(n - i):
            row_str += f"{table[i][j]:8.4f} | "
        print(row_str)
    print()
