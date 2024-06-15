FROM golang:latest

WORKDIR /app

ENV GO111MODULE=on

RUN apt install make

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . .

EXPOSE 7095

CMD [ "air", "-build.cmd", "make compile", "-build.bin", "/app/bin/emc", "-build.include_ext", "['.go']", "-build.include_dir", "cmd,pkg,internal,api" ]