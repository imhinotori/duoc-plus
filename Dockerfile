FROM ubuntu:latest
LABEL authors="hinotori"

ENTRYPOINT ["top", "-b"]