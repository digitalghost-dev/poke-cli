FROM ubuntu:latest
LABEL authors="cs"

ENTRYPOINT ["top", "-b"]