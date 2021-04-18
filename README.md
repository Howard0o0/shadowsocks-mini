
# shadowsocks-mini 

shadowsocks的精简实现 

暂时没有兼容现有的shadowsocks客户端,与shadowsocks的加密隧道协议有出入

## Install
```bash
go get -u -v github.com/Howard0o0/shadowsocks-mini
```

## Usage

### Sever 
配置文件通过绝对路径指定 
```bash
shadowsocks-mini -s -conf /etc/ssmini.json
```

### Client 
```bash
shadowsocks-mini -c -conf /etc/ssmini.json
```

## ssmini.json
```json
{
        "server": "server_ip",
        "server_port": "server_port",
        "local_address": "127.0.0.1",
        "local_port": "local_port",
        "password": "your_password",
        "method": "chacha20",
        "timeout": 300
}
```




