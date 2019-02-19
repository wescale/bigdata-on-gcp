FROM golang:1.11-alpine AS build-env
RUN apk add --no-cache bash ca-certificates git gcc g++ libc-dev git
RUN mkdir -p $GOPATH/src/app
ENV GO111MODULE=on
WORKDIR $GOPATH/src/app
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download

FROM build-env AS build-stage
WORKDIR $GOPATH/src/app
COPY src .
RUN go install

# final stage
FROM alpine
RUN apk add --no-cache curl ca-certificates
VOLUME /secret
ENV PROM_PORT 8080
ENV TOPIC_NAME "projects/slavayssiere-sandbox/topics/mastodon-raw"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"
ENV HASHTAG "#GiletsJaunes"
ENV MASTODON_SERVER "https://linuxrocks.online"
ENV MASTODON_CLIENT_ID "***"
ENV MASTODON_CLIENT_SECRET "***"
ENV MASTODON_LOGIN "***"
ENV MASTODON_PASSWORD "***"

WORKDIR /app
COPY --from=build-stage /go/bin/mastodon /app/
HEALTHCHECK --interval=5s --retries=3 --timeout=5s CMD curl -f http://localhost:8080/healthz || exit 1
ENTRYPOINT ./mastodon
EXPOSE 8080
