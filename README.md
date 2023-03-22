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

Saf in persian means Queue. One the problem that we face is deploying rabbitmq on the cloud and it brings us many challengers.
I want to see how nats with jetstream can work as queue to replace rabbitmq.
