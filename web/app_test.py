from streamlit.testing.v1 import AppTest

at = AppTest.from_file("app.py").run(timeout=10)
at.run()


def test_header():
    assert at.header[0].value == "Pokémon TCG Tournament Results"


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
    assert "player_quantity" not in cols
    assert "location" not in cols
    assert "start_date" not in cols
    assert "end_date" not in cols
    assert "text_date" not in cols
    assert "type" not in cols


def test_metrics():
    assert at.metric[0].label == "Total Players"
    assert at.metric[1].label == "Winner"
    assert at.metric[0].value is not None
    assert at.metric[1].value is not None
