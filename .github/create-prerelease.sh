#!/bin/bash

# Exit on error
set -e

# Get the message from the latest commit and ask for confirmation
MESSAGE=$(git log -1 --pretty=%B)
echo "Do you want to create a pre-release on this commit?"
echo -e "\033[33m$MESSAGE\033[0m"
read -p "[y/N]: " -n 1 -r

if [[ ! $REPLY =~ ^[Yy]$ ]]; then
  echo ""
  echo "Aborting."
  exit
fi

TAG_NAME="pre.$(date -u +'%Y.%m.%d.%H%M')"
echo "Creating and pushing tag: $TAG_NAME"

git tag -a "$TAG_NAME" -m "Create pre-release tag $TAG_NAME"

git push origin "$TAG_NAME"

echo "Tag $TAG_NAME has been successfully pushed!"
