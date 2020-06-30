#!/bin/bash

# providers=(azure)
providers=(aws knative openfaas kubeless fission)
tests=(BurstLvl1)
iterations=10
interval=5m

for provider in ${providers[*]}
do
  for test in ${tests[*]}
  do
    echo "Testing $provider $test"
    # export KUBECONFIG=/app/kubeconfigs/kubeconfig_$provider
    ./runtest.sh "./faasbenchmark run $provider $test -r results/$test/$provider" $iterations $interval
    echo "$provider $test done!"
    echo "Sleeping for $interval"
    sleep $interval
  done
done
