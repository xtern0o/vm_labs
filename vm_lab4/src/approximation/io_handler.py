"""Модуль для ввода-вывода данных согласно заданию"""

import json
import numpy as np
from typing import Tuple, Dict, Any
from pathlib import Path
from .result import ApproximationResult


class InputHandler:
    """Обработка входных данных"""
    
    @staticmethod
    def input_data(filepath: str = "") -> Tuple[np.ndarray, np.ndarray]:
        """
        Ввести данные из файла или консоли
        
        Args:
            filepath: Путь к JSON файлу. Если пусто - ввод вручную
            
        Returns:
            (x, y) - массивы значений
        """
        if filepath and filepath.strip():
            return InputHandler._read_from_file(filepath)
        else:
            return InputHandler._read_from_console()
    
    @staticmethod
    def _read_from_file(filepath: str) -> Tuple[np.ndarray, np.ndarray]:
        """Чтение данных из JSON файла"""
        try:
            path = Path(filepath)
            if not path.exists():
                raise FileNotFoundError(f"Файл не найден: {filepath}")
            
            with open(filepath, 'r', encoding='utf-8') as f:
                data = json.load(f)
            
            if not isinstance(data, dict) or 'x' not in data or 'y' not in data:
                raise ValueError("JSON должен содержать поля 'x' и 'y'")
            
            x = np.array(data['x'], dtype=float)
            y = np.array(data['y'], dtype=float)
            
            # Проверка на дублирующиеся x
            if len(x) != len(set(x)):
                raise ValueError("Ошибка: одинаковые x значения! Пример: (1, 2) и (1, 3)")
            
            if len(x) != len(y):
                raise ValueError(f"Размер x ({len(x)}) != размер y ({len(y)})")
            
            if len(x) < 4 or len(x) > 12:
                raise ValueError(f"Требуется 4-12 точек, получено {len(x)}")
            
            print(f"Загружено {len(x)} точек из {filepath}")
            return x, y
            
        except FileNotFoundError as e:
            raise e
        except json.JSONDecodeError as e:
            raise ValueError(f"Ошибка чтения JSON: {e}")
        except (ValueError, TypeError) as e:
            raise ValueError(f"Ошибка данных: {e}")
    
    @staticmethod
    def _read_from_console() -> Tuple[np.ndarray, np.ndarray]:
        """Ввод данных с консоли"""
        print("Введите данные для аппроксимации (8-12 точек)\n")
        
        while True:
            try:
                n = int(input("Количество точек (8-12): "))
                if n < 4 or n > 12:
                    print("Требуется 4-12 точек\n")
                    continue
                break
            except ValueError:
                print("Введите целое число\n")
        
        x = []
        y = []
        
        print(f"\nВведите {n} точек в формате: x y (через пробел)\n")
        
        for i in range(n):
            while True:
                try:
                    line = input(f"Точка {i+1}: ").strip()
                    if not line:
                        print("Введите оба значения")
                        continue
                    
                    parts = line.split()
                    if len(parts) != 2:
                        print("Введите два числа")
                        continue
                    
                    xi, yi = float(parts[0]), float(parts[1])
                    x.append(xi)
                    y.append(yi)
                    break
                    
                except ValueError:
                    print("Некорректные числа")
        
        # Проверка на дублирующиеся x
        if len(x) != len(set(x)):
            raise ValueError("Ошибка: одинаковые x значения! Пример: (1, 2) и (1, 3)")
        
        x = np.array(x, dtype=float)
        y = np.array(y, dtype=float)
        
        print(f"\nВведено {n} точек")
        return x, y


