from streamlit.testing.v1 import AppTest


def test_selectbox_default_value():
    at = AppTest.from_file("app.py").run(timeout=10)

    assert not at.exception
    assert len(at.selectbox) == 1
    assert at.selectbox[0].value == at.selectbox[0].options[0]


def test_selectbox_change():
    at = AppTest.from_file("app.py").run(timeout=10)

    second_option = at.selectbox[0].options[1]
    at.selectbox[0].select(second_option).run()

    assert at.selectbox[0].value == second_option
