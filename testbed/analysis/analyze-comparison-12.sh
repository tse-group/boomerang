#! /bin/bash -ve

N_nodes=100
SCENARIOS=(0 1 2 11 12 13 14 18 19 21)
NUMEXID="12"
NUMEXID_SCRIPTS="12"
SCID="02_nodes${N_nodes}_txs500_paths25_edges4.605170to6.907755"


EXID="${NUMEXID}_N${N_nodes}"
mkdir -p results/${EXID} || true


param_paths=25

DATAFILE="results/${EXID}/${SCID}_paths${param_paths}.data"
echo -n "u algo " > ${DATAFILE}
./compute-rvs.py >> ${DATAFILE}

for param_redundancy in 0 2 4 6 8 10 15 20 25 30 40 50 75 100 125 150; do
    for algo in "retry-02-amp-2" "redundancy-02-amp-2"; do
        params="${param_paths}-${param_redundancy}"
        echo "Reading combination: ${algo} ${params}"
        
        echo > results/${EXID}/${SCID}_paths${param_paths}_${algo}_${params}.data
        for i_scenario in ${SCENARIOS[@]}; do
            ./aggregate-nodes.py ../server/results/${EXID}/${SCID}_${algo}_${params}_${i_scenario}.data >> results/${EXID}/${SCID}_paths${param_paths}_${algo}_${params}.data
        done

        echo -n "${param_redundancy} ${algo} " >> ${DATAFILE}
        ./compute-rvs.py results/${EXID}/${SCID}_paths${param_paths}_${algo}_${params}.data >> ${DATAFILE}
    done

    algo="redundantretry-02-amp-2"

    for param_leeway in 10; do
        params="${param_paths}-${param_redundancy}-${param_leeway}"
        echo "Reading combination: ${algo} ${params}"

        echo > results/${EXID}/${SCID}_paths${param_paths}_${algo}_${params}.data
        for i_scenario in ${SCENARIOS[@]}; do
            ./aggregate-nodes.py ../server/results/${EXID}/${SCID}_${algo}_${params}_${i_scenario}.data >> results/${EXID}/${SCID}_paths${param_paths}_${algo}_${params}.data
        done

        echo -n "${param_redundancy} ${algo}-${param_leeway} " >> ${DATAFILE}
        ./compute-rvs.py results/${EXID}/${SCID}_paths${param_paths}_${algo}_${params}.data >> ${DATAFILE}
    done
done


gnuplot -e "DATAFILENAME='${DATAFILE}'" analyze-comparison-${NUMEXID_SCRIPTS}-01-volume_success.gnuplot
gnuplot -e "DATAFILENAME='${DATAFILE}'" analyze-comparison-${NUMEXID_SCRIPTS}-02-throughput_success.gnuplot
gnuplot -e "DATAFILENAME='${DATAFILE}'" analyze-comparison-${NUMEXID_SCRIPTS}-03-ttc_for_successful_tx.gnuplot
gnuplot -e "DATAFILENAME='${DATAFILE}'" analyze-comparison-${NUMEXID_SCRIPTS}-04-volume_for_successful_tx.gnuplot
gnuplot -e "DATAFILENAME='${DATAFILE}'" analyze-comparison-${NUMEXID_SCRIPTS}-05-runtime.gnuplot
gnuplot -e "DATAFILENAME='${DATAFILE}'" analyze-comparison-${NUMEXID_SCRIPTS}-06-count_success.gnuplot
