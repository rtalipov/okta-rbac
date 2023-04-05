# Okta RBAC Report Generator

This CLI tool generates a report of users in Okta organization with their group memebership.
It supports optional filtering using okta user query and exclusion of certain groups.
The report can be generated in either CSV of JSON format

## Usage

To use this program, you need to provide your Okta organization URL and API token as environment variables:

```text
export OKTA_ORG_URL=<https://your-okta-domain.com>
export OKTA_API_TOKEN=your-api-token
```

Run the program with the decired command-line options

To output in JSON:

```text
okta-rbac -output json -exclude group1,group2 -query "status eq \"ACTIVE\""
```

To output in CSV:

```text
okta-rbac -file report.csv -exclude group1,group2 -query "status eq \"ACTIVE\""
```

Available options:

- `-exclude`: a comma-separated list of excluded groups from the report (default is "Everyone")
- `-file`: the name of output csv file (default is "okta_users.csv")
- `-output`: output format CSV or JSON (default is "csv")
- `-query`: User query options in Okta Query Language

## License

This program is licensed under the MIT License.
