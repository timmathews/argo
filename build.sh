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
  go get github.com/op/go-logging
  go get gopkg.in/vmihailenco/msgpack.v2
  go get github.com/wsxiaoys/terminal
  go get github.com/gorilla/mux
  go get github.com/gorilla/websocket
  go get git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git
  go get github.com/burntsushi/toml
  go get github.com/imdario/mergo
  echo "# Compiling."
  go build
  echo "# Done."
fi

if [ "$1" == "--install" ]; then
  echo "# Installing argo."
  install -p -g root -o root -m 755 argo /usr/bin
  install -p -g root -o root -m 644 actisense/actisense.rules /etc/udev/rules.d
  install -p -g root -o root -m 644 canusb/canusb.rules /etc/udev/rules.d
  install -p -g root -o root -m 755 actisense/actisense.sh /lib/udev/actisense
  install -p -g root -o root -m 755 canusb/canusb.sh /lib/udev/canusb
fi

if [ "$1" == "--uninstall" ]; then
  echo "# Removing argo."
  rm /usr/bin/argo
  rm /etc/udev/rules.d/actisense.rules
  rm /etc/udev/rules.d/canusb.rules
  rm /lib/udev/actisense
  rm /lib/udev/canusb
fi
