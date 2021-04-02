#!/usr/bin/env bash
# bash scripts/gmcl.sh instead of sh scripts/gmcl.sh
set -e
shopt -s expand_aliases
alias time='date; time'

scriptdir=$(cd $(dirname $0); pwd -P)
sourcedir=$(cd $scriptdir/..; pwd -P)

go build && time ./go-mcl-benchmarks -test.benchtime 100x
