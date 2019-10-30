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
Environment="FLIPBOT_SERVER=chat.freenode.net"
Environment="FLIPBOT_PORT=6667"
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

## Usage

```text
FL1PBOT'S COMMANDS

#Basic
!flip :flips the table
!unflip :unflips the table
!ping :replies with pong (useful to see if your connection has stalled)
!random :get random advaita talk
!sleep :sleep emoticon
!cat :random cat emoticon
!shrug :shrug emoticon
!top :show top posters for past 7 days
!god :talk to god via random number generator
!debug :print stupid joke
!meditation : get spiritual advice
!toss : join the community of tossers

#Composite
!dict <word> : look up definitions of the given word
!syn <word> : look up synonyms of the given word
!calc <query> : calculate things
!flip <text> :flips the given text (ascii only)
!hug <nick> :hug someone
!seen <nick> :show the last time fl1pbot saw the given nick online
!define <word> :get Urban Dictionary definition for the given word
!time <location> :get time for the given location
!google <query> :get first result from google for the given query
!news <query> :search google news
!ducker <query> :search duckduckgo
!youtube <query> :get first result from youtube for the given query
!wiki <query> :get wikipedia summary for the given query
!memo <nick> <message> :send the given message to the given nick when they appear online on the channel
!remindme <duration> <message> :timer that reminds you after certain duration. Example !remindme 1 day 2 hours 3 minutes cook dinner

#Wheather
!w <location> :get weather info for the given location
!wf <location> :get 3 day weather forecast for the given location

#Translation
!trans <text in foreign language> : translate the given text to english (language detection may not work with single words, use full sentences)
!trans :<2_or_3_letter_lang_code> <text in english> :translate the given english text to the specified language
Example: !trans :ru How are you? :translates "How are you?" to russian

You can find the languages codes here https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes

-------------------------------------------------
PM bugs and feature requests to ugjka on freenode
```
