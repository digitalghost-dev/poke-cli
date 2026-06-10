import dagster as dg
import requests
from dagster import RetryPolicy, Backoff
from termcolor import colored

from ...utils.secret_retriever import fetch_n8n_webhook_secret


@dg.asset(
    kinds={"n8n"},
    retry_policy=RetryPolicy(max_retries=2, delay=5, backoff=Backoff.EXPONENTIAL),
)
def trigger_pikalytics_top_teams() -> None:
    """Trigger the n8n top-teams workflow (scrape → staging.pikalytics_top_teams) and
    block until it finishes. The webhook is configured to respond only when the last
    node completes, so a successful POST means staging is loaded."""
    webhook_url = fetch_n8n_webhook_secret("pikalytics-top-teams")

    resp = requests.post(webhook_url, timeout=600)
    resp.raise_for_status()
    print(colored(" ✓", "green"), "n8n top-teams workflow completed")
