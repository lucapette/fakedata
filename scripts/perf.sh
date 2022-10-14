#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

ROWS=${ROWS:-100000000}

echo "tick,rows_done" > template.perf

echo "{{Noun}}" | fakedata -l ${ROWS} | pv -b -l -a -t -n > /dev/null 2>> template.perf

sed -i -e 's/ /,/g' template.perf

echo "tick,rows_done" > generator.perf

fakedata noun -l ${ROWS} | pv -b -l -a -t -n > /dev/null 2>> generator.perf

sed -i -e 's/ /,/g' generator.perf

cat generator.perf | sqlite-utils insert -d --csv perf.db generator -
cat template.perf | sqlite-utils insert -d --csv perf.db template -

rm template.perf-e generator.perf-e

cat queries.csv | sqlite-utils insert --csv perf.db saved_queries  -
