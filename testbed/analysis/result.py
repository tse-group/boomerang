#! /usr/bin/env python3

import math

class Result(object):
    def __init__(self,
        volume,
        volume_success,
        volume_fail,
        count,
        count_success,
        count_fail,
        ttc,
        ttc_success,
        ttc_fail,
        runtime,
        runtime_success,
        runtime_fail,
        ldp,
        ldp_success,
        ldp_fail,
        ):

        self.volume = volume
        self.volume_success = volume_success
        self.volume_fail = volume_fail
        self.count = count
        self.count_success = count_success
        self.count_fail = count_fail
        self.ttc = ttc
        self.ttc_success = ttc_success
        self.ttc_fail = ttc_fail
        self.runtime = runtime
        self.runtime_success = runtime_success
        self.runtime_fail = runtime_fail
        self.ldp = ldp
        self.ldp_success = ldp_success
        self.ldp_fail = ldp_fail

    @classmethod
    def from_line(cls, line):
        values = line.split(' ')
        assert len(values) == 16
        assert values[0] == 'RESULTS:'
        return cls(
            float(values[1]),
            float(values[2]),
            float(values[3]),
            int(values[4]),
            int(values[5]),
            int(values[6]),
            float(values[7]),
            float(values[8]),
            float(values[9]),
            float(values[10]),
            float(values[11]),
            float(values[12]),
            float(values[13]),
            float(values[14]),
            float(values[15]),
            )

    @classmethod
    def sum(cls, results):
        # n = len(results)
        n = 1.0

        res_avg = cls(0.0, 0.0, 0.0, 0, 0, 0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
        for r in results:
            res_avg.volume += r.volume * 1./n
            res_avg.volume_success += r.volume_success * 1./n
            res_avg.volume_fail += r.volume_fail * 1./n
            res_avg.count += r.count * 1./n
            res_avg.count_success += r.count_success * 1./n
            res_avg.count_fail += r.count_fail * 1./n
            res_avg.ttc += r.ttc * 1./n
            res_avg.ttc_success += r.ttc_success * 1./n
            res_avg.ttc_fail += r.ttc_fail * 1./n
            res_avg.runtime += r.runtime * 1./n
            res_avg.runtime_success += r.runtime_success * 1./n
            res_avg.runtime_fail += r.runtime_fail * 1./n
            res_avg.ldp += r.ldp * 1./n
            res_avg.ldp_success += r.ldp_success * 1./n
            res_avg.ldp_fail += r.ldp_fail * 1./n

        return res_avg

    @classmethod
    def avgvar(cls, results):
        n = len(results)
        # n = 1.0

        res_avg = cls(0.0, 0.0, 0.0, 0, 0, 0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
        for r in results:
            res_avg.volume += r.volume * 1./n
            res_avg.volume_success += r.volume_success * 1./n
            res_avg.volume_fail += r.volume_fail * 1./n
            res_avg.count += r.count * 1./n
            res_avg.count_success += r.count_success * 1./n
            res_avg.count_fail += r.count_fail * 1./n
            res_avg.ttc += r.ttc * 1./n
            res_avg.ttc_success += r.ttc_success * 1./n
            res_avg.ttc_fail += r.ttc_fail * 1./n
            res_avg.runtime += r.runtime * 1./n
            res_avg.runtime_success += r.runtime_success * 1./n
            res_avg.runtime_fail += r.runtime_fail * 1./n
            res_avg.ldp += r.ldp * 1./n
            res_avg.ldp_success += r.ldp_success * 1./n
            res_avg.ldp_fail += r.ldp_fail * 1./n

        res_var = cls(0.0, 0.0, 0.0, 0, 0, 0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
        for r in results:
            res_var.volume += (r.volume - res_avg.volume)**2 * 1./n
            res_var.volume_success += (r.volume_success - res_avg.volume_success)**2 * 1./n
            res_var.volume_fail += (r.volume_fail - res_avg.volume_fail)**2 * 1./n
            res_var.count += (r.count - res_avg.count)**2 * 1./n
            res_var.count_success += (r.count_success - res_avg.count_success)**2 * 1./n
            res_var.count_fail += (r.count_fail - res_avg.count_fail)**2 * 1./n
            res_var.ttc += (r.ttc - res_avg.ttc)**2 * 1./n
            res_var.ttc_success += (r.ttc_success - res_avg.ttc_success)**2 * 1./n
            res_var.ttc_fail += (r.ttc_fail - res_avg.ttc_fail)**2 * 1./n
            res_var.runtime += (r.runtime - res_avg.runtime)**2 * 1./n
            res_var.runtime_success += (r.runtime_success - res_avg.runtime_success)**2 * 1./n
            res_var.runtime_fail += (r.runtime_fail - res_avg.runtime_fail)**2 * 1./n
            res_var.ldp += (r.ldp - res_avg.ldp)**2 * 1./n
            res_var.ldp_success += (r.ldp_success - res_avg.ldp_success)**2 * 1./n
            res_var.ldp_fail += (r.ldp_fail - res_avg.ldp_fail)**2 * 1./n

        return (res_avg, res_var)

    @classmethod
    def avgstd(cls, results):
        (res_avg, res_var) = cls.avgvar(results)

        res_std = cls(0.0, 0.0, 0.0, 0, 0, 0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 0.0)
        res_std.volume = math.sqrt(res_var.volume)
        res_std.volume_success = math.sqrt(res_var.volume_success)
        res_std.volume_fail = math.sqrt(res_var.volume_fail)
        res_std.count = math.sqrt(res_var.count)
        res_std.count_success = math.sqrt(res_var.count_success)
        res_std.count_fail = math.sqrt(res_var.count_fail)
        res_std.ttc = math.sqrt(res_var.ttc)
        res_std.ttc_success = math.sqrt(res_var.ttc_success)
        res_std.ttc_fail = math.sqrt(res_var.ttc_fail)
        res_std.runtime = math.sqrt(res_var.runtime)
        res_std.runtime_success = math.sqrt(res_var.runtime_success)
        res_std.runtime_fail = math.sqrt(res_var.runtime_fail)
        res_std.ldp = math.sqrt(res_var.ldp)
        res_std.ldp_success = math.sqrt(res_var.ldp_success)
        res_std.ldp_fail = math.sqrt(res_var.ldp_fail)

        return (res_avg, res_std)

    @classmethod
    def load_from_file(cls, filename):
        results = []
        fp = open(filename, 'r')
        for l in fp.readlines():
            l.strip()
            if l.startswith('RESULTS: '):
                r = Result.from_line(l)
                results.append(r)
        return results

    def render(self, prefix='RESULTS'):
        return '%s: %f %f %f %d %d %d %f %f %f %f %f %f %f %f %f'% (
            prefix,
            self.volume,
            self.volume_success,
            self.volume_fail,
            self.count,
            self.count_success,
            self.count_fail,
            self.ttc,
            self.ttc_success,
            self.ttc_fail,
            self.runtime,
            self.runtime_success,
            self.runtime_fail,
            self.ldp,
            self.ldp_success,
            self.ldp_fail,
        )

    def __repr__(self):
        return self.render()
