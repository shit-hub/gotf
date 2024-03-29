# GoTF

练手小项目：使用Go实现一个文件传输工具，可指定本地目录下所有文件，通过TCP传输到服务端，支持显示文件发送进度，支持多文件并发传输，并发数为CPU核数。

## 基本功能

- 读取命令行参数
- 递归读取文件夹下所有文件
- 支持发送大文件
- 用协程并发传输多文件
- 显示文件发送进度
- 支持AES加解密

## 程序用法

```txt
Usage of ./server:
  -aes string
        enable AES encrypt and set mode: CBC/ECB/CFB
  -aes-key string
        the key of AES encrypt (default "ABCDEFGHIJKLMNOP")
  -f string
        Download Path (default ".")
  -l string
        Listen Host (default "0.0.0.0:3000")
```

## 使用实例

- 服务端

```bash
./server -l 0.0.0.0:3000 -f ./dst -aes CBC
```

- 客户端

```bash
./client -c 127.0.0.1:3000 -f ./src -aes CBC
```

> 说明:
> 
> -l 指定服务器监听的地址
> 
> -c 指定服务器IP和端口
> 
> -f 指定文件夹路径
> 
> -aes 使能AES加密并制定AES加密模式
