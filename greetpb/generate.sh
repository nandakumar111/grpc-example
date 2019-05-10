#!/bin/bash

## Greetings
protoc ./greet.proto --go_out=plugins=grpc:.

## Calculator
#protoc greet/calculatorpb/calculator.proto --go_out=plugins=grpc:.