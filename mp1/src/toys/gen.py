import time
import sys
i = 0
for i in range(5):
    print(f"hello, world! #{i}")
    sys.stdout.flush()
    i += 1
    time.sleep(0.5)
