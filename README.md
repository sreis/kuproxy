# kuproxy

Proof of concept for a Kubernetes aware external load balancer. Uses
HAproxy under the hood for the actual load balancing.

The idea is to have an haproxy that self configures when a Kubernetes
pod goes online/offline.
