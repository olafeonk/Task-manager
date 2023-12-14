#!/bin/bash

#variables in env file
if [ -f ".env" ]; then
    export $(grep -v '^#' .env | xargs)
fi

if ! [ "${REMINDER_REGISTRY}" ]; then
    echo "SPACE_DAYS_REGISTRY is not set!"
    exit 1
fi


if ! [ "${SERVICE_ACCOUNT_ID}" ]; then
    echo "SERVICE_ACCOUNT_ID is not set!"
    exit 1
fi

REMINDER_REGISTRY=$(echo $REMINDER_REGISTRY | tr -d '\r')
REMINDER_BACKEND_CONTAINER_ID=$(echo $REMINDER_BACKEND_CONTAINER_ID | tr -d '\r')
SERVICE_ACCOUNT_ID=$(echo $SERVICE_ACCOUNT_ID | tr -d '\r')

new_image_name=$REMINDER_REGISTRY/task-server:prod;
echo $new_image_name;
docker build -t $new_image_name . ;
docker push $new_image_name;

yc sls container revisions deploy \
    --container-id ${REMINDER_BACKEND_CONTAINER_ID} \
    --memory 512M \
    --cores 1 \
    --execution-timeout 30s \
    --concurrency 8 \
    --min-instances 0 \
    --environment POSTGRES_PASSWORD=${POSTGRES_PASSWORD},GIN_MODE=${GIN_MODE},POSTGRES_HOST=${POSTGRES_HOST},POSTGRES_PORT=${POSTGRES_PORT},POSTGRES_USER=${POSTGRES_USER},POSTGRES_DB=${POSTGRES_DB}  \
    --service-account-id ${SERVICE_ACCOUNT_ID} \
    --image "$new_image_name";

