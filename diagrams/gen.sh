#!/bin/bash

find . -type f -name '*.mmd' -exec  mmdc -i '{}' -o '{}.png' -t dark -b transparent --scale 2 \;
