# ChainSQL/go-chainsql-api
  
使用请参考 [test/main.go](./test/main.go)

# 安装运行
1. windoss
要把./cgofuns/cdll/win下面的所有dll文件与可执行程序放在一起。
2. linux
- 版本低于ubuntu-16.04时，把./cgofuns/cdll/linx下面的所有so文件与可执行程序放在一起。
- 版本等于或者高于ubuntu-16.04，不做其它多余操作。 
3. arm
- 以下方式暂时方式
- 打开cgofun/cgo.go文件，删除这一行：#cgo linux LDFLAGS: -Wl,-RPATH="./" -L ./clib/linux/ -lsignature -lboost_regex -lcrypto -lssl -ldl -lstdc++
- 在相应位置增加这一行：#cgo LDFLAGS: -L ./clib/arm/ -lsignature -lboost_regex -lssl -lcrypto -lstdc++ -ldl
4. aarch64
- 以下方式暂时方式
- 打开cgofun/cgo.go文件，删除这一行：#cgo linux LDFLAGS: -Wl,-RPATH="./" -L ./clib/linux/ -lsignature -lboost_regex -lcrypto -lssl -ldl -lstdc++
- 在相应位置增加这一行：#cgo LDFLAGS: -L ./clib/aarch64/ -lsignature -lboost_regex -lssl -lcrypto -lstdc++ -ldl
