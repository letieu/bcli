#!/bin/bash
JIRA_URL="https://oneline.atlassian.net"
JIRA_API_ENDPOINT="$JIRA_URL/rest/api/3/search"

PROJECT="COM"
source ./_env.sh

JQL_QUERY="assignee=currentuser() AND project=COM AND createdDate>'$START_DATE' AND createdDate < '$END_DATE'"

ENCODED_JQL_QUERY=$(echo $JQL_QUERY | jq -sRr @uri)

# Save the response to a variable
response=$(curl -u "$USERNAME:$API_TOKEN" -X GET \
"$JIRA_API_ENDPOINT?jql=$ENCODED_JQL_QUERY" \
-H "Content-Type: application/json")

json_output=$(echo '{}' | jq '.')

# Read issues into an array
mapfile -t issues < <(echo "$response" | jq -c '.issues[]')

for issue in "${issues[@]}"; do
    PR_RESPONSE=$(curl -u "$USERNAME:$API_TOKEN" -X GET \
    "$JIRA_URL/rest/dev-status/1.0/issue/details?issueId=$(echo $issue | jq -r '.id')&applicationType=github&dataType=pullrequest" \
    -H "Content-Type: application/json")

    first_detail=$(echo "$PR_RESPONSE" | jq -c '.detail[0]')

    if [ "$first_detail" != "null" ]; then
        list_PR=$(echo $first_detail | jq -r '.pullRequests[] | select(.status == "MERGED" or .status == "OPEN") | .url')
        newest_PR=$(echo "$list_PR" | head -n 1)

        issue_key=$(echo $issue | jq -r '.key')

        if [ -n "$newest_PR" ]; then
            json_output=$(echo "$json_output" | jq --arg key "$issue_key" --arg url "$newest_PR" '. + {($key): $url}')
        fi
    fi
done

echo $json_output
