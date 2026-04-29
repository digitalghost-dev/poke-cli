import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import requests
from unittest.mock import patch, MagicMock

from pipelines.sensors import discord_success_sensor, discord_failure_sensor

# Access raw Python functions — the Dagster decorator wraps them in a SensorDefinition
_success_fn = discord_success_sensor._run_status_sensor_fn
_failure_fn = discord_failure_sensor._run_status_sensor_fn


def _make_context(run_id: str = "test-run-id", job_name: str = "test-job") -> MagicMock:
    ctx = MagicMock()
    ctx.dagster_run.run_id = run_id
    ctx.dagster_run.job_name = job_name
    return ctx


# ---------------------------------------------------------------------------
# discord_success_sensor()
# ---------------------------------------------------------------------------


@patch("pipelines.sensors.fetch_n8n_webhook_secret", return_value="https://n8n.example.com/hook")
@patch("pipelines.sensors.requests.post")
def test_discord_success_sensor_posts_webhook(mock_post, mock_secret, benchmark):
    mock_post.return_value.status_code = 200
    ctx = _make_context()

    benchmark(_success_fn, ctx)

    mock_post.assert_called_with(
        "https://n8n.example.com/hook",
        json={"job_name": "test-job", "status": "SUCCESS", "run_id": "test-run-id"},
        timeout=10,
    )


@patch("pipelines.sensors.fetch_n8n_webhook_secret", return_value="https://n8n.example.com/hook")
@patch("pipelines.sensors.requests.post", side_effect=requests.RequestException("connection refused"))
def test_discord_success_sensor_handles_request_exception(mock_post, mock_secret, benchmark):
    ctx = _make_context()

    benchmark(_success_fn, ctx)  # must not raise

    assert ctx.log.error.called  # nosec
    assert "connection refused" in ctx.log.error.call_args[0][0]  # nosec


@patch("pipelines.sensors.fetch_n8n_webhook_secret", return_value="https://n8n.example.com/hook")
@patch("pipelines.sensors.requests.post", side_effect=Exception("unexpected error"))
def test_discord_success_sensor_handles_generic_exception(mock_post, mock_secret, benchmark):
    ctx = _make_context()

    benchmark(_success_fn, ctx)  # must not raise

    assert ctx.log.error.called  # nosec
    assert "unexpected error" in ctx.log.error.call_args[0][0]  # nosec


# ---------------------------------------------------------------------------
# discord_failure_sensor()
# ---------------------------------------------------------------------------


@patch("pipelines.sensors.fetch_n8n_webhook_secret", return_value="https://n8n.example.com/hook")
@patch("pipelines.sensors.requests.post")
def test_discord_failure_sensor_posts_webhook(mock_post, mock_secret, benchmark):
    mock_post.return_value.status_code = 200
    ctx = _make_context()

    benchmark(_failure_fn, ctx)

    mock_post.assert_called_with(
        "https://n8n.example.com/hook",
        json={"job_name": "test-job", "status": "FAILURE", "run_id": "test-run-id"},
        timeout=10,
    )


@patch("pipelines.sensors.fetch_n8n_webhook_secret", return_value="https://n8n.example.com/hook")
@patch("pipelines.sensors.requests.post", side_effect=requests.RequestException("timeout"))
def test_discord_failure_sensor_handles_request_exception(mock_post, mock_secret, benchmark):
    ctx = _make_context()

    benchmark(_failure_fn, ctx)  # must not raise

    assert ctx.log.error.called  # nosec
    assert "timeout" in ctx.log.error.call_args[0][0]  # nosec


@patch("pipelines.sensors.fetch_n8n_webhook_secret", return_value="https://n8n.example.com/hook")
@patch("pipelines.sensors.requests.post", side_effect=Exception("service unavailable"))
def test_discord_failure_sensor_handles_generic_exception(mock_post, mock_secret, benchmark):
    ctx = _make_context()

    benchmark(_failure_fn, ctx)  # must not raise

    assert ctx.log.error.called  # nosec
    assert "service unavailable" in ctx.log.error.call_args[0][0]  # nosec
