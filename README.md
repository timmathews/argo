Argo
====

Argo is a [WebSocket](http://www.w3.org/TR/websockets/) and
[ZeroMQ](http://zeromq.org/) [NMEA
2000](http://en.wikipedia.org/wiki/NMEA_2000) server based on
[CANboat](https://github.com/canboat/canboat) and written in
[Go](http://golang.org). It is in the early stages of development and should be
considered a proof-of-concept at this point and not a mature product.

Installation
------------

If you don't have go, you'll need to install that. See the golang [Getting
Started](http://golang.org/doc/install) guide.

After that, you can install argo.

```
$ sudo apt-get install libzmq-dev
$ cd $GOPATH/src
$ git clone git@github.com:timmathews/argo
$ cd argo
$ ./build.sh --build
$ sudo ./build.sh --install
```

By default, argo will try to use /dev/ttyUSB0 for the Actisense NGT-1, but you
can change this by calling argo with a specific device like ```argo
/dev/ttyUSB3```.

TODO
----

* Daemonize argo and write an upstart script for it

