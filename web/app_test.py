from streamlit.testing.v1 import AppTest

at = AppTest.from_file("app.py").run(timeout=10)
at.run()


def test_header():
    assert at.header[0].value == "Pokémon TCG Tournament Results"


def test_selectbox_default_value():
    assert not at.exception
    assert len(at.selectbox) == 1
    assert at.selectbox[0].value == at.selectbox[0].options[0]


def test_selectbox_change():
    second_option = at.selectbox[0].options[1]
    at.selectbox[0].select(second_option).run()

    assert at.selectbox[0].value == second_option


def test_selectbox_label():
    assert at.selectbox[0].label == "Filter by Tournament *(ordered by date)*"


def test_tournament_info():
    assert not at.exception
    assert "•" in at.markdown[0].value


def test_metrics():
    assert at.metric[0].label == "Total Players"
    assert at.metric[1].label == "Winner"
    assert at.metric[0].value is not None
    assert at.metric[1].value is not None
