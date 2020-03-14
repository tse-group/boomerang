#! /bin/bash

if [ "$#" -lt 2 ]; then
    echo "Usage: ./run-own.sh <number of nodes> <algorithm> [<optional parameters for algorithm>]"
    exit -1
fi


args=("$@")

N_node=${args[0]}
algo=${args[1]}
algo_params=${args[@]:2}


json_suffix=".json"
neig_prefix="n"
tran_prefix="tr"
tran_suffix=".txt"
g_file="graph.txt"
log_suffix=".log"
log_prefix=""
path_prefix="pa"
path_suffix=".txt"

rm ./pid ./*.log || true

echo "-> Starting clients"
for((i=1;i<=N_node;i++))
do
    cmd="nohup ../../mockclient3/mockclient3 ${i}${json_suffix} ${neig_prefix}${i}${json_suffix} ${g_file} ${tran_prefix}${i}${tran_suffix} ${path_prefix}${i}${path_suffix} ${algo} ${algo_params}"
    log="${log_prefix}${i}${log_suffix}"
    echo $cmd
    # $cmd >out_${log} 2>err_${log} &
    $cmd >/dev/null 2>err_${log} &
    echo $! >> ./pid
done

echo "-> Waiting for clients to finish"
while [ `cat err_${log_prefix}*${log_suffix} | grep "FINISHED" | wc -l` -lt ${N_node} ]; do
    sleep 5
done

echo "-> Killing processes"
for line in `cat ./pid`
do
    echo $line
    kill -9 $line
done
