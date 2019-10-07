import networkx as nx
import csv
import random
import math
import sys
import os


def generate_scenario(tx_amounts, num_nodes, num_txs_per_node, num_paths, cap_log_min, cap_log_max, random_seed):
	random.seed(random_seed)


	# RANDOM NETWORK CONSTRUCTION

	GG = nx.watts_strogatz_graph(num_nodes, 8, 0.8, random_seed)   # num_nodes nodes, connected to nearest 8 neighbors in ring topology, then 0.8 probability of rewiring, random seed
	assert nx.is_connected(GG)
	G = nx.DiGraph()

	for e in GG.edges():
		cap_01_log = cap_log_min + random.random() * (cap_log_max - cap_log_min)
		cap_10_log = cap_log_min + random.random() * (cap_log_max - cap_log_min)
		cap_01 = math.exp(cap_01_log)
		cap_10 = math.exp(cap_10_log)
		G.add_edge(e[0], e[1], capacity=cap_01)
		G.add_edge(e[1], e[0], capacity=cap_10)


	# GENERATE PATHS

	paths = {}

	for src in G.nodes():
		for dst in G.nodes():
			if src == dst:
				continue

			path_set = list(nx.edge_disjoint_paths(G, src, dst, cutoff=num_paths))
			paths[(src, dst)] = path_set


	# GENERATE PAYMENTS

	payments = []

	for k in range(num_nodes * num_txs_per_node):
		amt = random.choice(tx_amounts)
		src = random.randrange(num_nodes)
		while True:
			dst = random.randrange(num_nodes)
			if dst != src:
				break

		payments.append((src, dst, amt))


	return (G, paths, payments, GG)



def main():
	# PARAMETERS

	# number of nodes in the network
	num_nodes = 100
	# number of payments to run through the network
	num_txs_per_node = 500 #1000000
	# how many paths to compute max. (for "mice" payments)
	num_paths = 25
	# lower and upper limit of log-uniform distribution for capacity
	cap_log_min = math.log(100)
	cap_log_max = math.log(1000)
	# number of scenarios to simulate
	num_scenarios = 100


	# BUILD PATH FOR STORAGE

	dir_prefix = '02_nodes%d_txs%d_paths%d_edges%fto%f'% (num_nodes, num_txs_per_node, num_paths, cap_log_min, cap_log_max)
	try:
		os.mkdir(dir_prefix)
	except FileExistsError:
		pass


	# LOAD TX AMOUNTS FROM RIPPLE DATASET

	tx_amounts = []

	with open('data/ripple_val.csv', 'r') as f: 
		csv_reader = csv.reader(f, delimiter=',')
		for row in csv_reader:
			if float(row[2]) > 0:
				tx_amounts.append(float(row[2]))


	# GENERATE SCENARIOS

	for i in range(num_scenarios):
		print('Scenario', i)

		(G, paths, payments, GG) = generate_scenario(tx_amounts, num_nodes, num_txs_per_node, num_paths, cap_log_min, cap_log_max, i)

		# write out graph
		f = open(dir_prefix + "/%d_graph.txt"% (i,), 'w+')
		for e in G.edges():
			f.write("%d,%d,%f\n" % (e[0], e[1], G.get_edge_data(*e)['capacity']))
		f.close()

		# write out paths
		f = open(dir_prefix + "/%d_paths.txt"% (i,), 'w+')
		for k, v in paths.items():
			for p in v:
				f.write("%d,%d,%s\n"% (k[0], k[1], ','.join([ str(p_) for p_ in p ])))
		f.close()

		# write out payments
		f = open(dir_prefix + "/%d_payments.txt"% (i,), 'w+')
		for (src, dst, amt) in payments:
			f.write("%d,%d,%f\n"% (src, dst, amt))
		f.close()



if __name__ == "__main__":
	main()
