import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import pytest
import polars as pl
import responses
from pydantic import ValidationError
from unittest.mock import patch

from pipelines.defs.extract.tcgcsv.extract_pricing import (
    build_dataframe,
    extract_card_name,
    get_card_number,
    is_card,
    normalize_card_number,
    pull_product_information,
    SET_PRODUCT_MATCHING,
)

_build_dataframe_fn = build_dataframe.op.compute_fn.decorated_fn


# ---------------------------------------------------------------------------
# is_card()
# ---------------------------------------------------------------------------


def test_is_card_with_number_field(benchmark):
    item = {"extendedData": [{"name": "Number", "value": "025/195"}]}
    assert benchmark(is_card, item) is True  # nosec


def test_is_card_no_extended_data(benchmark):
    assert benchmark(is_card, {}) is False  # nosec


def test_is_card_empty_extended_data(benchmark):
    assert benchmark(is_card, {"extendedData": []}) is False  # nosec


def test_is_card_no_number_field(benchmark):
    item = {"extendedData": [{"name": "Color", "value": "Yellow"}]}
    assert benchmark(is_card, item) is False  # nosec


# ---------------------------------------------------------------------------
# get_card_number()
# ---------------------------------------------------------------------------


def test_get_card_number_found(benchmark):
    card = {"extendedData": [{"name": "Number", "value": "025/195"}]}
    assert benchmark(get_card_number, card) == "025/195"  # nosec


def test_get_card_number_no_extended_data(benchmark):
    assert benchmark(get_card_number, {}) is None  # nosec


def test_get_card_number_no_number_field(benchmark):
    card = {"extendedData": [{"name": "HP", "value": "60"}]}
    assert benchmark(get_card_number, card) is None  # nosec


def test_get_card_number_no_value_key(benchmark):
    card = {"extendedData": [{"name": "Number"}]}
    assert benchmark(get_card_number, card) is None  # nosec


# ---------------------------------------------------------------------------
# normalize_card_number()
# ---------------------------------------------------------------------------


def test_normalize_card_number_single_digit(benchmark):
    assert benchmark(normalize_card_number, "1/149") == "001/149"  # nosec


def test_normalize_card_number_double_digit(benchmark):
    assert benchmark(normalize_card_number, "10/149") == "010/149"  # nosec


def test_normalize_card_number_triple_digit(benchmark):
    assert benchmark(normalize_card_number, "100/149") == "100/149"  # nosec


def test_normalize_card_number_non_numeric_parts(benchmark):
    assert benchmark(normalize_card_number, "GG01/GG70") == "GG01/GG70"  # nosec


def test_normalize_card_number_no_slash_passthrough(benchmark):
    assert benchmark(normalize_card_number, "SWSH001") == "SWSH001"  # nosec


def test_normalize_card_number_already_padded(benchmark):
    assert benchmark(normalize_card_number, "001/149") == "001/149"  # nosec


# ---------------------------------------------------------------------------
# extract_card_name()
# ---------------------------------------------------------------------------


def test_extract_card_name_plain(benchmark):
    assert benchmark(extract_card_name, "Pikachu") == "Pikachu"  # nosec


def test_extract_card_name_strip_dash_variant(benchmark):
    assert benchmark(extract_card_name, "Pikachu - 045/195") == "Pikachu"  # nosec


def test_extract_card_name_strip_parenthetical_number(benchmark):
    assert benchmark(extract_card_name, "Pikachu (010)") == "Pikachu"  # nosec


def test_extract_card_name_strip_full_art(benchmark):
    assert benchmark(extract_card_name, "Charizard (Full Art)") == "Charizard"  # nosec


def test_extract_card_name_strip_secret(benchmark):
    assert benchmark(extract_card_name, "Charizard (Secret)") == "Charizard"  # nosec


def test_extract_card_name_strip_reverse_holofoil(benchmark):
    assert benchmark(extract_card_name, "Pikachu (Reverse Holofoil)") == "Pikachu"  # nosec


def test_extract_card_name_strip_gold(benchmark):
    assert benchmark(extract_card_name, "Pikachu (Gold)") == "Pikachu"  # nosec


def test_extract_card_name_accented_characters(benchmark):
    assert benchmark(extract_card_name, "Flabébé - 088/195") == "Flabebe"  # nosec


def test_extract_card_name_dash_and_variant_suffix(benchmark):
    assert benchmark(extract_card_name, "Pikachu - 045/195 (Full Art)") == "Pikachu"  # nosec


def test_extract_card_name_unknown_variant_not_stripped(benchmark):
    # variants not in the hardcoded list are left in place
    assert benchmark(extract_card_name, "Pikachu (Shiny)") == "Pikachu (Shiny)"  # nosec


# ---------------------------------------------------------------------------
# pull_product_information()
# ---------------------------------------------------------------------------


def _make_product(product_id: int, name: str, card_number: str) -> dict:
    return {
        "productId": product_id,
        "name": name,
        "extendedData": [{"name": "Number", "value": card_number}],
    }


def _make_price(product_id: int, market_price: float | None, sub_type: str = "Normal") -> dict:
    return {"productId": product_id, "marketPrice": market_price, "subTypeName": sub_type}


