package main

import (
	"net"
	"sync"
	"time"
	"log"
)


// P2P MESSAGING PROTOCOL

// p2p messaging protocol

const P2PMSG_TYPE_P1_ABORT = 1
const P2PMSG_TYPE_P1_RESERVE = 2
const P2PMSG_TYPE_P1_RESERVE_RET = 3
const P2PMSG_TYPE_P2_POS_EXECUTE = 4
const P2PMSG_TYPE_P2_POS_EXECUTE_RET = 5
const P2PMSG_TYPE_P2_NEG_ROLLBACK = 6
const P2PMSG_TYPE_P2_NEG_ROLLBACK_RET = 7

type P2pMsg struct {
	Type					int				`json:"type"`
	Path					Path			`json:"path"`
	TxId					string			`json:"tid"`
	ParamReserveRequestAmount		float			`json:"p_res_req_amt"`
	ParamReserveResultSuccessful	bool			`json:"p_res_res_ok"`
	ParamReserveResultPath			Path			`json:"p_res_res_p"`
}


func is_fast(m P2pMsg) bool {
	if m.Type == P2PMSG_TYPE_P1_RESERVE_RET ||
		m.Type == P2PMSG_TYPE_P2_POS_EXECUTE ||
		m.Type == P2PMSG_TYPE_P2_NEG_ROLLBACK ||
		m.Type == P2PMSG_TYPE_P1_ABORT {
		return true
	} else if m.Type == P2PMSG_TYPE_P1_RESERVE ||
		m.Type == P2PMSG_TYPE_P2_POS_EXECUTE_RET ||
		m.Type == P2PMSG_TYPE_P2_NEG_ROLLBACK_RET {
		return false
	} else {
		log.Panicln("Unknown message type:", m)
		return true
	}
}

func is_slow(m P2pMsg) bool {
	return !is_fast(m)
}


// paths through the network

type Path struct {
	Ids			[]int		`json:"nids"`
}

func src(p Path) int {
	return myfirsti(p.Ids)
}

func dst(p Path) int {
	return mylasti(p.Ids)
}

func rev(p Path) Path {
	var newp Path
	newp.Ids = myreversei_outofplace(p.Ids)
	return newp
}

func next_hop(p Path, nid int) (int, bool, bool, int) {
	if pos, found := myfindi(p.Ids, nid); found {
		if pos == len(p.Ids) - 1 {
			// found, but no next hop
			return -1, true, false, pos
		} else {
			// found, next hop available
			return p.Ids[pos+1], true, true, pos
		}
	} else {
		// not found
		return -1, false, false, -1
	}
}

func prev_hop(p Path, nid int) (int, bool, bool, int) {
	if pos, found := myfindi(p.Ids, nid); found {
		if pos == 0 {
			// found, but no prev hop
			return -1, true, false, pos
		} else {
			// found, prev hop available
			return p.Ids[pos-1], true, true, pos
		}
	} else {
		// not found
		return -1, false, false, -1
	}
}

func until_here(p Path, nid int) Path {
	pos, found := myfindi(p.Ids, nid)
	myassert(found, "Node ID not found in path!")

	var p2 Path
	p2.Ids = p.Ids[0:pos+1]
	return p2
}


// HANDLING OF MESSAGES

// internal data of the client

type Session struct {
	Peer				NodeInfo
	QueueTransmit		chan P2pMsg
	Connection			*net.TCPConn

	Capacity			float
	CapacityLock		sync.Mutex

	// LogAmountRx			float
	// LogAmountTx			float
}


// planned transactions

type ReservedTx struct {
	TxId				string
	Amount				float
	Session				*Session
	TimeStarted			time.Time
}

type ReservedTxStore struct {
	Txs					map[string]ReservedTx
}

func has_txid(store *ReservedTxStore, txid string) bool {
	if _, found := store.Txs[txid]; found {
		return true
	} else {
		return false
	}
}

func txid_seen(store_out *ReservedTxStore, store_in *ReservedTxStore, msg P2pMsg) bool {
	txid := msg.TxId

	_, found_in := store_in.Txs[txid]
	_, found_out := store_out.Txs[txid]

	if found_in || found_out {
		return true
	} else {
		return false
	}
}

