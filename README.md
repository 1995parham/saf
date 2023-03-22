<<<<<<< HEAD

# nats-jetstream

Test NATS Jetstream features before using them on production

## How we can have a system account?

Check out the `cluster1.yaml` to see how we can have system account in helm values.
Please note that this doesn't affect applications and they can continue without authentication.

```yaml
auth:
  enabled: true
  systemAccount: admin
  basic:
    accounts:
      admin:
        users:
          - user: admin
            password: amdin
```

## Super-Cluster

First create a simple cluster without any gateway configuration and then create the following stream:

```sh
nats stream new rides --subjects 'ride.accepted, ride.finish' --max-age '5m' --max-bytes '10m' --replicas 2 --storage memory --retention limits --discard old

nats pub --count 10 ride.accepted 'ride {{ID}} started on {{Time}}'
```

Next upgrade its configuration to have gateway and also create a new cluster to form a super cluster and see how it works with jetstream.
Then you can see streams in both regions but each stream has its leader in its cluster.

```sh
nats stream report
nats server request gateways --user admin --password amdin | jq
```

Last step is to create a new stream in the new cluster to see it will be synced to the old cluster.

```sh
nats stream new murche --subjects 'ride.eta' --max-age '5m' --max-bytes '10m' --replicas 2 --storage memory --retention limits --discard old
```

# saf

<h1 align="center">Saf</h1>
<h6 align="center">Saf means Queue in Persian</h6>

<p align="center">

<img src="https://img.shields.io/github/actions/workflow/status/1995parham/saf/ci.yaml?label=ci&logo=github&style=for-the-badge&branch=main" alt="GitHub Workflow Status" />

<a href="https://codecov.io/gh/1995parham/saf">
<img src="https://img.shields.io/codecov/c/gh/1995parham/saf?logo=codecov&style=for-the-badge" alt="Codecov" />
</a>

</p>

## Introduction

One of the problems, that we face on projects with queues is deploying [RabbitMQ](https://www.rabbitmq.com/) on the cloud which brings us many challenges for CPU load, split brain, etc.
I want to see how [NATS](https://nats.io/) with [Jetstream](https://docs.nats.io/nats-concepts/jetstream) can work as a queue manager to replace RabbitMQ.
Saf project here have been developed to show the Jetstream usage as a queue manager and also load test it.

Meanwhile, we did some real-world tests and the results are not good, after the recent Erlang upgrades to make it cloud compatible,
I think RabbitMQ will work better.

## Description

Saf gets events from its producer side and publish them into NATS.
The consumer side gets events from NATS and do the process which may takes time.
The producer side here is an HTTP server which validate the given event request, then
after embedding trace information, it publishes event into the Jetstream.
[Stream](https://docs.nats.io/nats-concepts/jetstream/streams) created by the consumer side, and it is defined in `internal/cmq/cmq.go`.

Produce uses a same channel for all subjects and marshals data into JSON.

## Scenarios

On production, we need to handle two following scenarios when using a queue manager:

1. We have event on the producer side, but there isn't any available server, so we need to send an error.
2. There is no consumer available so events must be available on NATS until it gets back.

Jetstream stores messages for 1 hour in memory. So you can shut down the consumer
and send events happily and then after consumer starts again consumes these events.

The following description shows the stream that stores messages in memory:

```bash
# you can install natscli from https://github.com/nats-io/natscli
nats stream ls
```

```
╭─────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│                                                     Streams                                                     │
├────────┬────────────────────────────────────────────────┬─────────────────────┬──────────┬───────┬──────────────┤
│ Name   │ Description                                    │ Created             │ Messages │ Size  │ Last Message │
├────────┼────────────────────────────────────────────────┼─────────────────────┼──────────┼───────┼──────────────┤
│ events │ Saf's event channel contains only events topic │ 2023-03-22 14:46:02 │ 1        │ 188 B │ 4m9s         │
╰────────┴────────────────────────────────────────────────┴─────────────────────┴──────────┴───────┴──────────────╯
```

Also, you can see the description of the consumer:

```bash
nats consumer info events saf
```

```
Information for Consumer events > saf created 2023-03-22T14:46:02+03:30

Configuration:

                Name: saf
    Delivery Subject: _INBOX.u4ykdwReyQtmYjiZpxEJNM
      Filter Subject: events
      Deliver Policy: All
 Deliver Queue Group: saf
          Ack Policy: Explicit
            Ack Wait: 30s
       Replay Policy: Instant
     Max Ack Pending: 1,000
        Flow Control: false

State:

   Last Delivered Message: Consumer sequence: 1 Stream sequence: 1 Last delivery: 7m54s ago
     Acknowledgment floor: Consumer sequence: 1 Stream sequence: 1 Last Ack: 7m54s ago
         Outstanding Acks: 0 out of maximum 1,000
     Redelivered Messages: 0
     Unprocessed Messages: 0
          Active Interest: Active using Queue Group saf
```

Please note that we have the [Push consumer](https://natsbyexample.com/examples/jetstream/push-consumer/go)
and [Pull Consumer](https://natsbyexample.com/examples/jetstream/pull-consumer/go) in code and,
you can switch between them by commenting out the one you don't want.

## Up and Running

Everything you need to test the project and gather some results are available
in a single docker-compose:

```bash
cd deployment && docker compose up
```

Then you can build and run:

```bash
cd cmd/saf && go build
```

## APIs

The following API call publishes an event to subject `hello`.

```bash
curl -X POST -d '{ "subject": "hello" }' -H 'Content-Type: application/json' http://127.0.0.1:1378/api/event
```

## Clustering

One of the main issues about using Jetstream is how does its clustering work?
There are three RAFT groups exist in a single Jetstream cluster:

1. Meta Group: all servers join the Meta Group and the Jetstream API is managed by this group.
   A leader is elected and this owns the API and takes care of server placement.

2. Stream Group: each Stream creates a RAFT group, this group synchronizes state and data between its members.
   The elected leader handles ACKs and so forth, if there is no leader the stream will not accept messages.

3. Consumer Group: each Consumer creates a RAFT group, this group synchronizes consumer state between its members.
   The group will live on the machines where the Stream Group is and handle consumption ACKs etc. Each Consumer will have their own group.
