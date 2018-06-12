#!/home/artem/anaconda3/bin/python

import postprocessing.python.loaders as loaders
import postprocessing.python.plots as plots
import postprocessing.python.savers as savers
import matplotlib.pyplot as plt
import sys
import os.path


if __name__ == '__main__':
    args = sys.argv

    img_dir = args[1]
    data_dir = args[2]

    cycle_data_path = os.path.join(data_dir, "3n.csv")
    cycle_data = loaders.read_double_compressor_data(cycle_data_path)

    cycle_eta_plot_name = "cycle_eta_plot.png"
    cycle_eta_df = loaders.get_max_eta_df(cycle_data)
    plots.plot_scheme_characteristics(cycle_eta_df, 0.65)
    plt.savefig(os.path.join(img_dir, cycle_eta_plot_name))
    plt.close()

    cycle_power_plot_name = "cycle_power_plot.png"
    cycle_power_df = loaders.get_max_power_df(cycle_data)
    plots.plot_scheme_characteristics(cycle_power_df, 0.6)
    plt.savefig(os.path.join(img_dir, cycle_power_plot_name))
    plt.close()

    def save_profile(data_name_1, data_name_2, plot_name):
        df_1 = loaders.read_profile_data(os.path.join(data_dir, data_name_1))
        df_2 = loaders.read_profile_data(os.path.join(data_dir, data_name_2))
        plots.plot_profile(df_1.x, df_1.y)
        plots.plot_profile(df_2.x, df_2.y)
        plt.grid()
        plt.savefig(os.path.join(img_dir, plot_name))
        plt.close()

    profile_data = [
        ["stator_root_1.csv", "stator_root_2.csv", "stator_root.png"],
        ["stator_mid_1.csv", "stator_mid_2.csv", "stator_mid.png"],
        ["stator_top_1.csv", "stator_top_2.csv", "stator_top.png"],
        ["rotor_root_1.csv", "rotor_root_2.csv", "rotor_root.png"],
        ["rotor_mid_1.csv", "rotor_mid_2.csv", "rotor_mid.png"],
        ["rotor_top_1.csv", "rotor_top_2.csv", "rotor_top.png"]
    ]
    [save_profile(tup[0], tup[1], tup[2]) for tup in profile_data]

    inlet_profile_df = loaders.read_profile_angles(os.path.join(data_dir, "inlet_angle.csv"))
    plots.plot_profile_angles(inlet_profile_df, [r"$\alpha_1$", r"$\beta_1$"])
    plt.savefig(os.path.join(img_dir, "inlet_angle.png"))
    plt.close()

    rotor_profile_df = loaders.read_profile_angles(os.path.join(data_dir, "outlet_angle.csv"))
    plots.plot_profile_angles(rotor_profile_df, [r"$\alpha_2$", r"$\beta_2$"])
    plt.savefig(os.path.join(img_dir, "outlet_angle.png"))
    plt.close()

    ss_profile_no_front_df = loaders.read_cooling_data(os.path.join(data_dir, "cooling_2_no_front_ss.json"))
    ps_profile_no_front_df = loaders.read_cooling_data(os.path.join(data_dir, "cooling_2_no_front_ps.json"))

    plots.plot_cooling_temperature(ps_profile_no_front_df, ss_profile_no_front_df)
    plt.savefig(os.path.join(img_dir, "cooling_2_t_no_front.png"))
    plt.close()

    ss_profile_front_df = loaders.read_cooling_data(os.path.join(data_dir, "cooling_2_front_ss.json"))
    ps_profile_front_df = loaders.read_cooling_data(os.path.join(data_dir, "cooling_2_front_ps.json"))

    plots.plot_cooling_temperature(ps_profile_front_df, ss_profile_front_df)
    plt.savefig(os.path.join(img_dir, "cooling_2_t_front.png"))
    plt.close()

    plots.plot_cooling_alpha(ps_profile_no_front_df, ss_profile_no_front_df)
    plt.savefig(os.path.join(img_dir, "cooling_2_alpha_no_front.png"))
    plt.close()

    plots.plot_cooling_alpha(ps_profile_front_df, ss_profile_front_df)
    plt.savefig(os.path.join(img_dir, "cooling_2_alpha_front.png"))
    plt.close()

    inventor_data = savers.get_temperature_df(ps_profile_no_front_df, ss_profile_no_front_df)
    # savers.to_csv_air(os.path.join(data_dir, "cooling_t_air.csv"), inventor_data)
    # savers.to_csv_wall(os.path.join(data_dir, "cooling_t_wall.csv"), inventor_data)
    # savers.to_csv_film(os.path.join(data_dir, "cooling_t_film.csv"), inventor_data)
    savers.to_csv_complex(os.path.join(data_dir, "cooling_t_complex.csv"), inventor_data)
    savers.to_csv_efficiency(os.path.join(data_dir, "cooling_t_efficiency.csv"), inventor_data)