func abort_requested(aborted_txs map[string]bool, msg P2pMsg) bool {
	txid := msg.TxId
	if _, found := aborted_txs[txid]; found {
		return true
	} else {
		return false
	}
}

func reserve_inbound_tx(store *ReservedTxStore, session *Session, msg P2pMsg) {
	txid := msg.TxId
	myassert(!has_txid(store, txid), "TxId already known!")

	var res ReservedTx
	res.TxId = txid
	res.Amount = msg.ParamReserveRequestAmount
	res.Session = session
	res.TimeStarted = time.Now()
	store.Txs[txid] = res
}

func reserve_outbound_tx(store *ReservedTxStore, session *Session, msg P2pMsg) {
	txid := msg.TxId
	myassert(!has_txid(store, txid), "TxId already known!")

	session.CapacityLock.Lock()
	var res ReservedTx
	res.TxId = txid
	res.Amount = msg.ParamReserveRequestAmount
	res.Session = session
	res.TimeStarted = time.Now()
	myassert(session.Capacity >= res.Amount, "Insufficient capacity!")
	session.Capacity -= res.Amount
	store.Txs[txid] = res
	session.CapacityLock.Unlock()
}

func attempt_reserve_outbound_tx(store *ReservedTxStore, session *Session, msg P2pMsg) bool {
	var success bool

	txid := msg.TxId
	myassert(!has_txid(store, txid), "TxId already known!")

	session.CapacityLock.Lock()
	if session.Capacity >= msg.ParamReserveRequestAmount {
		// enough capacity available
		var res ReservedTx
		res.TxId = txid
		res.Amount = msg.ParamReserveRequestAmount
		res.Session = session
		res.TimeStarted = time.Now()
		session.Capacity -= res.Amount
		store.Txs[txid] = res

		success = true
	} else {
		// not enough capacity available
		// do nothing

		success = false		
	}
	session.CapacityLock.Unlock()

	return success
}

func rollback_inbound_tx(store *ReservedTxStore, msg P2pMsg) float {
	txid := msg.TxId
	myassert(has_txid(store, txid), "TxId not known!")

	res, _ := store.Txs[txid]
	delete(store.Txs, txid)

	return float(time.Now().Sub(res.TimeStarted).Seconds()) * res.Amount
}

func rollback_outbound_tx(store *ReservedTxStore, msg P2pMsg) float {
	txid := msg.TxId
	myassert(has_txid(store, txid), "TxId not known!")

	res, _ := store.Txs[txid]
	res.Session.CapacityLock.Lock()
	res.Session.Capacity += res.Amount
	delete(store.Txs, txid)
	res.Session.CapacityLock.Unlock()

	return float(time.Now().Sub(res.TimeStarted).Seconds()) * res.Amount
}

func execute_inbound_tx(store *ReservedTxStore, msg P2pMsg) float {
	return rollback_outbound_tx(store, msg)
}

func execute_outbound_tx(store *ReservedTxStore, msg P2pMsg) float {
	return rollback_inbound_tx(store, msg)
}


// MAIN SUB-ROUTINES THAT RUN IN PARALLEL

func session_by_next_hop(sessions_by_nid map[int]*Session, p Path, self_nid int) *Session {
	nid, is_found, is_forward, _ := next_hop(p, self_nid)
	myassert(is_found && is_forward, "session_by_next_hop failed!")
	return sessions_by_nid[nid]
}

func session_by_prev_hop(sessions_by_nid map[int]*Session, p Path, self_nid int) *Session {
	nid, is_found, is_forward, _ := prev_hop(p, self_nid)
	myassert(is_found && is_forward, "session_by_prev_hop failed!")
	return sessions_by_nid[nid]
}

// forward message (strictly outbound!)

func forward_msg_outbound(self NodeInfo, msg P2pMsg, sessions_outbound_by_nid map[int]*Session) {
	session := session_by_next_hop(sessions_outbound_by_nid, msg.Path, self.Id)
	log.Println("Forwarding", msg, "via", session)
	session.QueueTransmit <- msg
}

