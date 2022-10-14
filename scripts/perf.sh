#!/bin/bash

set -euo pipefail
IFS=$'\n\t'

ROWS=${ROWS:-100000000}

# https://stackoverflow.com/questions/59895/how-do-i-get-the-directory-where-a-bash-script-is-located-from-within-the-script
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

echo "tick,rows_done" > ${SCRIPT_DIR}/template.perf

echo "{{Noun}}" | ${SCRIPT_DIR}/../fakedata -l ${ROWS} | pv -b -l -a -t -n > /dev/null 2>> ${SCRIPT_DIR}/template.perf

sed -i -e 's/ /,/g' ${SCRIPT_DIR}/template.perf

echo "tick,rows_done" > ${SCRIPT_DIR}/generator.perf

${SCRIPT_DIR}/../fakedata noun -l ${ROWS} | pv -b -l -a -t -n > /dev/null 2>> ${SCRIPT_DIR}/generator.perf

sed -i -e 's/ /,/g' ${SCRIPT_DIR}/generator.perf

cat generator.perf | sqlite-utils insert -d --csv ${SCRIPT_DIR}/perf.db generator -
cat template.perf | sqlite-utils insert -d --csv ${SCRIPT_DIR}/perf.db template -

rm ${SCRIPT_DIR}/template.perf-e ${SCRIPT_DIR}/generator.perf-e

cat ${SCRIPT_DIR}/queries.csv | sqlite-utils insert --csv perf.db saved_queries  -
