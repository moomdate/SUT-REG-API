FROM golang:1.13-alpine3.10 AS build-env

ARG build_type
ENV BUILD_TYPE $build_type

RUN mkdir -p /src/moomdate/reg-api
WORKDIR /src/moomdate/reg-api
COPY . .
RUN go get ./


RUN go build .

#CMD if [ ${BUILD_TYPE} = prod ]; \
#    then \
#    ./reg-api; \
#    else \
#    go run main.go; \
#    fresh; \
#    fi
CMD ./reg-api
EXPOSE 8081
