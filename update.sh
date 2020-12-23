#!/bin/bash

cd /home/ugjka/sources/flipbot || exit
git checkout master || exit
git pull origin master || exit
go build || exit
git checkout bootybot || exit
git pull origin bootybot || exit
go build || exit
git checkout mysterybot || exit
git pull origin mysterybot || exit
go build || exit
mv flipbot /home/ugjka/flipbot/ || exit
mv bootybot /home/ugjka/bootybot/ || exit
mv mysterybot /home/ugjka/mysterybot/ || exit
sudo systemctl restart bootybot.service || exit
sudo systemctl start flipbotreload.service mysterybotreload.service || exit
