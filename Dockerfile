FROM docker.io/library/busybox:stable-uclibc

WORKDIR /src/adztbotv2

# Copy the executable build i nthe previous step
COPY adztbotv2 .

# Env variables
ENV DB none
ENV URL none
ENV CHANM none
ENV CHANV none
ENV TOKEN none
ENV ADMIN 0
ENV DEBUG false
ENV TIMER 3600000000000

CMD ["sh", "-c", "./adztbotv2 -db $DB -url $URL -chanm $CHANM -chanv $CHANV -token $TOKEN -admin $ADMIN -debug $DEBUG -timer $TIMER"]
