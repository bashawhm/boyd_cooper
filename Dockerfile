# Builder image
FROM golang:1.26-alpine as build

RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -v

# Server image
FROM gcr.io/distroless/static-debian13

COPY --from=build /app/boyd_cooper /app/boyd_cooper
WORKDIR /app
CMD [ "/app/boyd_cooper" ]
