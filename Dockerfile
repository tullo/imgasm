FROM golang:1.16.0-alpine3.13 AS go-builder
ENV CGO_ENABLED 1
RUN apk --no-cache add gcc musl-dev vips-dev
RUN go install honnef.co/go/tools/cmd/staticcheck@v0.1.2
RUN mkdir /build
WORKDIR /build
COPY . .
RUN staticcheck -go 1.16 \
		-tests ./backblaze/... ./db/... ./file/... ./models/... ./ui/templates/...
WORKDIR /build/app/imgasm
RUN go build -mod=vendor


FROM alpine:3.13.2
RUN apk --no-cache add ca-certificates vips-dev
RUN addgroup -g 3000 -S app && adduser -u 100000 -S app -G app --no-create-home --disabled-password \
    && mkdir -p /app/badger.db && chown app:app /app/badger.db
USER 100000
WORKDIR /app
COPY --from=go-builder --chown=app:app /build/app/imgasm/imgasm /app/imgasm
COPY --from=go-builder --chown=app:app /build/ui/static /app/ui/static
ENTRYPOINT ["./imgasm"]
