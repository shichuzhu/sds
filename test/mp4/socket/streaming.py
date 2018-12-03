from __future__ import print_function

import sys

from pyspark import SparkContext
from pyspark.streaming import StreamingContext

if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: network_wordcount.py <hostname> <port>", file=sys.stderr)
        sys.exit(-1)
    sc = SparkContext(appName="PythonStreamingNetworkStreamProcessing")
    ssc = StreamingContext(sc, 1)

    lines = ssc.socketTextStream(sys.argv[1], int(sys.argv[2]))
    lines.filter(lambda line: (len(line) > 15)).map(lambda line: line.upper()).filter(
        lambda line: ("ABCd" in line or "BCDe" in line)).saveAsTextFiles("test_folder")

    ssc.start()
    ssc.awaitTermination()
