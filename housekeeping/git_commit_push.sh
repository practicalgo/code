#!/bin/bash
set -ex

detect-os-executables -c "$1"
git add -A 
git commit -m "Update"
git push
