# AWS Node Termination Handler Daemon (nthd)

## Overview
`nthd` is a lightweight daemon that monitors the [EC2 Instance Metadata Service](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) (IMDS) (v1 or v2) for [Amazon EC2 Spot Interruption Termination Notices](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/spot-interruptions.html) (ITNs). When `nthd` discovers an ITN for that instance it is running, it will use the [dbus](https://www.freedesktop.org/wiki/Software/dbus/) to send a `power-off-multiple-sessions` signal to the operating system. 

`nthd`'s behavior implements a "shutdown on spot ITN" which applications can handle generically by registing a [systemd-inhibitor](https://www.freedesktop.org/software/systemd/man/systemd-inhibit.html) lock. 

## Kubernetes (k8s)

Beginning in Kubernetes >= v1.21, the kubelet can register a systemd-inhibitor lock that will delay an OS shutdown until the kubelet can gracefully drain the node (cordon + eviction of pods). K8s calls this feature "[Graceful Node Shutdown](https://kubernetes.io/blog/2021/04/21/graceful-node-shutdown-beta/)". In order for k8s pods to handle an OS shutdown gracefully, they need catch the SIGTERM signal (sent by the kubelet) and perform any data checkpointing or connection draining.

## Installation

```
$ sudo yum install https://github.com/bwagner5/nthd/releases/download/v0.0.10/nthd_0.0.10_linux_amd64.rpm
```

## Security

`nthd` runs as an unprivileged user (the `nthd` user) and uses the [`polkit`](https://www.freedesktop.org/software/polkit/docs/latest/polkit.8.html) to authorize a `power-off-multiple-sessions`.  