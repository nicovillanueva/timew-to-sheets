package main

import (
	"flag"
	"io/ioutil"
	"os"
	"text/template"
)

const (
	credentialsFile = "gsheets/creds/credentials.json"
	generatedFile   = "gsheets/credentials.go"
)

type CredentialsData struct {
	RawData string
}

func abortIf(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	var erase = flag.Bool("erase", false, "if present, it erases the credentials data from the file")
	flag.Parse()

	tpl, err := template.New("credentials").Parse(fileTemplate)
	abortIf(err)
	f, err := ioutil.ReadFile(credentialsFile)
	abortIf(err)
	var c CredentialsData
	if !*erase {
		c = CredentialsData{string(f)}
	} else {
		c = CredentialsData{""}
	}
	generated, err := os.Create(generatedFile)
	abortIf(err)
	err = tpl.Execute(generated, c)
	abortIf(err)
}

var fileTemplate = `/*
Autogenerated file, using creds/credgen.go. Do not modify.

In order to run this utility using your own Google API credentials, place
your "credentials.json" file in the "creds/" folder.
Obtain such file following these instructions:
https://developers.google.com/identity/protocols/OAuth2InstalledApp
*/

package gsheets

import (
	"io/ioutil"
)

var credentialsRaw = ` + "`{{.RawData}}`" + `

const credentialsFile string = "creds/credentials.json"

func GetCredentials() []byte {
	if b, err := ioutil.ReadFile(credentialsFile); err == nil {
		return b
	}
	return []byte(credentialsRaw)
}
`
