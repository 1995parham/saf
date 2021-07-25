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

Queue with NATS Jetstream to remove all the erlangs from cloud

## Introduction

Saf in Persian means Queue. One of the problems, that we face on projects with queues is deploying RabbitMQ on the cloud which brings us many challenges for CPU load, etc.
I want to see how NATS with Jetstream can work as the queue to replace RabbitMQ.

## Description

Saf gets events from its producer side and publish them into NATS. The consumer side gets events from NATS and do the process which may takes time.

## Scenarios

1. We have event on the producer side but there isn't any available server so we need to send an error.
2. There is no consumer available so events must be available on NATS until they gets back.
