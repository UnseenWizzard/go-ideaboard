FROM golang:1.19-alpine

RUN addgroup ideator && adduser -s /bin/false -G ideator -D ideator

WORKDIR /app/

COPY go.mod .
RUN go mod download

COPY cmd/ ./cmd
COPY internal/ ./internal
COPY web/ ./web

RUN go build -o go-ideaboard ./cmd/go-ideaboard

RUN chown -R ideator:ideator /app/go-ideaboard
RUN chmod +x /app/go-ideaboard
USER ideator

EXPOSE 8080

CMD [ "/app/go-ideaboard" ]