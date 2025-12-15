import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import pytest
import polars as pl
import responses
from pipelines.defs.extract.extract_data import extract_series_data

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
def test_extract_series_data_success(mock_api_response):
    """Test successful extraction and filtering"""
    # Mock the API call
    responses.add(
        responses.GET,
        "https://api.tcgdex.net/v2/en/series",
        json=mock_api_response,
        status=200
    )

    result = extract_series_data()

    # Assertions
    assert isinstance(result, pl.DataFrame)
    assert len(result) == 3  # Only swsh, sv, me
    assert set(result["id"].to_list()) == {"swsh", "sv", "me"}
    assert "name" in result.columns
    assert "logo" in result.columns