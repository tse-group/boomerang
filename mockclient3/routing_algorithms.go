package main

import (
    "time"
    "log"
    "math/rand"
    "os"
    "fmt"
)



type RetryCtrl struct {
	// provided from outside, overall setting
	Self						NodeInfo
	TfReq						TfReq
	TfReqNo						int
	TfReqSuccessful				bool
	CandidatePaths				[]Path

	// state
	TxReqs						map[string]*TxReq
	NumRemainingRetries			int

	// statistics
	TimeStart					time.Time
	TimeStopTTC					time.Time
	TimeStopRuntime				time.Time
	Stats						StatisticsIncrement

	// internals, completely invisible from outside (hopefully)
	NextTxIdNo					int
}

func (ctrl *RetryCtrl) NextTxId() string {
	ctrl.NextTxIdNo += 1
	return fmt.Sprintf("node-%d-txreq-%d-part-%d", ctrl.Self.Id, ctrl.TfReqNo, ctrl.NextTxIdNo)
}

func (ctrl *RetryCtrl) DrawRandomPath() Path {
	return ctrl.CandidatePaths[rand.Intn(len(ctrl.CandidatePaths))]
}

func (ctrl *RetryCtrl) GetTxReqs() map[string]*TxReq {
	return ctrl.TxReqs
}

func (ctrl *RetryCtrl) GetTxReq(txid string) (*TxReq, bool) {
	txreq, is_found := ctrl.GetTxReqs()[txid]
	return txreq, is_found
}

func (ctrl *RetryCtrl) AddTxReq(txreq *TxReq) {
	txid := txreq.TxId
	ctrl.TxReqs[txid] = txreq
}

func (ctrl *RetryCtrl) UpdateStats() {
	ctrl.Stats.amount = ctrl.TfReq.Amount
	ctrl.Stats.count = 1
	ctrl.Stats.ttc = ctrl.TimeStopTTC.Sub(ctrl.TimeStart)
	ctrl.Stats.runtime = ctrl.TimeStopRuntime.Sub(ctrl.TimeStart)
	ctrl.Stats.ldp = 0.0
	for _, txreq := range ctrl.GetTxReqs() {
		ctrl.Stats.ldp += txreq.LiquidityDelayProduct()
	}
}



