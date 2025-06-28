# Dockerfile.prod
FROM python:3.12-slim AS builder

WORKDIR /build

RUN pip install mkdocs mkdocs-material

COPY ../mkdocs.yml /build/mkdocs.yml
COPY docs/index.md /build/docs/index.md

RUN mkdocs build

# --- Serve with lightweight HTTP server ---
FROM python:3.12-slim

WORKDIR /site

COPY --from=builder /build/site /site

EXPOSE 8080

CMD ["python3", "-m", "http.server", "8080"]