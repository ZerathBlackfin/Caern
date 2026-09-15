#!/bin/sh
set -e

uid="${PUID:-0}"
gid="${PGID:-0}"

if [ "$uid" != "0" ] || [ "$gid" != "0" ]; then
	chown -R "$uid:$gid" /config 2>/dev/null || true
	exec su-exec "$uid:$gid" caern "$@"
fi

exec caern "$@"
