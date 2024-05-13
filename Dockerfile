#  base image sa gooom
FROM golang:1.17-alpine

# postavlja se radni direktorijum
WORKDIR /app

# kopira se cijeli projekat
COPY . .

# Builduje se aplikacija
RUN go build -o main .

# otvara se port na kom slusa aplikacija zahtjeve
EXPOSE 8000

# Definisemo komandu sa kojom pokrecemo
CMD ["./main"]
