FROM golang:alpine AS build-env

ARG app_env
ENV APP_ENV $app_env

RUN mkdir -p /src/moomdate/reg-api
WORKDIR /src/moomdate/reg-api
COPY . .
RUN apk add git

#RUN go build /go/src/github.com/moomdate/reg-api/main.go
# RUN go get github.com/gocolly/colly
# RUN go get github.com/gorilla/mux
# RUN go get github.com/rs/cors
RUN go get -d -v

RUN go build .

CMD if [ ${APP_ENV} = production ]; \
    then \
    ./reg-api; \
    ["flask", "run"]\
    else \
    go run main.go \
    ["flask", "run"] \
    fresh; \
    fi 

EXPOSE 8081
