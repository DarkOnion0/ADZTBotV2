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

---

FROM alpine:latest

WORKDIR /src/adztbotv2

COPY --from=builder /src/adztbotv2/ADZTBotV2 /src/adztbotv2/

# Env variables
ENV DB none
ENV URL none
ENV CHANM none
ENV CHANV none
ENV TOKEN none

CMD ["sh", "-c", "ADZTBotV2 -db $DB -url $URL -chanm $CHANM -chanv $CHANV -token $TOKEN"]