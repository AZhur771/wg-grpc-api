#!/bin/bash
echo "dlv wrapper script to run debugger as sudo"
if [ "$DEBUG_AS_ROOT" = "true" ]; then
	exec /usr/bin/sudo -S -E $HOME/go/bin/dlv --only-same-user=false "$@"
else
	exec $HOME/go/bin/dlv "$@"
fi
