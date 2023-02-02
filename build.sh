#!/usr/bin/env bash


CURR_PATH=$(cd `dirname $0`;pwd)

function _info(){
	local msg=$1
	local now=`date '+%Y-%m-%d %H:%M:%S'`
	echo  "$now $msg"
}


function main() {
  local platform=$1
	case $platform in
	"linux")
		_info "开始构建Linux平台版本 ..."
		GOOS=linux GOARCH=amd64 \
		CGO_ENABLED=0 go build $main_file
		;;
	*)
		_info "开始本地构建 ..."
		CGO_ENABLED=0 go build $main_file
		;;
	esac
	if [ $? -ne 0 ];then
	    _info "构建失败"
		exit 1
	fi
	_info "程序构建完成: $CURR_PATH"
}

main ${1:-local}  .