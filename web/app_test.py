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
