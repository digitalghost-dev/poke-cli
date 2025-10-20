# Dockerfile.prod
FROM python:3.12-slim AS builder

RUN groupadd -r -g 10001 docsuser && \
    useradd -r -u 10001 -g docsuser -m -s /sbin/nologin docsuser

WORKDIR /build

RUN pip install --no-cache-dir mkdocs mkdocs-material mkdocs-nav-weight

COPY --chown=docsuser:docsuser mkdocs.yml /build/mkdocs.yml
COPY --chown=docsuser:docsuser docs/ /build/docs/

RUN chown -R docsuser:docsuser /build

USER docsuser

RUN mkdocs build

# --- Serve with lightweight HTTP server ---
FROM python:3.12-slim

# Create non-root user for runtime
RUN groupadd -r -g 10001 docsuser && \
    useradd -r -u 10001 -g docsuser -m -s /sbin/nologin docsuser

WORKDIR /site

COPY --from=builder --chown=docsuser:docsuser /build/site /site

USER docsuser

EXPOSE 8080

CMD ["python3", "-m", "http.server", "8080"]