func process_tfreqs_retry_02_amp_2(
		self NodeInfo,
		tfreqs chan TfReq,
		paths map[int][]Path,
		peers []NodeInfo,
		sessions_outbound []Session,
		queue_routing chan P2pMsg,
		queue_local chan P2pMsg,
		algo_param_paths int,
		algo_param_retries int,
	) {

	var stats Statistics


	// attempt to satisfy transaction requests

	i := -1
	for tfreq := range tfreqs {
		i++

		log.Println("TfReq", i, tfreq)
		// spew.Dump(tfreq)

		var ctrl RetryCtrl
		ctrl.Self = self
		ctrl.TfReq = tfreq
		ctrl.TfReqNo = i
		ctrl.TfReqSuccessful = false
		ctrl.TimeStart = time.Now()
		ctrl.TxReqs = make(map[string]*TxReq)
		ctrl.NumRemainingRetries = algo_param_retries


		// randomize candidate paths

		candidate_paths := paths[tfreq.Dst]
		n_paths := len(candidate_paths)
		rand.Shuffle(n_paths, func(i int, j int) { candidate_paths[i], candidate_paths[j] = candidate_paths[j], candidate_paths[i] })
		ctrl.CandidatePaths = candidate_paths


		// send payment attempts

		path_amount := tfreq.Amount / float(algo_param_paths)

		for j := 0; j < algo_param_paths; j++ {
			p := ctrl.DrawRandomPath()
			txreq := NewTxReq(&ctrl, path_amount, p)
			txreq.SendReserve(queue_routing)
		}


		// wait for results (as long as there is a chance for this to be successful)

		for len(FilterTxReqsPositiveOutcome(&ctrl)) < algo_param_paths && len(FilterTxReqsUnknownOutcome(&ctrl)) > 0 && len(FilterTxReqsPositiveOutcome(&ctrl)) + len(FilterTxReqsUnknownOutcome(&ctrl)) + ctrl.NumRemainingRetries >= algo_param_paths {
			// receive a response
			ret := <- queue_local
			HandleRetPacket(&ctrl, ret, false)
			// HandleRetPacket(&ctrl, ret, true)
			for _, txreq := range FilterTxReqsIsPhase2ReadyNeg(&ctrl) {
				txreq.SendRollback(queue_routing)
			}
			for len(ctrl.GetTxReqs()) - len(FilterTxReqsIsFinalNeg(&ctrl)) < algo_param_paths && ctrl.NumRemainingRetries > 0 {
				p := ctrl.DrawRandomPath()
				txreq := NewTxReq(&ctrl, path_amount, p)
				txreq.SendReserve(queue_routing)
				ctrl.NumRemainingRetries--
			}
		}


		// abort outstanding attempts

		for _, txreq := range FilterTxReqsUnknownOutcome(&ctrl) {
			txreq.SendAbort(queue_routing)
		}


		// check whether we have enough positive responses

		if len(FilterTxReqsPositiveOutcome(&ctrl)) == algo_param_paths {
			// yes, enough positive responses, EXECUTE all
			for _, txreq := range FilterTxReqsPositiveOutcome(&ctrl) {
				myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD, "State unexpectedly is not TXREQ_STATE_POS_RESERVE_RET_RECVD")
				txreq.SendExecute(queue_routing)
			}
			ctrl.TfReqSuccessful = true
		} else {
			// no, not enough positive responses, ROLLBACK all
			for _, txreq := range FilterTxReqsPositiveOutcome(&ctrl) {
				myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD, "State unexpectedly is not TXREQ_STATE_POS_RESERVE_RET_RECVD")
				txreq.SendRollback(queue_routing)
			}
			ctrl.TfReqSuccessful = false
		}


		// tx completed

		ctrl.TimeStopTTC = time.Now()


		// cancel remaining attempts (irrespective of outcome)

		for len(FilterTxReqsUnknownOutcome(&ctrl)) > 0 {
		// for len(ctrl.GetTxReqs()) - len(FilterTxReqsIsFinal(&ctrl)) > 0 {
			// receive a response
			ret := <- queue_local
			HandleRetPacket(&ctrl, ret, false)
			// HandleRetPacket(&ctrl, ret, true)
			for _, txreq := range FilterTxReqsIsPhase2Ready(&ctrl) {
				txreq.SendRollback(queue_routing)
			}
		}


		// how did we do? perform some statistics bookkeeping ...

		ctrl.TimeStopRuntime = time.Now()
		ctrl.UpdateStats()

		log_stats(&stats, &ctrl.Stats, ctrl.TfReqSuccessful)


		// output checkpoint of state

		if i % 100 == 0 {
			dump_stats(os.Stderr, &stats, fmt.Sprintf("CHECKPOINT-%d", i))
		}
	}


	// dump statistics

	dump_stats(os.Stderr, &stats, "RESULTS")
	fmt.Fprintln(os.Stderr, "\nFINISHED!")
}



type RedundancyCtrl struct {
	// provided from outside, overall setting
	Self						NodeInfo
	TfReq						TfReq
	TfReqNo						int
	TfReqSuccessful				bool
	CandidatePaths				[]Path

	// state
	TxReqs						map[string]*TxReq

	// statistics
	TimeStart					time.Time
	TimeStopTTC					time.Time
	TimeStopRuntime				time.Time
	Stats						StatisticsIncrement

	// internals, completely invisible from outside (hopefully)
	NextTxIdNo					int
}

func (ctrl *RedundancyCtrl) NextTxId() string {
	ctrl.NextTxIdNo += 1
	return fmt.Sprintf("node-%d-txreq-%d-part-%d", ctrl.Self.Id, ctrl.TfReqNo, ctrl.NextTxIdNo)
}

func (ctrl *RedundancyCtrl) DrawRandomPath() Path {
	return ctrl.CandidatePaths[rand.Intn(len(ctrl.CandidatePaths))]
}

func (ctrl *RedundancyCtrl) GetTxReqs() map[string]*TxReq {
	return ctrl.TxReqs
}

