"""
Главный файл программы аппроксимации функции
методом наименьших квадратов
"""

import sys
from approximation import ApproximationAnalyzer, InputHandler, OutputHandler, ApproximationPlotter


def print_header():
    """Вывести заголовок"""
    print("\n" + "=" * 70)
    print("ЛР 4: АППРОКСИМАЦИЯ ФУНКЦИИ МЕТОДОМ НАИМЕНЬШИХ КВАДРАТОВ")
    print("=" * 70 + "\n")


def main():
    """Главная программа"""
    
    print_header()
    
    input_file = input("Путь к JSON файлу (Enter для ввода вручную): ").strip()
    
    try:
        x, y = InputHandler.input_data(input_file if input_file else "")
    except Exception as e:
        print(f"Ошибка ввода: {e}")
        return 1
    
    print("Выполняется аппроксимация:")
    print("  1. Линейная:       f(x) = a*x + b")
    print("  2. Полином 2:      f(x) = a0 + a1*x + a2*x^2")
    print("  3. Полином 3:      f(x) = a0 + a1*x + a2*x^2 + a3*x^3")
    print("  4. Экспоненциальная: f(x) = a*exp(b*x)")
    print("  5. Логарифмическая: f(x) = a*ln(x) + b")
    print("  6. Степенная:      f(x) = a*x^b\n")
    
    try:
        analyzer = ApproximationAnalyzer(x, y)
        results = analyzer.fit_all()
    except Exception as e:
        print(f"Ошибка аппроксимации: {e}")
        return 1
    
    print("\nВЫБОР ЛУЧШЕЙ АППРОКСИМАЦИИ")
    print("-" * 70)
    
    try:
        best_name, best_result = analyzer.get_best_approximation()
    except Exception as e:
        print(f"Ошибка выбора: {e}")
        return 1
    
    print(f"Лучшая функция: {best_name}")
    print(f"  S (сумма квадратов) = {best_result.deviation_sum:.6e}")
    print(f"  sigma (сигма) = {best_result.sigma:.6e}")
    print(f"  R2 (детерминация) = {best_result.r_squared:.6f}")
    
    print("\n\nВЫВОД РЕЗУЛЬТАТОВ")
    print("-" * 70)
    
    output_file = input("Путь для сохранения JSON (Enter для консоли): ").strip()
    
    try:
        OutputHandler.output_results(
            x, y, results, best_name, best_result,
            output_file if output_file else ""
        )
    except Exception as e:
        print(f"Ошибка вывода: {e}")
        return 1
    
    # Построение графиков
    try:
        plotter = ApproximationPlotter(x, y, results, best_name)
        single_plot, build_best_fit, filename = plotter.ask_plot_preferences()
        plotter.build_plots(single_plot, build_best_fit, filename)
    except Exception as e:
        print(f"Ошибка при построении графиков: {e}")
        import traceback
        traceback.print_exc()
        return 1
    
    return 0


if __name__ == "__main__":
    try:
        exit_code = main()
        sys.exit(exit_code)
    except KeyboardInterrupt:
        print("\n\nПрограмма прервана")
        sys.exit(1)
    except Exception as e:
        print(f"\nКритическая ошибка: {e}")
        sys.exit(1)
