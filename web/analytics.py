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


def _truncate_ip(ip: str | None) -> str | None:
    """Zero out the last octet (IPv4) or last 80 bits (IPv6) so the stored
    IP can still resolve to a country but cannot uniquely identify a visitor."""
    if not ip:
        return None
    if ":" in ip:
        parts = ip.split(":")[:3]
        return ":".join(parts) + "::"
    parts = ip.split(".")
    if len(parts) == 4:
        return ".".join(parts[:3] + ["0"])
    return None


def track_visit() -> None:
    if st.session_state.get("ph_visited"):
        return

    posthog = init_posthog()
    distinct_id = st.session_state.setdefault("ph_distinct_id", str(uuid.uuid4()))
    forwarded_for = st.context.headers.get("X-Forwarded-For", "")
    raw_ip = forwarded_for.split(",")[0].strip() if forwarded_for else None
    ip = _truncate_ip(raw_ip)

    # Events are anonymous: no persistent person profile is created, and the
    # IP is truncated before being sent so PostHog can resolve country but
    # cannot uniquely identify or track users across visits.
    properties: dict = {
        "$current_url": str(st.context.url),
        "$process_person_profile": False,
    }
    if ip:
        properties["$ip"] = ip

    posthog.capture("$pageview", distinct_id=distinct_id, properties=properties)
    st.session_state["ph_visited"] = True
