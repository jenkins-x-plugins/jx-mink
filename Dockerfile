FROM debian:buster-slim as builder

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

RUN apt-get update -y \
	&& apt-get -yy -q install --no-install-recommends --no-install-suggests --fix-missing \
		bash-static curl

RUN cp /bin/bash-static /sh

RUN echo using jx-mink version $VERSION and OS $TARGETOS arch $TARGETARCH && \
  cd /tmp && \
  curl -L https://github.com/jenkins-x-plugins/jx-mink/releases/download/v$VERSION/jx-mink-$TARGETOS-$TARGETARCH.tar.gz | tar xzv && \
  mv jx-mink /sh

FROM gcr.io/jenkinsxio/mink/mink:v20201123-local-7fd0bff2-dirty

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

COPY --from=0 /sh /usr/bin

ADD minx.sh kaniko.sh /usr/bin/

ENV PATH /usr/local/bin:/bin:/usr/bin:/kaniko

ENTRYPOINT ["minx.sh"]
