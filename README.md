# archof

Remote docker image arch checker

## Install

```console
$ go install github.com/udzura/archof
```

## Usage

```console
$ archof docker.io/amd64/ubuntu:latest
amd64

$ archof docker.io/arm64v8/ubuntu:latest
arm64

$ archof gcr.io/udzura-dev/uimage:latest --bearer "$(gcloud auth print-access-token)"
amd64
```
