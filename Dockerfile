# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.21-alpine3.18 AS build

ARG PACKAGE="mriyam.dev/cni-quickstart"
#ENV PACKAGE=$PACKAGE
WORKDIR "/go/src/${PACKAGE}"

RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download -x

ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -o cni-quickstart
#--mount=type=ssh \


FROM --platform=$BUILDPLATFORM gcr.io/distroless/static-debian12
USER nonroot:nonroot
COPY --from=build /bin/cni-quickstart /cni-quickstart
ENTRYPOINT ["/cni-quickstart"]