package main

import (
	"time"
	"fmt"
	"os"
)


type Statistics struct {
	amount				float
	amount_success		float
	amount_fail			float
	count				int
	count_success		int
	count_fail			int
	ttc					time.Duration
	ttc_success			time.Duration
	ttc_fail			time.Duration
	runtime				time.Duration
	runtime_success		time.Duration
	runtime_fail		time.Duration
	ldp					float
	ldp_success			float
	ldp_fail			float
}

type StatisticsIncrement struct {
	amount				float
	count				int
	ttc					time.Duration
	runtime				time.Duration
	ldp					float
}

func dump_stats(fp *os.File, s *Statistics, caption string) {
	fmt.Fprintln(fp, "\n" + caption + ":",
		s.amount,
		s.amount_success,
		s.amount_fail,
		s.count,
		s.count_success,
		s.count_fail,
		s.ttc.Seconds(),
		s.ttc_success.Seconds(),
		s.ttc_fail.Seconds(),
		s.runtime.Seconds(),
		s.runtime_success.Seconds(),
		s.runtime_fail.Seconds(),
		s.ldp,
		s.ldp_success,
		s.ldp_fail,
	)
}

func _log_stats_all(s *Statistics, d *StatisticsIncrement) {
	s.amount += d.amount
	s.count += d.count
	s.ttc += d.ttc
	s.runtime += d.runtime
	s.ldp += d.ldp
}

func _log_stats_success(s *Statistics, d *StatisticsIncrement) {
	s.amount_success += d.amount
	s.count_success += d.count
	s.ttc_success += d.ttc
	s.runtime_success += d.runtime
	s.ldp_success += d.ldp
}

func _log_stats_fail(s *Statistics, d *StatisticsIncrement) {
	s.amount_fail += d.amount
	s.count_fail += d.count
	s.ttc_fail += d.ttc
	s.runtime_fail += d.runtime
	s.ldp_fail += d.ldp
}

func log_stats_success(s *Statistics, d *StatisticsIncrement) {
	_log_stats_all(s, d)
	_log_stats_success(s, d)
}

func log_stats_fail(s *Statistics, d *StatisticsIncrement) {
	_log_stats_all(s, d)
	_log_stats_fail(s, d)
}

func log_stats(s *Statistics, d *StatisticsIncrement, successful bool) {
	if successful {
		log_stats_success(s, d)
	} else {
		log_stats_fail(s, d)
	}
}
