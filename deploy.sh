#!/bin/sh

cdk deploy MicrobatchLambdaStack --outputs-file outputs.json

yq eval-all --inplace --input-format yaml --output-format yaml '
    select(fileIndex == 0).region = select(fileIndex == 1).MicrobatchLambdaStack.ServerRegion |
    select(fileIndex == 0).lambda-url = select(fileIndex == 1).MicrobatchLambdaStack.ServerUrl |
    select(fileIndex == 0)
    '\
    config.yaml \
    outputs.json
