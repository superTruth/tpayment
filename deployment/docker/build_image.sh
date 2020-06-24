#!/bin/bash
commit_sha=$(git rev-parse HEAD)
export IMAGE_NAME=bindo123/tpayment:$commit_sha

docker build . -t $IMAGE_NAME
if [ $? -eq 0 ]; then
    docker push $IMAGE_NAME
    echo $IMAGE_NAME
else
    error "Docker image build failed !"
fi
