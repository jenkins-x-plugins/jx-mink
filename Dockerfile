FROM gcr.io/jenkins-x-labs-bdd/mink/mink:v20201121-local-d4b0e10e-dirty@sha256:f95d488e1436ecc4bee3dbe59ba2e783e2e367791f1f41201cc9046a88623a32

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

RUN echo using jx-mink version $VERSION and OS $TARGETOS arch $TARGETARCH && \
  cd /tmp && \
  curl -L https://github.com/jenkins-x/jx-mink/releases/download/v$VERSION/jx-mink-$TARGETOS-$TARGETARCH.tar.gz | tar xzv && \
  mv jx-mink /usr/bin

ENTRYPOINT ["jx-mink"]
