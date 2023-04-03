package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/okta/okta-sdk-golang/v2/okta"
)

var userNames []string
var userIds []string

type oktaUser struct {
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	ID        string   `json:"id"`
	Login     string   `json:"login"`
	Status    string   `json:"status"`
	Groups    []string `json:"groups"`
}

func main() {

	orgUrl := "https://dev-94800730.okta.com"
	token := "00nBYYHhu3K_gs_5mtwudGAiElM-7CEEligKIfL37A"

	client, err := createOktaClient(orgUrl, token)
	if err != nil {
		fmt.Printf("Error creating Okta client: %v\n", err)
		return
	}

	users, err := getAllUsers(client)

	for _, user := range users {
		user.Groups = getUserGroups(user.ID, client)
		//fmt.Printf("User with ID %s and name %s is in groups: %v\n", user.id, user.login, user.groups)
		user.print()
	}
}

func (u oktaUser) print() {
	jsonData, err := json.MarshalIndent(u, "", "    ")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Printf("%+v\n", jsonData)
	fmt.Println(string(jsonData))
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

func getAllUsers(client *okta.Client) ([]oktaUser, error) {
	ctx := context.TODO()
	users, resp, err := client.User.ListUsers(ctx, nil)
	if err != nil {
		fmt.Printf("Error getting all users: %v\n", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Failed to retrieve users: %v", resp.Status)
	}

	var oktaUsers []oktaUser

	for _, user := range users {
		if login, ok := (*user.Profile)["login"].(string); ok {
			oktaUser := oktaUser{
				FirstName: (*user.Profile)["firstName"].(string),
				LastName:  (*user.Profile)["lastName"].(string),
				ID:        user.Id,
				Login:     login,
				Status:    user.Status,
			}
			oktaUsers = append(oktaUsers, oktaUser)
		}
	}
	return oktaUsers, nil
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
