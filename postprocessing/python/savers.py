
import pandas as pd


def to_csv_film(path: str, data: pd.DataFrame):
    to_csv_selector(path, data, ['l', 't_film'])


def to_csv_air(path: str, data: pd.DataFrame):
    to_csv_selector(path, data, ['l', 't_air'])


def to_csv_wall(path: str, data: pd.DataFrame):
    to_csv_selector(path, data, ['l', 't_wall'])


def to_csv_efficiency(path: str, data: pd.DataFrame):
    to_csv_selector(path, data, ['l', 'film_efficiency'])


def to_csv_selector(path: str, data: pd.DataFrame, selector):
    data[selector].to_csv(path, header=False, index=False)


def to_csv_complex(path: str, data: pd.DataFrame):
    data.to_csv(path, header=False, index=False)


def get_temperature_df(ps_data: pd.DataFrame, ss_data: pd.DataFrame) -> pd.DataFrame:
    data = get_united_df(ps_data, ss_data)

    result = data[['l', 't_film', 't_wall', 't_air', 'film_efficiency']]
    return result


def get_united_df(ps_data: pd.DataFrame, ss_data: pd.DataFrame) -> pd.DataFrame:
    data = pd.concat([ps_data, ss_data], ignore_index=True)
    data.l = pd.concat([ps_data.l, -ss_data.l], ignore_index=True)
    data.sort_values(by='l', inplace=True)
    data.l *= 1000
    data = data[data.index % 5 == 0]
    return data