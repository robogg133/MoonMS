#! /bin/sh

shopt -s nullglob

for dir in "$1"/*/; do
    echo "$dir"

    rm "${dir}exec"
done