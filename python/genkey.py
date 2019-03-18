#!/usr/bin/env python
# -*- coding: utf-8 -*-

import argparse

from mnemonic import Mnemonic
from pycoin.key.BIP32Node import BIP32Node
from pycoin.key.key_from_text import key_from_text

def genWords(sec):
    m = Mnemonic("english")
    code = m.generate(sec)
    return code

def genRoot(words, passphrase, netcode, private=False):
    seed = Mnemonic.to_seed(words, passphrase)
    root = BIP32Node.from_master_secret(seed)
    root._netcode = netcode
    rootkey = root.wallet_key(as_private=private)
    return rootkey

def genAccount(root, path, netcode):
    master = key_from_text(root)
    key = master.subkey_for_path(path)
    key._netcode = netcode
    pubkey = key.wallet_key()
    prvkey = key.wallet_key(as_private=True)
    return pubkey, prvkey

if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('-n', dest='netcode', default='BTC', help='BTC/XTN/LTC/DASH/BCH/BSV, default is BTC')
    args = parser.parse_args()

    codePath = {
        'BTC' : '44H/0H/0H',
        'XTN' : '44H/1H/0H',
        'LTC' : '44H/2H/0H',
        'DASH': '44H/5H/0H',
        'BCH' : '44H/145H/0H',
        'BSV' : '44H/236H/0H',
    }

    netcode = args.netcode
    if netcode not in codePath:
        parser.print_help()
        exit()

    words = genWords(256)
    passphrase = ''
    path = codePath[netcode]
    if netcode == 'BCH' or netcode == 'BSV':
        netcode = 'BTC'

    master = genRoot(words, passphrase, netcode, True)
    account, prvkey = genAccount(master, path, netcode)

    print("Seed:", words)
    print("Rootkey:", master)
    print("AccountPath:", path)
    print("AccountPrivate:", prvkey)
    print("Account:", account)
