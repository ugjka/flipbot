#!/bin/bash

#This prints ten random words.

shuf -n 10 words --random-source=/dev/urandom | tr '\n' ' '
