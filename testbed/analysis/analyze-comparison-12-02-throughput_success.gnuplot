load("analyze-comparison.gnuplot")
set output DATAFILENAME."-".ARG0.".png";
set title ARG0;

set ylabel "Success Throughput";
set xlabel "Number of Retries/Redundant Payments";

plot \
    DATAFILENAME every 3::0 using (column("throughput_success-mean")):(column("throughput_success-std")):xtic(1) title 'Retry-02-AMP' linecolor rgb COLOR_RETRY_02_AMP, \
    DATAFILENAME every 3::1 using (column("throughput_success-mean")):(column("throughput_success-std")):xtic(1) title 'Redundancy-02-AMP' linecolor rgb COLOR_REDUNDANCY_02_AMP, \
    DATAFILENAME every 3::2 using (column("throughput_success-mean")):(column("throughput_success-std")):xtic(1) title 'Redundantretry-02-AMP(10)' linecolor rgb COLOR_REDUNDANTRETRY_02_AMP_10, \
    ;
