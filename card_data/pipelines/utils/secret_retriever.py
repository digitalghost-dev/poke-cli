import botocore
import botocore.session
from aws_secretsmanager_caching import SecretCache, SecretCacheConfig

import json
from typing import TypedDict, cast


class SupabaseSecret(TypedDict):
    database_uri: str


def fetch_secret() -> str:
    client = botocore.session.get_session().create_client("secretsmanager")
    cache_config = SecretCacheConfig()
    cache = SecretCache(config=cache_config, client=client)

    secret = cast(str, cache.get_secret_string("supabase"))

    # convert to dictionary
    secret_dict: SupabaseSecret = json.loads(secret)

    return secret_dict["database_uri"]
