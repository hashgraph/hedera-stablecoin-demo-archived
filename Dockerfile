# staging
FROM golang:1.14-alpine AS staging

COPY ./go.mod /srv/stable-coin/go.mod
COPY ./go.sum /srv/stable-coin/go.sum

RUN cd /srv/stable-coin && go mod download

# build
FROM staging AS build
ARG target

COPY ./ /srv/stable-coin/

RUN cd /srv/stable-coin && go build -o stable-coin-${target} ./${target}

# run
FROM alpine
ARG target

COPY --from=build /srv/stable-coin/stable-coin-${target} /srv/stable-coin-${target}
COPY --from=build /srv/stable-coin/network /srv/network
COPY --from=build /srv/stable-coin/data/schema /srv/data/schema

WORKDIR "/srv"
ENTRYPOINT "/srv/stable-coin-${TARGET}"
