FROM golang:1.16

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
RUN go install -v ./...

# Env variables
ENV DB none
ENV URL none
ENV CHANM none
ENV CHANV none
ENV TOKEN none

CMD ["sh", "-c", "ADZTBotV2 -db $DB -url $URL -chanm $CHANM -chanv $CHANV -token $TOKEN"]