#!/bin/bash

make build

if [ $? -eq 0 ];then
  echo "build success"
else
  echo "build fail"
  exec 1
fi

