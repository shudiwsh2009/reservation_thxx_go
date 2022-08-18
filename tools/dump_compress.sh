#!/bin/bash

export LC_ALL="en_US.UTF-8"
cd /app/dump
tar -zcvf mongodb.dump.$(date -d '1 days ago' +%Y_%m_%d).tar.gz mongodb.dump.$(date -d '1 days ago' +%Y.%m.%d)* && rm -rf mongodb.dump.$(date -d '1 days ago' +%Y.%m.%d)* && rm mongodb.dump.$(date -d '30 days ago' +%Y_%m_%d).tar.gz
