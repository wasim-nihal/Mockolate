FROM golang:1.19 as build

WORKDIR /go/src/app
COPY . ./

RUN go mod download
RUN go vet -v

RUN CGO_ENABLED=0 go build -o /mock-server

FROM gcr.io/distroless/static-debian11

COPY --from=build /mock-server /

CMD ["/mock-server"]