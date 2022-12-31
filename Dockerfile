FROM golang:1.19-alpine

RUN addgroup ideator && adduser -s /bin/false -G ideator -D ideator
USER ideator

WORKDIR /app/

COPY go.mod .
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY web/ ./web

RUN go build -o go-ideaboard ./cmd/go-ideaboard

EXPOSE 8080

CMD [ "./go-ideaboard" ]