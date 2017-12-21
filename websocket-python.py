#!/usr/bin/env python

import websocket
import ssl
import sys
import os
import time

registrations = []

def on_message(ws, message):
    print(message)

def on_error(ws, error):
    print(error)

def on_close(ws):
    print("### closed ###")

def on_open(ws):
    if len(registrations) < 1:
        print "Registering for all events"
        ws.send("register *.*.*")
    else:
        for s in registrations:
            print "Registering for %s" % (s)
            ws.send("register %s" % (s))

if __name__ == "__main__":
    key = os.getenv('RS_KEY', 'rocketskates:r0cketsk8ts')
    token = os.getenv('RS_TOKEN', '')
    endpoint = os.getenv('RS_ENDPOINT', 'https://127.0.0.1:8092')

    if (token == '') & (key == ''):
        print "Please provide a key or token."
        sys.exit(1)

    if token == '':
        token = key

    if len(sys.argv) > 1:
        registrations = sys.argv[1:]

    if (endpoint != '') & endpoint.startswith('https://'):
        endpoint = endpoint[len('https://'):]

    ws = websocket.WebSocketApp(
            "wss://%s/api/v3/ws?token=%s" % (endpoint, token),
            on_message = on_message,
            on_error = on_error,
            on_close = on_close,
            on_open = on_open)
    ws.run_forever(sslopt = {"cert_reqs": ssl.CERT_NONE})
