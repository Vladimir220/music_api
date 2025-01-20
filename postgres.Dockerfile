FROM postgres

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

EXPOSE 1234

CMD go run ${PWD}