#!/bin/bash

# check if running inside this directory
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
if [ "$script_dir" != "$(pwd)" ]; then
    echo "Error: Please run this script from the directory where it is located"
    exit 1
fi

# check if user's variables are set
if [ -z "$FABRIC_CFG_PATH" ]; then
    echo "Error: Please set the FABRIC_CFG_PATH environment variable in user level"
    exit 1
fi

source .env

peer lifecycle chaincode package ${NOFEEPAY_CC_NAME}.tar.gz \
--path . \
--lang golang \
--label $NOFEEPAY_CC_NAME

echo "Chaincode package is generated in $(realpath ${NOFEEPAY_CC_NAME}.tar.gz)"