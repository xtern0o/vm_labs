import math
import sys

from io_handler import (
    read_from_console, 
    read_from_file, 
    generate_from_function, 
    print_difference_table
)
from interpolation import (
    interpolate_lagrange,
    interpolate_newton_divided,
    interpolate_newton_divided2,
    calculate_divided_differences,
    interpolate_newton_finite,
    calculate_finite_differences,
)
from plotter import plot_graph

def f1(x: float) -> float:
    return math.sin(x)

def f2(x: float) -> float:
    return x**2 + 3 * x - 1

FUNCTIONS = {
    "1": ("sin(x)", f1),
    "2": ("x^2 + 3x - 1", f2)
}

def check_evenly_spaced(x_nodes: list[float]) -> bool:
    """Проверка, являются ли узлы равноотстоящими."""
    if len(x_nodes) < 2:
        return True
    h = x_nodes[1] - x_nodes[0]
    for i in range(2, len(x_nodes)):
        if not math.isclose(x_nodes[i] - x_nodes[i - 1], h, rel_tol=1e-5, abs_tol=1e-5):
            return False
    return True

def main():
    print("--- Лабораторная работа №5. Интерполяция ---")
    print("Выберите способ ввода данных:")
    print("1. С клавиатуры")
    print("2. Из файла")
    print("3. На основе функции")
    
    choice = input("Ваш выбор: ").strip()
    x_nodes, y_nodes = [], []
    
    if choice == "1":
        x_nodes, y_nodes = read_from_console()
    elif choice == "2":
        while True:
            filepath = input("Введите путь к файлу: ").strip()
            try:
                x_nodes, y_nodes = read_from_file(filepath)
                break
            except Exception as e:
                print(f"Ошибка: {e}")
    elif choice == "3":
        print("Доступные функции:")
        for k, v in FUNCTIONS.items():
            print(f"{k}. {v[0]}")
        f_choice = input("Выберите функцию: ").strip()
        if f_choice not in FUNCTIONS:
            print("Неверный выбор функции.")
            return
            
        start = float(input("Введите начало интервала a: ").replace(',', '.'))
        end = float(input("Введите конец интервала b: ").replace(',', '.'))
        n = int(input("Введите количество точек: "))
        
        try:
            x_nodes, y_nodes = generate_from_function(FUNCTIONS[f_choice][1], start, end, n)
        except Exception as e:
            print(f"Ошибка: {e}")
            return
    else:
        print("Неверный выбор.")
        return

    # Вывод исходной таблицы
    print("\nИсходные узлы:")
    for x, y in zip(x_nodes, y_nodes):
        print(f"x: {x:8.4f} | y: {y:8.4f}")

    if len(set(x_nodes)) != len(x_nodes):
        print("\nОшибка: Все значения x (узлы интерполяции) должны быть уникальными.")
        return

    is_even = check_evenly_spaced(x_nodes)
    if is_even:
        print("\nУзлы равноотстоящие. Будут рассчитаны конечные разности.")
        finite_table = calculate_finite_differences(y_nodes)
        print_difference_table(x_nodes, finite_table, is_finite=True)
    else:
        print("\nУзлы неравноотстоящие. Будут рассчитаны разделенные разности.")
        divided_table = calculate_divided_differences(x_nodes, y_nodes)
        print_difference_table(x_nodes, divided_table, is_finite=False)

    try:
        target_x = float(input("\nВведите значение x для интерполирования: ").replace(',', '.'))
    except ValueError:
        print("Некорректное значение x.")
        return

    if target_x < min(x_nodes) or target_x > max(x_nodes):
        print(f"\n[ВНИМАНИЕ] Выполняется экстраполирование, так как {target_x} выходит за пределы [{min(x_nodes)}, {max(x_nodes)}]")
    else:
        print(f"\nВыполняется интерполирование.")

    # Вычисление Лагранжем
    res_lagrange = interpolate_lagrange(x_nodes, y_nodes, target_x)
    print(f"Значение по многочлену Лагранжа: {res_lagrange:.6f}")

    # Вычисление Ньютоном
    if is_even:
        res_newton, newton_method = interpolate_newton_finite(x_nodes, y_nodes, target_x)
        print(f"Использована {newton_method}")
        print(f"Значение по многочлену Ньютона (конечные разности): {res_newton:.6f}")
        
    else:
        res_newton = interpolate_newton_divided(x_nodes, y_nodes, target_x)
        print(f"Значение по многочлену Ньютона (разделенные разности): {res_newton:.6f}")
        res_newton2 = interpolate_newton_divided2(x_nodes, y_nodes, target_x)
        print(f"Значение по многочлену Ньютона ДВА (разделенные разности): {res_newton2:.6f}")


    print(f"Разница между методами: {abs(res_lagrange - res_newton):.6e}")

    # Отрисовка графика
    # В зависимости от равномерности сетки передаем нужную реализацию Ньютона
    func_newton = interpolate_newton_finite if is_even else interpolate_newton_divided
    plot_graph(x_nodes, y_nodes, target_x, res_lagrange, interpolate_lagrange, func_newton)

if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("\nВыход из программы.")
        sys.exit(0)
