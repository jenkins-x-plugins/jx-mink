FROM gcr.io/jenkinsxio/mink/mink:v20201123-local-7fd0bff2-dirty

ARG BUILD_DATE
ARG VERSION
ARG REVISION
ARG TARGETARCH
ARG TARGETOS

LABEL maintainer="jenkins-x"

RUN addgroup -S app \
    && adduser -S -g app app \
    && apk --no-cache add \
    curl git

ADD minx.sh kaniko.sh /user/bin/

RUN echo using jx-mink version $VERSION and OS $TARGETOS arch $TARGETARCH && \
  cd /tmp && \
  curl -L https://github.com/jenkins-x-plugins/jx-mink/releases/download/v$VERSION/jx-mink-$TARGETOS-$TARGETARCH.tar.gz | tar xzv && \
  mv jx-mink /usr/bin

ENTRYPOINT ["minx.sh"]
