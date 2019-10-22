#!/bin/bash

echo "Photomaton..."
echo "Starting"
/home/pi/go/src/go-usbmuxd/server &>/home/pi/photomaton.log &
echo "End launch"