## Build
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /emprestimo .

FROM gcr.io/distroless/static-debian11

WORKDIR /

COPY --from=build /emprestimo /emprestimo

EXPOSE 8080

ENTRYPOINT ["/emprestimo"]