#!/usr/bin/python

import os
import subprocess
import sys

failed = 0
if len(sys.argv) != 2:
    sys.exit('Must specify a directory path to test')

for root, dirs, files in os.walk(sys.argv[1]):
    src_dir = None
    if root.startswith('./.git'):
        continue
    for f in files:
        if '.go' in f:
            src_dir = root
            break
    if not src_dir:
        continue
    try:
        subprocess.check_output(["go", "test", "-v"], cwd=src_dir,
                stderr=subprocess.PIPE)
    except subprocess.CalledProcessError as e:
        print('Test failure: {0}'.format(src_dir))
        failed = 1

sys.exit(failed)
