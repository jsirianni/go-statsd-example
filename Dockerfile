FROM golang:1.19-alpine as stage

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o app

FROM scratch
COPY --from=stage /app/app /app
ENTRYPOINT /app
