package gsheets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	tokenFilename = "./token.json"
)

type Header []interface{}

var TableHeader Header = []interface{}{"Start", "End", "Delta", "Tags"}

func readToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenFilename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(&token)
	return token, err
}

func createToken(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline) // Token type
	log.Printf("Auth URL: %s", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("could not read auth code: %v", err)
	}
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("could not retrieve new token: %v", err)
	}
	return token
}

func saveToken(token *oauth2.Token) {
	log.Printf("saving token")
	f, err := os.OpenFile(tokenFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("cannot write %s", tokenFilename)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("could not encode token: %v", err)
	}
}

func getClient() *http.Client {
	creds, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("could not read credentials.json: %v", err)
	}
	conf, err := google.ConfigFromJSON(creds, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("could not fetch config: %v", err)
	}

	t, err := readToken()
	if err != nil {
		log.Printf("could not read token.json; creating new one")
		t = createToken(conf)
		saveToken(t)
	} else {
		log.Printf("reusing token.json")
	}
	return conf.Client(context.Background(), t)
}

func GetSheetsService() *sheets.Service {
	srv, err := sheets.New(getClient())
	if err != nil {
		log.Fatalf("could not create service: %v", err)
	}
	return srv
}

func DoBatchUpdate(spreadsheetID string, requests []*sheets.Request) error {
	_, err := GetSheetsService().Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	return err
}

func addNewSheets(srv *sheets.Service, spreadsheetID string, names []string) error {
	requests := make([]*sheets.Request, len(names))
	for i, n := range names {
		requests[i] = &sheets.Request{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: n,
				},
			},
		}
	}

	_, err := srv.Spreadsheets.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return err
	}
	return nil
}
