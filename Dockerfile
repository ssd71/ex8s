# Build environment
# -----------------
FROM golang:1.14-alpine as build-env
WORKDIR /ex8s

RUN apk update && apk add --no-cache gcc musl-dev git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags '-w -s' -a -o ./bin/app


# Deployment environment
# ----------------------
FROM alpine

COPY --from=build-env /ex8s/bin/app /ex8s/

EXPOSE 8080
CMD ["/ex8s/app"]