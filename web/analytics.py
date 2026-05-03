import uuid

import streamlit as st
from posthog import Posthog

POSTHOG_API_KEY = "phc_qLvoCFJ5U9qgMS4p6LuJPgc3ZcrCRZYBNHLueHE9MU4C"
POSTHOG_HOST = "https://us.i.posthog.com"


@st.cache_resource
def init_posthog() -> Posthog:
    return Posthog(
        POSTHOG_API_KEY,
        host=POSTHOG_HOST,
        disable_geoip=False,
    )


def track_visit() -> None:
    if st.session_state.get("ph_visited"):
        return

    posthog = init_posthog()
    distinct_id = st.session_state.setdefault("ph_distinct_id", str(uuid.uuid4()))
    forwarded_for = st.context.headers.get("X-Forwarded-For", "")
    ip = forwarded_for.split(",")[0].strip() if forwarded_for else None

    properties: dict = {"$current_url": str(st.context.url)}
    if ip:
        properties["$ip"] = ip

    posthog.capture("$pageview", distinct_id=distinct_id, properties=properties)
    st.session_state["ph_visited"] = True
