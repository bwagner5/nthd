#!/usr/bin/env bash
set -euo pipefail

if ! [[ $(getent passwd nthd) ]]; then
    adduser --system -s /bin/false --user-group nthd
fi
