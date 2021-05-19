#!/usr/bin/python

import os
import subprocess
import sys

failed = 0

if len(sys.argv) != 2:
    sys.exit('Must specify a directory path to lint')

for root, dirs, files in os.walk(sys.argv[1]):
    src_dir = None
    if root.startswith('./.git'):
        continue
    for f in files:
        if f.endswith('.go'):
            src_dir = root
            break
    if not src_dir:
        continue
    print('Linting: {0}\n-------'.format(src_dir))
    try:
        print(subprocess.check_output(["golint"], cwd=src_dir,
                stderr=subprocess.PIPE).decode("utf-8"))
    except subprocess.CalledProcessError as e:
        print('Lint failure: {0}'.format(src_dir))
        failed = 1

sys.exit(failed)
