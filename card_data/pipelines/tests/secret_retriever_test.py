import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import json
import pytest
from unittest.mock import patch, MagicMock
from pipelines.utils.secret_retriever import fetch_secret, fetch_n8n_webhook_secret

@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_secret_success(mock_get_session, mock_secret_cache_cls):
    """Test successful retrieval of the Supabase database URI."""
    secret_payload = json.dumps({"database_uri": "postgresql://user:pass@host/db"})

    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = secret_payload
    mock_secret_cache_cls.return_value = mock_cache_instance

    result = fetch_secret()

    assert result == "postgresql://user:pass@host/db"  # nosec
    mock_cache_instance.get_secret_string.assert_called_once_with("supabase")


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_secret_missing_key(mock_get_session, mock_secret_cache_cls):
    """Test KeyError when the secret JSON is missing 'database_uri'."""
    secret_payload = json.dumps({"some_other_key": "value"})

    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = secret_payload
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(KeyError, match="database_uri"):
        fetch_secret()


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_secret_invalid_json(mock_get_session, mock_secret_cache_cls):
    """Test that invalid JSON in the secret raises JSONDecodeError."""
    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = "not valid json"
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(json.JSONDecodeError):
        fetch_secret()


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_secret_empty_json_object(mock_get_session, mock_secret_cache_cls):
    """Test KeyError when the secret is an empty JSON object."""
    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = "{}"
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(KeyError, match="database_uri"):
        fetch_secret()


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_secret_cache_raises(mock_get_session, mock_secret_cache_cls):
    """Test that an exception from SecretCache propagates."""
    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.side_effect = Exception(
        "Secret not found"
    )
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(Exception, match="Secret not found"):
        fetch_secret()


# ---------------------------------------------------------------------------
# fetch_n8n_webhook_secret()
# ---------------------------------------------------------------------------


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_n8n_webhook_secret_success(mock_get_session, mock_secret_cache_cls):
    """Test successful retrieval of the n8n webhook URL."""
    secret_payload = json.dumps({"n8n_webhook": "https://n8n.example.com/hook/abc"})

    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = secret_payload
    mock_secret_cache_cls.return_value = mock_cache_instance

    result = fetch_n8n_webhook_secret()

    assert result == "https://n8n.example.com/hook/abc"  # nosec
    mock_cache_instance.get_secret_string.assert_called_once_with("n8n_webhook")


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_n8n_webhook_secret_missing_key(mock_get_session, mock_secret_cache_cls):
    """Test KeyError when the secret JSON is missing 'n8n_webhook'."""
    secret_payload = json.dumps({"wrong_key": "value"})

    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = secret_payload
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(KeyError, match="n8n_webhook"):
        fetch_n8n_webhook_secret()


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_n8n_webhook_secret_invalid_json(
    mock_get_session, mock_secret_cache_cls
):
    """Test that invalid JSON in the secret raises JSONDecodeError."""
    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = "{broken"
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(json.JSONDecodeError):
        fetch_n8n_webhook_secret()


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_n8n_webhook_secret_empty_json_object(
    mock_get_session, mock_secret_cache_cls
):
    """Test KeyError when the secret is an empty JSON object."""
    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.return_value = "{}"
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(KeyError, match="n8n_webhook"):
        fetch_n8n_webhook_secret()


@pytest.mark.benchmark
@patch("pipelines.utils.secret_retriever.SecretCache")
@patch("pipelines.utils.secret_retriever.botocore.session.get_session")
def test_fetch_n8n_webhook_secret_cache_raises(
    mock_get_session, mock_secret_cache_cls
):
    """Test that an exception from SecretCache propagates."""
    mock_cache_instance = MagicMock()
    mock_cache_instance.get_secret_string.side_effect = Exception(
        "Access denied"
    )
    mock_secret_cache_cls.return_value = mock_cache_instance

    with pytest.raises(Exception, match="Access denied"):
        fetch_n8n_webhook_secret()
        