// receive the next P1 RET, discarding other RETs on the way

func receive_P1_RET(queue_local chan P2pMsg) (P2pMsg, int) {
	num_discarded_P2_RETs := 0
	for {
		ret := <- queue_local
		if ret.Type == P2PMSG_TYPE_P1_RESERVE_RET {
			return ret, num_discarded_P2_RETs
		} else if ret.Type == P2PMSG_TYPE_P2_NEG_ROLLBACK_RET || ret.Type == P2PMSG_TYPE_P2_POS_EXECUTE_RET {
			num_discarded_P2_RETs++
		} else {
			log.Panicln("Unexpected msg received:", ret)
		}
	}
}

// debug

func _debug_sessions_capacities(sessions []Session) ([]int, []float) {
	lst_nid := make([]int, len(sessions))
	lst_cap := make([]float, len(sessions))
	for i, s := range sessions {
		lst_nid[i] = s.Peer.Id
		lst_cap[i] = s.Capacity
	}
	return lst_nid, lst_cap
}

// handle all the routing

func handle_routing(self NodeInfo, queue_routing chan P2pMsg, queue_local chan P2pMsg, sessions_outbound []Session) {

	// store pending and aborted txs

	var pending_txs_in ReservedTxStore
	var pending_txs_out ReservedTxStore
	pending_txs_in.Txs = make(map[string]ReservedTx)
	pending_txs_out.Txs = make(map[string]ReservedTx)

	var aborted_txs map[string]bool
	aborted_txs = make(map[string]bool)


	// index outbound sessions by NodeID of the respective peer (to simplify forwarding)

	sessions_outbound_by_nid := make(map[int]*Session)
	for i, session := range sessions_outbound {
		sessions_outbound_by_nid[session.Peer.Id] = &(sessions_outbound[i])
	}


	// process the msgs queued up for routing

	for msg := range queue_routing {

		log.Println("Routing msg:", msg)
		// spew.Dump(msg)
		// log.Println("Pre routing:", "In:", pending_txs_in, "Out:", pending_txs_out)
		// log.Println("Channel capacities:", _debug_sessions_capacities(sessions_outbound))

		if msg.Type == P2PMSG_TYPE_P1_ABORT {

			// was the txid seen already?
			if txid_seen(&pending_txs_out, &pending_txs_in, msg) {
				// yes, ie, the RESERVE has passed this node already
				// -> forward, unless this is the destination
				if dst(msg.Path) == self.Id {
					// the RESERVE_RET is already on the way back, nothing can be done
				} else {
					// forward
					forward_msg_outbound(self, msg, sessions_outbound_by_nid)
				}
			} else {
				// no, ie, it either returned early because of insufficient
				// funds, or it will still arrive here
				// -> note ABORT request, do not forward
				aborted_txs[msg.TxId] = true
			}

		} else if msg.Type == P2PMSG_TYPE_P1_RESERVE {

			if !(src(msg.Path) == self.Id) {

				// whatever happens, an inbound transaction needs to be registered
				session_prev := session_by_prev_hop(sessions_outbound_by_nid, msg.Path, self.Id)
				reserve_inbound_tx(&pending_txs_in, session_prev, msg)

			}

			// abort was requested, fail propagation here
			if abort_requested(aborted_txs, msg) {

				delete(aborted_txs, msg.TxId)

				msg.Type = P2PMSG_TYPE_P1_RESERVE_RET
				msg.Path = rev(msg.Path)
				msg.ParamReserveResultSuccessful = false
				msg.ParamReserveResultPath = until_here(rev(msg.Path), self.Id)

				queue_routing <- msg

			// abort was not requested, normal operation
			} else {

				if dst(msg.Path) == self.Id {

					msg.Type = P2PMSG_TYPE_P1_RESERVE_RET
					msg.Path = rev(msg.Path)
					msg.ParamReserveResultSuccessful = true
					msg.ParamReserveResultPath = until_here(rev(msg.Path), self.Id)

					queue_routing <- msg

				} else {

					session_next := session_by_next_hop(sessions_outbound_by_nid, msg.Path, self.Id)
					successful := attempt_reserve_outbound_tx(&pending_txs_out, session_next, msg)

					if successful {

						forward_msg_outbound(self, msg, sessions_outbound_by_nid)

					} else {

						msg.Type = P2PMSG_TYPE_P1_RESERVE_RET
						msg.Path = rev(msg.Path)
						msg.ParamReserveResultSuccessful = false
						msg.ParamReserveResultPath = until_here(rev(msg.Path), self.Id)
						
						queue_routing <- msg

					}

				}

			}

		} else if msg.Type == P2PMSG_TYPE_P1_RESERVE_RET {

			if dst(msg.Path) == self.Id {

				// this node is dst of the msg

				queue_local <- msg

			} else {

				// this node is src or intermediary in this msg

				forward_msg_outbound(self, msg, sessions_outbound_by_nid)

			}

		} else if msg.Type == P2PMSG_TYPE_P2_NEG_ROLLBACK {

			if dst(msg.Path) == self.Id {

				// this node is dst of the msg

				msg.Type = P2PMSG_TYPE_P2_NEG_ROLLBACK_RET
				msg.Path = rev(msg.Path)
				queue_routing <- msg

			} else {

				// this node is src or intermediary in this msg
				
				forward_msg_outbound(self, msg, sessions_outbound_by_nid)

			}

		} else if msg.Type == P2PMSG_TYPE_P2_NEG_ROLLBACK_RET {

			if src(msg.Path) == self.Id && dst(msg.Path) == self.Id {

				myassert(len(msg.Path.Ids) == 1, "This node is both src and dst of the path, it better be a path of length 1")

				queue_local <- msg

			} else if src(msg.Path) == self.Id {

				// this node is src of the msg

				rollback_inbound_tx(&pending_txs_in, msg)

				forward_msg_outbound(self, msg, sessions_outbound_by_nid)

			} else if dst(msg.Path) == self.Id {

				// this node is dst of the msg

				rollback_outbound_tx(&pending_txs_out, msg)

				queue_local <- msg

			} else {

				// this node is intermediary in this msg
				
				rollback_outbound_tx(&pending_txs_out, msg)
				rollback_inbound_tx(&pending_txs_in, msg)

				forward_msg_outbound(self, msg, sessions_outbound_by_nid)

			}

		} else if msg.Type == P2PMSG_TYPE_P2_POS_EXECUTE {

			if dst(msg.Path) == self.Id {

				// this node is dst of the msg

				msg.Type = P2PMSG_TYPE_P2_POS_EXECUTE_RET
				msg.Path = rev(msg.Path)
				queue_routing <- msg

			} else {

				// this node is src or intermediary in this msg
				
				forward_msg_outbound(self, msg, sessions_outbound_by_nid)

			}

		} else if msg.Type == P2PMSG_TYPE_P2_POS_EXECUTE_RET {

			if src(msg.Path) == self.Id && dst(msg.Path) == self.Id {

				myassert(len(msg.Path.Ids) == 1, "This node is both src and dst of the path, it better be a path of length 1")

				queue_local <- msg

			} else if src(msg.Path) == self.Id {

				// this node is src of the msg

				execute_inbound_tx(&pending_txs_in, msg)

				forward_msg_outbound(self, msg, sessions_outbound_by_nid)

			} else if dst(msg.Path) == self.Id {

				// this node is dst of the msg

				execute_outbound_tx(&pending_txs_out, msg)

				queue_local <- msg


			} else {

				// this node is intermediary in this msg

				execute_outbound_tx(&pending_txs_out, msg)
				execute_inbound_tx(&pending_txs_in, msg)

				forward_msg_outbound(self, msg, sessions_outbound_by_nid)

			}

		} else {

			log.Panicln("Unknown P2pMsg type:", msg.Type)

		}


		// log.Println("Post routing:", "In:", pending_txs_in, "Out:", pending_txs_out)
		// log.Println("Channel capacities:", _debug_sessions_capacities(sessions_outbound))
	}
}

