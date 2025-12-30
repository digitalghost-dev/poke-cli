import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import pytest
import polars as pl
import responses
from pipelines.defs.extract.tcgdex.extract_sets import extract_sets_data


@pytest.fixture
def mock_api_response():
    """Sample API responses matching tcgdex series format with sets"""
    return {
        "https://api.tcgdex.net/v2/en/series/me": {
            "id": "me",
            "name": "Mega Evolution",
            "sets": [
                {
                    "id": "me01",
                    "name": "Mega Evolution",
                    "cardCount": {"official": 12, "total": 12},
                    "logo": "https://example.com/me01.png",
                    "symbol": "https://example.com/me01-symbol.png",
                },
                {
                    "id": "me02",
                    "name": "Phantasmal Flames",
                    "cardCount": {"official": 25, "total": 25},
                    "logo": "https://example.com/me02.png",
                    "symbol": "https://example.com/me02-symbol.png",
                },
            ],
        },
        "https://api.tcgdex.net/v2/en/series/sv": {
            "id": "sv",
            "name": "Scarlet & Violet",
            "sets": [
                {
                    "id": "sv01",
                    "name": "Scarlet & Violet",
                    "cardCount": {"official": 198, "total": 258},
                    "logo": "https://example.com/sv01.png",
                    "symbol": "https://example.com/sv01-symbol.png",
                },
                {
                    "id": "sv02",
                    "name": "Paldea Evolved",
                    "cardCount": {"official": 193, "total": 279},
                    "logo": "https://example.com/sv02.png",
                    "symbol": None,
                },
            ],
        },
        "https://api.tcgdex.net/v2/en/series/swsh": {
            "id": "swsh",
            "name": "Sword & Shield",
            "sets": [
                {
                    "id": "swsh1",
                    "name": "Sword & Shield",
                    "cardCount": {"official": 202, "total": 216},
                    "logo": None,
                    "symbol": "https://example.com/swsh1-symbol.png",
                },
            ],
        },
    }


@pytest.mark.benchmark
@responses.activate
def test_extract_sets_data_success(mock_series_responses):
    """Test successful extraction of sets from multiple series"""
    # Mock all API calls
    for url, response_data in mock_series_responses.items():
        responses.add(
            responses.GET,
            url,
            json=response_data,
            status=200,
        )

    result = extract_sets_data()

    # Assertions
    assert isinstance(result, pl.DataFrame)  # nosec
    assert len(result) == 5  # nosec (2 + 2 + 1 sets)
    assert set(result.columns) == {  # nosec
        "series_id",
        "set_id",
        "set_name",
        "official_card_count",
        "total_card_count",
        "logo",
        "symbol",
    }
    assert set(result["series_id"].to_list()) == {"me", "sv", "swsh"}  # nosec
    assert set(result["set_id"].to_list()) == {"me01", "me02", "sv01", "sv02", "swsh1"}  # nosec


@pytest.mark.benchmark
@responses.activate
def test_extract_sets_data_empty_sets(mock_series_responses):
    """Test extraction when a series has no sets"""
    # Modify one response to have empty sets
    mock_series_responses["https://api.tcgdex.net/v2/en/series/me"]["sets"] = []

    for url, response_data in mock_series_responses.items():
        responses.add(
            responses.GET,
            url,
            json=response_data,
            status=200,
        )

    result = extract_sets_data()

    assert isinstance(result, pl.DataFrame)  # nosec
    assert len(result) == 3  # nosec (0 + 2 + 1 sets)
    assert "me" not in result["series_id"].to_list()  # nosec


@pytest.mark.benchmark
@responses.activate
def test_extract_sets_data_null_card_counts():
    """Test extraction with null card counts"""
    mock_responses = {
        "https://api.tcgdex.net/v2/en/series/me": {
            "id": "me",
            "name": "Mega Evolution",
            "sets": [],
        },
        "https://api.tcgdex.net/v2/en/series/sv": {
            "id": "sv",
            "name": "Scarlet & Violet",
            "sets": [
                {
                    "id": "sv01",
                    "name": "Scarlet & Violet",
                    "cardCount": {},
                    "logo": None,
                    "symbol": None,
                },
            ],
        },
        "https://api.tcgdex.net/v2/en/series/swsh": {
            "id": "swsh",
            "name": "Sword & Shield",
            "sets": [],
        },
    }

    for url, response_data in mock_responses.items():
        responses.add(
            responses.GET,
            url,
            json=response_data,
            status=200,
        )

    result = extract_sets_data()

    assert isinstance(result, pl.DataFrame)  # nosec
    assert len(result) == 1  # nosec
    assert result["official_card_count"].to_list()[0] is None  # nosec
    assert result["total_card_count"].to_list()[0] is None  # nosec
