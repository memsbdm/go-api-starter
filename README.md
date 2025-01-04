# Starter

This is Go API project starter implementing net/http with hexagonal architecture.

Migrations are performed using pressly/goose.

## Getting started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Register DB container
```bash
make docker-up
```

Shutdown DB Container
```bash
make docker-down
```

DB Integration Tests:
```bash
make itest
```

Run the application
```bash
make run
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

Seed the database:
```bash
make seed
```