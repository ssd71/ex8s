package sheetutil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

var secretDir string

var sheetService *sheets.Service
var spreadsheetID string
var vRange string

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := path.Join(secretDir, "token.json")
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		// saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// Init is
func Init(sheetID string) {
	log.Println("Initialising Google Sheets Service...")
	secretDir = os.Getenv("SECRET_DIR")
	if secretDir == "" {
		secretDir = "."
		log.Println("SECRET_DIR Environment variable is empty. Trying to read secret files from current directory...")
	} else {
		log.Println("SECRET_DIR Environment found. Trying to read secret files...")
	}
	b, err := ioutil.ReadFile(path.Join(secretDir, "credentials.json"))
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	sheetService, err = sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	spreadsheetID = sheetID
	vRange = "Sheet1!A3:H"
	resp, err := sheetService.Spreadsheets.Values.Get(spreadsheetID, vRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	} else {
		log.Println("Google Sheets Service successfully initialized")
	}

	if len(resp.Values) == 0 {
		// fmt.Println("No data found.")
	} else {
		// fmt.Println("Name, Major:")
		// for _, row := range resp.Values {
		// 	// Print columns A and E, which correspond to indices 0 and 4.
		// 	fmt.Printf("%v\n", row)
		// }
	}
}

// UpdateOrInsert is
func UpdateOrInsert(data []string) error {
	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, len(data))
	for i, v := range data {
		values[0][i] = v
	}

	log.Println("Checking if received data is already present...")
	// Check if data for a date is already present
	isPresent := true
	resp, err := sheetService.Spreadsheets.Values.Get(spreadsheetID, vRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		currentRow := resp.Values[len(resp.Values)-1]
		for c := range data {
			if (c != 0) && (currentRow[c] != data[c]) {
				isPresent = false
			}
		}
	}
	rb := &sheets.ValueRange{
		Values: values,
	}

	if isPresent == false {
		log.Println("Received data distinct from current data. Pushing Updates")
		_, err = sheetService.Spreadsheets.Values.Append(spreadsheetID, vRange, rb).ValueInputOption("USER_ENTERED").Do()
		if err != nil {
			log.Fatalln("Error in pushing update to Google Sheets: ", err)
		} else {
			log.Println("Successfully pushed update")
		}
	} else {
		log.Println("Received data already present in Google Sheets")
	}
	return err
}
