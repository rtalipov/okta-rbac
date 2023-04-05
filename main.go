package main

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type oktaUser struct {
	ID        string   `json:"id"`
	Login     string   `json:"login"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Status    string   `json:"status"`
	Groups    []string `json:"groups"`
}

func main() {
	orgUrl := getEnv("OKTA_ORG_URL")
	token := getEnv("OKTA_API_TOKEN")

	fileName := flag.String("file", "", "Output file name")
	outputFormat := flag.String("output", "csv", "Output format (csv|json)")
	excludedGroups := flag.String("exclude", "Everyone", "Excluded groups from reporting")
	userQuery := flag.String("query", "", "User query options")
	flag.Parse()

	if *outputFormat == "csv" && *fileName == "" {
		*fileName = "okta_users.csv"
	} else if *outputFormat == "json" && *fileName == "" {
		*fileName = "okta_users.json"
	}

	_, client, err := createOktaClient(orgUrl, token)
	if err != nil {
		fmt.Printf("Error creating Okta client: %v\n", err)
		return
	}

	filter := query.NewQueryParams(query.WithFilter(*userQuery))
	users, err := getUsers(client, filter)
	if err != nil {
		fmt.Printf("Error getting users: %v\n", err)
		return
	}

	if *outputFormat == "csv" {
		file, err := os.Create(*fileName)
		if err != nil {
			fmt.Printf("Error creating csv file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		err = writer.Write([]string{"ID", "Login", "FistName", "LastName", "Status", "Groups"})
		if err != nil {
			fmt.Printf("Error writing header to file: %v\n", err)
			os.Exit(1)
		}

		for _, user := range users {
			user.Groups = getUserGroups(user.ID, client)
			filteredGroups := excludeGroups(user.Groups, *excludedGroups)

			row := []string{user.ID, user.Login, user.FirstName, user.LastName, user.Status}
			row = append(row, strings.Join(filteredGroups, ","))

			err := writer.Write(row)
			if err != nil {
				fmt.Printf("Error writing rows to file: %v\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Saved results to: %v\n", *fileName)
	} else if *outputFormat == "json" {
		file, err := os.Create(*fileName)
		if err != nil {
			fmt.Printf("Error creating json file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		var allUsers []oktaUser
		for _, user := range users {
			userGroups := getUserGroups(user.ID, client)
			filteredGroups := excludeGroups(userGroups, *excludedGroups)
			user.Groups = filteredGroups
			allUsers = append(allUsers, user)
		}
		jsonData, err := json.MarshalIndent(allUsers, "", "    ")
		if err != nil {
			fmt.Printf("Error identing json data: %v\n", err)
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Printf("Error writing json to file: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Saved results to: %v\n", *fileName)
	} else {
		fmt.Printf("Unsupported output format: %v\n", *outputFormat)
	}
}

func getEnv(envVar string) string {
	value, exists := os.LookupEnv(envVar)
	if !exists {
		fmt.Printf("Environmental variable is not set: %s\n", envVar)
		os.Exit(1)
	}
	return value
}

func getUserGroups(u string, client *okta.Client) []string {
	ctx := context.TODO()
	groups, resp, err := client.User.ListUserGroups(ctx, u)
	if err != nil {
		fmt.Printf("Error getting group list for user: %v\n", err)
		return []string{}
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Failed to retrieve groups: %v", resp.Status)
		return []string{}
	}

	var groupNames []string

	for _, group := range groups {
		groupNames = append(groupNames, group.Profile.Name)
	}
	return groupNames
}

func getUsers(client *okta.Client, filter *query.Params) ([]oktaUser, error) {
	ctx := context.TODO()
	users, resp, err := client.User.ListUsers(ctx, filter)
	if err != nil {
		fmt.Printf("Error getting all users: %v\n", err)
	}

	if resp.StatusCode != 200 {
		fmt.Printf("Failed to retrieve users: %v", resp.Status)
	}

	var oktaUsers []oktaUser

	for _, user := range users {
		if login, ok := (*user.Profile)["login"].(string); ok {
			oktaUser := oktaUser{
				ID:        user.Id,
				Login:     login,
				FirstName: (*user.Profile)["firstName"].(string),
				LastName:  (*user.Profile)["lastName"].(string),
				Status:    user.Status,
			}
			oktaUsers = append(oktaUsers, oktaUser)
		}
	}
	return oktaUsers, err
}

func excludeGroups(allGroups []string, excludedGroups string) []string {

	filteredGroups := []string{}
	for _, group := range allGroups {
		excluded := false
		for _, excludedGroup := range strings.Split(excludedGroups, ",") {
			if group == excludedGroup {
				excluded = true
				break
			}
		}
		if !excluded {
			filteredGroups = append(filteredGroups, group)
		}
	}
	return filteredGroups
}

func createOktaClient(orgUrl string, token string) (ctx context.Context, client *okta.Client, err error) {
	ctx, client, err = okta.NewClient(
		context.TODO(),
		okta.WithOrgUrl(orgUrl),
		okta.WithToken(token),
	)

	if err != nil {
		fmt.Printf("Error creating Okta client: %v\n", err)
		return nil, nil, err
	}

	return ctx, client, err
}
