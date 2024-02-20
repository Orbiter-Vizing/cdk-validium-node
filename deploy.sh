#!/bin/bash

harbor="codefamer/cdk-validium-node"
tag="v0.0.2"

echo $harbor:$tag

docker build --no-cache -t ${harbor}:$tag .
docker push ${harbor}:${tag}
docker rmi ${harbor}:${tag}


echo "docker build success."