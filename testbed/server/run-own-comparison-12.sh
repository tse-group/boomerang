#! /bin/bash -ve

N_nodes=100
N_scenarios=100
EXID="12_N${N_nodes}"
SCID="02_nodes${N_nodes}_txs500_paths25_edges4.605170to6.907755"

mkdir -p results/${EXID} || true


param_paths=25

for((i_scenario=0;i_scenario<N_scenarios;i_scenario++))
do
    cd ../parse_graph
    ./load_cfg.sh ${SCID} ${i_scenario} ${N_nodes}
    cd ../server
    ./load_cfg.sh

    for param_redundancy in 0 2 4 6 8 10 15 20 25 30 40 50 75 100 125 150
    do
        for algo in "retry-02-amp-2" "redundancy-02-amp-2"
        do
            params="${param_paths}-${param_redundancy}"
            echo "Running combination: ${algo} ${param_paths} ${param_redundancy} ${params}"
            ./run-own3.sh ${N_nodes} ${algo} ${param_paths} ${param_redundancy}
            cat err_*.log | grep "RESULTS:" > results/${EXID}/${SCID}_${algo}_${params}_${i_scenario}.data
            cat err_*.log | grep "CHECKPOINT" >> results/${EXID}/${SCID}_${algo}_${params}_${i_scenario}.data
        done

        algo="redundantretry-02-amp-2"

        for param_leeway in 10
        do
            params="${param_paths}-${param_redundancy}-${param_leeway}"
            echo "Running combination: ${algo} ${param_paths} ${param_redundancy} ${param_leeway} ${params}"
            ./run-own3.sh ${N_nodes} ${algo} ${param_paths} ${param_redundancy} ${param_leeway}
            cat err_*.log | grep "RESULTS:" > results/${EXID}/${SCID}_${algo}_${params}_${i_scenario}.data
            cat err_*.log | grep "CHECKPOINT" >> results/${EXID}/${SCID}_${algo}_${params}_${i_scenario}.data
        done

    done
done
