FROM golang:alpine AS build

RUN apk --no-cache add git ca-certificates

RUN mkdir rtmp-auth
WORKDIR /rtmp-auth
COPY main.go /rtmp-auth
RUN go env -w GO111MODULE=off; \
    go get -u github.com/gorilla/mux; \
    go get -u github.com/lib/pq; \
    CGO_ENABLED=0 \
    go build \
    -installsuffix "static" \
    -o rtmp-auth

FROM scratch AS final
COPY --from=build /rtmp-auth/rtmp-auth /rtmp-auth
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/rtmp-auth"]