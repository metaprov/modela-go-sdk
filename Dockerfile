# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-modelagoclient"
LABEL REPO="https://github.com/metaprov/modelagoclient"

ENV PROJPATH=/go/src/github.com/metaprov/modelagoclient

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/metaprov/modelagoclient
WORKDIR /go/src/github.com/metaprov/modelagoclient

RUN make build-alpine

# Final Stage
FROM metaprov/modela-go-client

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/metaprov/modelagoclient"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/modelagoclient/bin

WORKDIR /opt/modelagoclient/bin

COPY --from=build-stage /go/src/github.com/metaprov/modelagoclient/bin/modelagoclient /opt/modelagoclient/bin/
RUN chmod +x /opt/modelagoclient/bin/modelagoclient

# Create appuser
RUN adduser -D -g '' modelagoclient
USER modelagoclient

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/modelagoclient/bin/modelagoclient"]
