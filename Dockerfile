FROM docker.io/library/golang:1.16 as builder

WORKDIR /src/adztbotv2

# Project file
COPY main.go .
COPY db/ ./db/
COPY config/ ./config/
COPY commands/ ./commands/
# Go settings
COPY go.mod .
COPY go.sum .

# Donwload and Install project
RUN go get -d -v ./...
RUN go build -o ADZTBotV2 main.go

FROM docker.io/library/ubuntu:20.04

WORKDIR /src/adztbotv2

COPY --from=builder /src/adztbotv2/ADZTBotV2 /src/adztbotv2/

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
RUN update-ca-certificates

# Env variables
ENV DB none
ENV URL none
ENV CHANM none
ENV CHANV none
ENV TOKEN none

CMD ["sh", "-c", "./ADZTBotV2 -db $DB -url $URL -chanm $CHANM -chanv $CHANV -token $TOKEN"]