
# shadowsocks-mini 

a simple and mini implementation of shadowsocks-server

rewrite shadowsocks in go to learn how it works

thanks [go-shadowsocks2](https://github.com/shadowsocks/go-shadowsocks2)


## Features
- very simple and easy, just a few files
- capture GFW's active probe and relay attack

## Install
```bash
go get -u -v github.com/Howard0o0/shadowsocks-mini@latest
```

## Usage

1. Create ssmini's workspace
```bash
#if not in root, use sudo 
touch /etc/ssmini
```

2. Create a config.json save as /etc/ssmini/config.json 
{
        "listen_port": "1998",
        "password": "your_password",
        "method": "AEAD_CHACHA20_POLY1305",
        "logdir": "path/to/save/log"
}
```

3. Run 
```bash
nohup $GOPATH/bin/shadowsocks-mini  /etc/ssmini/config.json &
```




