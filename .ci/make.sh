#!/bin/bash

git_status_output=`git status`

make -f Makefile || exit

git_status_output_after=`git status`

if [[ "$git_status_output" != "$git_status_output_after" ]]; then
   echo ERROR: make command modified code
   echo "$git_status_output"
   echo "$git_status_output_after"
    exit 1
fi

