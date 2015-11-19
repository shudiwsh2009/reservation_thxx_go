#!/bin/bash
echo "#update git"
go get -u all

echo "#deploy website"
supervisorctl stop reservation_thxx_go
sleep 5
cd $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go
go build
supervisorctl start reservation_thxx_go

echo "#deploy reminder"
cd $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go_reminder
go build
echo "restart cron"
service cron restart

echo "###Deployment Completed"
