#!/usr/bin/env bash
set -euo pipefail

systemctl daemon-reload
systemctl enable nthd.service
systemctl start nthd.service
