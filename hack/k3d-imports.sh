#!/bin/bash

CERT_MANAGER_VERSION=${CERT_MANAGER_VERSION:-"v1.9.1"}
K3S_CONTEXT=${K3S_CONTEXT:-"k3s-default"}
IMAGE_REPO=${IMAGE_REPO:-"quay.io/jetstack"}
IMAGES=$(cat <<-END
  cert-manager-webhook
  cert-manager-controller
  cert-manager-cainjector
END
)

for img in ${IMAGES}; do
  _full_img=${IMAGE_REPO}/${img}:${CERT_MANAGER_VERSION}
  echo -e "\n\nImporting image ${_full_img}"
  k3d -c ${K3S_CONTEXT} image import ${_full_img}
done
