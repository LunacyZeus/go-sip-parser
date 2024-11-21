#!/usr/bin/env python
# -*- coding: utf-8 -*-


import telnetlib

HOST = "127.0.0.1"
PORT = 4320



def command(con, flag, str_=""):
    #data = con.read_until(flag.encode())
    data = con.read_all()
    print(data.decode(errors='ignore'))
    con.write(str_.encode() + b"\n")
    return data

tn = telnetlib.Telnet(HOST,PORT)
data = command(tn, "\n", "login")
print(data)