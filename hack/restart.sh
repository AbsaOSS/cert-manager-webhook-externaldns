#!/bin/bash
DIR="${DIR:-$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )}"

k3d cluster delete
k3d cluster create --no-lb --k3s-arg "--disable=traefik,servicelb,metrics-server,local-storage@server:*" --image docker.io/rancher/k3s:v1.23.7-k3s1

# install cert-manager
${DIR}/k3d-imports.sh
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.9.1/cert-manager.yaml

sleep 15

# use either this or the following line
kubectl apply -f ${DIR}/../manifest/
#tilt up --stream