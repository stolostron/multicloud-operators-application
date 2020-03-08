FROM registry.access.redhat.com/ubi8/ubi-minimal:8.1-398

ENV OPERATOR=/usr/local/bin/multicloud-operators-application \
    USER_UID=1001 \
    USER_NAME=multicloud-operators-application

# copy applcation CRD
COPY deploy/crds/*.yaml  /usr/local/etc/application/crds/

# install operator binary
COPY build/_output/bin/multicloud-operators-application ${OPERATOR}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/entrypoint"]

USER ${USER_UID}
