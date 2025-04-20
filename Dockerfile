# Builder image
FROM golang:1.24-alpine as build

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -v

# Server image
FROM gcr.io/distroless/static-debian12

COPY --from=build /app/boyd_cooper /app/boyd_cooper
WORKDIR /app
CMD [ "/app/boyd_cooper" ]
