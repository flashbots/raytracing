#!/bin/bash

function mine() {
  curl \
    -X POST \
    -H 'content-type:application/json' \
    --data '{"id":1,"jsonrpc":"2.0","method":"evm_mine","params":[]}' \
    127.0.0.1:8545
}

mine
