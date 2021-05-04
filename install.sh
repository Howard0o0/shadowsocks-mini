#!/bin/bash

# global param
gport=9797
gpasswd=default_password
workspace=/etc/ssmini

function getConfParam() {

	while true; do

		read -p "listen port:" port
		read -p "password:" passwd
		echo
		echo "+++++++++++++++++++++++"
		echo "listen port: $port"
		echo "password : $passwd"
		echo "-----------------------"
		echo

		read -r -p "Are You Sure? [Y/n] " input

		case $input in
		[yY][eE][sS] | [yY])
			echo "Yes"
			gport=$port
			gpasswd=$passwd
			# ./test -port $port -passwd $passwd
			return 0
			;;

		[nN][oO] | [nN])
			echo "restart \n"
			;;
		*)
			echo "Invalid input..."
			;;
		esac
	done
}

function genSSminiConf() {

	if [ ! -d "$workspace" ]; then
		mkdir $workspace
	fi
	if [ ! -d "$workspace/log" ]; then
		mkdir $workspace/log
	fi

	chmod -R 777 $workspace

	touch conf.json
	echo "{" >>conf.json
	echo "  \"listen_port\":\"$gport\",   " >>conf.json
	echo "  \"password\":\"$gpasswd\",   " >>conf.json
	echo "  \"method\":\"AEAD_CHACHA20_POLY1305\",   " >>conf.json
	echo "  \"logdir\":\"$workspace/log\"   " >>conf.json
	echo "}" >>conf.json

	mv -f conf.json $workspace
}

function downloadExec() {

	if [ -f "/usr/local/bin/ssmini" ]; then
		rm "/usr/local/bin/ssmini"
	fi

	wget --tries=5 --timeout=3 -cO /usr/local/bin/ssmini https://github.com/Howard0o0/shadowsocks-mini/releases/download/v3.0/ssmini
	if [ $? != 0 ]; then
		echo "download ssmini failed, please check your network situation and try again"
		exit 1
	fi

	chmod +x /usr/local/bin/ssmini
	echo "export PATH=$PATH:/usr/local/bin" >>/etc/profile
	source /etc/profile
}

function genSysctlConf() {

	touch ssmini.service

	echo "[Unit]" >>ssmini.service
	echo "Description=\"ssmini daemonize systemd\"" >>ssmini.service
	echo "After=network.target" >>ssmini.service
	echo "After=network-online.target" >>ssmini.service
	echo "Wants=network-online.target" >>ssmini.service
	echo >>ssmini.service
	echo "[Service]" >>ssmini.service
	echo "ExecStart=/usr/local/bin/ssmini /etc/ssmini/conf.json" >>ssmini.service
	echo >>ssmini.service
	echo "[Install]" >>ssmini.service
	echo "WantedBy=multi-user.target" >>ssmini.service

	mv -f ssmini.service /etc/systemd/system/
}

# 查看当前用户是否为root
function checkUser() {

	if [ "$(whoami)" == "root" ]; then
		return 0
	fi

	echo "current user is not root, please login as root and try again"
	exit 1
}

checkUser
getConfParam
genSSminiConf
downloadExec
genSysctlConf
systemctl enable ssmini
systemctl start ssmini

echo "ssmini installed"
ssmini -uri -conf $workspace/conf.json
