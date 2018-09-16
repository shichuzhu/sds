#!/bin/sh

# Copy the log files to VMs
for i in $(seq 1 10)
do
	scp ../../data/vm$i.log szhu28@fa18-cs425-g44-$(printf %02d $i).cs.illinois.edu:~/data/mp1/
done
