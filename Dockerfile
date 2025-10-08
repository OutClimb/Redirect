FROM golang:1.25-alpine AS redirect-builder

COPY . /app
WORKDIR /app

RUN go mod download && go mod verify
RUN go build -v -o /app/redirect cmd/service/main.go
RUN go build -v -o /app/redirect_create_user cmd/create_user/main.go

FROM alpine:latest AS redirect

ENV GIN_MODE release

WORKDIR /app

COPY --from=redirect-builder /app/redirect /app/redirect
COPY --from=redirect-builder /app/redirect_create_user /app/redirect_create_user
COPY --from=redirect-builder /app/LICENSE.md /app/LICENSE.md
COPY --from=redirect-builder /app/README.md /app/README.md

RUN apk --no-cache add curl

ENTRYPOINT ["/app/redirect"]
