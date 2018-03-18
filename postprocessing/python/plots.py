
import matplotlib.pyplot as plt
import numpy as np
import pandas as pd


def plot_cooling_alpha(ps_data, ss_data):
    data = concat_profiles(ps_data, ss_data)

    plt.grid()

    plt.plot(data.l * 1e3, data.alpha_air, color='blue')
    plt.plot(data.l * 1e3, data.alpha_gas, color='red')

    plt.legend([r'$\alpha_{в}$', r'$\alpha_{пл}$'], loc='best')

    x_min = min(data.l) * 1e3
    x_max = max(data.l) * 1e3
    plt.xlim(x_min, x_max)

    t_1 = min(data.alpha_air)
    t_2 = max(data.alpha_gas)

    t_min = min([t_1, t_2])
    t_max = 825

    t_min -= (t_max - t_min) * 0.05

    plt.ylim(0, 2e3)

    plt.plot([0, 0], [t_min, t_max], color='black', lw=2)

    plt.text(0.6 * x_min, t_max - 50, r'$спинка$', fontsize=16)
    plt.text(0.4 * x_max, t_max - 50, r'$корыто$', fontsize=16)
    plt.xlabel(r'$x,\ мм$', fontsize=14)
    plt.ylabel(r'$\alpha,\ Вт/\left(м^2 \cdot К \right)$', fontsize=14)


def plot_cooling_temperature(ps_data, ss_data):
    data = concat_profiles(ps_data, ss_data)

    plt.grid()

    plt.plot(data.l * 1e3, data.t_wall, color='green')
    plt.plot(data.l * 1e3, data.t_air, color='blue')
    plt.plot(data.l * 1e3, data.t_film, color='red')

    plt.legend(['$T_{ст}$', '$T_{возд}$', '$T_{пл}$'], loc='best')

    x_min = min(data.l) * 1e3
    x_max = max(data.l) * 1e3
    plt.xlim(x_min, x_max)

    t_min = min(data.t_air)
    t_max = max(data.t_film)
    t_max += (t_max - t_min) * 0.05
    plt.ylim(t_min, t_max)

    plt.plot([0, 0], [t_min, t_max], color='black', lw=2)

    plt.text(0.6 * x_min, t_max - 50, r'$спинка$', fontsize=16)
    plt.text(0.4 * x_max, t_max - 50, r'$корыто$', fontsize=16)
    plt.xlabel(r'$x,\ мм$', fontsize=14)
    plt.ylabel(r'$T,\ К$', fontsize=14)


def concat_profiles(ps_data, ss_data) -> pd.DataFrame:
    data = pd.concat([ps_data, ss_data], ignore_index=True)
    data.l = pd.concat([ps_data.l, -ss_data.l], ignore_index=True)
    data.sort_values(by='l', inplace=True)
    return data


def plot_profile_angles(data, angle_names):
    plt.plot(np.rad2deg(data.angle_in), data.h)
    plt.plot(np.rad2deg(data.angle_out), data.h)
    plt.grid()
    plt.legend(angle_names, loc="best")
    plt.ylabel('$\overline{h}$')


def plot_scheme_characteristics(data, y_min=0.8, y_max=1.02):
    local_data = data[data.pi <= 40]
    plt.title('$Приведенные \ характеристики \ установки \ (\overline{f} = f / f_{max})$')
    plt.plot(local_data.pi, local_data.G / local_data.G.max(), '-bo', markevery=[20])
    plt.plot(local_data.pi, local_data.N_e / local_data.N_e.max(), '-go', markevery=[20])
    plt.plot(local_data.pi, local_data.eta / local_data.eta.max(), '-ro', markevery=[20])
    plt.plot([20, 20], [y_min, y_max], color='black')
    plt.ylim([y_min, y_max])
    plt.xlabel('$\pi_\Sigma$', fontsize=20)
    plt.grid()
    plt.legend(['$\overline{G}$', '$\overline{L}$', '$\overline{\eta}$'], loc='lower right')


def plot_height_parameter(parameter, parameter_name):
    h_rel = np.linspace(0, 1, len(parameter))
    plt.plot(parameter, h_rel)
    plt.grid()
    plt.xlabel(parameter_name)
    plt.ylabel('$\overline{h}$')


def plot_profile(profile_x, profile_y):
    plt.plot(profile_x, profile_y, color='blue')
    plt.axis('equal')