# Deployment Guide

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Deployment Guide](#deployment-guide)
    - [RBAC](#rbac)
        - [Deployment](#deployment)
    - [General process](#general-process)
<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## RBAC

The service account is `multicluster-operators-application`.

The role `multicluster-operators-application` is binded to that service account.

### Deployment

```shell
cd multicloud-operators-application
kubectl apply -f deploy/crds/standalone
kubectl apply -f deploy/crds
kubectl apply -f deploy
```

## General process

Application CR:

```yaml
apiVersion: app.k8s.io/v1beta1
kind: Application
metadata:
  name: subscription-app
spec:
  componentKinds:
  - group: apps.open-cluster-management.io
    kind: Subscription
  descriptor: {}
  selector:
    matchExpressions:
    - key: app
      operator: In
      values:
      - subscription-app
```
