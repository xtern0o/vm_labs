"""Модуль для построения графиков аппроксимации"""

import matplotlib
# Используем неинтерактивный backend для работы без дисплея
matplotlib.use('Agg')

import os
import numpy as np
import matplotlib.pyplot as plt
from typing import Dict, Tuple, Optional
from datetime import datetime
from .result import ApproximationResult


class ApproximationPlotter:
    """Класс для построения графиков функций аппроксимации"""
    
    FUNCTION_NAMES = {
        'linear': 'Линейная: f(x) = a*x + b',
        'polynomial_2': 'Полином 2: f(x) = a_0 + a_1*x + a_2*x^2',
        'polynomial_3': 'Полином 3: f(x) = a_0 + a_1*x + a_2*x^2 + a_3*x^3',
        'exponential': 'Экспоненциальная: f(x) = a*exp(b*x)',
        'logarithmic': 'Логарифмическая: f(x) = a*ln(x) + b',
        'power': 'Степенная: f(x) = a*x^b'
    }
    
    def __init__(self, x: np.ndarray, y: np.ndarray, 
                 results: Dict[str, ApproximationResult],
                 best_name: str):
        """
        Args:
            x: Массив значений x
            y: Массив значений y
            results: Словарь с результатами аппроксимации
            best_name: Название лучшей аппроксимации
        """
        self.x = np.array(x, dtype=float)
        self.y = np.array(y, dtype=float)
        self.results = results
        self.best_name = best_name
        
        plt.style.use('seaborn-v0_8-darkgrid')
    
    def ask_plot_preferences(self) -> Tuple[bool, bool, str]:
        """
        Спросить пользователя о предпочтениях для построения графиков
        
        Returns:
            Tuple[bool, bool, str]: (все_на_одном_графике, нужен_best_fit, имя_файла)
        """
        print("\n\nПОСТРОЕНИЕ ГРАФИКОВ")
        print("-" * 70)
        
        single_plot = None
        while single_plot is None:
            choice = input("Все графики на одном рисунке? (y/n): ").strip().lower()
            if choice in ['y', 'yes', 'д', 'да', 'yes', 'yes']:
                single_plot = True
            elif choice in ['n', 'no', 'н', 'нет']:
                single_plot = False
            else:
                print("  WARNING: Пожалуйста, ответьте 'y' (да) или 'n' (нет)")
        
        best_fit = None
        while best_fit is None:
            choice = input("Построить отдельно график лучшей аппроксимации? (y/n): ").strip().lower()
            if choice in ['y', 'yes', 'д', 'да']:
                best_fit = True
            elif choice in ['n', 'no', 'н', 'нет']:
                best_fit = False
            else:
                print("  WARNING: Пожалуйста, ответьте 'y' (да) или 'n' (нет)")
        
        filename = input("Имя файла для сохранения графика (Enter для дефолтного): ").strip()
        
        if not filename:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            filename = f"approximation_plot_{timestamp}"
        
        plots_dir = "plots"
        if not os.path.exists(plots_dir):
            os.makedirs(plots_dir)
            print(f"  OK: Создана папка '{plots_dir}/'")
        
        filename = os.path.join(plots_dir, filename)
        
        return single_plot, best_fit, filename
    
    def plot_all_single(self) -> str:
        """
        Построить все графики на одном рисунке
        
        Returns:
            str: Путь к сохраненному файлу
        """
        print("  Построение графиков на одном рисунке...")
        
        fig, ax = plt.subplots(figsize=(14, 8))
        
        ax.scatter(self.x, self.y, color='black', s=80, alpha=0.6, 
                  label='Исходные данные', zorder=5)
        
        x_smooth = np.linspace(self.x.min(), self.x.max(), 300)
        
        colors = plt.cm.Set2(np.linspace(0, 1, len(self.results)))
        
        for idx, (func_name, result) in enumerate(self.results.items()):
            try:
                y_smooth = self._evaluate_function(func_name, x_smooth)
                ax.plot(x_smooth, y_smooth, 
                       color=colors[idx], linewidth=2.5, 
                       label=self.FUNCTION_NAMES.get(func_name, func_name),
                       alpha=0.8)
            except Exception as e:
                print(f"    Ошибка при построении {func_name}: {e}")
        
        ax.set_xlabel('x', fontsize=12, fontweight='bold')
        ax.set_ylabel('y', fontsize=12, fontweight='bold')
        ax.set_title('Результаты аппроксимации (все функции)', 
                    fontsize=14, fontweight='bold')
        ax.legend(loc='best', fontsize=10)
        ax.grid(True, alpha=0.3)
        
        plt.tight_layout()
        
        return self._save_plot(fig, "all_functions")
    
    def plot_all_separate(self) -> str:
        """
        Построить все графики на отдельных подграфиках (grid)
        
        Returns:
            str: Путь к сохраненному файлу
        """
        print("  Построение графиков на отдельных подграфиках...")
        
        num_plots = len(self.results) + 1  # +1 для лучшей аппроксимации
        cols = 3
        rows = (num_plots + cols - 1) // cols
        
        fig, axes = plt.subplots(rows, cols, figsize=(16, 4 * rows))
        axes = axes.flatten()
        
        x_smooth = np.linspace(self.x.min(), self.x.max(), 300)
        
        for idx, (func_name, result) in enumerate(self.results.items()):
            ax = axes[idx]
            
            try:
                # Исходные данные
                ax.scatter(self.x, self.y, color='black', s=60, alpha=0.6, 
                          label='Исходные данные', zorder=5)
                
                # Аппроксимация
                y_smooth = self._evaluate_function(func_name, x_smooth)
                
                color = 'red' if func_name == self.best_name else 'blue'
                linewidth = 3 if func_name == self.best_name else 2
                
                ax.plot(x_smooth, y_smooth, color=color, 
                       linewidth=linewidth, label='Аппроксимация', alpha=0.8)
                
                # Информация о функции
                title = self.FUNCTION_NAMES.get(func_name, func_name)
                if func_name == self.best_name:
                    title += " [BEST]"
                
                ax.set_title(title, fontsize=11, fontweight='bold')
                ax.set_xlabel('x', fontsize=10)
                ax.set_ylabel('y', fontsize=10)
                ax.legend(fontsize=9)
                ax.grid(True, alpha=0.3)
                
                # Добавить статистику на график
                stats_text = f"S = {result.deviation_sum:.2e}\n"
                stats_text += f"sigma = {result.sigma:.2e}\n"
                stats_text += f"R^2 = {result.r_squared:.4f}"
                ax.text(0.02, 0.98, stats_text, transform=ax.transAxes,
                       fontsize=8, verticalalignment='top',
                       bbox=dict(boxstyle='round', facecolor='wheat', alpha=0.5))
                
            except Exception as e:
                ax.text(0.5, 0.5, f"Ошибка: {e}", 
                       ha='center', va='center', transform=ax.transAxes)
                ax.set_title(self.FUNCTION_NAMES.get(func_name, func_name),
                           fontsize=11, fontweight='bold', color='red')
        
        for idx in range(len(self.results), len(axes)):
            fig.delaxes(axes[idx])
        
        fig.suptitle('Результаты аппроксимации (отдельные функции)',
                    fontsize=14, fontweight='bold', y=0.995)
        plt.tight_layout()
        
        return self._save_plot(fig, "all_separate")
    
    def plot_best_comparison(self) -> str:
        """
        Построить график лучшей аппроксимации с исходными данными
        
        Returns:
            str: Путь к сохраненному файлу
        """
        print("  Построение графика лучшей аппроксимации...")
        
        fig, ax = plt.subplots(figsize=(12, 7))
        
        ax.scatter(self.x, self.y, color='black', s=100, alpha=0.7,
                  label='Исходные данные', zorder=5, edgecolors='black', linewidth=1.5)
        
        x_smooth = np.linspace(self.x.min(), self.x.max(), 300)
        
        try:
            y_smooth = self._evaluate_function(self.best_name, x_smooth)
            ax.plot(x_smooth, y_smooth, color='red', linewidth=3,
                   label='Аппроксимация', alpha=0.8)
        except Exception as e:
            print(f"  Ошибка при построении лучшей функции: {e}")
        
        best_result = self.results[self.best_name]
        
        ax.set_xlabel('x', fontsize=12, fontweight='bold')
        ax.set_ylabel('y', fontsize=12, fontweight='bold')
        
        title = f"Лучшая аппроксимация: {self.FUNCTION_NAMES.get(self.best_name, self.best_name)}"
        ax.set_title(title, fontsize=14, fontweight='bold')
        
        stats_text = f"S (сумма квадратов) = {best_result.deviation_sum:.6e}\n"
        stats_text += f"sigma = {best_result.sigma:.6e}\n"
        stats_text += f"R^2 (детерминация) = {best_result.r_squared:.6f}"
        if best_result.pearson_r is not None:
            stats_text += f"\nr (корреляция) = {best_result.pearson_r:.6f}"
        
        ax.text(0.02, 0.98, stats_text, transform=ax.transAxes,
               fontsize=11, verticalalignment='top',
               bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))
        
        ax.legend(fontsize=11, loc='best')
        ax.grid(True, alpha=0.3)
        
        plt.tight_layout()
        
        return self._save_plot(fig, "best_fit")
    
    def _evaluate_function(self, func_name: str, x_values: np.ndarray) -> np.ndarray:
        """
        Вычислить значения функции для заданных x
        
        Args:
            func_name: Название функции
            x_values: Массив значений x
            
        Returns:
            np.ndarray: Вычисленные значения y
        """
        result = self.results[func_name]
        coef = result.coefficients
        
        if func_name == 'linear':
            return coef[0] * x_values + coef[1]
        
        elif func_name == 'polynomial_2':
            return coef[0] + coef[1] * x_values + coef[2] * x_values**2
        
        elif func_name == 'polynomial_3':
            return coef[0] + coef[1] * x_values + coef[2] * x_values**2 + coef[3] * x_values**3
        
        elif func_name == 'exponential':
            return coef[0] * np.exp(coef[1] * x_values)
        
        elif func_name == 'logarithmic':
            return coef[0] * np.log(x_values) + coef[1]
        
        elif func_name == 'power':
            return coef[0] * x_values**coef[1]
        
        else:
            raise ValueError(f"Неизвестная функция: {func_name}")
    
    def _save_plot(self, fig, suffix: str = "") -> str:
        """
        Сохранить график в файл
        
        Args:
            fig: Figure объект matplotlib
            suffix: Суффикс для добавления к имени файла
            
        Returns:
            str: Путь к сохраненному файлу
        """
        filename = getattr(self, '_filename', None)
        if not filename:
            # timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            filename = os.path.join("plots", f"approximation_plot")
        
        if suffix:
            dirname = os.path.dirname(filename)
            basename = os.path.basename(filename)
            filename = os.path.join(dirname, f"{basename}_{suffix}")
        
        filepath_png = f"{filename}.png"
        filepath_pdf = f"{filename}.pdf"
        
        try:
            fig.savefig(filepath_png, dpi=300, bbox_inches='tight')
            fig.savefig(filepath_pdf, bbox_inches='tight')
            print(f"  OK: Сохранено: {filename}.*")
        except Exception as e:
            print(f"  ERROR: Ошибка при сохранении: {e}")
            filepath_png = None
        
        plt.close(fig)
        
        return filepath_png if filepath_png else filename
    
    def build_plots(self, single_plot: bool, build_best_fit: bool, filename: str) -> None:
        """
        Построить графики согласно предпочтениям пользователя
        
        Args:
            single_plot: Все ли графики на одном рисунке
            build_best_fit: Нужно ли строить график лучшей аппроксимации
            filename: Имя файла для сохранения
        """
        self._filename = filename
        
        print(f"\nPLOTTING: (имя: {filename})...")
        
        try:
            if single_plot:
                print("  Режим: Все функции на одном графике")
                self.plot_all_single()
            else:
                print("  Режим: Отдельные подграфики в grid")
                self.plot_all_separate()
            
            if build_best_fit:
                print("  Дополнительно: График лучшей аппроксимации")
                self.plot_best_comparison()
            
            print("\nOK: Все графики успешно построены и сохранены!")
        
        except Exception as e:
            print(f"\nERROR: Ошибка при построении графиков: {e}")
            import traceback
            traceback.print_exc()
