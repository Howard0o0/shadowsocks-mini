
# shadowsocks-mini 

a simple and mini implementation of shadowsocks-server

rewrite shadowsocks in go to learn how it works

thanks [go-shadowsocks2](https://github.com/shadowsocks/go-shadowsocks2)


## Features
- very simple and easy, just a few files
- capture GFW's active probe and relay attack

## Install

```bash
# login as root first
su root
# one-click setup
bash <(curl -s https://raw.githubusercontent.com/Howard0o0/shadowsocks-mini/dev/install.sh)
```

once installed, ssmini is ready and running, which will self-boot after reboot

you can check log in /etc/ssmini/log

or check ssmini's status by usage below

## Usage

check status 
```bash
systemctl status ssmini
```

show URI 
```bash
ssmini -uri
```

stop 
```bash
systemctl stop ssmini
```

start 
```bash
systemctl start ssmini
```






