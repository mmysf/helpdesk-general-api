# builder
FROM golang:1.22-alpine AS builder

# change workdir
WORKDIR /home

# add go modules lockfiles
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build a fully standalone binary with zero dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o build-app main.go

# final image
FROM alpine

# copy the binary from the builder image
COPY --from=builder /home/build-app .

# run the binary
ENTRYPOINT ./build-app
