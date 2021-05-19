#!/bin/bash

cat $1 | expand -t 8 | pbcopy
