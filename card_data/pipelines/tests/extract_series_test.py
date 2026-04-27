import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import pytest
import polars as pl
import requests
import responses
from pydantic import ValidationError
from pipelines.defs.extract.tcgdex.extract_series import extract_series_data


@pytest.fixture
def mock_api_response():
    """Sample API response matching tcgdex format"""
    return [
        {"id": "sv", "name": "Scarlet & Violet", "logo": "https://example.com/sv.png"},
        {"id": "swsh", "name": "Sword & Shield", "logo": "https://example.com/swsh.png"},
        {"id": "xy", "name": "XY", "logo": "https://example.com/xy.png"},
        {"id": "me", "name": "McDonald's Collection", "logo": "https://example.com/me.png"},
        {"id": "sm", "name": "Sun & Moon", "logo": None},
    ]


@responses.activate
def test_extract_series_data_success(benchmark, mock_api_response):
    """Test successful extraction and filtering"""
    responses.add(
        responses.GET,
        "https://api.tcgdex.net/v2/en/series",
        json=mock_api_response,
        status=200
    )

    result = benchmark(extract_series_data)

    assert isinstance(result, pl.DataFrame)  # nosec
    assert len(result) == 4  # nosec
    assert set(result["id"].to_list()) == {"swsh", "sv", "me", "sm"}  # nosec
    assert "name" in result.columns  # nosec
    assert "logo" in result.columns  # nosec


@responses.activate
def test_extract_series_data_validation_error(benchmark):
    """Test that Pydantic ValidationError propagates when a required field is missing."""
    responses.add(
        responses.GET,
        "https://api.tcgdex.net/v2/en/series",
        json=[{"logo": "https://example.com/test.png"}],  # missing required 'id' and 'name'
        status=200,
    )

    def run():
        with pytest.raises(ValidationError):
            extract_series_data()

    benchmark(run)


@responses.activate
def test_extract_series_data_http_error(benchmark):
    """Test that an HTTP 500 from the API propagates as HTTPError."""
    responses.add(
        responses.GET,
        "https://api.tcgdex.net/v2/en/series",
        json={"error": "internal server error"},
        status=500,
    )

    def run():
        with pytest.raises(requests.exceptions.HTTPError):
            extract_series_data()

    benchmark(run)


@responses.activate
def test_extract_series_data_all_filtered_out(benchmark):
    """Test that an empty DataFrame is returned when no series match the allowed IDs."""
    responses.add(
        responses.GET,
        "https://api.tcgdex.net/v2/en/series",
        json=[{"id": "bw", "name": "Black & White", "logo": None}],
        status=200,
    )

    result = benchmark(extract_series_data)

    assert isinstance(result, pl.DataFrame)  # nosec
    assert result.is_empty()  # nosec
