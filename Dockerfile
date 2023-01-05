FROM golang:1.19-alpine

RUN addgroup 1000 && adduser -s /bin/false -G 1000 -D 1000

WORKDIR /app/

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY web/ ./web

RUN go build -o go-ideaboard ./cmd/go-ideaboard

RUN chmod a+x /app/go-ideaboard
USER 1000

EXPOSE 8080

CMD [ "/app/go-ideaboard" ]