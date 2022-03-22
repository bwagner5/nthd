#!/usr/bin/env bash
set -euo pipefail

systemctl stop nthd.service || :
systemctl disable nthd.service || : 
systemctl daemon-reload || :
systemctl reset-failed || :

if [[ $(getent passwd nthd) ]]; then
    userdel -r nthd
fi
