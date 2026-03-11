import polars as pl
from app import PlayersCountrySection, RawStandingsSection
from streamlit.testing.v1 import AppTest

at = AppTest.from_file("app.py").run(timeout=10)


def test_header():
    assert at.header[0].value == "Pokémon TCG Tournament Data"


def test_tabs():
    assert len(at.tabs) == 2
    tab_labels = [t.label for t in at.tabs]
    assert "Overview" in tab_labels
    assert "Tournaments" in tab_labels


def test_overview_tab_header():
    assert at.header[1].value == "Tournament Locations"


def test_selectbox_default_value():
    assert not at.exception
    assert len(at.selectbox) == 1
    assert at.selectbox[0].index == 0


def test_selectbox_label():
    assert at.selectbox[0].label == "Filter by Tournament *(ordered by date)*"


def test_tournament_info():
    assert not at.exception
    assert "•" in at.markdown[0].value
    assert "flagcdn.com" in at.markdown[0].value


def test_dataframe_renders():
    assert not at.exception
    assert len(at.dataframe) == 1


def test_dataframe_columns():
    cols = at.dataframe[0].value.columns.tolist()
    # expected columns
    assert "rank" in cols
    assert "name" in cols
    assert "points" in cols
    assert "record" in cols
    assert "opp_win_percent" in cols
    assert "opp_opp_win_percent" in cols
    assert "deck" in cols
    assert "decklist" in cols
    assert "player_country" in cols
    # dropped columns
    assert "country_code" not in cols
    assert "iso_code" not in cols
    assert "logo" not in cols
    assert "player_quantity" not in cols
    assert "location" not in cols
    assert "start_date" not in cols
    assert "end_date" not in cols
    assert "text_date" not in cols
    assert "type" not in cols


def test_metrics():
    assert at.metric[0].label == "Total Players"
    assert at.metric[1].label == "Winner"
    assert at.metric[2].label == "Winning Deck"
    assert at.metric[0].value is not None
    assert at.metric[1].value is not None
    assert at.metric[2].value is not None


def test_dataframe_sorted_by_rank():
    ranks = at.dataframe[0].value["rank"].to_list()
    assert ranks == sorted(ranks)


def test_raw_standings_subheader():
    assert not at.exception
    subheader_values = [s.value for s in at.subheader]
    assert "Raw Standings" in subheader_values


SAMPLE_RAW_DF = pl.DataFrame(
    {
        "rank": [2, 1, 3],
        "name": ["Alice", "Bob", "Charlie"],
        "points": [20, 30, 10],
        "record": ["4-2", "5-1", "3-3"],
        "opp_win_percent": [0.6, 0.7, 0.5],
        "opp_opp_win_percent": [0.55, 0.65, 0.45],
        "deck": ["DeckA", "DeckB", "DeckC"],
        "decklist": ["http://a.com", "http://b.com", "http://c.com"],
        "player_country": ["US", "JP", "KR"],
    }
)


def test_raw_standings_init_stores_df():
    section = RawStandingsSection(SAMPLE_RAW_DF, "some-tourney-id")
    assert section.df.shape == SAMPLE_RAW_DF.shape


def test_raw_standings_init_stores_tourney_filter():
    section = RawStandingsSection(SAMPLE_RAW_DF, "some-tourney-id")
    assert section.tourney_filter == "some-tourney-id"


SAMPLE_DF = pl.DataFrame(
    {
        "player_country": ["US", "US", "US", "JP", "JP", "KR", "KR", "AU"],
        "points": [30, 20, 10, 25, 15, 18, 12, 5],
    }
)


def test_top_countries_by_count_order():
    section = PlayersCountrySection(SAMPLE_DF)
    assert section.countries[0] == "US"  # 3 players
    assert section.countries[1] in ("JP", "KR")  # 2 players each


def test_top_countries_by_count_limit():
    section = PlayersCountrySection(SAMPLE_DF)
    section.n = 2
    section.countries = section._top_countries_by_count()
    assert len(section.countries) == 2


def test_median_order_sorted_descending():
    section = PlayersCountrySection(SAMPLE_DF)
    median_map = (
        SAMPLE_DF.filter(pl.col("player_country").is_in(section.countries))
        .group_by("player_country")
        .agg(pl.median("points").alias("median_points"))
        .to_pandas()
        .set_index("player_country")["median_points"]
        .to_dict()
    )
    ordered_medians = [median_map[c] for c in section.median_order]
    assert ordered_medians == sorted(ordered_medians, reverse=True)
    assert section.median_order[-1] == "AU"  # AU median 5 — always last


def test_median_order_excludes_countries_not_in_top():
    section = PlayersCountrySection(SAMPLE_DF)
    assert all(c in section.countries for c in section.median_order)


def test_bar_chart_title():
    section = PlayersCountrySection(SAMPLE_DF)
    assert section._build_bar_chart().to_dict()["title"] == "Player Count by Country"


def test_bar_chart_encoding():
    section = PlayersCountrySection(SAMPLE_DF)
    spec = section._build_bar_chart().to_dict()
    assert spec["encoding"]["x"]["field"] == "player_country"
    assert spec["encoding"]["y"]["field"] == "player_count"


def test_box_chart_title():
    section = PlayersCountrySection(SAMPLE_DF)
    assert section._build_box_chart().to_dict()["title"] == "Points Spread by Country"


def test_box_chart_encoding():
    section = PlayersCountrySection(SAMPLE_DF)
    spec = section._build_box_chart().to_dict()
    assert spec["encoding"]["x"]["field"] == "player_country"
    assert spec["encoding"]["y"]["field"] == "points"


def test_box_chart_sort_order():
    section = PlayersCountrySection(SAMPLE_DF)
    spec = section._build_box_chart().to_dict()
    assert spec["encoding"]["x"]["sort"] == section.median_order


def test_box_chart_filters_to_top_countries():
    section = PlayersCountrySection(SAMPLE_DF)
    chart = section._build_box_chart()
    data_countries = set(chart.data["player_country"].tolist())
    assert data_countries == set(section.countries)
