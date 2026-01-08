FROM golang:1.26-rc-alpine3.23 AS builder-dev
WORKDIR /build
COPY go.mod go.sum ./
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download && \
    # go install github.com/swaggo/swag/cmd/swag@latest && \
    # cd ./app && \
    # swag init --parseDependency --parseInternal -g index.go -dir api/routers/v1 --instanceName v1 && \
    # cd .. && \
    CGO_ENABLED=0 GOOS=linux go build -o ./dist/pets ./cmd/server/main.go

    
FROM golang:1.26-rc-alpine3.23 AS builder-prod
WORKDIR /build
COPY go.mod go.sum ./
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN rm -rf ./app/docs && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -o ./dist/pets ./cmd/server/main.go


FROM alpine:3.20 AS test
COPY --from=builder-dev /build/dist /app
COPY ./configs/settings.test.yaml ./app/settings.yaml
RUN apk add tzdata && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone
WORKDIR /app
CMD ["./pets"]


FROM alpine:3.20 AS prod
COPY --from=builder-prod /build/dist /app
COPY ./configs/settings.prod.yaml ./app/settings.yaml
RUN apk add tzdata && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' > /etc/timezone
WORKDIR /app
CMD ["./pets"]