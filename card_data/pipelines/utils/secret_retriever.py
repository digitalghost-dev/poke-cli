import botocore
import botocore.session
from aws_secretsmanager_caching import SecretCache, SecretCacheConfig

import json


def fetch_secret() -> str:
    client = botocore.session.get_session().create_client("secretsmanager")
    cache_config = SecretCacheConfig()
    cache = SecretCache(config=cache_config, client=client)

    secret = cache.get_secret_string("supabase")

    # convert to dictionary
    secret_dict = json.loads(secret)

    return secret_dict["database_uri"]
