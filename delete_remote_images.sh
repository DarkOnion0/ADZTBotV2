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

ls -larth

if [ "${ids_to_delete}" = "" ]
then
	echo "There are no dangling images to remove for this package"
	exit 0
fi

echo -e "\nDeleting dangling images..."
while read -r line; do
	id="$line"
	url="https://api.github.com/user/packages/container/$container/versions/$id"

	#echo $id $container $temp_file $url
	#echo -e "/user/packages/container/${container}/versions/${id}"

	curl -X DELETE -u $auth -H "Accept: application/vnd.github.v3+json" $url
	echo Dangling image with ID $id deleted successfully
done <<< $ids_to_delete

rm -rf $temp_file
echo -e "\nAll the dangling images have been removed successfully"
exit 0
