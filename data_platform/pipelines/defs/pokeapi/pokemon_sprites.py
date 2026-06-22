import dagster as dg
import polars as pl
from dagster import RetryPolicy, Backoff
from sqlalchemy.exc import OperationalError
from termcolor import colored

from ...utils.json_retriever import fetch_json
from ...utils.secret_retriever import fetch_secret

API_BASE: str = "https://api.github.com/repos/PokeAPI/sprites"
RAW_BASE: str = (
    "https://raw.githubusercontent.com/PokeAPI/sprites/refs/heads/master/sprites/pokemon"
)


def list_sprite_urls(parent_path: str, dir_name: str, ext: str, url_dir: str) -> dict[int, str]:
    """Map each Pokémon id to its raw sprite URL, for files named '<id>.<ext>' in a
    sprites subdirectory. The id is only used to line up each gif with its matching
    png; it is not written to the table (dbt derives that from the URL columns).
    Uses the Git Trees API because the Contents API caps listings at 1000 entries."""

    parent = fetch_json(f"{API_BASE}/contents/{parent_path}?ref=master")
    sha = next(entry["sha"] for entry in parent if entry["name"] == dir_name)
    tree = fetch_json(f"{API_BASE}/git/trees/{sha}")["tree"]

    suffix = f".{ext}"
    return {
        int(stem): f"{url_dir}/{node['path']}"
        for node in tree
        if node["type"] == "blob"
        and node["path"].endswith(suffix)
        and (stem := node["path"].removesuffix(suffix)).isdigit()
    }


def create_dataframe() -> pl.DataFrame:
    gifs = list_sprite_urls(
        "sprites/pokemon/other", "showdown", "gif", f"{RAW_BASE}/other/showdown"
    )
    pngs = list_sprite_urls("sprites", "pokemon", "png", RAW_BASE)
    ids = sorted(i for i in set(gifs) | set(pngs) if i >= 1)

    return pl.DataFrame(
        {
            "gif_sprite_url": [gifs.get(i) for i in ids],
            "png_sprite_url": [pngs.get(i) for i in ids],
        }
    )


@dg.asset(
    kinds={"Supabase", "Postgres"},
    retry_policy=RetryPolicy(max_retries=3, delay=2, backoff=Backoff.EXPONENTIAL),
)
def load_vg_pokemon_sprites() -> None:
    database_url: str = fetch_secret()
    table_name: str = "staging.vg_pokemon_sprites"
    df = create_dataframe()

    try:
        df.write_database(
            table_name=table_name, connection=database_url, if_table_exists="replace"
        )
        print(
            colored(" ✓", "green"),
            f"Data loaded into {table_name} ({df.height} rows)",
        )
    except OperationalError as e:
        print(colored(" ✖", "red"), "Connection error in load_vg_pokemon_sprites():", e)
        raise
