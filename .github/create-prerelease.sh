#!/bin/bash

# Exit on error
set -e

TAG_NAME="pre.$(date +'%Y.%m.%d.%H%M')"
echo "Creating and pushing tag: $TAG_NAME"

git tag -a "$TAG_NAME" -m "Create pre-release tag $TAG_NAME"

git push origin "$TAG_NAME"

echo "Tag $TAG_NAME has been successfully pushed!"
