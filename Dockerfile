FROM golang:1.15.2-alpine3.12 AS go-builder
ENV CGO_ENABLED 1
RUN apk --no-cache add gcc musl-dev vips-dev
COPY . /build/
WORKDIR /build
RUN go build -o main


FROM alpine:3.12
RUN apk --no-cache add ca-certificates vips-dev
RUN addgroup -g 1000 -S app && adduser -u 1000 -S app -G app --no-create-home --disabled-password \
    && mkdir -p /app/badger.db && chown app:app /app/badger.db
USER app
WORKDIR /app
COPY --from=go-builder --chown=app:app /build/main /app/imgasm
COPY --from=go-builder --chown=app:app /build/static /app/static
ENTRYPOINT ["./imgasm"]
