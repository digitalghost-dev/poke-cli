import dagster as dg
import requests
from dagster import RetryPolicy, Backoff
from termcolor import colored

from ...utils.secret_retriever import fetch_n8n_webhook_secret
from .usage import trigger_pikalytics_usage


@dg.asset(
    kinds={"n8n"},
    deps=[trigger_pikalytics_usage],
    retry_policy=RetryPolicy(max_retries=1, delay=5, backoff=Backoff.EXPONENTIAL),
)
def trigger_pikalytics_pokemon_comp_info() -> None:
    """Trigger the n8n pokemon-comp-info workflow and block until it finishes. It reads the
    fresh top-50 from staging.pikalytics_usage (hence the dependency on the usage trigger),
    scrapes each Pokémon's AI pokedex page, and loads staging.pikalytics_pokemon_comp_info.
    The webhook responds only when the last node completes; this is a long run (~50 page
    scrapes), so the timeout is generous and retries are limited (a retry re-runs all 50)."""
    webhook_url = fetch_n8n_webhook_secret("pikalytics-pokemon-comp-info")

    resp = requests.post(webhook_url, timeout=1200)
    resp.raise_for_status()
    print(colored(" ✓", "green"), "n8n pokemon-comp-info workflow completed")
