package main


type TfReqRoutingController interface {
    NextTxId()                  string
    DrawRandomPath()            Path
    GetTxReqs()                 map[string]*TxReq
    GetTxReq(string)            (*TxReq, bool)
    AddTxReq(*TxReq)
    UpdateStats()
}

func FilterTxReqsUnknownOutcome(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.HasUnknownOutcome() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsPositiveOutcome(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.HasPositiveOutcome() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsNegativeOutcome(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.HasNegativeOutcome() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsHasPhase1Sent(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.HasPhase1Sent() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsHasPhase1RetRecvd(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.HasPhase1RetRecvd() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsHasPhase2Sent(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.HasPhase2Sent() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsHasPhase2RetRecvd(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.HasPhase2RetRecvd() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsIsFinal(ctrl TfReqRoutingController) map[string]*TxReq {
    return FilterTxReqsHasPhase2RetRecvd(ctrl)
}

func FilterTxReqsIsFinalNeg(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.IsFinalNeg() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsIsFinalPos(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.IsFinalPos() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsIsPhase2Ready(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.IsPhase2Ready() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsIsPhase2ReadyNeg(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.IsPhase2ReadyNeg() {
            res[txid] = txreq
        }
    }
    return res
}

func FilterTxReqsIsPhase2ReadyPos(ctrl TfReqRoutingController) map[string]*TxReq {
    res := make(map[string]*TxReq)
    for txid, txreq := range ctrl.GetTxReqs() {
        if txreq.IsPhase2ReadyPos() {
            res[txid] = txreq
        }
    }
    return res
}

func NewTxReq(ctrl TfReqRoutingController, amt float, p Path) *TxReq {
    var txreq TxReq
    txreq.TxId = ctrl.NextTxId()
    txreq.Path = p
    txreq.Amount = amt
    txreq.State = TXREQ_STATE_NEW
    ctrl.AddTxReq(&txreq)
    return &txreq
}

func HandleRetPacket(ctrl TfReqRoutingController, ret P2pMsg, is_phase2_blocking bool) {
    myassert(ret.Type == P2PMSG_TYPE_P1_RESERVE_RET || ret.Type == P2PMSG_TYPE_P2_NEG_ROLLBACK_RET || ret.Type == P2PMSG_TYPE_P2_POS_EXECUTE_RET, "Unexpected non-RET response from network!")

    if ret.Type == P2PMSG_TYPE_P1_RESERVE_RET {

        txreq, is_found := ctrl.GetTxReq(ret.TxId)
        myassert(is_found, "RESERVE_RET from current TfReq, but TxId not found among TxReqs?!")
        if ret.ParamReserveResultSuccessful {
            txreq.PositiveReserveRetReceived(ret.ParamReserveResultPath)
        } else {
            txreq.NegativeReserveRetReceived(ret.ParamReserveResultPath)
        }

    } else if ret.Type == P2PMSG_TYPE_P2_NEG_ROLLBACK_RET {

        txreq, is_found := ctrl.GetTxReq(ret.TxId)
        if is_found {
            txreq.RollbackRetReceived()
        } else if is_phase2_blocking {
            myassert(is_found, "NEG_ROLLBACK_RET from current TfReq, but TxId not found among TxReqs?!")
        }

    } else if ret.Type == P2PMSG_TYPE_P2_POS_EXECUTE_RET {

        txreq, is_found := ctrl.GetTxReq(ret.TxId)
        if is_found {
            txreq.ExecuteRetReceived()
        } else if is_phase2_blocking {
            myassert(is_found, "POS_EXECUTE_RET from current TfReq, but TxId not found among TxReqs?!")
        }

    }

}

