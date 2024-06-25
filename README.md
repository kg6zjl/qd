[![Build & Release](https://github.com/kg6zjl/qd/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/kg6zjl/qd/actions/workflows/goreleaser.yml) [![Tests](https://github.com/kg6zjl/qd/actions/workflows/lint-and-check.yml/badge.svg)](https://github.com/kg6zjl/qd/actions/workflows/lint-and-check.yml)

<p align="center">
  <img src="images/qd.jpg?raw=true" alt="Happy Containers" width="40%">
</p>

# K8s Quick Deploy
A super simple way to deploy a container to K8s.

## Overview
```
qd run alpine:latest

qd run centos

qd exec ubuntu:20.04

qd list

qd stop
```

## Run
```
$ qd run alpine:latest
$ qd run node
$ qd run python
$ qd run ubuntu:20.04
```

## List
```
# list only qd deployments 
$ qd list
alpine-qd-740167019386
node-qd-249575108461
python-qd-208428483955
ubuntu-qd-650186231124

$ kubectl get pods
NAME                                      READY   STATUS    RESTARTS   AGE
alpine-qd-740167019386-85bb955857-82fbf   1/1     Running   0          48s
node-qd-249575108461-7964cf9c7b-v929q     1/1     Running   0          44s
python-qd-208428483955-66cfc9ccf4-hjsfh   1/1     Running   0          39s
ubuntu-qd-650186231124-57b69c4559-6zkb8   1/1     Running   0          85s
```

## Stop
```
$ qd stop
Stopped deployment alpine-qd-740167019386
Stopped deployment node-qd-249575108461
Stopped deployment python-qd-208428483955
Stopped deployment ubuntu-qd-650186231124
```

## Exec
```
$ qd exec ubuntu
```