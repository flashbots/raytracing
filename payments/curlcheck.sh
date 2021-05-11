#!/bin/bash


function check() {
  curl \
    -X POST \
    -H 'content-type:application/json' \
    --data '{"id":1,"jsonrpc":"2.0","method":"eth_getCode","params":["0xcc71dfd118f723049ab3d79ea6d1b34b8b4c928a"]}' \
    http://138.68.75.41:8545
}

check
