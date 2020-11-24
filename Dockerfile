FROM debian:buster-slim as builder

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

RUN apt-get update -y \
	&& apt-get -yy -q install --no-install-recommends --no-install-suggests --fix-missing \
		bash-static curl tar gzip git ca-certificates netcat-openbsd

RUN cp /bin/bash-static /out/sh

RUN echo using jx-mink version $VERSION and OS $TARGETOS arch $TARGETARCH && \
  cd /tmp && \
  curl -k -L https://github.com/jenkins-x-plugins/jx-mink/releases/download/v$VERSION/jx-mink-$TARGETOS-$TARGETARCH.tar.gz | tar xzv && \
  mv jx-mink /out/jx-mink

FROM gcr.io/jenkinsxio/mink/mink:v20201124-local-6ea9cba4-dirty@sha256:5cb24ad8efffc82c6ed4f6a95b292fa5e068b6cd85743230d8b5c6179c49460e

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

COPY --from=0 /out /bin

ADD minx.sh kaniko.sh /usr/bin/

ENV PATH /usr/local/bin:/bin:/usr/bin:/kaniko

ENTRYPOINT ["minx.sh"]
