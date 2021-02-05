cd /www/bin/server/wechat_server

echo 更新代码

git fetch --all

git reset --hard origin/master

echo 开始构建

export GOPROXY=https://goproxy.cn

/usr/local/go/bin/go build -o wechat_server cmd/main.go

echo 重启服务

supervisorctl restart wechat_server