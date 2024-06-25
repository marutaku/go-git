#!/usr/bin/env python3
import zlib
import sys

if len(sys.argv) > 1:
    print(f"File in was: {sys.argv[1]}")
    file = sys.argv[1]
else:
    print("File in was BLANK, please pass a zlibbed file path...")
    exit(1)

try:
    compressed_str = open(file, "rb").read()
    decompressed = zlib.decompress(compressed_str)
    print(decompressed)
except Exception as e:
    print("Something went wrong: ", e)
