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
        \pi = %.1f && \eta = %.3f && L_e = %.3f \ МДж/кг && G = %.1f \ кг/с
    \end{align}
    ''' % (pi, eta, power, mass_rate))


def get_3_shaft_nominal_parameters_note(df):
    return get_2_shaft_nominal_parameters_note(df)


# функция выводит график зависимости относительного КПД, относительного расхода и относительной мощности
# от степени повышения давления в компрессоре на номинальном режиме
def plot_nom_characteristic(df):
    norm_df = df / df.max()
    norm_df.pi = df.pi
    plt.title(
        '$Приведенная \ характеристика \ установки \ на \ номинальном \ режиме\ (\overline{f} = f / f_{max})$',
        fontsize=24
    )
    plt.plot(norm_df.pi, norm_df.mass_rate)
    plt.plot(norm_df.pi, norm_df.efficiency)
    plt.plot(norm_df.pi, norm_df.specific_power)
    plt.xlabel('$\pi$', fontsize=20)
    plt.xticks(fontsize=15)
    plt.yticks(fontsize=15)
    plt.grid()
    plt.ylim([0.83, 1.005])
    plt.legend(
        ['$\overline{G}$', '$\overline{\eta}$', '$\overline{L_e}$'], fontsize=20, loc='lower right',
    )


# функция выводит зависимость относительного КПД и относительного расхода установки в зависимости от относительной
# мощности установки
def plot_common_characteristics(df):
    norm_df = df / df.max()
    plt.title('$Приведенные \ характеристики \ установки \ (\overline{f} = f / f_{max})$', fontsize=24)
    plt.plot(norm_df.power, norm_df.mass_rate)
    plt.plot(norm_df.power, norm_df.eta)
    plt.xlabel('$\overline{N_e}$', fontsize=20)
    plt.xticks(fontsize=15)
    plt.yticks(fontsize=15)
    plt.xlim(0.3, 1.05)
    plt.ylim(0.5, 1.01)
    plt.grid()
    plt.legend(['$\overline{G}$', '$\overline{\eta}$'], loc='lower right', fontsize=20)


# функция строит сравнение некоторого параметра по относительной мощности
def plot_comparison(dfs, y_selector):
    for df in dfs:
        plt.plot(df.power / df.power.max(), df[y_selector])
    plt.xlabel('$\overline{N_e}$', fontsize=20)
    plt.xticks(fontsize=15)
    plt.yticks(fontsize=15)
    plt.xlim(0.3, 1.05)
    plt.grid()