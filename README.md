# Telegram Service

Telegram integration for Amadeus. The service listens for incoming Telegram messages, publishes `new_message` events to RabbitMQ, and exposes an MCP tool for sending messages.

## Requirements

- Go 1.26.3 or newer
- Docker with Docker Compose
- Telegram API credentials: `APP_ID` and `APP_HASH`
- The [`dev-infra`](https://github.com/AmadeusAI-dev/dev-infra) repository

## Install

Install Go dependencies:

```bash
go mod download
```

Create the local configuration:

```bash
cp .env.example .env
```

Fill in `.env` with your Telegram credentials, phone number, optional two-factor authentication password, and the desired MCP host and port.

## Configure RabbitMQ

This step is required. Without the exchange, queue, and binding, local message delivery will not work.

1. Start the local infrastructure from the `dev-infra` repository:

   ```bash
   cd ../dev-infra
   make up
   ```

2. Open the [RabbitMQ Management UI](http://localhost:15672).
3. Sign in with username `amadeus` and password `amadeus`.
4. Open **Exchanges** and create a `direct` exchange named `main`.
5. Open **Queues and Streams** and create a queue, for example `telegram_messages`.
6. Open the queue and add a binding from exchange `main` with routing key `new_message`.

## Authorize Telegram

Create the local Telegram session before the first service start:

```bash
make authorize
```

Enter the confirmation code sent by Telegram when prompted. The session will be stored in the file configured by `SESSION_FILE` and must not be committed to Git.

## Run

Make sure the local infrastructure is running and then start the service:

```bash
cd ../telegram-service
make up
```

With the example configuration, the MCP HTTP server listens on [http://localhost:8080](http://localhost:8080).

## Checks

```bash
make lint
make test
```
