FROM registry.ci.openshift.org/stolostron/builder:go1.25-linux AS builder

WORKDIR /go/src/open-cluster-management.io/multicloud-operators-application
COPY . .
RUN make -f Makefile build

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

RUN microdnf update -y && \
     microdnf clean all

ENV OPERATOR=/usr/local/bin/multicluster-operators-application \
    USER_UID=1001 \
    USER_NAME=multicluster-operators-application

# copy applcation CRD
COPY deploy/crds/*.yaml  /usr/local/etc/application/crds/

# install operator binary
COPY --from=builder go/src/open-cluster-management.io/multicloud-operators-application/build/_output/bin/multicluster-operators-application ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
