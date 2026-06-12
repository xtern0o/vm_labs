import numpy as np
import matplotlib.pyplot as plt
from typing import Callable

def plot_graph(
    x_nodes: list[float], 
    y_nodes: list[float], 
    target_x: float, 
    target_y: float, 
    func_lagrange: Callable[[list[float], list[float], float], float],
    func_newton: Callable[[list[float], list[float], float], float]
):
    """
    Построение графика исходных точек, интерполяционной кривой и искомой точки.
    """
    plt.figure(figsize=(10, 6))
    
    # Исходные узлы
    plt.scatter(x_nodes, y_nodes, color='blue', label='Узлы интерполяции', zorder=5)
    
    # Искомая точка
    plt.scatter([target_x], [target_y], color='red', s=100, label=f'Искомая точка ({target_x}, {target_y:.4f})', zorder=6)
    
    # Гладкая кривая
    min_x, max_x = min(x_nodes), max(x_nodes)
    margin = (max_x - min_x) * 0.1
    # Если margin 0 (только одна точка), зададим дефолтный
    if margin == 0:
        margin = 1.0
        
    start_plot_x = min(min_x - margin, target_x)
    end_plot_x = max(max_x + margin, target_x)
        
    plot_x = np.linspace(start_plot_x, end_plot_x, 200)
    
    # Вычисляем Y для разных методов
    plot_y_lagrange = [func_lagrange(x_nodes, y_nodes, px) for px in plot_x]
    
    # Обработка для Ньютона (теперь возвращает кортеж, если это конечные разности)
    plot_y_newton = []
    for px in plot_x:
        res = func_newton(x_nodes, y_nodes, px)
        if isinstance(res, tuple):
            plot_y_newton.append(res[0])
        else:
            plot_y_newton.append(res)
    
    # Лагранж
    plt.plot(plot_x, plot_y_lagrange, color='green', linewidth=4, alpha=0.5, label='Лагранж', linestyle='-')
    # Ньютон
    plt.plot(plot_x, plot_y_newton, color='orange', linewidth=2, label='Ньютон', linestyle='--')
    
    plt.title('Сравнение методов интерполяции')
    plt.xlabel('x')
    plt.ylabel('y')
    plt.grid(True)
    plt.legend()
    plt.savefig("plots/graph.pdf")
    plt.savefig("plots/graph.png")

    # plt.show()
