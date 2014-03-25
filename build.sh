#!/bin/bash

if [ "$#" == 0 ]; then
  echo "Usage: ./build.sh <flag>"
  echo "  --build     Builds argo."
  echo "  --install   Installs argo. Must be run as root."
  echo "  --uninstall Removes argo. Must be run as root."
  exit 1
fi

if [ "$1" == "--build" ]; then
  echo "# Building argo."
  echo "# Installing dependencies."
  go get github.com/schleibinger/sio
  go get -tags zmq_2_1 github.com/alecthomas/gozmq
  go get github.com/vmihailenco/msgpack
  go get github.com/gorilla/mux
  go get github.com/gorilla/websocket
  echo "# Compiling."
  go build
  echo "# Done."
fi

if [ "$1" == "--install" ]; then
  echo "# Installing argo."
  install -p -g root -o root -m 755 argo /usr/bin
  install -p -g root -o root -m 644 actisense.rules /etc/udev/rules.d
  install -p -g root -o root -m 755 actisense.sh /lib/udev/actisense
fi

if [ "$1" == "--uninstall" ]; then
  echo "# Removing argo."
  rm /usr/bin/argo
  rm /etc/udev/rules.d/actisense.rules
  rm /lib/udev/actisense
fi
