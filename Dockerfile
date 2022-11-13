FROM golang:1.18.7-buster

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY css ./css
COPY html ./html
COPY js ./js
COPY internal ./internal
COPY *.go ./

RUN go build main.go

ENTRYPOINT [ "./main" ]
