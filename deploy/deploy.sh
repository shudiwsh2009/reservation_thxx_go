#!/bin/zsh
echo "#update git"
cd $GOPATH/src/github.com/shudiwsh2009/reservation_thxx_go
echo "git fetch..."
git checkout master
git fetch --all
echo "git update..."
git reset --hard origin/master
git pull origin master

echo "#set system env"
source ./profile

echo "#deploy website"
echo "go build"
go build
echo "supervisorctl restart reservation_thxx_go"
supervisorctl restart reservation_thxx_go
sleep 5

echo "#deploy reminder"
echo "go build"
# TODO
echo "restart cron"
service cron restart

echo "###Deployment Completed"