func (ctrl *RedundancyCtrl) GetTxReq(txid string) (*TxReq, bool) {
	txreq, is_found := ctrl.GetTxReqs()[txid]
	return txreq, is_found
}

func (ctrl *RedundancyCtrl) AddTxReq(txreq *TxReq) {
	txid := txreq.TxId
	ctrl.TxReqs[txid] = txreq
}

func (ctrl *RedundancyCtrl) UpdateStats() {
	ctrl.Stats.amount = ctrl.TfReq.Amount
	ctrl.Stats.count = 1
	ctrl.Stats.ttc = ctrl.TimeStopTTC.Sub(ctrl.TimeStart)
	ctrl.Stats.runtime = ctrl.TimeStopRuntime.Sub(ctrl.TimeStart)
	ctrl.Stats.ldp = 0.0
	for _, txreq := range ctrl.GetTxReqs() {
		ctrl.Stats.ldp += txreq.LiquidityDelayProduct()
	}
}



func process_tfreqs_redundancy_02_amp_2(
		self NodeInfo,
		tfreqs chan TfReq,
		paths map[int][]Path,
		peers []NodeInfo,
		sessions_outbound []Session,
		queue_routing chan P2pMsg,
		queue_local chan P2pMsg,
		algo_param_paths int,
		algo_param_redundancy int,
	) {

	var stats Statistics


	// attempt to satisfy transaction requests

	i := -1
	for tfreq := range tfreqs {
		i++

		log.Println("TfReq", i, tfreq)
		// spew.Dump(tfreq)

		var ctrl RedundancyCtrl
		ctrl.Self = self
		ctrl.TfReq = tfreq
		ctrl.TfReqNo = i
		ctrl.TfReqSuccessful = false
		ctrl.TimeStart = time.Now()
		ctrl.TxReqs = make(map[string]*TxReq)


		// randomize candidate paths

		candidate_paths := paths[tfreq.Dst]
		n_paths := len(candidate_paths)
		rand.Shuffle(n_paths, func(i int, j int) { candidate_paths[i], candidate_paths[j] = candidate_paths[j], candidate_paths[i] })
		ctrl.CandidatePaths = candidate_paths


		// send payment attempts

		path_amount := tfreq.Amount / float(algo_param_paths)

		for j := 0; j < algo_param_paths + algo_param_redundancy; j++ {
			p := ctrl.DrawRandomPath()
			txreq := NewTxReq(&ctrl, path_amount, p)
			txreq.SendReserve(queue_routing)
		}


		// wait for results (as long as there is a chance for this to be successful)

		for len(FilterTxReqsPositiveOutcome(&ctrl)) < algo_param_paths && len(FilterTxReqsUnknownOutcome(&ctrl)) > 0 && len(FilterTxReqsPositiveOutcome(&ctrl)) + len(FilterTxReqsUnknownOutcome(&ctrl)) >= algo_param_paths {
			// receive a response
			ret := <- queue_local
			HandleRetPacket(&ctrl, ret, false)
			for _, txreq := range FilterTxReqsIsPhase2ReadyNeg(&ctrl) {
				txreq.SendRollback(queue_routing)
			}
		}


		// abort outstanding attempts

		for _, txreq := range FilterTxReqsUnknownOutcome(&ctrl) {
			txreq.SendAbort(queue_routing)
		}


		// check whether we have enough positive responses

		if len(FilterTxReqsPositiveOutcome(&ctrl)) == algo_param_paths {
			// yes, enough positive responses, EXECUTE all
			for _, txreq := range FilterTxReqsPositiveOutcome(&ctrl) {
				myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD, "State unexpectedly is not TXREQ_STATE_POS_RESERVE_RET_RECVD")
				txreq.SendExecute(queue_routing)
			}
			ctrl.TfReqSuccessful = true
		} else {
			// no, not enough positive responses, ROLLBACK all
			for _, txreq := range FilterTxReqsPositiveOutcome(&ctrl) {
				myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD, "State unexpectedly is not TXREQ_STATE_POS_RESERVE_RET_RECVD")
				txreq.SendRollback(queue_routing)
			}
			ctrl.TfReqSuccessful = false
		}


		// tx completed

		ctrl.TimeStopTTC = time.Now()


		// cancel remaining attempts (irrespective of outcome)

		for len(FilterTxReqsUnknownOutcome(&ctrl)) > 0 {
			// receive a response
			ret := <- queue_local
			HandleRetPacket(&ctrl, ret, false)
			for _, txreq := range FilterTxReqsIsPhase2Ready(&ctrl) {
				txreq.SendRollback(queue_routing)
			}
		}


		// how did we do? perform some statistics bookkeeping ...

		ctrl.TimeStopRuntime = time.Now()
		ctrl.UpdateStats()

		log_stats(&stats, &ctrl.Stats, ctrl.TfReqSuccessful)


		// output checkpoint of state

		if i % 100 == 0 {
			dump_stats(os.Stderr, &stats, fmt.Sprintf("CHECKPOINT-%d", i))
		}
	}


	// dump statistics

	dump_stats(os.Stderr, &stats, "RESULTS")
	fmt.Fprintln(os.Stderr, "\nFINISHED!")
}



