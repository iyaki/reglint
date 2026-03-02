#!/usr/bin/env bash

action="${1:-start}"
port="${2:-4096}"

case "$action" in
	start)
		# Start in background so the script returns immediately.
		nohup opencode serve --port "$port" > /tmp/opencode-serve.log 2>&1 &
		;;
	stop)
		if [[ -n "$2" ]]; then
			pkill -f "opencode serve --port $port"
		else
			pkill -f "opencode serve"
		fi
		;;
	*)
		echo "Usage: $0 [start|stop] [port]"
		exit 1
		;;
esac
