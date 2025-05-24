# build backend
FROM golang:1.23-alpine as buildback

WORKDIR /app/backend

COPY backend .
RUN go mod download
RUN go build -o app cmd/panel/panel.go 
RUN go build -o cli cmd/cli/cli.go 

# build app
FROM alpine:latest

WORKDIR /opt/sslpanel

COPY --from=buildback /app/backend/app /app/backend/cli ./
COPY --from=buildback /app/backend/fixtures ./fixtures
COPY --from=buildback /app/backend/migrations ./migrations
COPY --from=buildback /app/backend/var ./var

CMD [ "/opt/sslpanel/app" ]
