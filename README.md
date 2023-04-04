# Okta RBAC Report Generator

This CLI tool generates a report of users in Okta organization with their group memebership.
It supports optional filtering using okta user query and exclusion of certain groups.
The report can be generated in either CSV of JSON format

## Usage

To use this program, you need to provide your Okta organization URL and API token as environment variables:

```
export OKTA_ORG_URL=<https://your-okta-domain.com>
export OKTA_API_TOKEN=your-api-token
```

Run the program with the decired command-line options

To output in JSON:

```
okta-rbac -o json -e group1,group2 -q "status eq \"ACTIVE\""
```

To output in CSV:

```
okta-rbac -f report.csv -e group1,group2 -q "status eq \"ACTIVE\""
```

Available options:

- `-e`: a comma-separated list of excluded groups from the report (default is "Everyone")
- `-f`: the name of output csv file (default is "okta_users.csv")
- `-o`: output format CSV or JSON (default is "csv")
- `-q`: User query options in Okta Query Language

## License

This program is licensed under the MIT License.
