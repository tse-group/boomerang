package main

import (
    "time"
)


const TXREQ_STATE_NEW = 1
const TXREQ_STATE_RESERVE_SENT = 2
const TXREQ_STATE_POS_RESERVE_RET_RECVD = 3
const TXREQ_STATE_NEG_RESERVE_RET_RECVD = 4
const TXREQ_STATE_EXECUTE_SENT = 5
const TXREQ_STATE_ROLLBACK_SENT = 6
const TXREQ_STATE_EXECUTE_RET_RECVD = 7
const TXREQ_STATE_ROLLBACK_RET_RECVD = 8

type TxReq struct {
    TxId                        string
    Path                        Path
    Amount                      float
    State                       int
    ReserveResultPath           Path
    TimeReserveSent             time.Time
    TimeReserveRetRecvd         time.Time
    TimePhase2Sent              time.Time
    TimePhase2RetRecvd          time.Time
}

func (txreq *TxReq) SendReserve(queue_routing chan P2pMsg) {
    myassert(txreq.State == TXREQ_STATE_NEW, "Requires state TXREQ_STATE_NEW for TxReq")

    var msg P2pMsg
    msg.Type = P2PMSG_TYPE_P1_RESERVE
    msg.Path = txreq.Path
    msg.TxId = txreq.TxId
    msg.ParamReserveRequestAmount = txreq.Amount
    queue_routing <- msg

    txreq.TimeReserveSent = time.Now()
    txreq.State = TXREQ_STATE_RESERVE_SENT
}

func (txreq *TxReq) PositiveReserveRetReceived(p Path) {
    myassert(txreq.State == TXREQ_STATE_RESERVE_SENT, "Requires state TXREQ_STATE_RESERVE_SENT for TxReq")
    txreq.ReserveResultPath = p
    txreq.TimeReserveRetRecvd = time.Now()
    txreq.State = TXREQ_STATE_POS_RESERVE_RET_RECVD
}

func (txreq *TxReq) NegativeReserveRetReceived(p Path) {
    myassert(txreq.State == TXREQ_STATE_RESERVE_SENT, "Requires state TXREQ_STATE_RESERVE_SENT for TxReq")
    txreq.ReserveResultPath = p
    txreq.TimeReserveRetRecvd = time.Now()
    txreq.State = TXREQ_STATE_NEG_RESERVE_RET_RECVD
}

func (txreq *TxReq) SendExecute(queue_routing chan P2pMsg) {
    myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD, "Requires state TXREQ_STATE_POS_RESERVE_RET_RECVD for TxReq")

    var msg P2pMsg
    msg.Type = P2PMSG_TYPE_P2_POS_EXECUTE
    msg.Path = txreq.ReserveResultPath
    msg.TxId = txreq.TxId
    queue_routing <- msg

    txreq.TimePhase2Sent = time.Now()
    txreq.State = TXREQ_STATE_EXECUTE_SENT
}

func (txreq *TxReq) SendRollback(queue_routing chan P2pMsg) {
    myassert(txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD || txreq.State == TXREQ_STATE_NEG_RESERVE_RET_RECVD, "Requires state TXREQ_STATE_POS_RESERVE_RET_RECVD or TXREQ_STATE_NEG_RESERVE_RET_RECVD for TxReq")

    var msg P2pMsg
    msg.Type = P2PMSG_TYPE_P2_NEG_ROLLBACK
    msg.Path = txreq.ReserveResultPath
    msg.TxId = txreq.TxId
    queue_routing <- msg

    txreq.TimePhase2Sent = time.Now()
    txreq.State = TXREQ_STATE_ROLLBACK_SENT
}

func (txreq *TxReq) ExecuteRetReceived() {
    myassert(txreq.State == TXREQ_STATE_EXECUTE_SENT, "Requires state TXREQ_STATE_EXECUTE_SENT for TxReq")
    txreq.TimePhase2RetRecvd = time.Now()
    txreq.State = TXREQ_STATE_EXECUTE_RET_RECVD
}

func (txreq *TxReq) RollbackRetReceived() {
    myassert(txreq.State == TXREQ_STATE_ROLLBACK_SENT, "Requires state TXREQ_STATE_ROLLBACK_SENT for TxReq")
    txreq.TimePhase2RetRecvd = time.Now()
    txreq.State = TXREQ_STATE_ROLLBACK_RET_RECVD
}

func (txreq *TxReq) SendAbort(queue_routing chan P2pMsg) {
    var msg P2pMsg
    msg.Type = P2PMSG_TYPE_P1_ABORT
    msg.Path = txreq.Path
    msg.TxId = txreq.TxId
    queue_routing <- msg
}

func (txreq *TxReq) HasUnknownOutcome() bool {
    return txreq.State == TXREQ_STATE_NEW || txreq.State == TXREQ_STATE_RESERVE_SENT
}

func (txreq *TxReq) HasPositiveOutcome() bool {
    return txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD || txreq.State == TXREQ_STATE_EXECUTE_SENT || txreq.State == TXREQ_STATE_EXECUTE_RET_RECVD
}

func (txreq *TxReq) HasNegativeOutcome() bool {
    return txreq.State == TXREQ_STATE_NEG_RESERVE_RET_RECVD || txreq.State == TXREQ_STATE_ROLLBACK_SENT || txreq.State == TXREQ_STATE_ROLLBACK_RET_RECVD
}

func (txreq *TxReq) HasPhase1Sent() bool {
    return !(txreq.State == TXREQ_STATE_NEW)
}

func (txreq *TxReq) HasPhase1RetRecvd() bool {
    return !(txreq.State == TXREQ_STATE_NEW || txreq.State == TXREQ_STATE_RESERVE_SENT)
}

func (txreq *TxReq) HasPhase2Sent() bool {
    return txreq.State == TXREQ_STATE_EXECUTE_SENT || txreq.State == TXREQ_STATE_ROLLBACK_SENT || txreq.State == TXREQ_STATE_EXECUTE_RET_RECVD || txreq.State == TXREQ_STATE_ROLLBACK_RET_RECVD
}

func (txreq *TxReq) HasPhase2RetRecvd() bool {
    return txreq.State == TXREQ_STATE_EXECUTE_RET_RECVD || txreq.State == TXREQ_STATE_ROLLBACK_RET_RECVD
}

func (txreq *TxReq) IsFinal() bool {
    return txreq.HasPhase2RetRecvd()
}

func (txreq *TxReq) IsFinalNeg() bool {
    return txreq.State == TXREQ_STATE_ROLLBACK_RET_RECVD
}

func (txreq *TxReq) IsFinalPos() bool {
    return txreq.State == TXREQ_STATE_EXECUTE_RET_RECVD
}

func (txreq *TxReq) IsPhase2Ready() bool {
    return txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD || txreq.State == TXREQ_STATE_NEG_RESERVE_RET_RECVD
}

func (txreq *TxReq) IsPhase2ReadyNeg() bool {
    return txreq.State == TXREQ_STATE_NEG_RESERVE_RET_RECVD
}

func (txreq *TxReq) IsPhase2ReadyPos() bool {
    return txreq.State == TXREQ_STATE_POS_RESERVE_RET_RECVD
}

func (txreq *TxReq) LiquidityDelayProduct() float {
    var reftime_stop time.Time

    if txreq.State == TXREQ_STATE_NEW {
        return 0.0
    }

    if !txreq.HasPhase2Sent() {
        reftime_stop = time.Now()
    } else {
        reftime_stop = txreq.TimePhase2Sent
    }

    return txreq.Amount * float(reftime_stop.Sub(txreq.TimeReserveSent).Seconds())
}
