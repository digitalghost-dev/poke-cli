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
    result = supabase.table("standings").select("location").execute()
    return sorted({row["location"] for row in result.data})


st.set_page_config(page_title="Pokémon Tournament Results", layout="wide")


def header() -> str:
    with st.container():
        col1, col2 = st.columns(2)

        with col1:
            st.header("Pokémon TCG Tournament Results")

        with col2:
            tourney_list = unique_locations()
            tourney_filter = st.selectbox(
                "Filter by tournament",
                tourney_list,
            )

        st.divider()

    return tourney_filter


def display_latest_tournament() -> None:
    tourney_filter = header()

    df = pl.from_dicts(run_query(tourney_filter))
    st.dataframe(
        df,
        column_config={
            "decklist": st.column_config.LinkColumn(
                label="Decklist", display_text=":material/open_in_new:"
            )
        },
    )


display_latest_tournament()
