#!/bin/bash
echo "#update git"
go get -u github.com/shudiwsh2009/reservation_thxx_go
go get -u github.com/shudiwsh2009/reservation_thxx_go_reminder

echo "#deploy website"
supervisorctl stop reservation_thxx_go
sleep 5
source /etc/profile
cd $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go
go build
supervisorctl start reservation_thxx_go

echo "#deploy reminder"
cd $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go_reminder
go build
echo "restart cron"
service cron restart

echo "###Deployment Completed"
