#!/bin/sh
set -ex
os=$(go env GOOS)
arch=$(go env GOARCH)

curl -L https://go.kubebuilder.io/dl/2.3.1/${os}/${arch} | tar -xz -C /tmp/

export TARGET=$HOME/kubebuilder
sudo mv /tmp/kubebuilder_2.3.1_${os}_${arch}/* ${TARGET}
ls -la ${TARGET}
ls -la ${TARGET}/bin
export PATH=$PATH:$TARGET/bin
