#!/usr/bin/env bash

scripts_path=$1

GIT_DIR=$(git rev-parse --git-dir)
echo  -e $"Installing hooks..."
# simlinking the prepush script to hooks.
ln -s ${scripts_path}/scripts/pre-push.sh ${GIT_DIR}/hooks/pre-push
echo -e $"Hook configured successfully"