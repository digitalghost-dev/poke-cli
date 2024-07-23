import requests
import polars as pl


def call_api(pokemon_id) -> tuple[int, str, str, str, list[str], list[str]]:
    url = f"https://pokeapi.co/api/v2/pokemon/{pokemon_id}"
    r = requests.get(url)

    id = r.json()["id"]
    name = r.json()["name"]
    height = r.json()["height"]
    weight = r.json()["weight"]

    types = [t["type"]["name"] for t in r.json()["types"]]
    abilities = [a["ability"]["name"] for a in r.json()["abilities"]]

    return id, height, weight, name, types, abilities


def build_dataframe():
    data = []
    for i in range(906, 1026):
        id, height, weight, name, types, abilities = call_api(i)

        types += [None] * (2 - len(types))
        abilities += [None] * (3 - len(abilities))

        data.append({
            "id": id,
            "height": height,
            "weight": weight,
            "name": name,
            "type_1": types[0],
            "type_2": types[1],
            "ability_1": abilities[0],
            "ability_2": abilities[1],
            "ability_3": abilities[2],
        })

    df = pl.DataFrame(data)

    with pl.Config(tbl_cols=-1):
        return df


def upload_dataframe():
    dataframe = build_dataframe()

    uri = "postgresql://pokemon-db-user:n3w-database-pw-1234@127.0.0.1:5432/pokemon-db"

    dataframe.write_database(table_name="pokemon-schema.pokemon"
                             , connection=uri, if_table_exists="append")

if __name__ == "__main__":
    upload_dataframe()