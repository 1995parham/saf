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

nats pub ride.accepted 'i have a ride'
```

Next upgrade its configuration to have gateway and also create a new cluster to form a super cluster and see how it works with jetstream.
