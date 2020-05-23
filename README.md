# nats101
[![Drone (cloud)](https://img.shields.io/drone/build/1995parham/nats101.svg?style=flat-square)](https://cloud.drone.io/1995parham/nats101)

## Introduction
NATS is message broker that works as a dial tone between services.
It can be used for broadcasting messages among services or it can be used for synchronous request/response.
In this repository we want to grasp on some of it features and try it on the cloud.

## Up and Running
Here we use the nats [chart](https://github.com/nats-io/k8s/tree/master/helm/charts)
to setup a cluster on the cloud. Based on this setup you can nats101 kubernetes files in `k8s/`.
