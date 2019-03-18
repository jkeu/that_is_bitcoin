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

NET_CODE = 'XTN'

def signTxLegacy(txHex, fromWif):
    print('TxRaw:', txHex)
    # print('fromWif:', fromWif)

    tx = Tx.from_hex(txHex)
    sign_tx(tx, wifs=[fromWif])
    print('TxFinal:', repr(tx))
    # print('TxSigned:', tx.as_hex())

    return tx.as_hex()

def pkh_segwit_address_from_wif(wif):
    """
    The P2SH redeemScript is always 22 bytes.
    It starts with a OP_0, followed by a canonical push of the keyhash (i.e. 0x0014{20-byte keyhash})
    Same as any other P2SH, the scriptPubKey is OP_HASH160 hash160(redeemScript) OP_EQUAL
    """
    my_key = Key.from_text(wif)
    script = ScriptPayToAddressWit(b'\0', my_key.hash160()).script()
    script_hex = binascii.hexlify(script).decode()
    return address_for_pay_to_script(script, netcode=NET_CODE), script_hex

def spend_sh_fund(tx_ins, wif_keys, tx_outs):
    """
    spend script hash fund
    the key point of an input comes from multisig address is that,
    its sign script is combined with several individual signs
    :param tx_ins: list with tuple(tx_id, idx, balance, address, redeem_script)
    :param wif_keys: private keys in wif format,
        technical should be the same order with the pubkey in redeem script,
        but pycoin has inner control, so here order is not mandatory
    :param tx_outs: balance, receiver_address
    :return: raw hex and tx id
    """
    _txs_in = []
    _un_spent = []
    for tx_id, idx, balance, address, _ in tx_ins:
        # must h2b_rev NOT h2b
        tx_id_b = h2b_rev(tx_id)
        _txs_in.append(TxIn(tx_id_b, idx))

        _un_spent.append(Spendable(balance, script_obj_from_address(address, netcodes=[NET_CODE]).script(),
                                   tx_id_b, idx))

    _txs_out = []
    for balance, receiver_address in tx_outs:
        # _txs_out.append(TxOut(balance, script_obj_from_address(receiver_address, netcodes=[NET_CODE]).script()))
        _txs_out.append(TxOut(balance, receiver_address))

    version, lock_time = 1, 0
    tx = Tx(version, _txs_in, _txs_out, lock_time)
    tx.set_unspents(_un_spent)

    # construct hash160_lookup[hash160] = (secret_exponent, public_pair, compressed) for each individual key
    hash160_lookup = build_hash160_lookup([Key.from_text(wif_key).secret_exponent() for wif_key in wif_keys])

    for i in range(0, len(tx_ins)):
        # you can add some conditions that if the input script is not p2sh type, not provide p2sh_lookup,
        # so that all kinds of inputs can work together
        p2sh_lookup = build_p2sh_lookup([binascii.unhexlify(tx_ins[i][-1])])
        tx.sign_tx_in(hash160_lookup, i, tx.unspents[i].script, hash_type=SIGHASH_ALL, p2sh_lookup=p2sh_lookup)

    return tx.as_hex(), tx.id()

def signTxSegwit(txHex, wif_key):
    tx = Tx.from_hex(txHex)
    address, redeem = pkh_segwit_address_from_wif(wif_key)
    print(address, redeem)
    print('unspents:', tx.unspents)
    amounts = []
    for uns in tx.unspents:
        # print('\tuns:', uns)
        # amount += uns.coin_value
        amounts.append(uns.coin_value)
    
    print('amounts:', amounts)
    print('tx:', repr(tx))
    print('txIn:', tx.TxIn)
    tx_ins = []
    idx = 0
    for txIn in tx.txs_in:
        print(txIn)
        print(b2h_rev(txIn.previous_hash), txIn.previous_index)
        amount = amounts[idx]
        newTxIn = (b2h_rev(txIn.previous_hash), txIn.previous_index, amount, address, redeem)
        tx_ins.append(newTxIn)
        idx += 1

    print(tx_ins)
    # tx_ins = [('501fb0d89c3ae82a0c8971c1ff9c8f79f3235be2049e3bf45aa4b7099347bf4a', 0, 1000000,
    #           '2N7x1K4xpHdazzWhSnSNfoqj369bfu5gqF7', '0014035568a6faa755c17337a30896db68837ff49731')
    #          ]
    in_keys = [wif_key]

    tx_outs = []
    for txOut in tx.txs_out:
        newTxOut = (txOut.coin_value, txOut.script)
        # newTxOut = (txOut.coin_value, tools.disassemble(txOut.script))
        # newTxOut = (txOut.coin_value, '2N7x1K4xpHdazzWhSnSNfoqj369bfu5gqF7')
        tx_outs.append(newTxOut)

    # tx_outs = [(990000, 'mqswKFc8VTMtF986xfQKBHaQoh42P4c1gi')]
    print(tx_outs)

    raw_hex, tx_id = spend_sh_fund(tx_ins, in_keys, tx_outs)
    print('signed raw hex:')
    print(raw_hex)
    print('txn id/hash:')
    print(tx_id)

    return raw_hex

def signTxSegwit2(txHex, wif_key):
    tx = Tx.from_hex(txHex)

    address, redeem = pkh_segwit_address_from_wif(wif_key)
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
        signedTx = signTxSegwit2(txHex, fromWif)

    print('TxSigned:', signedTx)

