#!/usr/bin/env bash

set -e

./build.sh

scp build/mrsans pi@alexchi-raspi.local:~

ssh pi@alexchi-raspi.local "sudo mv mrsans /opt/mrsans/ && sudo systemctl restart mrsans"
