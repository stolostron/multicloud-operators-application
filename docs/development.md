# Development guide

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Development Guide](#development-guide)
    - [Launch dev mode](#launch-dev-mode)
    - [Build a local image](#build-a-local-image)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Launch dev mode

Run the following command to launch developer mode:

```shell
git clone git@github.com:stolostron/multicloud-operators-application.git
cd multicloud-operators-application
export GITHUB_USER=<github_user>
export GITHUB_TOKEN=<github_token>
make
make build
kubectl apply -f deploy/crds/standalone
export POD_NAMESPACE=<pod namespace to wire up webhook>
./build/_output/bin/multicluster-operators-application --application-crd-file deploy/crds/app.k8s.io_applications_crd_v1.yaml
```

## Build a local image

Build a local image by running the following command:

```shell
git clone git@github.com:stolostron/multicloud-operators-application.git
cd multicloud-operators-application
export GITHUB_USER=<github_user>
export GITHUB_TOKEN=<github_token>
make
make build-images
```
