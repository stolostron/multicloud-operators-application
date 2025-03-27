FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

RUN microdnf update && \
    microdnf clean all

ENV OPERATOR=/usr/local/bin/multicluster-operators-application \
    USER_UID=1001 \
    USER_NAME=multicluster-operators-application

# copy applcation CRD
COPY deploy/crds/*.yaml  /usr/local/etc/application/crds/

# install operator binary
COPY build/_output/bin/multicluster-operators-application ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}