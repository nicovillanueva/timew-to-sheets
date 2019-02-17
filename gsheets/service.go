package gsheets

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"github.com/mattn/go-tty"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

var (
	tokenPath = "~/.timewarrior/to-sheets/token.json"
)

func init() {
	var err error
	p, err := homedir.Expand(tokenPath[:strings.LastIndex(tokenPath, "/")])
	if err != nil {
		log.Fatalf("could not pull home directory: %v", err)
	}
	if _, err := os.Stat(p); os.IsNotExist(err) {
		log.Debugf("creating folder %s", p)
		if err := os.MkdirAll(p, 0755); err != nil {
			log.Fatalf("could not create directory %s: %v", p, err)
		}
	}
	tokenPath, err = homedir.Expand(tokenPath)
	if err != nil {
		log.Fatalf("cannot determine token location: %v", err)
	}
}

func readToken() (*oauth2.Token, error) {
	f, err := os.Open(tokenPath)
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
	log.Infof("\nOpen this URL in your browser and paste the generated code here:\n%s\n", authURL)
	inputTty, err := tty.Open()
	if err != nil {
		log.Fatalf("could not open new TTY to input auth code")
	}
	defer inputTty.Close()
	authCode, err := inputTty.ReadString()
	if err != nil {
		log.Fatalf("could not read auth code from new inputTty")
	}
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("could not retrieve new token: %v", err)
	}
	return token
}

func saveToken(token *oauth2.Token) {
	log.Info("saving token to ", tokenPath)
	f, err := os.OpenFile(tokenPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal("cannot write", tokenPath)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Fatalf("could not encode token: %v", err)
	}
}

func getClient() *http.Client {
	conf, err := google.ConfigFromJSON(GetCredentials(), "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("could not build API config: %v", err)
	}

	t, err := readToken()
	if err != nil {
		log.Infof("could not read token.json; creating new one")
		t = createToken(conf)
		saveToken(t)
	} else {
		log.Debug("reusing token.json")
	}
	return conf.Client(context.Background(), t)
}

// GetSheetsService spawns a new client to GSheets
func GetSheetsService() *sheets.Service {
	srv, err := sheets.New(getClient())
	if err != nil {
		log.Fatalf("could not create service: %v", err)
	}
	return srv
}

// DoBatchUpdate is a convenience function. Spawns a session to GSheets and performs all `requests`
// If a AddSheet request is found to try to add a pre-existing sheet, it is removed from the batch.
// TODO: Fill the `.Fields()` method to save on traffic
func DoBatchUpdate(spreadsheetID string, requests []*sheets.Request) error {
	srv := GetSheetsService().Spreadsheets
	currentStatus, err := srv.Get(spreadsheetID).Fields().Do()
	if err != nil {
		log.Fatalf("could not pull current spreadsheet status: %v", err)
	}
	requests = cancelPreExistingSheets(requests, currentStatus.Sheets)
	log.Infof("will dispatch %d requests to google API", len(requests))
	_, err = srv.BatchUpdate(spreadsheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	return err
}

// cancelPreExistingSheets filters out the existing sheet IDs from a list of candidate requests, preserving order.
// Currently performs two allocations: one of the same size of the input candidates, which results in <nil> entries when
// a candidate is found to exist already; and another one of the size of the candidates minus the nil entries.
// The latter is returned after removing nil entries.
func cancelPreExistingSheets(candidates []*sheets.Request, existing []*sheets.Sheet) []*sheets.Request {
	confirmed := make([]*sheets.Request, len(candidates))
	var skip = false
	var lastIndex int
	for _, req := range candidates {
		if req.AddSheet != nil && req.AddSheet.Properties != nil {
			for _, cS := range existing {
				if req.AddSheet.Properties.SheetId == cS.Properties.SheetId {
					skip = true
					break
				}
			}
			if !skip {
				confirmed[lastIndex] = req
				lastIndex++
			}
		} else {
			confirmed[lastIndex] = req
			lastIndex++
		}
		skip = false
	}
	n := make([]*sheets.Request, lastIndex)
	for i, r := range confirmed {
		if r == nil {
			break
		}
		n[i] = r
	}
	return n
}
