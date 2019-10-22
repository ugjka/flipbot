#!/bin/bash

#This prints ten random words.

echo "$(shuf -n 10 words --random-source=/dev/urandom | tr '\n' ' ')"
