#! /usr/bin/env python3

import sys
import os
from result import Result

assert len(sys.argv) == 2

results = Result.load_from_file(sys.argv[1])
total = Result.sum(results)

print(total.render())
