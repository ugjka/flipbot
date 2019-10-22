# flipbot

The bot for #rschizophrenia on freenode

## depends

Needs [translate-shell](https://github.com/soimort/translate-shell) and [googler](https://github.com/jarun/googler) package

## systemd configuration example

```ini
[Service]
Type=idle
ExecStart=/home/flipbot/flipbot
User=flipbot
Group=flipbot
Restart=always
RestartSec=60
WorkingDirectory=/home/flipbot/
Environment="FLIPBOT_EMAIL=<your email>"
Environment="FLIPBOT_YOUTUBE=<youtube api key>"
Environment="FLIPBOT_SUB=/r/someredddit/new"
Environment="FLIPBOT_SERVER=chat.freenode.net:6667"
Environment="FLIPBOT_NICK=<botnick>"
Environment="FLIPBOT_PASS=<nick pass>"
Environment="FLIPBOT_CHAN=#examplechan"
Environment="FLIPBOT_OW=<open weather map api key>"
Environment="FLIPBOT_OP=<op nick>"
Environment="FLIPBOT_SERVER_MAIL=server email"
Environment="FLIPBOT_WOLF=<wolfram alpha api key>"
[Install]
WantedBy=network-online.target
```
