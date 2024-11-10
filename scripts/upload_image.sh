#!/bin/bash
JIRA_URL="https://oneline.atlassian.net"
JIRA_API_ENDPOINT="$JIRA_URL/rest/api/3/search"

PROJECT="COM"
USERNAME="xx"
API_TOKEN="xxx"

JQL_QUERY="assignee=currentuser() AND project=COM AND createdDate>'2024-10-30' AND createdDate < '2024-11-30'"
ENCODED_JQL_QUERY=$(echo $JQL_QUERY | jq -sRr @uri)

JIRA_TO_BLP_MAPPING_FILE=/home/tieu/.bcli/files/jira_to_blp.json
JIRA_TO_PR_MAPPING_FILE=/home/tieu/.bcli/files/jira_to_pr.json
PR_IMAGES_FOLDER=/home/tieu/code/personal/pr-images/screenshots/

# Save the response to a variable
response=$(curl -u "$USERNAME:$API_TOKEN" -X GET \
"$JIRA_API_ENDPOINT?jql=$ENCODED_JQL_QUERY" \
-H "Content-Type: application/json")

# Read issues into an array
mapfile -t issues < <(echo "$response" | jq -c '.issues[]')

for issue in "${issues[@]}"; do
  issue_key=$(echo "$issue" | jq -r '.key')
  blue_print_id=$(jq -r --arg key "$issue_key" '.[$key]' $JIRA_TO_BLP_MAPPING_FILE)
  pr_link=$(jq -r --arg key "$issue_key" '.[$key]' $JIRA_TO_PR_MAPPING_FILE)

  if [ -z "$blue_print_id" ]; then
    echo -e "\e[31mNo mapping found for $issue_key\e[0m"
    continue
  fi

  if [ -z "$pr_link" ] || [ "$pr_link" = "null" ]; then
    echo -e "\e[31mNo PR link found for $issue_key\e[0m"
    continue
  fi

  pr_id=$(echo $pr_link | awk -F'/' '{print $NF}')
  pr_images_folder="$PR_IMAGES_FOLDER$pr_id"
  echo $pr_images_folder

  if [ ! -d "$pr_images_folder" ]; then
    echo -e "\e[31mNo images found for $issue_key\e[0m"
    continue
  fi

  for image in $pr_images_folder/*; do
    echo "bcli task add-file $blue_print_id -f $image"
    bcli task add-file $blue_print_id -f $image
  done
done
