# Builder image
FROM golang:1.21-alpine as build

RUN mkdir /app 
COPY . /app
WORKDIR /app
RUN go build

# Server image
FROM gcr.io/distroless/static-debian12

COPY --from=build /app/boyd_cooper /app/boyd_cooper
CMD [ "/app/boyd_cooper" ]
