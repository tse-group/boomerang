#! /bin/bash -ve

if [ "$#" -ne 3 ]; then
    echo "Usage: ./load_cfg.sh <scenario ID> <scenario instance> <num nodes>"
    exit -1
fi


rm *.txt *.json || true
cp ../gen_trace/$1/$2_graph.txt ./graph.txt
cp ../gen_trace/$1/$2_payments.txt ./payments.txt
cp ../gen_trace/$1/$2_paths.txt ./paths.txt

./parse_graph graph.txt payments.txt paths.txt 127.0.0.1 $3
