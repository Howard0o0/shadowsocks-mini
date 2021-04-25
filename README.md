
# shadowsocks-mini 

## Install
```bash
go get -u -v github.com/Howard0o0/shadowsocks-mini
```

## Usage

```bash
ssmini  /etc/ssmini.json
```

## ssmini.json
```json
{
        "listen_port": "1998",
        "password": "your_password",
        "method": "AEAD_CHACHA20_POLY1305",
        "logdir": "path/to/save/log"
}
```




