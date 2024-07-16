#!/usr/bin/bash

MINWAIT=2
MAXWAIT=300
DELAY=$((MINWAIT+RANDOM % (MAXWAIT-MINWAIT)))

sleep $DELAY

echo $DELAY
./bcli auth login -f ~/.bcli/.credential.txt; ./bcli punch
