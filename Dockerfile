FROM golang:1.12.9 AS build

RUN mkdir -p /memento
WORKDIR /memento

ADD go.mod go.mod
ADD go.sum go.sum
RUN go mod download

ADD . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM scratch
COPY --from=build /memento/memento .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["./memento", "run", "--config=/config/config.yml"]