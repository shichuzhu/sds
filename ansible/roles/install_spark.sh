#!/bin/bash

mkdir soft
cd soft
wget http://mirrors.ibiblio.org/apache/spark/spark-2.4.0/spark-2.4.0-bin-hadoop2.7.tgz
tar xf spark-2.4.0-bin-hadoop2.7.tgz

echo 'export PATH=/home/szhu28/soft/spark-2.4.0-bin-hadoop2.7/bin:$PATH' >> ~/.bashrc
