JIRA_URL="https://oneline.atlassian.net"
JIRA_API_ENDPOINT="$JIRA_URL/rest/api/3/search"

PROJECT="COM"

USERNAME="xxx"
API_TOKEN="XXX"

JQL_QUERY="assignee=currentuser() AND project=COM AND createdDate>'2024-10-30' AND createdDate < '2024-11-30'"
ENCODED_JQL_QUERY=$(echo $JQL_QUERY | jq -sRr @uri)

TASK_TEMPLATE=~/.bcli/templates/COM-code.json

# Save the response to a variable
response=$(curl -u "$USERNAME:$API_TOKEN" -X GET \
"$JIRA_API_ENDPOINT?jql=$ENCODED_JQL_QUERY" \
-H "Content-Type: application/json")

# Use jq to parse the response and calculate the total number of issues
total_issues=$(echo "$response" | jq '.total')

# Print the total number of issues
echo "Total issues: $total_issues"

json_output=$(echo '{}' | jq '.')

# Read issues into an array
mapfile -t issues < <(echo "$response" | jq -c '.issues[]')

# Iterate over each issue and create blueprints task
for issue in "${issues[@]}"; do
  issue_key=$(echo "$issue" | jq -r '.key')
  issue_summary=$(echo "$issue" | jq -r '.fields.summary')
  parent_summary=$(echo "$issue" | jq -r '.fields.parent.fields.summary // "No parent summary"')
  url="$JIRA_URL/browse/$issue_key"
  encoded_url=$(echo $url | jq -sRr @uri)

  content="<p><strong>Task Ref:</strong> <a href=\"$encoded_url\">$url</a></p>
  <p><strong>US description:</strong> $parent_summary</p>
  <p><strong>Task:</strong> $issue_summary</p>"

  bluePrintId=$(bcli task create -T $TASK_TEMPLATE -t "$issue_summary" -c "$content")

  json_output=$(echo "$json_output" | jq --arg key "$issue_key" --arg url "$bluePrintId" '. + {($key): $url}')
done

echo $json_output
