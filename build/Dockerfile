FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV OPERATOR=/usr/local/bin/cluster-api-server \
    USER_UID=root \
    USER_NAME=cluster-api-server

COPY _output/cluster-api-server ${OPERATOR}

ENTRYPOINT ["/usr/local/bin/cluster-api-server"]
