import polars as pl
from streamlit.testing.v1 import AppTest

from app import build_bar_chart, build_box_chart, median_order_by_country, top_countries_by_count

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


SAMPLE_DF = pl.DataFrame({
    "player_country": ["US", "US", "US", "JP", "JP", "KR", "KR", "AU"],
    "points": [30, 20, 10, 25, 15, 18, 12, 5],
})


def test_top_countries_by_count_order():
    result = top_countries_by_count(SAMPLE_DF, n=3)
    assert result[0] == "US"  # 3 players
    assert result[1] in ("JP", "KR")  # 2 players each
    assert len(result) == 3


def test_top_countries_by_count_limit():
    result = top_countries_by_count(SAMPLE_DF, n=2)
    assert len(result) == 2


def test_median_order_sorted_descending():
    countries = ["US", "JP", "KR", "AU"]
    result = median_order_by_country(SAMPLE_DF, countries)
    # Rebuild medians in the returned order to verify descending sort
    median_map = (
        SAMPLE_DF.filter(pl.col("player_country").is_in(countries))
        .group_by("player_country")
        .agg(pl.median("points").alias("median_points"))
        .to_pandas()
        .set_index("player_country")["median_points"]
        .to_dict()
    )
    ordered_medians = [median_map[c] for c in result]
    assert ordered_medians == sorted(ordered_medians, reverse=True)
    assert "AU" == result[-1]  # AU median 5 — always last


def test_median_order_filters_to_given_countries():
    result = median_order_by_country(SAMPLE_DF, ["US", "JP"])
    assert "KR" not in result
    assert "AU" not in result


def test_bar_chart_title():
    chart = build_bar_chart(SAMPLE_DF)
    assert chart.to_dict()["title"] == "Player Count by Country"


def test_bar_chart_encoding():
    spec = build_bar_chart(SAMPLE_DF).to_dict()
    assert spec["encoding"]["x"]["field"] == "player_country"
    assert spec["encoding"]["y"]["field"] == "player_count"


def test_box_chart_title():
    order = ["JP", "US", "KR", "AU"]
    chart = build_box_chart(SAMPLE_DF, ["US", "JP", "KR", "AU"], order)
    assert chart.to_dict()["title"] == "Points Spread by Country"


def test_box_chart_encoding():
    order = ["JP", "US", "KR", "AU"]
    spec = build_box_chart(SAMPLE_DF, ["US", "JP", "KR", "AU"], order).to_dict()
    assert spec["encoding"]["x"]["field"] == "player_country"
    assert spec["encoding"]["y"]["field"] == "points"


def test_box_chart_sort_order():
    order = ["JP", "US", "KR"]
    spec = build_box_chart(SAMPLE_DF, ["US", "JP", "KR"], order).to_dict()
    assert spec["encoding"]["x"]["sort"] == order


def test_box_chart_filters_to_given_countries():
    order = ["JP", "US"]
    chart = build_box_chart(SAMPLE_DF, ["JP", "US"], order)
    data_countries = set(chart.data["player_country"].tolist())
    assert data_countries == {"JP", "US"}
    assert "KR" not in data_countries
