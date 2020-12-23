FROM golang:1.15

WORKDIR /app
COPY . .

RUN go install -v ./cmd/fattingo/

# CMD ["fattingo"]
CMD ["go", "run", "./cmd/fattingo/"]
