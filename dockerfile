FROM golang:alpine as builder

RUN mkdir app
WORKDIR /app

COPY . . 

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ocp_quiz

FROM alpine:latest
RUN mkdir /app
WORKDIR /app

COPY --from=builder /app/ocp_quiz .

CMD ["./ocp_quiz"]