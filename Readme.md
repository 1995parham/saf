[![Build Status](https://cloud.drone.io/api/badges/elahe-dastan/NATS/status.svg)](https://cloud.drone.io/elahe-dastan/NATS)

## NATS
In this repository I want to learn NATS and write a very simple code

## What is it
NATS messaging enables the exchange of data that is segmented into messages among computer applications and services.<br/>
These messages are addressed by subjects and do not depend on network location. This provides an abstraction layer<br/>
between the application or service and the underlying physical network. Data is encoded and framed as a message and sent<br/>
by a publisher. The message is received, decoded, and processed by one or more subscribers.

## Protocol
In this project I tried to code a publish subscribe protocol, there is a subject in common, publisher publishes a <br/>
message to that subject and all the ones who have subscribed that subject get the message.<br/>
the I used Queue Subscriber in which the ones who have subscribed a subject form a group and each time the publisher<br/>
publishes a message one the subscribers is randomly chosen to get the message so as you have probably guessed this <br/>
works as a load balancer

## Usage
To show the usage of nats as a load balancer I created a GET endpoint, each time it's called a heavy calculation<br/>
(simulated by 5 second sleeping) will be done 