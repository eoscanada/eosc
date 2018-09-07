#!/bin/bash -e

TAG_OUT=${1}

if [[ ${TAG_OUT} == "" ]]; then
   TAG_OUT=`date +"%Y-%m-%d-%H-%M-%S"`
fi

export CLOUDSDK_CORE_PROJECT=eoscanada-shared-services


gcloud builds submit . \
        --config cloudbuild_docs.yaml \
        --timeout 15m \
        --substitutions=SHORT_SHA=${TAG_OUT}

echo "TAG_OUT: ${TAG_OUT}"