type RedundantRetryCtrl struct {
	// provided from outside, overall setting
	Self						NodeInfo
	TfReq						TfReq
	TfReqNo						int
	TfReqSuccessful				bool
	CandidatePaths				[]Path

	// state
	TxReqs						map[string]*TxReq
	NumRemainingRetries			int

	// statistics
	TimeStart					time.Time
	TimeStopTTC					time.Time
	TimeStopRuntime				time.Time
	Stats						StatisticsIncrement

	// internals, completely invisible from outside (hopefully)
	NextTxIdNo					int
}

func (ctrl *RedundantRetryCtrl) NextTxId() string {
	ctrl.NextTxIdNo += 1
	return fmt.Sprintf("node-%d-txreq-%d-part-%d", ctrl.Self.Id, ctrl.TfReqNo, ctrl.NextTxIdNo)
}

func (ctrl *RedundantRetryCtrl) DrawRandomPath() Path {
	return ctrl.CandidatePaths[rand.Intn(len(ctrl.CandidatePaths))]
}

func (ctrl *RedundantRetryCtrl) GetTxReqs() map[string]*TxReq {
	return ctrl.TxReqs
}

func (ctrl *RedundantRetryCtrl) GetTxReq(txid string) (*TxReq, bool) {
	txreq, is_found := ctrl.GetTxReqs()[txid]
	return txreq, is_found
}

func (ctrl *RedundantRetryCtrl) AddTxReq(txreq *TxReq) {
	txid := txreq.TxId
	ctrl.TxReqs[txid] = txreq
}

func (ctrl *RedundantRetryCtrl) UpdateStats() {
	ctrl.Stats.amount = ctrl.TfReq.Amount
	ctrl.Stats.count = 1
	ctrl.Stats.ttc = ctrl.TimeStopTTC.Sub(ctrl.TimeStart)
	ctrl.Stats.runtime = ctrl.TimeStopRuntime.Sub(ctrl.TimeStart)
	ctrl.Stats.ldp = 0.0
	for _, txreq := range ctrl.GetTxReqs() {
		ctrl.Stats.ldp += txreq.LiquidityDelayProduct()
	}
}



