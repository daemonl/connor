FROM debian:jessie
COPY build/connor /connor
WORKDIR /
ENTRYPOINT ["/connor"]
