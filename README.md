## Prerequisites

Before getting started, ensure you have done the following things:

1. Install `Go`
2. Install `Docker`

## Getting Started

To start this repo on your machine, do the following:

1.  Clone this repo
2.  Go to the repo directory on your machine
3.  Execute `go mod tidy && go mod vendor`
4.  Setup the environment variables in `./.env` file (see below for details)
5.  Start `Docker` on your machine
6.  Run `docker compose up -d`
7.  Finally, `go run ./module-name` to start the Go app

## Environment Variables

Setup the following environment variables in `./.env` file:

```
MONGODB_URI="atlas mongoDB URI"
PUBSUB_EMULATOR_HOST="localhost:8085"
```
