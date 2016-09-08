#!/bin/bash
echo "#update git"
cd $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go
git reset --hard
git fetch $1
git checkout $1/$2

echo "#deploy website"
cd ~/thxxfzzx_go
go install github.com/shudiwsh2009/reservation_thxx_go/server
kill -9 $(lsof -t -i:8080)
sleep 5
now=$(date +"%Y_%m_%d_%T")
mv ~/thxxfzzx_go/log/server.log ~/thxxfzzx_go/log/server-${now}.log
cp -f $GOPATH/bin/server ./server/server.run
cp -rf $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go/templates ./templates
cp -rf $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go/assets ./assets
cd ~/thxxfzzx_go/server
chmod a+x ./server.run
nohup ./server.run --app-env="ONLINE" --sms-uid="shudiwsh2009" --sms-key="946fee2e7ad699b065f1" > ~/thxxfzzx_go/log/server.log 2>&1 & echo $! > ~/thxxfzzx_go/run.pid &

echo "#deploy reminder"
cd ~/thxxfzzx_go
go install github.com/shudiwsh2009/reservation_thxx_go/reminder
cp -f $GOPATH/bin/reminder ./server/reminder.run
chmod a+x ./server/reminder.run
#0 20 * * * /root/thxxfzzx_go/server/reminder.run --app-env="ONLINE" --sms-uid="shudiwsh2009" --sms-key="946fee2e7ad699b065f1"
echo "restart cron"
service cron restart

echo "###Deployment Completed"
