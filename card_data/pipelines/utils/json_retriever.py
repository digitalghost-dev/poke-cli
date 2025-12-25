import requests
import json

def fetch_json(url: str, timeout: int = 30) -> dict:
    response = requests.get(url, timeout=timeout)
    response.raise_for_status()

    return response.json()
