#!/bin/bash

repo=$1
token=$2
name=$3


docker pull opensecurity/nodejsscan:cli > /dev/null
docker pull zricethezav/gitleaks > /dev/null
if [[ ! -d results ]]; then
  mkdir results
fi

# format the url correctly
cloneUrl=$(echo $repo | sed "s/https:\/\//https:\/\/${token}@/")


# clone
cloneDir="/tmp/$name-$RANDOM"
pkgFile="$cloneDir/package.json"
git clone $cloneUrl $cloneDir


if [[ -f "$pkgFile" ]]; then
    # perform ssca
    echo Performing source code analysis
    docker run -v $cloneDir:/src opensecurity/nodejsscan:cli -d /src > ./results/$name-ssca-results.json


    # perform leak analysis
    echo Checking for pushed tokens 
    docker run -v $cloneDir:/src zricethezav/gitleaks -r /src > ./results/$name-key-results.log

    # perform audit
    echo Auditing packages
    yarn audit --json > ./results/$name-audit-results.json
fi