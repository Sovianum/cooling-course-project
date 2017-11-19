
import matplotlib.pyplot as plt
import numpy as np


def plot_scheme_characteristics(data, y_min=0.8, y_max=1.02):
    plt.title('$Приведенные \ характеристики \ установки \ (\overline{f} = f / f_{max})$')
    plt.plot(data.pi, data.G / data.G.max())
    plt.plot(data.pi, data.N_e / data.N_e.max())
    plt.plot(data.pi, data.eta / data.eta.max())
    plt.ylim([y_min, y_max])
    plt.grid()
    plt.legend(['$\overline{G}$', '$\overline{N_e}$', '$\overline{\eta}$'], loc='lower right')


def plot_height_parameter(parameter, parameter_name):
    h_rel = np.linspace(0, 1, len(parameter))
    plt.plot(parameter, h_rel)
    plt.grid()
    plt.xlabel(parameter_name)
    plt.ylabel('$\overline{h}$')


def plot_profile(profile_x, profile_y):
    plt.plot(profile_x, profile_y, color='blue')
    plt.grid()
    plt.axis('equal')