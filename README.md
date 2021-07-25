# saf

## Introduction

Saf in Persian means Queue. One of the problems, that we face on projects with queues is deploying RabbitMQ on the cloud which brings us many challenges for CPU load, etc.
I want to see how NATS with Jetstream can work as the queue to replace RabbitMQ.

## Description

Saf gets events from its producer side and publish them into NATS. The consumer side gets events from NATS and do the process which may takes time.

## Scenarios

1. We have event on the producer side but there isn't any available server so we need to send an error.
2. There is no consumer available so events must be available on NATS until they gets back.

## APIs

```sh
curl -X POST -d '{ "subject": "hello" }' -H 'Content-Type: application/json' http://127.0.0.1:1378/api/event
```
