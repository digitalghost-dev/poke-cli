import altair as alt
import polars as pl
import pydeck
import streamlit as st

from supabase import Client, create_client


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
    result = (
        supabase.table("standings")
        .select("location, text_date")
        .order("start_date")
        .execute()
    )
    return list(
        dict.fromkeys((row["location"], row["text_date"]) for row in result.data)
    )


st.set_page_config(page_title="Pokémon Tournament Results", layout="wide")


def data_table(tourney_filter: str) -> pl.DataFrame:
    standings_table = pl.from_dicts(run_query(tourney_filter))

    return standings_table


def header() -> str:
    tourney_list = unique_locations()
    tournament_filter = st.selectbox(
        "Filter by Tournament *(ordered by date)*",
        tourney_list,
        format_func=lambda x: f"{x[0]} - {x[1]}",
    )
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


class PlayersCountrySection:
    def __init__(self, df: pl.DataFrame):
        self.df = df
        self.n = 10
        self.countries = self._top_countries_by_count()
        self.median_order = self._median_order_by_country()

    def _top_countries_by_count(self) -> list[str]:
        return (
            self.df.group_by("player_country")
            .agg(pl.len().alias("player_count"))
            .sort("player_count", descending=True)
            .head(self.n)
            .get_column("player_country")
            .to_list()
        )

    def _median_order_by_country(self) -> list[str]:
        return (
            self.df.filter(pl.col("player_country").is_in(self.countries))
            .group_by("player_country")
            .agg(pl.median("points").alias("median_points"))
            .sort("median_points", descending=True)
            .get_column("player_country")
            .to_list()
        )

    def _build_bar_chart(self) -> alt.Chart:
        countries_df = (
            self.df.group_by("player_country")
            .agg(pl.len().alias("player_count"))
            .head(15)
            .sort("player_count", descending=True)
        )
        return (
            alt.Chart(countries_df.to_pandas())
            .mark_bar()
            .encode(
                x=alt.X("player_country:N", sort="-y", title="Country"),
                y=alt.Y("player_count:Q", title="Players"),
            )
            .properties(title="Player Count by Country")
        )

    def _build_box_chart(self) -> alt.Chart:
        box_df = self.df.filter(pl.col("player_country").is_in(self.countries))
        return (
            alt.Chart(box_df.to_pandas())
            .mark_boxplot(extent="min-max")
            .encode(
                x=alt.X("player_country:N", sort=self.median_order, title="Country"),
                y=alt.Y("points:Q", title="Points"),
            )
            .properties(title="Points Spread by Country")
        )

    def render(self) -> None:
        st.divider()

        st.header("Stats per Country")
        st.subheader("Top 512 players", divider="blue")

        col1, col2 = st.columns(2)

        with col1:
            st.altair_chart(self._build_bar_chart(), width="stretch")

        with col2:
            st.altair_chart(self._build_box_chart(), width="stretch")


def tournament_locations() -> None:
    st.header("Tournament Locations")

    tournaments = (
        supabase.table("standings")
        .select(
            "location, tournament_latitude, tournament_longitude, type, text_date, player_quantity"
        )
        .eq("rank", 1)
        .order("start_date")
        .execute()
        .data
    )
    type_colors = {
        "International": [220, 50, 50, 200],
        "Regional": [255, 204, 0, 200],
        "Special Event": [50, 100, 220, 200],
        "World": [50, 200, 100, 200],
    }
    for t in tournaments:
        t["color"] = type_colors.get(t["type"], [200, 200, 200, 200])

    point_layer = pydeck.Layer(
        "ScatterplotLayer",
        data=tournaments,
        id="tournament-layer",
        get_position=["tournament_longitude", "tournament_latitude"],
        get_radius=200000,
        get_fill_color="color",
        pickable=True,
        auto_highlight=True,
    )

    view_state = pydeck.ViewState(latitude=15, longitude=10, zoom=1.3, controller=True)
    deck = pydeck.Deck(
        point_layer,
        initial_view_state=view_state,
        tooltip={"text": "{location}\n{type}\n{text_date}\nPlayers: {player_quantity}"},
    )

    st.pydeck_chart(deck, on_select="rerun", selection_mode="single-object")


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
    df = df.drop(
        [
            "country_code",
            "player_quantity",
            "iso_code",
            "logo",
            "location",
            "start_date",
            "end_date",
            "text_date",
            "type",
            "tournament_latitude",
            "tournament_longitude",
        ]
    )

    df = df.sort("rank")

    st.header("Raw Standings - Top 512")

    st.dataframe(
        df,
        column_config={
            "rank": st.column_config.NumberColumn(
                label="Rank",
                help="The player's placement in the tournament.",
            ),
            "name": st.column_config.TextColumn(
                label="Name",
            ),
            "points": st.column_config.NumberColumn(
                label="Points",
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
            ),
        },
    )


def main():
    st.header("Pokémon TCG Tournament Data")

    overview_tab, regionals_tab = st.tabs(["Overview", "Tournaments"])

    with overview_tab:
        tournament_locations()

    with regionals_tab:
        tourney_filter = header()
        tournament_info(tourney_filter)
        tournament_stats(tourney_filter)
        section = PlayersCountrySection(data_table(tourney_filter))
        section.render()
        display_latest_tournament(tourney_filter)


main()
