# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-modeldgoclient"
LABEL REPO="https://github.com/metaprov/modeldgoclient"

ENV PROJPATH=/go/src/github.com/metaprov/modeldgoclient

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/metaprov/modeldgoclient
WORKDIR /go/src/github.com/metaprov/modeldgoclient

RUN make build-alpine

# Final Stage
FROM metaprov/modeld-go-client

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/metaprov/modeldgoclient"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/modeldgoclient/bin

WORKDIR /opt/modeldgoclient/bin

COPY --from=build-stage /go/src/github.com/metaprov/modeldgoclient/bin/modeldgoclient /opt/modeldgoclient/bin/
RUN chmod +x /opt/modeldgoclient/bin/modeldgoclient

# Create appuser
RUN adduser -D -g '' modeldgoclient
USER modeldgoclient

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/modeldgoclient/bin/modeldgoclient"]
