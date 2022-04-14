# This script come from this gist --> https://gist.github.com/ferferga/93ca1ab3056257d05f6e33af0d6ead49

#!/bin/bash
set -e

# Simple script to remove dangling images from GHCR.
# You need to have installed jq for this script to work properly
container="adztbotv2"
temp_file_delete="ghcr_prune_delete.ids"
temp_file_exclude="ghcr_prune_exclude.ids"

rm -rf $temp_file_exclude
rm -rf $temp_file_delete

echo "Fetching dangling images from GHCR..."
curl -u $auth -H "Accept: application/vnd.github.v3+json" https://api.github.com/user/packages/container/$container/versions > $temp_file_exclude

ids_to_delete=$(cat "$temp_file_delete" | jq -r '.[] | select(.metadata.container.tags==[]) | .id')

if [ "${ids_to_delete}" = "" ]
then
	echo "There are no dangling images to remove for this package"
	exit 0
fi

ids_to_exclude=$(cat "$temp_file_exclude" | jq -r '.[] | select(.metadata.container.tags!=[]) | .name')

while read -r line; do

done <<< $ids_to_exclude

if [ "${ids_to_exclude}" = "" ]
then
	echo "There are no dangling images to exclude for this package"
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

rm -rf $temp_file_exclude
rm -rf $temp_file_delete
echo -e "\nAll the dangling images have been removed successfully"
exit 0
