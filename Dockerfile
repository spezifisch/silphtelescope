FROM alpine:latest as base
EXPOSE 8000

# set default timezone
ARG TZ=Europe/Berlin
ENV DEFAULT_TZ ${TZ}

# build stage
FROM golang:1.16-alpine as build
RUN mkdir -p /build /build/bin /build/data
WORKDIR /build

# cache go deps
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/silpht ./cmd/silpht
RUN go build -o bin/pokedexgen ./cmd/pokedexgen
RUN go build -o bin/geodexgen ./cmd/geodexgen

# final stage
FROM base as final

RUN apk upgrade --update \
  && apk add -U tzdata \
  && ln -snf /usr/share/zoneinfo/${DEFAULT_TZ} /etc/localtime \
  && rm -rf /var/cache/apk/*

# copy built binaries
RUN mkdir /app
WORKDIR /app
COPY --from=build /build/bin/* ./

# copy runtime assets
COPY --from=build /build/data/* ./data/

RUN mkdir -p /data
VOLUME /data
CMD ["/app/silpht", "--bind=0.0.0.0:8000"]

