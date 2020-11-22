#!/bin/sh
echo "Collecting data for the different cells"
for k in $(seq 1 10)    
do
    for i in $(seq 1 100)
    do
        echo "*****************************"
        echo "ROUND : $k"
        echo "--------------------------------"
        echo "CellID $i"
        mkdir tcpdump/cellID$i
        tcpdump src host 172.22.0.2 or dst host 172.22.0.2 -i eth0 -w tcpdump/cellID$i/round$k.pcap &
        PID_TCPDUMP=$!
        echo "The PID of tcpdump is :$PID_TCPDUMP"
        python3 client.py grid -p key-client.pub -c attr.cred -r 'a' -t  $i
        kill $PID_TCPDUMP
    done
done