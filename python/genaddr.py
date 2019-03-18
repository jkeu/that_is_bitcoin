#!/usr/bin/env python
# -*- coding: utf-8 -*-

import argparse

from pycoin.key.key_from_text import key_from_text
from pycoin.tx.pay_to.ScriptPayToAddressWit import ScriptPayToAddressWit
from pycoin.ui import address_for_pay_to_script

def genAddress(pubkey, path, netcode, is_segwit=False):
    master = key_from_text(pubkey)
    key = master.subkey_for_path(path)
    key._netcode = netcode
    if not is_segwit:  # legacy
        address = key.address()
    elif is_segwit:    # segwit
        hash160_c = key.hash160(use_uncompressed=False)
        script = ScriptPayToAddressWit(b'\0', hash160_c).script()
        address = address_for_pay_to_script(script, key._netcode)

    return address

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-n', dest='netcode', default='BTC', help='BTC/XTN/LTC/DASH, default is BTC')
    parser.add_argument('pubkey', help='account public key')
    parser.add_argument('-s', dest='segwit', action='store_true', help='segwit address for BTC/XTN/LTC')
    parser.add_argument('-i', dest='fromIndex', default=0, help='path from 0/i, default 0')
    args = parser.parse_args()

    codePath = {
        'BTC' : '44H/0H/0H',
        'XTN' : '44H/1H/0H',
        'LTC' : '44H/2H/0H',
        'DASH': '44H/5H/0H',
        'BCH' : '44H/145H/0H',
        'BSV' : '44H/236H/0H',
    }

    segwit = args.segwit
    account = args.pubkey
    netcode = args.netcode
    index = int(args.fromIndex)
    if netcode not in codePath:
        parser.print_help()
        exit()
    
    if netcode == 'BCH' or netcode == 'BSV':
        netcode = 'BTC'

    print("Account:", account)
    
    for i in range(index, index + 10):
        path = '0/' + str(i)
        address = genAddress(account, path, netcode, segwit)
        print(path, ":", address)
