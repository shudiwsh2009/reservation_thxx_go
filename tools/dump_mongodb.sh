#!/bin/bash

export LC_ALL="en_US.UTF-8"
now=$(date +"%Y.%m.%d.%H")
mkdir -p /app/dump/mongodb.dump.${now}
mongodump --host "$MONGO_HOST" -u "$MONGO_USERNAME" -p "$MONGO_PASSWORD" --authenticationDatabase "$MONGO_DATABASE" --db "$MONGO_DATABASE" -o /app/dump/mongodb.dump.${now}
