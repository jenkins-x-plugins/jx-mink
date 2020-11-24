FROM alpine

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add \
    ca-certificates curl git make netcat-openbsd
    
RUN echo using jx-mink version $VERSION and OS $TARGETOS arch $TARGETARCH && \
  cd /tmp && \
  curl -k -L https://github.com/jenkins-x-plugins/jx-mink/releases/download/v$VERSION/jx-mink-$TARGETOS-$TARGETARCH.tar.gz | tar xzv && \
  mv jx-mink /jx-mink

FROM gcr.io/jenkinsxio/mink/mink:v20201124-local-6ea9cba4-dirty@sha256:d9d545c2aafe88afffd8021adc455f86616f9712c3059ceb992e6a6049f512f1

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

COPY --from=0 /jx-mink /usr/bin/jx-mink

ADD minx.sh kaniko.sh /usr/bin/

ENV PATH /usr/local/bin:/bin:/usr/bin:/kaniko:/ko-app

ENTRYPOINT ["jx-mink", "resolve"]
