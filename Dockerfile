FROM golang:1.15.2-alpine3.12 AS go-builder
ENV CGO_ENABLED 1
RUN apk --no-cache add gcc musl-dev vips-dev
RUN mkdir /build
WORKDIR /build
COPY . .
WORKDIR /build/app/imgasm
RUN go build


FROM alpine:3.12
RUN apk --no-cache add ca-certificates vips-dev
RUN addgroup -g 3000 -S app && adduser -u 100000 -S app -G app --no-create-home --disabled-password \
    && mkdir -p /app/badger.db && chown app:app /app/badger.db
USER 100000
WORKDIR /app
COPY --from=go-builder --chown=app:app /build/app/imgasm/imgasm /app/imgasm
COPY --from=go-builder --chown=app:app /build/ui/static /app/ui/static
ENTRYPOINT ["./imgasm"]
