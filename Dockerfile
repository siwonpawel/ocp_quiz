FROM alpine:latest
RUN mkdir /app
WORKDIR /app

COPY build/ocp_quiz /app/ocp_quiz

CMD ["./ocp_quiz"]