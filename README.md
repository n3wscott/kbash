# kbash

Run a bash script on Kubernetes with Knative.

_Note: this is a work in progress. Please give feedback by filing an issue._

## Usage

This tool requires [`ko`](https://github.com/google/ko) installed locally and
[`Knative`](https://knative.dev) in the target Kubernetes cluster.

To deploy the simple `now.sh` demo,

```shell
ko apply -f ./config/service.yaml
```

Get the service url,

```shell
$ kubectl get ksvc
NAME      URL                                 LATESTCREATED      LATESTREADY    READY   REASON
kbash     http://kbash.default.example.com    kbash-xyz          kbash-xyz      True
```

Invoke the service url,

```shell
$ curl  http://kbash.default.example.com
Mon Mar 23 05:48:48 UTC 2020
```

## Next

Change the script in `./kodata/now.sh` to your needs. The base image can be
adjusted by editing `.ko.yaml`. The script name can be changed in
`./config/service.yaml` (see: `value: now.sh`).

## Local development.

You can launch the server locally and test,

```shell
FILE_PATH=./kodata go run ./main.go
```
