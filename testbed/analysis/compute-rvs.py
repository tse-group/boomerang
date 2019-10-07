#! /usr/bin/env python3

import sys
import os
import math

from collections import OrderedDict

from result import Result


f_rv01_volume_success = lambda r: r.volume_success
f_rv02_throughput_success = lambda r: r.volume_success / r.runtime
f_rv03_ttc_for_successful_tx = lambda r: r.ttc_success / r.count_success
f_rv04_volume_for_successful_tx = lambda r: r.volume_success / r.count_success
f_rv05_runtime = lambda r: r.runtime
f_rv06_count_success = lambda r: r.count_success

rvs = OrderedDict()
rvs['SAMPLE'] = lambda r: 1
rvs['volume_success'] = f_rv01_volume_success
rvs['throughput_success'] = f_rv02_throughput_success
rvs['ttc_for_successful_tx'] = f_rv03_ttc_for_successful_tx
rvs['volume_for_successful_tx'] = f_rv04_volume_for_successful_tx
rvs['runtime'] = f_rv05_runtime
rvs['count_success'] = f_rv06_count_success


def my_sum(lst):
    return sum(lst)

def my_min(lst):
    return min(lst)

def my_max(lst):
    return max(lst)

def my_mean(lst):
    n = len(lst)
    return 1.0 / n * sum(lst)

def my_var(lst):
    n = len(lst)
    mu = my_mean(lst)
    if n < 2:
        return 0.0
    else:
        return 1.0 / (n-1) * sum([ (l - mu)**2 for l in lst ])

def my_std(lst):
    return math.sqrt(my_var(lst))

funcs = OrderedDict()
funcs['sum'] = my_sum
funcs['min'] = my_min
funcs['max'] = my_max
funcs['mean'] = my_mean
funcs['var'] = my_var
funcs['std'] = my_std


if len(sys.argv) == 1:
    for k, f in rvs.items():
        print(*[ '%s-%s'% (k, k2) for (k2, f2) in funcs.items() ], end=' ')

    print()

elif len(sys.argv) == 2:
    results = Result.load_from_file(sys.argv[1])

    for k, f in rvs.items():
        vals = list(map(f, results))
        print(*[ f2(vals) for (k2, f2) in funcs.items() ], end=' ')

    print()

else:
    assert len(sys.argv) in [1, 2]

