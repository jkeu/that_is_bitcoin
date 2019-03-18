#!/usr/bin/env python
# -*- coding: utf-8 -*-

import argparse
import binascii

from pycoin.tx.tx_utils import sign_tx
from pycoin.tx.Tx import Tx, SIGHASH_ALL
from pycoin.key.Key import Key
from pycoin.tx.pay_to.ScriptPayToAddressWit import ScriptPayToAddressWit
from pycoin.ui import address_for_pay_to_script, script_obj_from_address
from pycoin.serialize import h2b_rev, b2h_rev
from pycoin.tx.TxIn import TxIn
from pycoin.tx.TxOut import TxOut
from pycoin.tx.Spendable import Spendable
from pycoin.tx.pay_to import build_hash160_lookup, build_p2sh_lookup
from pycoin.tx.script import tools

def signTxLegacy(txHex, fromWif):
    tx = Tx.from_hex(txHex)
    sign_tx(tx, wifs=[fromWif])
    return tx.as_hex()

def signTxSegwit(txHex, wif_key):
    tx = Tx.from_hex(txHex)

    my_key = Key.from_text(wif_key)
    script = ScriptPayToAddressWit(b'\0', my_key.hash160()).script()
    redeem = binascii.hexlify(script).decode()

    in_keys = [wif_key]
    hash160_lookup = build_hash160_lookup([Key.from_text(wif).secret_exponent() for wif in in_keys])
    for i in range(0, len(tx.txs_in)):
        p2sh_lookup = build_p2sh_lookup([binascii.unhexlify(redeem)])
        tx.sign_tx_in(hash160_lookup, i, tx.unspents[i].script, hash_type=SIGHASH_ALL, p2sh_lookup=p2sh_lookup)
    return tx.as_hex()

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('txHex', help='unsigned tx hex')
    parser.add_argument('-w', dest='fromWif', help='from address wif')
    parser.add_argument('-s', dest='segwit', action='store_true', help='segwit signature for BTC/XTN/LTC')
    args = parser.parse_args()

    txHex = args.txHex
    fromWif = args.fromWif

    if not args.segwit:
        signedTx = signTxLegacy(txHex, fromWif)
    else:
        signedTx = signTxSegwit(txHex, fromWif)

    print('TxSigned:', signedTx)

