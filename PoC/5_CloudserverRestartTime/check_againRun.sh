#!/bin/bash

while true; do echo "`date -u` `curl -v http://10.0.0.180:31113/status`" >> check_againRunLog.txt; sleep 1; clear; done
