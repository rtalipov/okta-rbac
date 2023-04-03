package main

import (
	"context"
	"fmt"
	"log"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var userNames []string
var userIds []string

type oktaUser struct {
	id     string
	login  string
	groups []string
}

func main() {

	orgUrl := "https://dev-94800730.okta.com"
	token := "00nBYYHhu3K_gs_5mtwudGAiElM-7CEEligKIfL37A"

	client, err := createOktaClient(orgUrl, token)
	if err != nil {
		fmt.Printf("Error creating Okta client: %v\n", err)
		return
	}

	userNames, userIds = getAllUsers(client)

	for i, userid := range userIds {
		userGroups := getUserGroups(userid, client)
		fmt.Printf("User with ID %s and name %s is in groups: %v\n", userid, userNames[i], userGroups)
	}
}

func getUserGroups(u string, client *okta.Client) []string {
	ctx := context.TODO()
	groups, resp, err := client.User.ListUserGroups(ctx, u)
	if err != nil {
		fmt.Printf("Error getting group list for user: %v\n", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Failed to retrieve groups: %v", resp.Status)
	}

	var groupNames []string

	for _, group := range groups {
		groupNames = append(groupNames, group.Profile.Name)
	}
	return groupNames
}

func getAllUsers(client *okta.Client) ([]string, []string) {
	ctx := context.TODO()
	users, resp, err := client.User.ListUsers(ctx, nil)
	if err != nil {
		fmt.Printf("Error getting all users: %v\n", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Failed to retrieve users: %v", resp.Status)
	}

	for _, user := range users {
		if login, ok := (*user.Profile)["login"].(string); ok {
			userNames = append(userNames, login)
		}
		userIds = append(userIds, user.Id)
	}
	return userNames, userIds
}

func createOktaClient(orgUrl string, token string) (*okta.Client, error) {
	ctx := context.TODO()

	ctx, client, err := okta.NewClient(
		ctx,
		okta.WithOrgUrl(orgUrl),
		okta.WithToken(token),
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	return client, nil
}

// for index, user := range users {
// 	fmt.Printf("User %d: %+v\n", index, (*user.Profile)["login"])
// }
//fixed following the article: https://devforum.okta.com/t/get-login-details-from-a-user-profile-using-go-lang-sdk/14398
