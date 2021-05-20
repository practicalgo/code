#!/usr/bin/python

import os
import subprocess
import sys

failed = 0

if len(sys.argv) != 2:
    sys.exit('Must specify a directory path to vet')

for root, dirs, files in os.walk(sys.argv[1]):
    src_dir = None
    if root.startswith('./.git'):
        continue
    for f in files:
        if '.go' in f:
            src_dir = root
            break
    if not src_dir or "parked" in src_dir or "solutions" in src_dir or "service" in src_dir:
        print("Ignoring: {0}".format(src_dir))
        continue
    try:        
        subprocess.check_output(["go", "vet"], cwd=src_dir,
                stderr=subprocess.PIPE)
    except subprocess.CalledProcessError as e:
        print('Vet failure: {0}'.format(src_dir))
        failed = 1
sys.exit(failed)