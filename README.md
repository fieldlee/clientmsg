# CLIENT MESSAGE
the client of message

# build window dll file
env GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -buildmode=c-shared -o clientmsg.dll

# modify config
utils/config.go 修改config文件的服务器ip 端口
和修改 客户端的端口等

# download
git clone https://github.com/fieldlee/clientmsg.git

go mod download
go mod vendor

# cgo

bridge.h bridge.c go 和 c的交互struct结构定义和函数定义
