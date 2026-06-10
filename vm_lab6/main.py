import math
import os
import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt


EQUATIONS = [
    {
        "name": "y' = y",
        "f": lambda x, y: y,
        "exact": lambda x, x0, y0: y0 * math.exp(x - x0),
    },
    {
        "name": "y' = x + y",
        "f": lambda x, y: x + y,
        "exact": lambda x, x0, y0: (y0 + x0 + 1) * math.exp(x - x0) - x - 1,
    },
    {
        "name": "y' = x - y",
        "f": lambda x, y: x - y,
        "exact": lambda x, x0, y0: (y0 - x0 + 1) * math.exp(-(x - x0)) + x - 1,
    },
    {
        "name": "y' = x * y",
        "f": lambda x, y: x * y,
        "exact": lambda x, x0, y0: y0 * math.exp((x**2 - x0**2) / 2),
    },
    {
        "name": "y' = x^2 - y",
        "f": lambda x, y: x**2 - y,
        "exact": lambda x, x0, y0: (y0 - x0**2 + 2*x0 - 2) * math.exp(-(x - x0)) + x**2 - 2*x + 2,
    },
]


def euler_improved(f, x0, y0, xn, h):
    """Улучшенный метод Эйлера 2 порядка"""
    xs = [x0]
    ys = [y0]
    x, y = x0, y0
    while x + h <= xn + 1e-9:
        y_pred = y + h * f(x, y)
        y = y + h / 2 * (f(x, y) + f(x + h, y_pred))
        x = x + h
        xs.append(x)
        ys.append(y)
    return xs, ys


def runge_kutta4(f, x0, y0, xn, h):
    """Рунге-Кутта 4 порядка"""
    xs = [x0]
    ys = [y0]
    x, y = x0, y0
    while x + h <= xn + 1e-9:
        k1 = h * f(x, y)
        k2 = h * f(x + h/2, y + k1/2)
        k3 = h * f(x + h/2, y + k2/2)
        k4 = h * f(x + h, y + k3)
        y = y + (k1 + 2*k2 + 2*k3 + k4) / 6
        x = x + h
        xs.append(x)
        ys.append(y)
    return xs, ys


def milne(f, x0, y0, xn, h, eps):
    if xn - x0 < 3 * h:
        print("ошибка!!! интервал слишком мал для метода Милна: нужно xn - x0 >= 3*h")
        return [x0], [y0]

    xs, ys = runge_kutta4(f, x0, y0, x0 + 3*h, h)

    if len(xs) < 4:
        print("ошибка!!! интервал слишком мал для метода Милна (нужно минимум 4 точки)")
        return xs, ys

    fs = [f(xs[i], ys[i]) for i in range(len(xs))]

    x = xs[-1]
    while x + h <= xn + 1e-9:
        f3, f2, f1 = fs[-3], fs[-2], fs[-1]

        y_pred = ys[-4] + 4*h/3 * (2*f3 - f2 + 2*f1)
        x_new = x + h
        f_pred = f(x_new, y_pred)

        # корректор итерирую до сходимости
        y_corr = ys[-2] + h/3 * (f2 + 4*f1 + f_pred)
        while abs(y_corr - y_pred) >= eps:
            y_pred = y_corr
            f_pred = f(x_new, y_pred)
            y_corr = ys[-2] + h/3 * (f2 + 4*f1 + f_pred)
            if abs(y_corr - y_pred) >= eps:
                print(f"[Милн x={x_new}]: не хватило точности для сходимости: |y_corr - y_pred| >= eps {y_corr:.4f} - {y_pred:.4f} = {abs(y_corr - y_pred):.4f}  >= {eps}")
                print(f"[Милг x={x_new}]: проводим повторную коррекцию")
        
        print(f"[Милн x={x_new}]: коррекция завершена: y_corr - y_pred < eps: {y_corr:.4f} - {y_pred:.4f} = {abs(y_corr - y_pred):.4f}  < {eps}")
        print()

        xs.append(x_new)
        ys.append(y_corr)
        fs.append(f(x_new, y_corr))
        x = x_new

    return xs, ys

# --- error estimation ---

def runge_error(method, f, x0, y0, xn, h, p):
    """Правило Рунге: R = |y^h - y^(h/2)| / (2^p - 1) в x_n"""
    _, ys_h = method(f, x0, y0, xn, h)
    _, ys_h2 = method(f, x0, y0, xn, h / 2)

    print("--- Вычисление правила Рунге ---")

    print(f"значение в y_n при h: {ys_h[-1]}")
    print(f"значение в y_n при h/2: {ys_h2[-1]} <- более точное решение (порядок {p} при шаге h/2 = {h/2})")
    print(f"p={p}")
    print(f"Формула: |{ys_h[-1]} - {ys_h2[-1]}| / 2^({p}-1)")

    return abs(ys_h[-1] - ys_h2[-1]) / (2**p - 1)


def exact_error(xs, ys, exact, x0, y0):
    """max |y_точн - y_i| по всем узлам."""
    return max(abs(exact(xs[i], x0, y0) - ys[i]) for i in range(len(xs)))


# --- output ---

