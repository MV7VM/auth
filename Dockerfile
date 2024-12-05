FROM golang:1.22.4-alpine as builder

ENV APP_HOME="/app" \
    APP_MAIN="cmd/main.go" \
    CGO_ENABLED=0 \
    GO111MODULE="on" 

WORKDIR $APP_HOME

COPY . .
RUN go mod download && go mod verify
RUN go build -o main $APP_MAIN

FROM golang:1.22.4-alpine

ENV APP_HOME="/app"

WORKDIR $APP_HOME

COPY --from=builder $APP_HOME/main $APP_HOME/main
COPY --from=builder $APP_HOME/config $APP_HOME/config


EXPOSE 3000
EXPOSE 5000
CMD ["./main"]
