<h1 align="center">Saf</h1>
<h6 align="center">Saf means Queue in Persian</h6>

<p align="center">
  <img src="./.github/assets/banner.jpg" />
  <br>
  <img src="https://img.shields.io/github/actions/workflow/status/1995parham/saf/ci.yml?label=ci&logo=github&style=for-the-badge&branch=main" alt="GitHub Workflow Status" />
  <a href="https://codecov.io/gh/1995parham/saf">
    <img src="https://img.shields.io/codecov/c/gh/1995parham/saf?logo=codecov&style=for-the-badge" alt="Codecov" />
  </a>
  <a href="https://1995parham-teaching.github.io/ie-lecture/lectures/nats101/#/1">
    <img alt="Static Badge" src="https://img.shields.io/badge/Presentation-nats101-orange?style=for-the-badge&logo=revealdotjs">
  </a>
</p>

## Introduction

One of the problems, that we face on projects with queues is deploying [RabbitMQ](https://www.rabbitmq.com/) on the cloud which brings us many challenges for CPU load, split brain, etc.
I want to see how [NATS](https://nats.io/) with [Jetstream](https://docs.nats.io/nats-concepts/jetstream) can work as a queue manager to replace RabbitMQ.
Saf project here have been developed to show the Jetstream usage as a queue manager and also load test it.

Meanwhile, we did some real-world tests and the results are not good, after the recent Erlang upgrades to make it cloud compatible,
I think RabbitMQ will work better.

### What is NATS?

NATS messaging enables the exchange of data that is segmented into messages among computer applications and services.
These messages are addressed by subjects and do not depend on network location. This provides an abstraction layer
between the application or service and the underlying physical network. Data is encoded and framed as a message and sent
by a publisher. The message is received, decoded, and processed by one or more subscribers.

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
2. There is no consumer available, so events must be available on NATS until it gets back.

Jetstream stores messages for 1 hour in memory. So you can shut down the consumer
and send events happily and then after consumer starts again consumes these events.

The following description shows the stream that stores messages in memory:

```bash
# you can install natscli from https://github.com/nats-io/natscli
nats stream ls
```

```text
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

```text
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

Please note that we have the [Pull Consumer](https://natsbyexample.com/examples/jetstream/pull-consumer/go) using
new JetStream client API in code.

In case of having not important consumers, you can use NATS core consume type from your subject
(in this case, you don't need any stream). These consumers lose messages in their downtime
(please note that, NATS use TCP and you will not lose any messages in case of healthy producer and consumer)
but they have a little footprint on server.

### Channels

After consuming messages from NATS, Saf send events into Channels.
Channels are homemade concept of the Saf project. They are actually an interface defined as:

```go
type Channel interface {
 Init(*zap.Logger, trace.Tracer, interface{}, <-chan model.ChanneledEvent)
 Run()
 Name() string
}
```

Each channel has a name by which is referenced in configuration. They can accept arbitrary configuration (as an empty
interface) and they must validate it in their initiation process. Each channel has a way to receive events, and it will
be set by `SetChannel`. So we can describe a channel lifecycle as follows:

```go
Init() -> Run()
```

## Up and Running

Everything you need to test the project and gather some results are available
in a single docker-compose:

```bash
just dev up
```

## on Kubernetes

You can also deploy NATS on Kubernetes cluster using [NATS helm chart](https://github.com/nats-io/k8s).
Values for deploying two clusters are available in `./deployments/k8s/`.
Official chart support NATS exporter by default, and it can set up a `ServiceMonitor` too.

Also, Saf itself has charts (for producer and consumer) in `./charts` that you
can use to deploy it.

After deployment, you can easily use port-forwarding to send events and test your environment.

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

## Super-Cluster (with Jetstream)

First create a simple cluster without any gateway configuration and then create the following stream:

```bash
nats stream new rides --subjects 'ride.accepted, ride.finish' --max-age '5m' --max-bytes '10m' --replicas 2 --storage memory --retention limits --discard old

nats pub --count 10 ride.accepted 'ride {{ID}} started on {{Time}}'
```

Next upgrade its configuration to have gateway and also create a new cluster to form a super cluster and see how it works with Jetstream.
Then you can see streams in both regions, but each stream has its leader in its cluster.

```bash
nats stream report
nats server request gateways --user admin --password amdin | jq
```

Last step is to create a new stream in the new cluster to see it will be synced to the old cluster.

```bash
nats stream new murche --subjects 'ride.eta' --max-age '5m' --max-bytes '10m' --replicas 2 --storage memory --retention limits --discard old
```

## Enabling Authentication

In a scenario where the NATS cluster is already set up and users are actively producing and consuming from it,
the goal is to enable authentication without causing any downtime.

To enable authentication with zero downtime in a NATS cluster that currently lacks any form of authentication,
you can follow these steps:

1. Define users and passwords: Begin by creating a list of users and their corresponding passwords.
   Determine the access privileges and permissions for each user based on your requirements.
   This step can be completed offline without affecting the existing cluster.
2. Communicate new authentication requirements to clients: Once the user credentials are established,
   inform all clients (producers and consumers) about the upcoming authentication process.
   Request them to modify their configurations to include the username and password when connecting to the NATS cluster.
   (They can send usernames and passwords into the cluster without having any issue before enabling the authentication)
3. Enable authentication
4. Monitor and troubleshoot

By following these steps, you can effectively enable authentication on the NATS cluster without causing any downtime
for the existing producers and consumers.

### System Account

By enabling authentication, you can define an account (again, accounts are like tenants to NATS) to be a system account.
Using the system account you can read system events and APIs, so you can have access into the systematic commands of `natscli` too.

## Remove duplicate messages

Clients may send a same request multiple times, Jetstream can remove duplicate message based on their ID.
Each message has an ID header, and you can use your application logic to provide that ID, and ask Jetstream
to remove those messages.

Consider the following request:

```bash
curl -X POST -d '{ "subject": "hello", "data": "Hello World", "id": "1" }' -H 'Content-Type: application/json' http://127.0.0.1:1378/api/event
```

It sends a message with identification equals to 1, if your send another request you will not see that request on consumers.
