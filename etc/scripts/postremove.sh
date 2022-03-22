#!/usr/bin/env bash
set -euo pipefail

if [[ $(getent passwd nthd) ]]; then
    userdel -r nthd
fi
