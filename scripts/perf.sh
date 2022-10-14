#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

ROWS=100000000

echo "tick,rows_done" > template.perf

echo "{{Noun}}" | fakedata -l ${ROWS} | pv -b -l -a -t -n > /dev/null 2>> template.perf

sed -i -e 's/ /,/g' template.perf

echo "tick,rows_done" > generator.perf

fakedata noun -l ${ROWS} | pv -b -l -a -t -n > /dev/null 2>> generator.perf

sed -i -e 's/ /,/g' generator.perf

rm template.perf-e generator.perf-e