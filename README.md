# kuproxy

Proof of concept for a Kubernetes aware external load balancer. Uses
HAproxy under the hood for the actual load balancing.

The idea is to have an haproxy that self configures when a Kubernetes
pod goes online/offline.

For now it just prints a bunch of debug info when a Kubernetes pod
goes online/offline.

## Build

```
go get github.com/coreos/go-etcd/etc
go build
```

## Run

`kuproxy` depends on a
**[Kubernetes](https://github.com/GoogleCloudPlatform/kubernetes)**
cluster.

For now you can use this **[Kubernetes example
cluster](https://github.com/pires/kubernetes-vagrant-coreos-cluster)**.

Find out where the Kubernetes master is running.

```
vagrant ssh master
kubectl cluster-info
```

Then start the load balancer with.

```
kuproxy --master="http://<kubernetes.master.ip>:2379"
```

Check everything is running by launching a Kubernetes pod. For
example:

```
vagrant ssh master
kubectl run-container nginx --image=nginx --replicas=2 --port=80
```

Stop the pod with.

```
vagrant ssh master
kubectl stop rc nginx
```

