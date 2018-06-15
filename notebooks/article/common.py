import matplotlib.pyplot as plt
from IPython.display import Math


def get_2_shaft_nominal_parameters_note(df):
    m = df.max()
    eta = m.get('efficiency')
    mass_rate = m.get('mass_rate')
    pi = m.get('pi')
    power = m.get('specific_power') / 1e6
    return Math(r'''
    \begin{align}
        \pi = %.1f && \eta_e = %.3f && L_e = %.3f \ МДж/кг && G = %.1f \ кг/с
    \end{align}
    ''' % (pi, eta, power, mass_rate))


def get_3_shaft_nominal_parameters_note(df):
    return get_2_shaft_nominal_parameters_note(df)


# функция выводит график зависимости относительного КПД, относительного расхода и относительной мощности
# от степени повышения давления в компрессоре на номинальном режиме
def plot_nom_characteristic(df, ymin=0.83, ymax=1.005):
    norm_df = df / df.max()
    norm_df.pi = df.pi
    # plt.title(
    #     '$Приведенная \ характеристика \ установки \ на \ номинальном \ режиме\ (\overline{f} = f / f_{max})$',
    #     fontsize=24
    # )
    plt.plot(norm_df.pi, norm_df.mass_rate, color='red', linewidth=2.0)
    plt.plot(norm_df.pi, norm_df.efficiency, color='green', linewidth=2.0)
    plt.plot(norm_df.pi, norm_df.specific_power, color='blue', linewidth=2.0)
    plt.xlabel('$\pi_\Sigma$', fontsize=35, position=(0.95, 0))
    plt.xticks(fontsize=25)
    plt.yticks(fontsize=25)
    plt.grid()
    plt.ylim([ymin, ymax])
    plt.legend(
        ['$\overline{G}$', '$\overline{\eta_e}$', '$\overline{L_e}$'], fontsize=30, loc='best',
    )


# функция выводит зависимость относительного КПД и относительного расхода установки в зависимости от относительной
# мощности установки
def plot_common_characteristics(df, ymin=0.5, ymax=1.01):
    norm_df = df / df.max()
    # plt.title('$Приведенные \ характеристики \ установки \ (\overline{f} = f / f_{max})$', fontsize=24)
    plt.plot(norm_df.power, norm_df.mass_rate, color='red', linewidth=2.0)
    plt.plot(norm_df.power, norm_df.eta, color='blue', linewidth=2.0)
    # plt.xlabel('$\overline{N_e}$', fontsize=30)
    plt.xticks(fontsize=25)
    plt.yticks(fontsize=25)
    plt.xlim(0.3, 1.05)
    plt.ylim(ymin, ymax)
    plt.grid()
    plt.legend(['$\overline{G}$', '$\overline{\eta_e}$'], loc='lower right', fontsize=30)

# функция строит сравнение некоторого параметра по относительной мощности
def plot_rel_comparison(dfs, y_selector):
    max_y = -1e10
    for df in dfs:
        m = df[y_selector].max()
        if m > max_y:
            max_y = m

    colors = ['red', 'green', 'blue']
    for df, id in zip(dfs, range(len(dfs))):
        plt.plot(df.power / df.power.max(), df[y_selector] / max_y, linewidth=2, color=colors[id % len(colors)])
    plt.xlabel('$\overline{N_e}$', fontsize=30)
    plt.xticks(fontsize=25)
    plt.yticks(fontsize=25)
    plt.xlim(0.3, 1.05)
    plt.grid()

# функция строит сравнение некоторого параметра по относительной мощности
def plot_comparison(dfs, y_selector, ymin=0.4, ymax=1):
    colors = ['red', 'green', 'blue']
    for df, id in zip(dfs, range(len(dfs))):
        plt.plot(df.power / df.power.max(), df[y_selector], linewidth=2, color=colors[id % len(colors)])
    # plt.xlabel('$\overline{N_e}$', fontsize=30)
    plt.xticks(fontsize=25)
    plt.yticks(fontsize=25)
    plt.xlim(0.3, 1.05)
    plt.ylim(ymin, ymax)
    plt.grid()