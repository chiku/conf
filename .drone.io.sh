#/bin/bash

set -e

setup_environment() {
  export GOROOT=$HOME/go
  export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
}

install_os_packages() {
  echo sudo apt-get install -y wget
}

install_golang() {
  gofile=go1.6.2.linux-amd64.tar.gz
  mkdir -p $GOROOT
  pushd $HOME
  wget --continue "https://storage.googleapis.com/golang/$gofile" -O "$HOME/$gofile"
  tar -zxvf $gofile
  popd
}

install_golang_deps() {
  go get -v github.com/mattn/goveralls
}

build() {
  make
  goveralls -v -service drone.io -repotoken $COVERALLS_TOKEN
}

run() {
  commnd=$1
  echo "--- Start $commnd ---"
  "$commnd"
  echo "--- End $commnd ---"
}

main() {
  run setup_environment
  run install_os_packages
  run install_golang
  run install_golang_deps
  run build
}

main
