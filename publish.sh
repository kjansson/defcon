#!/bin/bash
# Usage: publish major|minor|patch
# Will automatically increment from last git tag and tag/push/publish

CURRENT_VER=$(git tag | sort | tail -n 1)
if ! [[ "$CURRENT_VER" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Could not get current version from git."
    exit
fi

CURRENT_MAJOR=$(echo $CURRENT_VER | sed -E 's/v([0-9]+)\.[0-9]+\.[0-9]+/\1/g')
CURRENT_MINOR=$(echo $CURRENT_VER | sed -E 's/v[0-9]+\.([0-9]+)\.[0-9]+/\1/g')
CURRENT_PATCH=$(echo $CURRENT_VER | sed -E 's/v[0-9]+\.[0-9]+\.([0-9]+)/\1/g')

if [[ "$1" == "major" ]]
then
    let CURRENT_MAJOR=$CURRENT_MAJOR+1
    CURRENT_MINOR=0
    CURRENT_PATCH=0
elif [[ "$1" == "minor" ]]
then
    let CURRENT_MINOR=$CURRENT_MINOR+1
    CURRENT_PATCH=0
elif [[ "$1" == "patch" ]]
then
    let CURRENT_PATCH=$CURRENT_PATCH+1
else
    echo "Not valid"
fi

echo "Publishing v$CURRENT_MAJOR.$CURRENT_MINOR.$CURRENT_PATCH"

MODULE_NAME=$(cat go.mod | grep module | awk {'print $2'})
echo $MODULE_NAME
if [[ -z $MODULE_NAME ]];then 
    echo "Could not get module name from go.mod"
    exit
fi

git tag v$CURRENT_MAJOR.$CURRENT_MINOR.$CURRENT_PATCH
if [[ $? != 0 ]]; then
    echo "Git tag error"
fi

git push origin v$CURRENT_MAJOR.$CURRENT_MINOR.$CURRENT_PATCH
if [[ $? != 0 ]]; then
    echo "Git push error"
fi

GOPROXY=proxy.golang.org go list -m $MODULE_NAME@v$CURRENT_MAJOR.$CURRENT_MINOR.$CURRENT_PATCH
if [[ $? != 0 ]]; then
    echo "Go list error"
fi

