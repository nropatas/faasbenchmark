#!/bin/bash

if [ -z "$1" ]
then
  echo "Please specify a command"
  exit 1
fi

if [[ -z "$2" || $(($2)) == 0 ]]
then
  echo "Please specify the number of iterations"
  exit 1
fi

if [[ -z "$3" && $(($2)) > 1 ]]
then
  echo "Please specify the time interval between iterations (e.g., 30s, 5m)"
  exit 1
fi

cmd=$1
iter=$(($2))
interval=$3

for i in `seq 1 $iter`
do
  echo "$i of $iter"
  eval $cmd

  if (( i < iter ))
  then
    echo "$i of $iter done!"
    echo "Sleeping for $interval..."
    sleep $interval
  fi
done
