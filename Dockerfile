# The base go-image
FROM golang:1.21-alpine

RUN mkdir /app 
COPY . /app

WORKDIR /app

RUN go build

# Run the server executable
CMD [ "/app/boyd_cooper" ]
