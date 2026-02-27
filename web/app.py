import polars as pl
import streamlit as st
from supabase import create_client, Client


@st.cache_resource
def init_connection() -> Client:
    url = st.secrets["SUPABASE_URL"]
    key = st.secrets["SUPABASE_KEY"]
    return create_client(url, key)


supabase = init_connection()


@st.cache_data(ttl=86400)
def run_query(location: str) -> list:
    return (
        supabase.table("standings").select("*").eq("location", location).execute().data
    )


@st.cache_data(ttl=86400)
def unique_locations() -> list:
    result = supabase.table("standings").select("location, text_date").order("start_date").execute()
    return list(dict.fromkeys((row["location"], row["text_date"]) for row in result.data))


st.set_page_config(page_title="Pokémon Tournament Results", layout="wide")


def data_table(tourney_filter: str) -> pl.DataFrame:
    standings_table = pl.from_dicts(run_query(tourney_filter))

    return standings_table


def header() -> str:
    with st.container():
        col1, col2 = st.columns(2)

        with col1:
            st.header("Pokémon TCG Tournament Results")

        with col2:
            tourney_list = unique_locations()
            tournament_filter = st.selectbox(
                "Filter by Tournament *(ordered by date)*",
                tourney_list,
                format_func=lambda x: f"{x[0]} - {x[1]}",
            )

        st.divider()

    return tournament_filter[0]


def tournament_info(tourney_filter: str):
    df = data_table(tourney_filter)

    date_text = df["text_date"][0]
    location = df["location"][0]
    st.markdown(f"### {location} • {date_text}")


def tournament_stats(tourney_filter: str) -> None:
    df = data_table(tourney_filter)

    players = df["player_quantity"].to_list()
    winner = df.filter(pl.col("rank") == 1)["name"][0]

    with st.container():

        col1, col2, col3, col4 = st.columns(4, border=True)

        with col1:
            st.metric(label="Total Players", value=players[0])

        with col2:
            st.metric(label="Winner", value=winner)


def display_latest_tournament(tourney_filter: str) -> None:
    df = data_table(tourney_filter)
    df = df.drop(["country_code", "player_quantity", "location", "start_date", "end_date", "text_date", "type"])

    df = df.sort("rank")

    st.dataframe(
        df,
        column_config={
            "decklist": st.column_config.LinkColumn(
                label="Decklist", display_text=":material/open_in_new:"
            )
        },
    )


def main():
    tourney_filter = header()
    tournament_info(tourney_filter)
    tournament_stats(tourney_filter)
    display_latest_tournament(tourney_filter)

main()

