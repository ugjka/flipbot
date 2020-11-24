#!/bin/bash
git checkout master || exit
git checkout bootybot || exit
git merge master --no-edit || exit
git checkout mysterybot || exit
git merge master --no-edit || exit
git checkout master || exit
git push --all || exit