@responses.activate
def test_pull_product_information_success(benchmark):
    product_id = "22873"  # sv01
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/products",
        json={"results": [
            _make_product(1001, "Pikachu", "025/198"),
            _make_product(1002, "Charizard", "006/198"),
        ]},
        status=200,
    )
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/prices",
        json={"results": [
            _make_price(1001, 1.50),
            _make_price(1002, None),
        ]},
        status=200,
    )

    df = benchmark(pull_product_information, "sv01")

    assert isinstance(df, pl.DataFrame)  # nosec
    assert len(df) == 2  # nosec
    assert df.filter(pl.col("name") == "Pikachu")["market_price"].to_list() == [1.50]  # nosec
    assert df.filter(pl.col("name") == "Charizard")["market_price"].to_list() == [None]  # nosec


@responses.activate
def test_pull_product_information_skips_variants(benchmark):
    product_id = "22873"  # sv01
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/products",
        json={"results": [
            _make_product(1001, "Pikachu", "025/198"),
            _make_product(1002, "Pikachu (Poke Ball Pattern)", "025/198"),
            _make_product(1003, "Pikachu (Master Ball Pattern)", "025/198"),
        ]},
        status=200,
    )
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/prices",
        json={"results": [
            _make_price(1001, 2.00),
            _make_price(1002, 3.00),
            _make_price(1003, 4.00),
        ]},
        status=200,
    )

    df = benchmark(pull_product_information, "sv01")

    assert len(df) == 1  # nosec
    assert df["name"].to_list() == ["Pikachu"]  # nosec


@responses.activate
def test_pull_product_information_skips_non_cards(benchmark):
    product_id = "22873"  # sv01
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/products",
        json={"results": [
            _make_product(1001, "Pikachu", "025/198"),
            {"productId": 1002, "name": "Booster Pack", "extendedData": []},
        ]},
        status=200,
    )
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/prices",
        json={"results": [_make_price(1001, 1.00)]},
        status=200,
    )

    df = benchmark(pull_product_information, "sv01")

    assert len(df) == 1  # nosec
    assert df["name"].to_list() == ["Pikachu"]  # nosec


@responses.activate
def test_pull_product_information_sm_normalizes_card_number(benchmark):
    product_id = "1863"  # sm1
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/products",
        json={"results": [_make_product(2001, "Rowlet", "9/149")]},
        status=200,
    )
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/prices",
        json={"results": [_make_price(2001, 0.50)]},
        status=200,
    )

    df = benchmark(pull_product_information, "sm1")

    assert df["card_number"].to_list() == ["009/149"]  # nosec


@responses.activate
def test_pull_product_information_excludes_reverse_holofoil_prices(benchmark):
    product_id = "22873"  # sv01 — both Normal and Reverse Holofoil entries for same card
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/products",
        json={"results": [_make_product(1001, "Pikachu", "025/198")]},
        status=200,
    )
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/prices",
        json={"results": [
            _make_price(1001, 5.00, "Reverse Holofoil"),
            _make_price(1001, 1.50, "Normal"),
        ]},
        status=200,
    )

    df = benchmark(pull_product_information, "sv01")

    # Reverse Holofoil price entry is skipped; Normal price is used
    assert df["market_price"].to_list() == [1.50]  # nosec


@responses.activate
def test_pull_product_information_validation_error_raises(benchmark):
    product_id = "22873"  # sv01
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/products",
        json={"results": [
            {
                "productId": "not-an-integer",
                "name": "Bad Card",
                "extendedData": [{"name": "Number", "value": "999/198"}],
            }
        ]},
        status=200,
    )
    responses.add(
        responses.GET,
        f"https://tcgcsv.com/tcgplayer/3/{product_id}/prices",
        json={"results": []},
        status=200,
    )

    def run():
        with pytest.raises(ValidationError):
            pull_product_information("sv01")

    benchmark(run)


# ---------------------------------------------------------------------------
# build_dataframe()
# ---------------------------------------------------------------------------


@patch("pipelines.defs.extract.tcgcsv.extract_pricing.pull_product_information")
def test_build_dataframe_concatenates_all_sets(mock_pull, benchmark):
    sample_df = pl.DataFrame({
        "product_id": [1001],
        "name": ["Pikachu"],
        "card_number": ["025/198"],
        "market_price": [1.50],
    })
    mock_pull.return_value = sample_df

    result = benchmark(_build_dataframe_fn)

    assert isinstance(result, pl.DataFrame)  # nosec
    assert len(result) == len(SET_PRODUCT_MATCHING)  # one row per set  # nosec
    assert result.columns == ["product_id", "name", "card_number", "market_price"]  # nosec
    assert result.dtypes == sample_df.dtypes  # nosec


@patch("pipelines.defs.extract.tcgcsv.extract_pricing.pull_product_information")
def test_build_dataframe_raises_on_empty_dataframe(mock_pull, benchmark):
    mock_pull.return_value = pl.DataFrame()

    def run():
        with pytest.raises(ValueError, match="Empty DataFrame"):
            _build_dataframe_fn()

    benchmark(run)
