# defaults which may be overridden from the build command
ARG GO_VERSION=1.20
ARG DISTROLESS_VERSION=nonroot

# build stage
FROM golang:${GO_VERSION}-buster AS builder

COPY . /go/src/github.com/Walkbase/sapi
WORKDIR /go/src/github.com/Walkbase/sapi
RUN go install .

# final stage
FROM gcr.io/distroless/base-debian10:${DISTROLESS_VERSION}

COPY --from=builder /go/bin/sapi /usr/bin/sapi
# perform any further action as an unprivileged user
USER nonroot:nonroot
ENTRYPOINT [ "/usr/bin/sapi" ]
EXPOSE 8080
