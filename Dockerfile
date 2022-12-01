# syntax=docker/dockerfile:1
# first stage to install dependencies and create executable
FROM golang:1.18-buster as builder
LABEL description="Gritface forum project @grit:lab"
LABEL creators="adrian1206, dicarste, oluwatosin, petrus_ambrosius and tvntvn"
WORKDIR /
COPY . ./
RUN go mod download
RUN go build -o gritface main.go

# second stage with only the executable and necessary files for the web app
# this stage will not copy the already existing database by default, but will create a new one, if you want to keep the old one,
# uncomment the line starting with "COPY" below
FROM golang:1.18-buster as prod
WORKDIR /app
COPY --from=builder ./gritface ./
COPY log ./log
COPY server/public_html ./server/public_html
COPY localhost.crt localhost.csr localhost.key Readme.Md ./
# COPY forum-db.db ./
EXPOSE 443
CMD [ "./gritface" ]