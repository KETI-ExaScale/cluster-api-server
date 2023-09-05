#!/bin/bash
docker_id="ketidevit2"
image_name="exascale.cluster-api-server"
operator="cluster-api-server"
version="latest"

export GO111MODULE=on
go mod vendor
#kubectl config view >> `pwd`/build/bin/config

go build -o ../build/_output/$operator -mod=vendor ../cmd/main.go && \
docker build -t $docker_id/$image_name:$version ../build && \
docker push $docker_id/$image_name:$version