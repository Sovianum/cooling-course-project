#!/usr/bin/python

import postprocessing.python.loaders as loaders
import postprocessing.python.plots as plots
import matplotlib.pyplot as plt
import sys
import os.path


if __name__ == '__main__':
    args = sys.argv

    img_dir = args[1]
    data_dir = args[2]

    cycle_plot_name = "cycle_plot.png"
    cycle_data_path = os.path.join(data_dir, "3n.csv")
    cycle_df = loaders.get_max_eta_df(loaders.read_double_compressor_data(cycle_data_path))
    plots.plot_scheme_characteristics(cycle_df, 0.6)
    plt.savefig(os.path.join(img_dir, cycle_plot_name))
    plt.close()

    def save_profile(data_name, plot_name):
        df = loaders.read_profile_data(os.path.join(data_dir, data_name))
        plots.plot_profile(df.x, df.y)
        plt.savefig(os.path.join(img_dir, plot_name))
        plt.close()

    profile_data = [
        ["stator_root.csv", "stator_root.png"],
        ["stator_mid.csv", "stator_mid.png"],
        ["stator_top.csv", "stator_top.png"],
        ["rotor_root.csv", "rotor_root.png"],
        ["rotor_mid.csv", "rotor_mid.png"],
        ["rotor_top.csv", "rotor_top.png"]
    ]

    [save_profile(pair[0], pair[1]) for pair in profile_data]
