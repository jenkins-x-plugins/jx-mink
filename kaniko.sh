#!/bin/sh

# source the variables and then run the jx-mink binary

source .jx/variables.sh

mkdir -p /kaniko/.docker
cp /tekton/creds-secrets/tekton-container-registry-auth/.dockerconfigjson /kaniko/.docker/config.json

if [ -z "$CONTEXT" ]
then
    export CONTEXT="/workspace/source"
fi

if [ -z "$DOCKERFILE" ]
then
    export DOCKERFILE="/workspace/source/Dockerfile"
fi

if [ -z "$DESTINATION" ]
then
    export DESTINATION="$DOCKER_REGISTRY/$DOCKER_REGISTRY_ORG/$APP_NAME:$VERSION"
fi

/kaniko/executor $KANIKO_FLAGS --context=$CONTEXT --dockerfile=$DOCKERFILE --destination=$DESTINATION

