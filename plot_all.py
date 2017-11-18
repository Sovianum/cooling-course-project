#!/usr/bin/python

import postprocessing.python.loaders as loaders
import postprocessing.python.plots as plots
import matplotlib.pyplot as plt
import sys
import os.path


if __name__ == '__main__':
    args = sys.argv

    img_dir = args[1]
    cycle_data_path = args[2]

    cycle_plot_name = "cycle_plot.png"

    cycle_df = loaders.get_max_eta_df(loaders.read_double_compressor_data(cycle_data_path))

    plots.plot_scheme_characteristics(cycle_df, 0.6)
    plt.savefig(os.path.join(img_dir, cycle_plot_name))
