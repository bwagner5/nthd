# AWS Node Termination Handler Daemon (nthd)

## Setup:

```
$ adduser --system --user-group nthd

$ cat << EOF > /etc/systemd/system/nthd.service
[Unit]
Description=Node Termination Handler
Requires=dbus.socket
After=network.target

[Service]
Type=simple
User=nthd
Group=nthd
ExecStart=/usr/bin/nthd
TimeoutSec=30
Restart=on-failure
RestartSec=10
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=nthd

[Install]
WantedBy=multi-user.target
EOF

$ mkdir -p /etc/polkit-1/localauthority/50-local.d

$ cat << EOF > /etc/polkit-1/localauthority/50-local.d/allow_nthd_user_to_shutdown.pkla
[Allow nthd to shutdown]
Identity=unix-user:nthd
Action=org.freedesktop.login1.power-off-multiple-sessions
ResultActive=yes
EOF

$ systemctl daemon-reload
$ systemctl enable nthd.service
$ systemctl start nthd.service
```