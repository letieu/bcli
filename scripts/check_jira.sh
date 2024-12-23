#!/bin/bash
JIRA_URL="https://oneline.atlassian.net"
JIRA_API_ENDPOINT="$JIRA_URL/rest/api/3/search"

PROJECT="COM"
source _env.sh

JQL_QUERY="assignee=currentuser() AND project=COM AND createdDate>'$START_DATE' AND createdDate < '$END_DATE'"
ENCODED_JQL_QUERY=$(echo $JQL_QUERY | jq -sRr @uri)

# Save the response to a variable
response=$(curl -u "$USERNAME:$API_TOKEN" -X GET \
"$JIRA_API_ENDPOINT?jql=$ENCODED_JQL_QUERY" \
-H "Content-Type: application/json")

# Use jq to parse the response and calculate the total number of issues
total_issues=$(echo "$response" | jq '.total')

# Print the total number of issues
echo "Total issues: $total_issues"

# Read issues into an array
mapfile -t issues < <(echo "$response" | jq -c '.issues[]')

# Iterate over each issue and print details
for issue in "${issues[@]}"; do
    issue_key=$(echo "$issue" | jq -r '.key')
    issue_summary=$(echo "$issue" | jq -r '.fields.summary')
    parent_summary=$(echo "$issue" | jq -r '.fields.parent.fields.summary // "No parent summary"')
    url="$JIRA_URL/browse/$issue_key"

    echo "[${issue_key}] $issue_summary"
    echo "Parent Summary: $parent_summary"
    echo "URL: $url"
    echo "-----------------------------"
done
