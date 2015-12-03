#!/bin/bash
rm -rf /root/dump/$(date -d '30 days ago' +%Y%m%d)/
mkdir -p /root/dump/$(date +%Y%m%d)/
mongodump -h 127.0.0.1 -o /root/dump/$(date +%Y%m%d)/