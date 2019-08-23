FROM debian:stable-slim as builder

ARG DOCKER_TAG=stable
ARG DRP_VERSION=${DOCKER_TAG}

ENV LANG=C.UTF-8
ENV DEBIAN_FRONTEND=noninteractive

# digital rebar provision install starts here
WORKDIR /provision/
COPY tools/install.sh .
# install provision and its deps
RUN echo "DRP_VERSION=${DRP_VERSION}" && \
    apt-get update && \
    apt-get install -y sudo curl procps iproute2 ipmitool libarchive-tools p7zip && \
    ./install.sh --isolated install --drp-version=${DRP_VERSION}

# Copy binaries following symlinks. This is used for easier copying from builder image.
RUN mkdir /provision/binaries && \
    cp -L /provision/dr-provision /provision/drpcli /provision/drpjoin /provision/binaries/

# Build final container
FROM debian:stable-slim
ENV LANG=C.UTF-8
RUN apt-get update && \
    apt-get install --no-install-recommends -y ca-certificates curl iproute2 ipmitool jq libarchive-tools p7zip && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /provision/drp-data

COPY --from=builder /provision/binaries/ /usr/bin/
RUN chmod +x /usr/bin/dr*

# run the api server so we can install sledgehammer image
RUN dr-provision --version || true

EXPOSE 8091 8092 69 67 4011
VOLUME ["/provision/drp-data"]

ENTRYPOINT ["dr-provision", "--base-root=/provision/drp-data"]
CMD []