def print_table(xs, ys, exact, x0, y0, h_orig=None):
    if h_orig is not None:
        idxes = [i for i, x in enumerate(xs) if abs(round((x - x0) / h_orig) * h_orig - (x - x0)) < 1e-9]
    else:
        idxes = list(range(len(xs)))
    print(f"{'x':>10}  {'y_прибл':>14}  {'y_точн':>14}  {'|y_прибл - y_точн|':>12}")
    for i in idxes:
        ye = exact(xs[i], x0, y0)
        print(f"{xs[i]:>10.4f}  {ys[i]:>14.8f}  {ye:>14.8f}  {abs(ye - ys[i]):>12.2e}")


def plot_results(results, exact, x0, y0, eq_name, save_dir):
    fig, ax = plt.subplots(figsize=(8, 5))
    fig.suptitle(eq_name)

    xs_ref = results[0][1]
    n_dense = 300
    xs_dense = [xs_ref[0] + i * (xs_ref[-1] - xs_ref[0]) / n_dense for i in range(n_dense + 1)]
    ys_exact = [exact(x, x0, y0) for x in xs_dense]
    ax.plot(xs_dense, ys_exact, 'k-', linewidth=2, label='точное')

    colors = ['tab:blue', 'tab:orange', 'tab:green']
    for (label, xs, ys), color in zip(results, colors):
        ax.plot(xs, ys, 'o--', color=color, markersize=4, linewidth=1, label=label)

    ax.set_xlabel('x')
    ax.set_ylabel('y')
    ax.set_title('Решения')
    ax.legend()
    ax.grid(True, alpha=0.3)

    plt.tight_layout()

    os.makedirs(save_dir, exist_ok=True)
    path = os.path.join(save_dir, "plot.png")
    plt.savefig(path, dpi=150)
    plt.close()
    print(f"График сохранен: {os.path.abspath(path)}")


# --- main ---

def select_equation():
    print("Выберите уравнение:")
    for i, eq in enumerate(EQUATIONS):
        print(f"  {i+1}. {eq['name']}")
    while True:
        try:
            idx = int(input("Ваш выбор: ")) - 1
            if 0 <= idx < len(EQUATIONS):
                return EQUATIONS[idx]
            print("Неверный номер.")
        except ValueError:
            print("Введите число.")


def get_float(prompt):
    while True:
        try:
            return float(input(prompt))
        except ValueError:
            print("Введите число.")


def get_positive_float(prompt):
    while True:
        val = get_float(prompt)
        if val > 0:
            return val
        print("Значение должно быть положительным.")


def get_eps():
    raw = input("eps (точность) [по умолчанию 0.001]: ").strip()
    if raw == "":
        return 0.001
    while True:
        try:
            val = float(raw)
            if val > 0:
                return val
            print("Значение должно быть положительным.")
        except ValueError:
            print("Введите число.")
        raw = input("eps (точность) [по умолчанию 0.001]: ").strip()
        if raw == "":
            return 0.001


def main():
    print("ЛР6 - Численное решение ОДУ")
    print("Вариант 5: Улучшенный Эйлер, РК4, Милн")

    eq = select_equation()
    f = eq["f"]
    exact = eq["exact"]

    x0 = get_float("x0 (начальная точка): ")
    y0 = get_float("y0 = y(x0): ")

    while True:
        xn = get_float("xn (конец интервала): ")
        if xn > x0:
            break
        print("xn должен быть больше x0.")

    h = get_positive_float("h (шаг): ")
    eps = get_eps()

    save_dir = input("Папка для сохранения графика [по умолчанию ./results]: ").strip()
    if save_dir == "":
        save_dir = "./results"

    methods = [
        ("Мод. Эйлер (порядок 2)", euler_improved, 2),
        ("Рунге-Кутта 4 (порядок 4)", runge_kutta4,  4),
        ("Милн (порядок 4)", milne, 4),
    ]

    MAX_ITER = 20
    results = []

    for name, method, order in methods:
        print(f"\n--- {name} ---")

        h_cur = h

        if method is milne:
            xs, ys = method(f, x0, y0, xn, h_cur, eps)
            print_table(xs, ys, exact, x0, y0, h)
            err = exact_error(xs, ys, exact, x0, y0)
            print(f"макс. погрешность (vs точное): {err:.2e}")
            # pass
        else:
            for i in range(MAX_ITER):
                print(f"\nИтерация {i}, h = {h_cur}")
                xs, ys = method(f, x0, y0, xn, h_cur)
                r = runge_error(method, f, x0, y0, xn, h_cur, order)
                print(f"оценка по правилу Рунге: {r:.2e}")
                if r <= eps:
                    print(f"точность {eps:.2e} достигнута. ЗАВЕРШАЕМСЯ")
                    print_table(xs, ys, exact, x0, y0, h)
                    break
                if i < MAX_ITER - 1:
                    h_cur /= 2
                    print(f"В методе {name} не хватило точности - уменьшаем h до {h_cur:.6f}")
            else:
                print(f"точность {eps:.2e} для метода {name} НЕ достигнута за {MAX_ITER} итераций уменьшения h!!")

        results.append((name, xs, ys))

    plot_results(results, exact, x0, y0, eq["name"], save_dir)


if __name__ == "__main__":
    main()