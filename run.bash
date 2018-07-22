#!/bin/bash

while true; do
    go generate
    exit_code=$?

    if [ $exit_code != 0 ]
    then
        echo "[BASH SCRIPT] go generate failed with exit code $exit_code"
        break
    fi

    go build -o bin/bot && bin/bot
    exit_code=$?

    if [ $exit_code == 0 ] || [ $exit_code == 1 ]
    then
        echo "[BASH SCRIPT] Terminate exit code recieved: $exit_code"
        break
    elif [ $exit_code == 3 ]
    then
        seconds=1
        echo "[BASH SCRIPT] Restarting without delay"
    else
        seconds=15
    fi

    # git reset --hard
    git pull origin master

    for ((second=$seconds; second > 0; second--))
    do
        echo -ne "[BASH SCRIPT] Restarting in $second seconds..\r"
        sleep 1
    done
done
