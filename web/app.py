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

    iso_code = df["iso_code"][0]
    flag = f'<img src="https://flagcdn.com/w40/{iso_code}.png"> ' if iso_code else ""
    logo = df["logo"][0]
    date_text = df["text_date"][0]
    location = df["location"][0]

    with st.container(horizontal=True):
        st.markdown(f"### {flag} • {location}\n{date_text}", unsafe_allow_html=True)
        st.space("stretch")
        if logo:
            st.image(logo, width=100)

def tournament_stats(tourney_filter: str) -> None:
    df = data_table(tourney_filter)

    players = df["player_quantity"].to_list()
    winner = df.filter(pl.col("rank") == 1)["name"][0]
    winning_deck = df.filter(pl.col("rank") == 1)["deck"][0]

    with st.container():

        col1, col2, col3 = st.columns(3, border=True)

        with col1:
            st.metric(label="Total Players", value=players[0])

        with col2:
            st.metric(label="Winner", value=winner)

        with col3:
            st.metric(label="Winning Deck", value=winning_deck.capitalize())


def display_latest_tournament(tourney_filter: str) -> None:
    df = data_table(tourney_filter)
    df = df.drop(["country_code", "player_quantity", "iso_code", "logo", "location", "start_date", "end_date", "text_date", "type"])

    df = df.sort("rank")

    st.dataframe(
        df,
        column_config={
            "rank": st.column_config.NumberColumn(
                label="Rank",
                format="plain",
                help="The player's placement in the tournament.",
            ),
            "name": st.column_config.TextColumn(
                label="Name",
            ),
            "points": st.column_config.NumberColumn(
                label="Points",
                format="plain",
                help="The player's total points in the tournament.",
            ),
            "record": st.column_config.TextColumn(
                label="Record",
                help="The player's record in the tournament.",
            ),
            "opp_win_percent": st.column_config.NumberColumn(
                label="OPW%",
                format="percent",
                help="The player's opponent's win percentage in the tournament.",
            ),
            "opp_opp_win_percent": st.column_config.NumberColumn(
                label="OOPW%",
                format="percent",
                help="The player's opponent's opponent's win percentage in the tournament.",
            ),
            "deck": st.column_config.TextColumn(
                label="Deck",
                help="The player's deck in the tournament.",
            ),
            "decklist": st.column_config.LinkColumn(
                label="Decklist", display_text=":material/open_in_new:"
            ),
            "player_country": st.column_config.TextColumn(
                label="Country",
                help="The player's home country.",
            )
        },
    )


def main():
    tourney_filter = header()
    tournament_info(tourney_filter)
    tournament_stats(tourney_filter)
    display_latest_tournament(tourney_filter)

main()

