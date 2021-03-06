# Boomerang: Redundancy Improves Latency and Throughput in Payment-Channel Networks


**Disclaimer:**
This codebase comes with absolutely no warranty.
Use strictly at your own risk!


## Abstract

In multi-path routing schemes for payment-channel networks,
Alice transfers funds to Bob by splitting them into partial payments
and routing them along multiple paths.
Undisclosed channel balances and mismatched transaction fees
cause delays and failures on some payment paths.
For atomic transfer schemes,
these straggling paths stall the whole transfer.
We show that the latency of transfers reduces when redundant payment paths are added.
This frees up liquidity in payment channels
and hence increases the throughput of the network.
We devise *Boomerang*, a generic
technique to be used on top of multi-path routing schemes to construct
redundant payment paths free of counterparty risk.
In our experiments, applying Boomerang to a baseline routing scheme
leads to 40% latency reduction and 2x throughput increase.
We build on ideas from publicly verifiable secret sharing, such that
Alice learns a secret of Bob
iff Bob overdraws funds from the redundant paths.
Funds are forwarded using Boomerang contracts,
which allow Alice to revert the transfer
iff she has learned Bob's secret.
We implement the Boomerang contract in Bitcoin Script.



## Paper

* Boomerang: Redundancy Improves Latency and Throughput in Payment-Channel Networks<br/>
  Vivek Bagaria, Joachim Neu, David Tse<br/>
  24th International Conference on Financial Cryptography and Data Security (FC'20), February 2020<br/>
  [https://arxiv.org/abs/1910.01834](https://arxiv.org/abs/1910.01834)

For an excellent summary of the paper, check out the [Bitcoin Optech Newsletter #86](https://bitcoinops.org/en/newsletters/2020/02/26/#boomerang-redundancy-improves-latency-and-throughput-in-payment-channel-networks).



## Getting Started

The following sections provide basic familiarity with this codebase.
For questions, please contact [Joachim Neu](https://www.jneu.net/).


### Dependencies

The Boomerang experiments were developed on Linux and have
the following major dependencies (we used the version in brackets):

* Python (v3.8) + dependencies listed in `requirements.txt`
* Go (v1.13.5)
* Gnuplot (v5.2)
* Bash (v5.0.11)


### Anatomy of the Repository

The anatomy of the repository is inspired by that of the [Flash routing protocol](https://github.com/NetX-lab/Offchain-routing-traces-and-code).

* `mockclient3`:
  The payment-channel network node software simulator/prototype, written in Go.
* `testbed`:
  The environment to run simulation experiments in.
  * `analysis`:
    Gnuplot and Bash scripts to turn simulation results into plots.
  * `gen_trace`:
    Python code `test-own-02-uniform-srcdst-ripple-amount-multiple.py` to setup simulation scenarios (randomly sample network topologies, initial channel balances, and transfer demands from certain probability distributions).
  * `parse_graph`:
    Go code to prepare a simulation scenario generated by `gen_trace` for execution by `mockclient3`.
  * `server`:
    Where it all comes together: The Bash script `run-own-comparison-12.sh` orchestrates the simulations (`mockclient3`) of certain scenarios (`gen_trace`) for various parameters.

**Where to start?**
The best entry point to this repository is probably the Bash script `testbed/server/run-own-comparison-12.sh`.
From there you can find your way to all the parts of the codebase.


### Download & Setup

Download the repository and setup the Python virtual environment (`venv`):

```
git clone https://github.com/tse-group/boomerang.git
cd boomerang/
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
```


### Analysis

View the plots in `testbed/analysis/results/12_N100`
and recreate them using Gnuplot with:

```
cd testbed/analysis/
./analyze-comparison-12.sh
cd ../../
```


### Simulation Scenarios

Uncompress the simulation scenarios:

```
cd testbed/gen_trace/
unzip 02_nodes100_txs500_paths25_edges4.605170to6.907755.zip
cd ../../
```

To create new simulation scenarios,
uncompress the underlying [Ripple trace dataset](https://crysp.uwaterloo.ca/software/speedymurmurs/),
adjust the parameters in `test-own-02-uniform-srcdst-ripple-amount-multiple.py`,
and create the scenarios, using:

```
cd testbed/gen_trace/
cd data/
bzip2 -d ripple_val.csv.bz2
cd ../
python test-own-02-uniform-srcdst-ripple-amount-multiple.py
cd ../../
```


### Run Experiments

First, build the Go code in `mockclient3/` and `testbed/parse_graph/`:

```
cd mockclient3/
go build
cd ../
cd testbed/parse_graph/
go build
cd ../../
```

Then, run the simulations:

```
cd testbed/server/
./run-own-comparison-12.sh
cd ../../
```

To adjust the parameters or which scenarios are being simulated,
modify `run-own-comparison-12.sh`.



## References

* Boomerang: Redundancy Improves Latency and Throughput in Payment-Channel Networks<br/>
  Vivek Bagaria, Joachim Neu, David Tse<br/>
  24th International Conference on Financial Cryptography and Data Security (FC'20), February 2020<br/>
  [https://arxiv.org/abs/1910.01834](https://arxiv.org/abs/1910.01834) - [Bitcoin Optech Newsletter #86](https://bitcoinops.org/en/newsletters/2020/02/26/#boomerang-redundancy-improves-latency-and-throughput-in-payment-channel-networks)

* Flash: Efficient Dynamic Routing for Offchain Networks<br/>
  Peng Wang, Hong Xu, Xin Jin, Tao Wang<br/>
  ACM CoNEXT 2019<br/>
  [https://arxiv.org/abs/1902.05260](https://arxiv.org/abs/1902.05260) - [Code on GitHub](https://github.com/NetX-lab/Offchain-routing-traces-and-code)

* Settling Payments Fast and Private: Efficient Decentralized Routing for Path-Based Transactions<br/>
  Stefanie Roos, Pedro Moreno-Sanchez, Aniket Kate, Ian Goldberg<br/>
  NDSS 2018<br/>
  [Paper](https://crysp.uwaterloo.ca/software/speedymurmurs/ndss.pdf) - [Code](https://crysp.uwaterloo.ca/software/speedymurmurs/)
