#!/bin/bash

NOFEEPAY_CHANNEL_ID=coinlog
NOFEEPAY_CC_NAME=nofeepay

# check if running inside this directory
validate_running_directory() {
    local script_dir
    script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
    if [ "$script_dir" != "$(pwd)" ]; then
        echo "Error: Please run this script from the directory where it is located"
        exit 1
    fi
}


# check if user's variables are set
validate_user_variables() {
    local required_vars=("$@")
    for var in "${required_vars[@]}"; do
        if [ -z "${!var}" ]; then
            echo "Error: Please set the $var environment variable in user level"
            exit 1
        fi
    done
}


validate_binary() {
    local required_bins=("$@")
    for bin in "${required_bins[@]}"; do
        if ! command -v $bin &> /dev/null; then
            echo "Error: '$bin' binary is required but not installed"
            exit 1
        fi
    done
}


package_chaincode() {
    validate_running_directory
    validate_user_variables "FABRIC_CFG_PATH"
    validate_binary "peer"

    peer lifecycle chaincode package "${NOFEEPAY_CC_NAME}.tar.gz" \
        --path . \
        --lang golang \
        --label $NOFEEPAY_CC_NAME

    echo "Chaincode package is generated in $(realpath ${NOFEEPAY_CC_NAME}.tar.gz)"
}


install_chaincode() {
    validate_running_directory
    validate_user_variables "FABRIC_CFG_PATH" "CORE_PEER_TLS_ENABLED" "CORE_PEER_TLS_ROOTCERT_FILE" "CORE_PEER_MSPCONFIGPATH" "CORE_PEER_ADDRESS" "CORE_PEER_LOCALMSPID"
    validate_binary "peer"

    filename=$NOFEEPAY_CC_NAME
    peer lifecycle chaincode install ${filename}.tar.gz
}

show_help() {
    cat << EOF
Usage: ./$(basename "$0") <command>

A CLI utility to manage the chaincode lifecycle for the NoFeePay

Commands:
  package   Package the chaincode into a .tar.gz file
  install   Install the chaincode on the organization's peer
  help      Display this help message
EOF
}


# --- Main CLI Router ---

# Show help if no arguments are provided
if [ $# -eq 0 ]; then
    show_help
    exit 0
fi

command="$1"
shift 

# Route to the appropriate function
case "$command" in
    package)
        package_chaincode
        ;;
    help|-h|--help)
        show_help
        ;;
    install)
        install_chaincode
        ;;
    *)
        echo "Error: Unknown command '$command'"
        show_help
        exit 1
        ;;
esac