#!/usr/bin/env bash

cd /app/log
tar -zcvf zip_log/server.log.$(date -d '1 days ago' +%Y_%m).tar.gz server.log.$(date -d '1 days ago' +%Y_%m)* && rm server.log.$(date -d '1 days ago' +%Y_%m)*
tar -zcvf zip_log/reminder.log.$(date -d '1 days ago' +%Y_%m).tar.gz reminder.log.$(date -d '1 days ago' +%Y_%m)* && rm reminder.log.$(date -d '1 days ago' +%Y_%m)*
tar -zcvf zip_log/feedback_reminder.log.$(date -d '1 days ago' +%Y_%m).tar.gz feedback_reminder.log.$(date -d '1 days ago' +%Y_%m)* && rm feedback_reminder.log.$(date -d '1 days ago' +%Y_%m)*
