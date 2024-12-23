#!/bin/bash
JIRA_URL="https://oneline.atlassian.net"
JIRA_API_ENDPOINT="$JIRA_URL/rest/api/3/search"

PROJECT="COM"
source ./_env.sh

JQL_QUERY="assignee=currentuser() AND project=COM AND createdDate>'$START_DATE' AND createdDate < '$END_DATE'"

ENCODED_JQL_QUERY=$(echo $JQL_QUERY | jq -sRr @uri)

MAPPING_FILE=/home/tieu/.bcli/files/jira_to_blp.json
ADD_TIME_TEMPLATE=/home/tieu/.bcli/templates/add-time-dev.json
ADD_POINT_TEMPLATE=/home/tieu/.bcli/templates/add-point-dev.json

# Save the response to a variable
response=$(curl -u "$USERNAME:$API_TOKEN" -X GET \
"$JIRA_API_ENDPOINT?jql=$ENCODED_JQL_QUERY" \
-H "Content-Type: application/json")

# Read issues into an array
mapfile -t issues < <(echo "$response" | jq -c '.issues[]')

for issue in "${issues[@]}"; do
  issue_key=$(echo "$issue" | jq -r '.key')

  time_tracking=$(echo "$issue" | jq -r '.fields.timespent')

  updated=$(echo "$issue" | jq -r '.fields.updated')
  updated=$(date -d "${updated}" +"%Y%m%d")
  echo $updated

  blue_print_id=$(jq -r --arg key "$issue_key" '.[$key]' $MAPPING_FILE)

  if [ -z "$blue_print_id" ]; then
    echo "No mapping found for $issue_key"
    continue
  fi

  hours=$((time_tracking / 3600))
  hours=$((hours + 2)) # Add 2 hour
  echo "bcli task add-time $blue_print_id -H $hours -d $updated -T $ADD_TIME_TEMPLATE"
  bcli task add-time $blue_print_id -H $hours -d $updated -T $ADD_TIME_TEMPLATE

  points=$((hours * 30))
  volume=$((points / 50))
  echo "bcli task add-point $blue_print_id -v $volume -T $ADD_POINT_TEMPLATE"
  bcli task add-point $blue_print_id -v $volume -T $ADD_POINT_TEMPLATE
done
