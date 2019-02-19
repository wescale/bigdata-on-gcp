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
ENV MESSAGE_SIZE 1024 
ENV TOPIC_NAME "projects/sandbox-perrier/topics/raw-twitter"
ENV SECRET_PATH "/secret/secret-sa-gcp-pubsub.json"

ENV CONSUMER_KEY "default"
ENV CONSUMER_SECRET "default"
ENV ACCESS_TOKEN "default"
ENV ACCESS_SECRET "default"
ENV HASHTAG "#CAC40"
ENV ZIPKIN_ENDPOINT "zipkin.istio-system.svc.cluster.local"

WORKDIR /app
COPY --from=build-stage /go/bin/twitter /app/
HEALTHCHECK --interval=5s --retries=3 --timeout=5s CMD curl -f http://localhost:8080/healthz || exit 1
ENTRYPOINT ./twitter
EXPOSE 8080
