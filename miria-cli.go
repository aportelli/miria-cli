package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"syscall"

	"golang.org/x/term"
)

type AuthRequest struct {
	Db        string `json:"dbName"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	SuperUser bool   `json:"superUser"`
}

type AuthToken struct {
	Db      string `json:"dbName"`
	Expire  int    `json:"expire"`
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

func main() {
	if len(os.Args[1:]) != 1 {
		fmt.Fprintln(os.Stderr, "usage:", path.Base(os.Args[0]), "<Miria host>")
		os.Exit(1)
	}
	host := os.Args[1]

	var auth AuthToken

	cacheDir, _ := os.UserCacheDir()
	appCacheDir := cacheDir + "/miria-cli"
	authPath := appCacheDir + "/auth.json"
	_, error := os.Open(authPath)
	if error == nil {
		authj, _ := os.ReadFile(authPath)
		json.Unmarshal(authj, &auth)
		fmt.Println("Token exists", auth.Access)
	} else {
		var authr AuthRequest

		fmt.Print("Enter username: ")
		fmt.Scanln(&authr.Name)
		fmt.Print("Enter password: ")
		bytepw, _ := term.ReadPassword(int(syscall.Stdin))
		authr.Password = string(bytepw)
		authr.Db = "ADA"
		authr.SuperUser = false
		authrj, _ := json.Marshal(authr)
		response, err := http.Post("http://"+host+"/restapi/auth/token/", "application/json", bytes.NewBuffer(authrj))
		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}
		defer response.Body.Close()
		jbuf, _ := io.ReadAll(response.Body)
		json.Unmarshal(jbuf, &auth)
		fileContent, _ := json.MarshalIndent(auth, "", "  ")
		os.MkdirAll(appCacheDir, 0700)
		os.WriteFile(authPath, fileContent, 0600)
		fmt.Println("Token created", auth.Access)
	}

	var bearer = "Bearer " + auth.Access

	request, err := http.NewRequest("GET", "http://"+host+"/restapi/datamanagement/repositories/", nil)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	request.Header.Add("Authorization", bearer)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))

	search := []byte(`{
		"rootObjectPath": "archive@dp207",
  	"resultType": "INST",
		"criteria": {
			"condition": "AND",
			"rules": [
				{
					"type": "FILE_NAME",
					"value": "",
					"value2": null,
					"operator": "contains"
				}
			]
		}
	}
`)

	request, err = http.NewRequest("POST", "http://"+host+"/restapi/files/advanced-search/", bytes.NewBuffer(search))
	// request, err = http.NewRequest("POST", "http://pie.dev/post", bytes.NewBuffer(search))
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	request.Header.Add("Authorization", bearer)
	request.Header.Set("Content-Type", "application/json")
	response, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	defer response.Body.Close()

	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
}
