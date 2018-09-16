#!/bin/sh

# Copy the log files to VMs
for i in $(seq 1 10)
do
	scp ../../data/vm$i.log $CS425NETID@fa18-cs425-g$CS425GROUPID-$(printf %02d $i).cs.illinois.edu:~/data/mp1/
done
