#!/bin/sh

for k in $(seq 1 10)    
do
    for i in $(seq 1 100)
    do
        echo "*****************************"
        echo "ROUND : $k"
        echo "--------------------------------"
        echo "CellID $i"
        echo "$i"
        mkdir csv/cellID$i
        tshark -r tcpdump/cellID$i/round$k.pcap -T fields -e frame.number -e frame.time_relative -e ip.src -e ip.dst -e frame.protocols -e frame.len -E separator=, > csv/cellID$i/dumpRAW$k.csv
        grep eth:ethertype:ip:tcp:ssl csv/cellID$i/dumpRAW$k.csv > csv/cellID$i/dump$k.csv
        rm csv/cellID$i/dumpRAW$k.csv

    done
done
