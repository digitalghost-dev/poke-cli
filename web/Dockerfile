# Dockerfile to build the Streamlit app.
FROM python:3.12-slim-bookworm AS builder

RUN groupadd -r -g 10001 streamlit_group && \
    useradd -r -u 10001 -g streamlit_group -m -s /sbin/nologin streamlit_user

WORKDIR /app

COPY --from=ghcr.io/astral-sh/uv:latest /uv /usr/local/bin/uv

COPY pyproject.toml uv.lock ./

RUN uv sync --frozen --no-dev --no-cache

COPY --chown=streamlit_user:streamlit_group .streamlit .streamlit
COPY --chown=streamlit_user:streamlit_group app.py .

RUN chown -R streamlit_user:streamlit_group /app

# --- Runtime ---
FROM python:3.12-slim-bookworm

RUN groupadd -r -g 10001 streamlit_group && \
    useradd -r -u 10001 -g streamlit_group -m -s /sbin/nologin streamlit_user

WORKDIR /app

COPY --from=builder --chown=streamlit_user:streamlit_group /app /app

RUN apt-get update && apt-get install -y --no-install-recommends curl && rm -rf /var/lib/apt/lists/*

USER streamlit_user

EXPOSE 8501

HEALTHCHECK CMD curl --fail http://localhost:8501/_stcore/health

ENTRYPOINT ["/app/.venv/bin/streamlit", "run", "app.py", "--server.port=8501", "--server.address=0.0.0.0"]
