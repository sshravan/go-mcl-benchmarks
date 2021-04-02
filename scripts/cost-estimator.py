import pandas as pd
import numpy as np
import json
import math

# ==============================================================================
NS_to_US = 10**3
US_to_MS = 10**3
MS_to_S = 10**3
NS_to_S = 10**9
US_to_S = 10**6
MS_to_S = 10**3
sep = "==============================================================="
# ==============================================================================


def IsPowOf2(x): return (x & (x - 1)) == 0


def next_pow_of_2(v):
    # Example input: output
    # 0 1
    # 1024 1024
    # 1023 1024
    # 1205 2048
    # 2048 2048
    # 491520 524288 (== 1 << 19)
    if v == 0:
        return 1

    return 1 << (v - 1).bit_length()


def std_out(msg, L, txn, total_cost): return "{:15}\t{:2}\t{:10}\t{:12.3f}\tseconds".format(
    msg, L, txn, total_cost)

# ==============================================================================


def naive_verification(db, L, txn):
    cost_per_txn = db['G1Mul'] + db['G1Sub'] + \
        (L + 1) * db['MillerLoopVec32Avg'] + \
        db['FinalExp'] + \
        db['GTIsEqual']
    total_cost = (txn * cost_per_txn)
    return total_cost


def naive_verification_driver(db, L, txn_count):
    for i in range(len(txn_count)):
        total_cost = naive_verification(db, L, txn_count[i])
        total_cost = total_cost / NS_to_S
        out_str = std_out("NaiveVerify", L, txn_count[i], total_cost)
        print(out_str)


def driver(db):
    txn_count = [1 << i for i in range(2, 15)]
    ell = [30]
    print(sep)
    for i in range(len(ell)):
        naive_verification_driver(db, ell[i], txn_count)
# ==============================================================================


if __name__ == "__main__":
    print("Hello, World!")
    filename = "benchmarking-results-nanoseconds.json"
    with open(filename) as f:
        db = json.load(f)

    print(*db.keys(), sep='\n')
    print()
    driver(db)
