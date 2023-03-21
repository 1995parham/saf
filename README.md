<h1 align="center">Saf</h1>
<h6 align="center">Saf means Queue in Persian</h6>

<p align="center">

<img src="https://img.shields.io/github/actions/workflow/status/1995parham/saf/ci.yaml?label=ci&logo=github&style=for-the-badge&branch=main" alt="GitHub Workflow Status" />

<a href="https://codecov.io/gh/1995parham/saf">
<img src="https://img.shields.io/codecov/c/gh/1995parham/saf?logo=codecov&style=for-the-badge" alt="Codecov" />
</a>

</p>

## Introduction

One of the problems, that we face on projects with queues is deploying RabbitMQ on the cloud which brings us many challenges for CPU load, split brain, etc.
I want to see how NATS with Jetstream can work as a queue manager to replace RabbitMQ.
Saf project here have been developed to show the Jetstream usage as a queue manager and also load test it.

Meanwhile, we did some real-world tests and the results are not good, after the recent Erlang upgrades to make it cloud compatible,
I think RabbitMQ will work better.

## Description

Saf gets events from its producer side and publish them into NATS.
The consumer side gets events from NATS and do the process which may takes time.
The producer side here is an HTTP server.

Produce uses a same channel for all subjects and marshals data into JSON.

## Scenarios

1. We have event on the producer side, but there isn't any available server, so we need to send an error.
2. There is no consumer available so events must be available on NATS until it gets back.

## Up and Running

Everything you need to test the project and gather some results are available
in a single docker-compose:

```bash
cd deployment && docker compose up
```

## APIs

The following API call publishes an event to subject `hello`.

```bash
curl -X POST -d '{ "subject": "hello" }' -H 'Content-Type: application/json' http://127.0.0.1:1378/api/event
```
