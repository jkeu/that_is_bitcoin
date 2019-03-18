#!/usr/bin/env python
# -*- coding: utf-8 -*-

import argparse

from pycoin.services.chain_so import ChainSoProvider
from pycoin.tx.tx_utils import create_tx, sign_tx
from pycoin.tx.Tx import Tx

def genTx(fromAddress, toAddress, netcode):
    b = ChainSoProvider(netcode)
    spendables = b.spendables_for_address(fromAddress)
    if not spendables:
        print("{} no spenables tx.".format(fromAddress))
        return

    print('FromAddress:', fromAddress)
    print('Spending:', spendables)
    print('ToAddress:', toAddress)

    tx = create_tx(spendables, [toAddress])
    print('TxRaw   :', tx.as_hex(include_unspents=True))

    return tx.as_hex(include_unspents=True)

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-n', dest='netcode', default='BTC', help='BTC/XTN/LTC/DASH, default is BTC')
    parser.add_argument('fromAddress', help='transfer from address')
    parser.add_argument('-t', dest='toAddress', help='transfer to address')
    args = parser.parse_args()

    codePath = {
        'BTC' : '44H/0H/0H',
        'XTN' : '44H/1H/0H',
        'LTC' : '44H/2H/0H',
        'DASH': '44H/5H/0H',
        'BCH' : '44H/145H/0H',
        'BSV' : '44H/236H/0H',
    }

    fromAddress = args.fromAddress
    toAddress = args.toAddress
    netcode = args.netcode
    if netcode not in codePath:
        parser.print_help()
        exit()
    
    if netcode == 'BCH' or netcode == 'BSV':
        netcode = 'BTC'

    rawTx = genTx(fromAddress, toAddress, netcode)
    print('Tx:', rawTx)
