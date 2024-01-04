
FROM golang:1.21.5-bullseye AS build

WORKDIR /app

COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/app


FROM gcr.io/distroless/static/debian

COPY --from=build /bin/app /bin
COPY .env.prod /bin
EXPOSE 3000


ENTRYPOINT ["/bin/app", "bin/.env.prod"]
