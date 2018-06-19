# build stage
FROM golang:alpine AS build-env

RUN mkdir -p /go/src/github.com/cloudnativecz/public-cloud-kubernetes-demo/backend
ADD ./ /go/src/github.com/cloudnativecz/public-cloud-kubernetes-demo/backend
WORKDIR /go/src/github.com/cloudnativecz/public-cloud-kubernetes-demo/backend

RUN go build -o backend cmd/server.go

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/cloudnativecz/public-cloud-kubernetes-demo/backend/backend /app/

ENTRYPOINT ./backend

