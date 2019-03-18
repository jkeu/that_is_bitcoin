#!/usr/bin/env python
# -*- coding: utf-8 -*-

import argparse

from pycoin.tx.tx_utils import sign_tx
from pycoin.tx.Tx import Tx

def signTx(txHex, fromWif):
    print('TxRaw:', txHex)
    print('fromWif:', fromWif)

    tx = Tx.from_hex(txHex)
    sign_tx(tx, wifs=[fromWif])
    print('TxFinal:', repr(tx))
    print('TxSigned:', tx.as_hex())

    return tx.as_hex()

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('txHex', help='unsigned tx hex')
    parser.add_argument('-w', dest='fromWif', help='from address wif')
    args = parser.parse_args()

    txHex = args.txHex
    fromWif = args.fromWif

    signedTx = signTx(txHex, fromWif)
    print(signedTx)
