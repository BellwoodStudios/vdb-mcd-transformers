#! /usr/bin/env bash

set -e

function message() {
    echo
    echo -----------------------------------
    echo "$@"
    echo -----------------------------------
    echo
}

ENVIRONMENT=$1
if [ "$ENVIRONMENT" == "prod" ]; then
TAG=latest
elif [ "$ENVIRONMENT" == "staging" ]; then
TAG=staging
else
   message UNKNOWN ENVIRONMENT
fi

if [ -z "$ENVIRONMENT" ]; then
    echo 'You must specifiy an envionrment (bash deploy.sh <ENVIRONMENT>).'
    echo 'Allowed values are "staging" or "prod"'
    exit 1
fi

message BUILDING EXECUTE DOCKER IMAGE
docker build -f dockerfiles/execute/Dockerfile . -t makerdao/vdb-execute:$TAG

message BUILDING BACKFILL-STORAGE DOCKER IMAGE
docker build -f dockerfiles/backfill_storage/Dockerfile . -t makerdao/vdb-backfill-storage:$TAG

message LOGGING INTO DOCKERHUB
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USER" --password-stdin

message PUSHING EXECUTE DOCKER IMAGE
docker push makerdao/vdb-execute:$TAG

message PUSHING BACKFILL-STORAGE DOCKER IMAGE
docker push makerdao/vdb-backfill-storage:$TAG

# service deploy
if [ "$ENVIRONMENT" == "prod" ]; then
  message DEPLOYING EXECUTE
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-execute-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$PROD_REGION.amazonaws.com --region $PROD_REGION
elif [ "$ENVIRONMENT" == "staging" ]; then
  message DEPLOYING EXECUTE
  aws ecs update-service --cluster vdb-cluster-$ENVIRONMENT --service vdb-execute-$ENVIRONMENT --force-new-deployment --endpoint https://ecs.$STAGING_REGION.amazonaws.com --region $STAGING_REGION
else
   message UNKNOWN ENVIRONMENT
fi
