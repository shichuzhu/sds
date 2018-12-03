from __future__ import print_function

import sys

from pyspark import SparkContext
from pyspark.streaming import StreamingContext

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("USSEAGE input file path", file=sys.stderr)
        sys.exit(-1)
    sc = SparkContext(appName="PythonStreamingNetworkStreamProcessing")
    ssc = StreamingContext(sc, 1)

    # lines = ssc.textFileStream(sys.argv[1])
    lines = ssc.textFileStream("test/mp4/data/input.txt")
    result = lines.filter(lambda line: (len(line) > 15)).map(
        lambda line: line.upper()).map(lambda line: line + "!!!")
    print("1234567890", result.collect())

    ssc.start()
    ssc.awaitTermination()
