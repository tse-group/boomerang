# Boomerang: Redundancy Improves Latency and Throughput in Payment-Channel Networks

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


## Usage

...


## References

* Boomerang: Redundancy Improves Latency and Throughput in Payment-Channel Networks<br/>
  Vivek Bagaria, Joachim Neu, David Tse<br/>
  24th International Conference on Financial Cryptography and Data Security (FC'20), February 2020<br/>
  [https://arxiv.org/abs/1910.01834](https://arxiv.org/abs/1910.01834) - [Bitcoin Optech Newsletter #86](https://bitcoinops.org/en/newsletters/2020/02/26/#boomerang-redundancy-improves-latency-and-throughput-in-payment-channel-networks)

* Flash: Efficient Dynamic Routing for Offchain Networks<br/>
  Peng Wang, Hong Xu, Xin Jin, Tao Wang<br/>
  ACM CoNEXT 2019<br/>
  [https://arxiv.org/abs/1902.05260](https://arxiv.org/abs/1902.05260) - [Code on GitHub](https://github.com/NetX-lab/Offchain-routing-traces-and-code)

