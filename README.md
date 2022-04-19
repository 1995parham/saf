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