class OutputHandler:
    """Обработка вывода результатов"""
    
    @staticmethod
    def output_results(
        x: np.ndarray,
        y: np.ndarray,
        results: Dict[str, ApproximationResult],
        best_name: str,
        best_result: ApproximationResult,
        filepath: str = ""
    ) -> None:
        """
        Вывести результаты в файл или консоль
        
        Args:
            x: Входные значения x
            y: Входные значения y
            results: Словарь результатов для всех функций
            best_name: Имя наилучшей функции
            best_result: Результат наилучшей функции
            filepath: Путь к выходному JSON файлу. Если пусто - консоль
        """
        output_data = OutputHandler._build_output_structure(
            x, y, results, best_name, best_result
        )
        
        if filepath and filepath.strip():
            OutputHandler._write_to_file(output_data, filepath)
        else:
            OutputHandler._print_to_console(output_data)
    
    @staticmethod
    def _build_output_structure(
        x: np.ndarray,
        y: np.ndarray,
        results: Dict[str, ApproximationResult],
        best_name: str,
        best_result: ApproximationResult
    ) -> Dict[str, Any]:
        """Построить структуру вывода"""
        
        functions_data = {}
        
        for name, result in results.items():
            func_data = {
                "name": name,
                "coefficients": result.coefficients.tolist(),
                "metrics": {
                    "deviation_sum": float(result.deviation_sum),
                    "sigma": float(result.sigma),
                    "r_squared": float(result.r_squared),
                }
            }
            
            # Коэффициент Пирсона только для линейной
            if result.pearson_r is not None:
                func_data["metrics"]["pearson_r"] = float(result.pearson_r)
            
            # Массивы значений для всех функций
            func_data["approximation_arrays"] = {
                "x_i": x.tolist(),
                "y_i": y.tolist(),
                "phi_x_i": result.approximated_y.tolist(),
                "errors_epsilon_i": result.errors.tolist()
            }
            
            functions_data[name] = func_data
        
        output_structure = {
            "input_data": {
                "points_count": len(x),
                "x_values": x.tolist(),
                "y_values": y.tolist()
            },
            "approximations": functions_data,
            "best_approximation": {
                "name": best_name,
                "coefficients": best_result.coefficients.tolist(),
                "metrics": {
                    "deviation_sum": float(best_result.deviation_sum),
                    "sigma": float(best_result.sigma),
                    "r_squared": float(best_result.r_squared),
                },
                "approximation_arrays": {
                    "x_i": x.tolist(),
                    "y_i": y.tolist(),
                    "phi_x_i": best_result.approximated_y.tolist(),
                    "errors_epsilon_i": best_result.errors.tolist()
                }
            }
        }
        
        # Добавим Пирсона для линейной в best если это линейная
        if best_name == "linear" and best_result.pearson_r is not None:
            output_structure["best_approximation"]["metrics"]["pearson_r"] = \
                float(best_result.pearson_r)
        
        return output_structure
    
    @staticmethod
    def _write_to_file(data: Dict[str, Any], filepath: str) -> None:
        """Сохранить результаты в JSON файл"""
        try:
            path = Path(filepath)
            path.parent.mkdir(parents=True, exist_ok=True)
            
            with open(filepath, 'w', encoding='utf-8') as f:
                json.dump(data, f, indent=2, ensure_ascii=False)
            
            print(f"Результаты сохранены в {filepath}")
            
        except IOError as e:
            print(f"Ошибка записи: {e}")
    
    @staticmethod
    def _print_to_console(data: Dict[str, Any]) -> None:
        """Вывести результаты в консоль"""
        
        print("\n=== РЕЗУЛЬТАТЫ АППРОКСИМАЦИИ ===\n")
        
        # Сводка входных данных
        input_data = data["input_data"]
        print(f"Количество точек: {input_data['points_count']}")
        
        # Таблица результатов
        print("\nРезультаты для всех функций:")
        print("-" * 70)
        print(f"{'Функция':<20} {'S':<15} {'sigma':<15} {'R2':<15}")
        print("-" * 70)
        
        approx = data["approximations"]
        best_name = data["best_approximation"]["name"]
        
        for name in sorted(approx.keys()):
            func = approx[name]
            s = func["metrics"]["deviation_sum"]
            sigma_val = func["metrics"]["sigma"]
            r2 = func["metrics"]["r_squared"]
            marker = " [ЛУЧШАЯ]" if name == best_name else ""
            print(f"{name:<20} {s:<15.6e} {sigma_val:<15.6e} {r2:<15.6f}{marker}")
        
        # Детали по каждой функции
        print("\n\nДетали по каждой функции:")
        print("-" * 70)
        
        for name in sorted(approx.keys()):
            func = approx[name]
            print(f"\n{name}:")
            
            # Коэффициенты
            coef_names = OutputHandler._get_coef_names(name, len(func["coefficients"]))
            for coef_name, coef_val in zip(coef_names, func["coefficients"]):
                print(f"  {coef_name} = {coef_val:.8f}")
            
            # Метрики
            metrics = func["metrics"]
            print(f"  S = {metrics['deviation_sum']:.6e}")
            print(f"  sigma = {metrics['sigma']:.6e}")
            print(f"  R2 = {metrics['r_squared']:.6f}")
            
            if "pearson_r" in metrics:
                print(f"  r = {metrics['pearson_r']:.6f}")
            
            # Таблица значений для этой функции
            if "approximation_arrays" in func:
                OutputHandler._print_table(func["approximation_arrays"])
        
        # Лучшая функция
        print("\n\nЛучшая аппроксимация: " + data["best_approximation"]["name"].upper())
        print("-" * 70)
        
        best = data["best_approximation"]
        coef_names = OutputHandler._get_coef_names(best['name'], len(best["coefficients"]))
        for coef_name, coef_val in zip(coef_names, best["coefficients"]):
            print(f"  {coef_name} = {coef_val:.8f}")
        
        metrics = best["metrics"]
        print(f"  S = {metrics['deviation_sum']:.6e}")
        print(f"  sigma = {metrics['sigma']:.6e}")
        print(f"  R2 = {metrics['r_squared']:.6f}")
        
        if "pearson_r" in metrics:
            print(f"  r = {metrics['pearson_r']:.6f}")
        
        # Таблица для лучшей аппроксимации
        if "approximation_arrays" in best:
            OutputHandler._print_table(best["approximation_arrays"])
        
        print()
    
    @staticmethod
    def _print_table(arrays: Dict[str, list]) -> None:
        """Вывести таблицу значений"""
        print("\n" + "=" * 80)
        print("Таблица значений")
        print("=" * 80)
        
        x_vals = arrays['x_i']
        y_vals = arrays['y_i']
        phi_vals = arrays['phi_x_i']
        eps_vals = arrays['errors_epsilon_i']
        
        # Заголовок
        print(f"{'i':<4} {'x':<12} {'y':<12} {'phi(x)':<12} {'eps_i':<12}")
        print("-" * 80)
        
        # Данные
        for i in range(len(x_vals)):
            print(f"{i+1:<4} {x_vals[i]:<12.6f} {y_vals[i]:<12.6f} {phi_vals[i]:<12.6f} {eps_vals[i]:<12.6f}")
        
        print("=" * 80)
    
    @staticmethod
    def _get_coef_names(func_name: str, count: int) -> list:
        """Get coefficient names for function"""
        if func_name == "linear":
            return ["a", "b"]
        elif func_name.startswith("polynomial"):
            return [f"a{i}" for i in range(count)]
        elif func_name == "exponential":
            return ["a", "b"]
        elif func_name == "logarithmic":
            return ["a", "b"]
        elif func_name == "power":
            return ["a", "b"]
        else:
            return [f"c{i}" for i in range(count)]