func process_tfreqs_redundantretry_02_amp_2(
		self NodeInfo,
		tfreqs chan TfReq,
		paths map[int][]Path,
		peers []NodeInfo,
		sessions_outbound []Session,
		queue_routing chan P2pMsg,
		queue_local chan P2pMsg,
		algo_param_paths int,
		algo_param_retries int,
		algo_param_redundancy int,
	) {

	var stats Statistics


	// attempt to satisfy transaction requests

	i := -1
	for tfreq := range tfreqs {
		i++

		log.Println("TfReq", i, tfreq)
		// spew.Dump(tfreq)

		var ctrl RedundantRetryCtrl
		ctrl.Self = self
		ctrl.TfReq = tfreq
		ctrl.TfReqNo = i
		ctrl.TfReqSuccessful = false
		ctrl.TimeStart = time.Now()
		ctrl.TxReqs = make(map[string]*TxReq)
		ctrl.NumRemainingRetries = algo_param_retries


		// randomize candidate paths

		candidate_paths := paths[tfreq.Dst]
		n_paths := len(candidate_paths)
		rand.Shuffle(n_paths, func(i int, j int) { candidate_paths[i], candidate_paths[j] = candidate_paths[j], candidate_paths[i] })
		ctrl.CandidatePaths = candidate_paths


		// send payment attempts

		path_amount := tfreq.Amount / float(algo_param_paths)

		for j := 0; j < algo_param_paths + mymini(algo_param_retries, algo_param_redundancy); j++ {
			p := ctrl.DrawRandomPath()
			txreq := NewTxReq(&ctrl, path_amount, p)
			txreq.SendReserve(queue_routing)
		}
		ctrl.NumRemainingRetries -= mymini(algo_param_retries, algo_param_redundancy)


		// wait for results (as long as there is a chance for this to be successful)

		for len(FilterTxReqsPositiveOutcome(&ctrl)) < algo_param_paths && len(FilterTxReqsUnknownOutcome(&ctrl)) > 0 && len(FilterTxReqsPositiveOutcome(&ctrl)) + len(FilterTxReqsUnknownOutcome(&ctrl)) + ctrl.NumRemainingRetries >= algo_param_paths {
			// receive a response
			ret := <- queue_local
			HandleRetPacket(&ctrl, ret, false)
			// HandleRetPacket(&ctrl, ret, true)
			for _, txreq := range FilterTxReqsIsPhase2ReadyNeg(&ctrl) {
				txreq.SendRollback(queue_routing)
			}
			for len(ctrl.GetTxReqs()) - len(FilterTxReqsIsFinalNeg(&ctrl)) < algo_param_paths + algo_param_redundancy && ctrl.NumRemainingRetries > 0 {
				p := ctrl.DrawRandomPath()
				txreq := NewTxReq(&ctrl, path_amount, p)
				txreq.SendReserve(queue_routing)
				ctrl.NumRemainingRetries--
			}
		}


		// abort outstanding attempts

		for _, txreq := range FilterTxReqsUnknownOutcome(&ctrl) {
			txreq.SendAbort(queue_routing)
		}


		// check whether we have enough positive responses

		if len(FilterTxReqsPositiveOutcome(&ctrl)) == algo_param_paths {
			// yes, enough positive responses, EXECUTE all
			for _, txreq := range FilterTxReqsPositiveOutcome(&ctrl) {
				myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD, "State unexpectedly is not TXREQ_STATE_POS_RESERVE_RET_RECVD")
				txreq.SendExecute(queue_routing)
			}
			ctrl.TfReqSuccessful = true
		} else {
			// no, not enough positive responses, ROLLBACK all
			for _, txreq := range FilterTxReqsPositiveOutcome(&ctrl) {
				myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD, "State unexpectedly is not TXREQ_STATE_POS_RESERVE_RET_RECVD")
				txreq.SendRollback(queue_routing)
			}
			ctrl.TfReqSuccessful = false
		}


		// tx completed

		ctrl.TimeStopTTC = time.Now()


		// cancel remaining attempts (irrespective of outcome)

		for len(FilterTxReqsUnknownOutcome(&ctrl)) > 0 {
		// for len(ctrl.GetTxReqs()) - len(FilterTxReqsIsFinal(&ctrl)) > 0 {
			// receive a response
			ret := <- queue_local
			HandleRetPacket(&ctrl, ret, false)
			// HandleRetPacket(&ctrl, ret, true)
			for _, txreq := range FilterTxReqsIsPhase2Ready(&ctrl) {
				txreq.SendRollback(queue_routing)
			}
		}


		// how did we do? perform some statistics bookkeeping ...

		ctrl.TimeStopRuntime = time.Now()
		ctrl.UpdateStats()

		log_stats(&stats, &ctrl.Stats, ctrl.TfReqSuccessful)


		// output checkpoint of state

		if i % 100 == 0 {
			dump_stats(os.Stderr, &stats, fmt.Sprintf("CHECKPOINT-%d", i))
		}
	}


	// dump statistics

	dump_stats(os.Stderr, &stats, "RESULTS")
	fmt.Fprintln(os.Stderr, "\nFINISHED!")
}
