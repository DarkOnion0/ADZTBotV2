# This script come from this gist --> https://gist.github.com/ferferga/93ca1ab3056257d05f6e33af0d6ead49

#!/bin/bash
set -e

# Simple script to remove dangling images from GHCR.
# You need to have installed gh cli and jq for this script to work properly
# You need to be logged to 'gh' first
container="adztbotv2"
temp_file="ghcr_prune.ids"
rm -rf $temp_file

echo "Fetching dangling images from GHCR..."
gh api /user/packages/container/${container}/versions --paginate > $temp_file

ids_to_delete=$(cat "$temp_file" | jq -r '.[] | select(.metadata.container.tags==[]) | .id')

if [ "${ids_to_delete}" = "" ]
then
	echo "There are no dangling images to remove for this package"
	exit 0
fi

echo -e "\nDeleting dangling images..."
while read -r line; do
	id="$line"
	## Workaround for https://github.com/cli/cli/issues/4286 and https://github.com/cli/cli/issues/3937
	echo -n | gh api --method DELETE /user/packages/container/${container}/versions/${id} --input -
	echo Dangling image with ID $id deleted successfully
done <<< $ids_to_delete

rm -rf $temp_file
echo -e "\nAll the dangling images have been removed successfully"
exit 0
