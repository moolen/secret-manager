#!/bin/sh
set -ex
os=$(go env GOOS)
arch=$(go env GOARCH)

curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

TARGET=${KUBEBUILDER_ASSETS:-/usr/local/kubebuilder}
sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch}/* ${TARGET}
export PATH=$PATH:/usr/local/kubebuilder/bin
