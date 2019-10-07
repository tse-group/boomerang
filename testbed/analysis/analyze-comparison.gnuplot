set terminal pngcairo dashed size 900,500 enhanced font "Verdana,10";

set style data histogram;
set style histogram errorbars gap 1 lw 2;   # set style histogram errorbars cluster gap 1;
set style fill solid;
set bars front
set boxwidth 0.9;
set xtics format "";
set grid ytics;
set yrange [0:];

set key on left top;


is_retry(x,y)= (x eq "retry-02-amp-2" ? y : 1/0)
is_redundancy(x,y)= (x eq "redundancy-02-amp-2" ? y : 1/0)
is_redundantretry(x,y)= (x eq "redundantretry-02-amp-2" ? y : 1/0)



COLOR_RETRY_02_AMP = '#77ac30'
COLOR_REDUNDANCY_02_AMP = '#4dbeee'
COLOR_REDUNDANTRETRY_02_AMP_10 = '#edb120'


## New default Matlab line colors, introduced together with parula (2014b)
#set style line 11 lt 1 lc rgb '#0072bd' pt 1 lw 2 dt 1    # blue
#set style line 12 lt 1 lc rgb '#0072bd' pt 1 lw 2 dt "-"  # blue
#set style line 13 lt 1 lc rgb '#0072bd' pt 1 lw 3 dt "."  # blue
#set style line 14 lt 1 lc rgb '#0072bd' pt 1 lw 2 dt 4    # blue
#set style line 15 lt 1 lc rgb '#0072bd' pt 1 lw 2 dt 5    # blue
#
#set style line 21 lt 1 lc rgb '#d95319' pt 2 lw 2 dt 1    # orange
#set style line 22 lt 1 lc rgb '#d95319' pt 2 lw 2 dt "-"  # orange
#set style line 23 lt 1 lc rgb '#d95319' pt 2 lw 3 dt "."  # orange
#set style line 24 lt 1 lc rgb '#d95319' pt 2 lw 2 dt 4    # orange
#set style line 25 lt 1 lc rgb '#d95319' pt 2 lw 2 dt 5    # orange
#
#set style line 31 lt 1 lc rgb '#edb120' pt 3 lw 2 dt 1    # yellow
#set style line 32 lt 1 lc rgb '#edb120' pt 3 lw 2 dt "-"  # yellow
#set style line 33 lt 1 lc rgb '#edb120' pt 3 lw 3 dt "."  # yellow
#set style line 34 lt 1 lc rgb '#edb120' pt 3 lw 2 dt 4    # yellow
#set style line 35 lt 1 lc rgb '#edb120' pt 3 lw 2 dt 5    # yellow
#
#set style line 41 lt 1 lc rgb '#7e2f8e' pt 4 lw 2 dt 1    # purple
#set style line 42 lt 1 lc rgb '#7e2f8e' pt 4 lw 2 dt "-"  # purple
#set style line 43 lt 1 lc rgb '#7e2f8e' pt 4 lw 3 dt "."  # purple
#set style line 44 lt 1 lc rgb '#7e2f8e' pt 4 lw 2 dt 4    # purple
#set style line 45 lt 1 lc rgb '#7e2f8e' pt 4 lw 2 dt 5    # purple
#
#set style line 51 lt 1 lc rgb '#77ac30' pt 5 lw 2 dt 1    # green
#set style line 52 lt 1 lc rgb '#77ac30' pt 5 lw 2 dt "-"  # green
#set style line 53 lt 1 lc rgb '#77ac30' pt 5 lw 3 dt "."  # green
#set style line 54 lt 1 lc rgb '#77ac30' pt 5 lw 2 dt 4    # green
#set style line 55 lt 1 lc rgb '#77ac30' pt 5 lw 2 dt 5    # green
#
#set style line 61 lt 1 lc rgb '#4dbeee' pt 6 lw 2 dt 1    # light-blue
#set style line 62 lt 1 lc rgb '#4dbeee' pt 6 lw 2 dt "-"  # light-blue
#set style line 63 lt 1 lc rgb '#4dbeee' pt 6 lw 3 dt "."  # light-blue
#set style line 64 lt 1 lc rgb '#4dbeee' pt 6 lw 2 dt 4    # light-blue
#set style line 65 lt 1 lc rgb '#4dbeee' pt 6 lw 2 dt 5    # light-blue
#
#set style line 71 lt 1 lc rgb '#a2142f' pt 7 lw 2 dt 1    # red
#set style line 72 lt 1 lc rgb '#a2142f' pt 7 lw 2 dt "-"  # red
#set style line 73 lt 1 lc rgb '#a2142f' pt 7 lw 3 dt "."  # red
#set style line 74 lt 1 lc rgb '#a2142f' pt 7 lw 2 dt 4    # red
#set style line 75 lt 1 lc rgb '#a2142f' pt 7 lw 2 dt 5    # red

