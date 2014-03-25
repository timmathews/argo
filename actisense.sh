#!/bin/sh

/sbin/modprobe ftdi_sio && echo 0403 d9aa > /sys/bus/usb-serial/drivers/ftdi_sio/new_id
