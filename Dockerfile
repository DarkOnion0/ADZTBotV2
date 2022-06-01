FROM docker.io/library/golang:1.17 as builder

WORKDIR /src/adztbotv2

# Copy project file
COPY . .

# Donwload and Install project
RUN go get -d -v ./...
RUN env CGO_ENABLED=0 go build -o ADZTBotV2 main.go

# Create a new very lightweight image for the runtime
FROM docker.io/library/alpine:latest

WORKDIR /src/adztbotv2

# Copy the executable build i nthe previous step
COPY --from=builder /src/adztbotv2/ADZTBotV2 /src/adztbotv2/

# Env variables
ENV DB none
ENV URL none
ENV CHANM none
ENV CHANV none
ENV TOKEN none
ENV ADMIN 0
ENV DEBUG false
ENV CRON "59 23 * * *"

CMD ["sh", "-c", "./ADZTBotV2 -db $DB -url $URL -chanm $CHANM -chanv $CHANV -token $TOKEN -admin $ADMIN -debug $DEBUG -cron $CRON"]
