import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import pytest
import responses
import requests
from pipelines.utils.json_retriever import fetch_json


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_success():
    """Test successful JSON retrieval."""
    responses.add(
        responses.GET,
        "https://api.example.com/data",
        json={"id": 1, "name": "Pikachu"},
        status=200,
    )

    result = fetch_json("https://api.example.com/data")

    assert isinstance(result, dict)  # nosec
    assert result["id"] == 1  # nosec
    assert result["name"] == "Pikachu"  # nosec


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_with_nested_data():
    """Test retrieval of nested JSON structures."""
    payload = {
        "results": [
            {"productId": 100, "name": "Card A"},
            {"productId": 200, "name": "Card B"},
        ],
        "totalItems": 2,
    }
    responses.add(
        responses.GET,
        "https://api.example.com/products",
        json=payload,
        status=200,
    )

    result = fetch_json("https://api.example.com/products")

    assert result["totalItems"] == 2  # nosec
    assert len(result["results"]) == 2  # nosec
    assert result["results"][0]["productId"] == 100  # nosec


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_http_404():
    """Test that a 404 response raises HTTPError."""
    responses.add(
        responses.GET,
        "https://api.example.com/missing",
        json={"error": "not found"},
        status=404,
    )

    with pytest.raises(requests.exceptions.HTTPError):
        fetch_json("https://api.example.com/missing")


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_http_500():
    """Test that a 500 response raises HTTPError."""
    responses.add(
        responses.GET,
        "https://api.example.com/error",
        json={"error": "internal server error"},
        status=500,
    )

    with pytest.raises(requests.exceptions.HTTPError):
        fetch_json("https://api.example.com/error")


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_connection_error():
    """Test that a connection error raises ConnectionError."""
    responses.add(
        responses.GET,
        "https://api.example.com/down",
        body=requests.exceptions.ConnectionError("Connection refused"),
    )

    with pytest.raises(requests.exceptions.ConnectionError):
        fetch_json("https://api.example.com/down")


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_timeout():
    """Test that a timeout raises an appropriate exception."""
    responses.add(
        responses.GET,
        "https://api.example.com/slow",
        body=requests.exceptions.ReadTimeout("Read timed out"),
    )

    with pytest.raises(requests.exceptions.ReadTimeout):
        fetch_json("https://api.example.com/slow")


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_empty_object():
    """Test retrieval of an empty JSON object."""
    responses.add(
        responses.GET,
        "https://api.example.com/empty",
        json={},
        status=200,
    )

    result = fetch_json("https://api.example.com/empty")

    assert result == {}  # nosec


@pytest.mark.benchmark
@responses.activate
def test_fetch_json_invalid_json():
    """Test that an invalid JSON body raises a ValueError (JSONDecodeError)."""
    responses.add(
        responses.GET,
        "https://api.example.com/bad",
        body="not valid json {{{",
        status=200,
        content_type="application/json",
    )

    with pytest.raises(requests.exceptions.JSONDecodeError):
        fetch_json("https://api.example.com/bad")
