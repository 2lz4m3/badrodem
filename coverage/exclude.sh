#!/bin/bash
# https://dev.to/talalyousif/excluding-files-from-code-coverage-in-go-291f
while read p || [ -n "$p" ]
do
  ep=${p//\//\\/}
  sed -i "/^${ep}/d" "`dirname $0`/out/coverage.out"
done < "`dirname $0`/exclusions.txt"
