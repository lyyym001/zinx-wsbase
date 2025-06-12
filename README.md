# zinx-wsbase
zinx-websocket基础
## msgID规定
    mId := msgID*1000 + subID


## 编译protobuf文件
    cd ./mdbc_server
    protoc --go_out=. ./pb/*.proto

## stun服务
### 1.

## 编译本地可执行程序
    1. windows
        go build -o ./app/server.exe ./server.go

### 本地测试
    1.GinService
        192.168.0.22:8080
    2.ZinxServer
        192.168.0.10106

### 参数说明
    1. Player字段说明
        Player.Status       0-未运行 1-登录服务器中 2-绑定中 3-数据获取中 4-服务器登录成功  
        Player.Role         1-红方 2-蓝方
        Player.AccountType  0-玩家 1-中控
    2. Att字段说明
        Att.Status          战斗状态 0-未开始 1-进行中 2-已结束
        Att.Flag            1-开始 2-重启 0-结束