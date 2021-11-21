package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type Token struct {
	Email        string
	PrivateKey   string
	PrivateKeyID string
	TokenURL     string
	Scopes       []string
}
type SpreadSheet struct {
	RawValues     map[string][]RawValue
	spreadsheetId string
	valueRange    string
	headerRange   string
	tokenFileName string
	titleCache    map[string]struct{}
	srv           *sheets.Service
	spreadSheet   *sheets.Spreadsheet
}

type RawValue struct {
	Id      uuid.UUID
	Content string
	Author  string
	Tags    []string
}

func Intersection(a []RawValue, b []RawValue) []RawValue {
	var ans []RawValue
	if len(a) == 0 || len(b) == 0 {
		return ans
	}
	hash := make(map[uuid.UUID]bool)
	for _, e := range a {
		hash[e.Id] = true
	}
	for _, e := range b {
		if hash[e.Id] {
			ans = append(ans, e)
		}
	}
	return ans
}

var ssInstance *SpreadSheet

func SpreadSheetInstance() *SpreadSheet {
	if ssInstance == nil {
		ssInstance = newSpreadSheet()
	}
	return ssInstance
}

func (s *SpreadSheet) TryAdd(value RawValue, guild string) bool {
	sheetTitle := fmt.Sprint(guild)
	if !s.checkSheetExist(sheetTitle) {
		if !s.tryAddNewSheet(sheetTitle) {
			return false
		}
	}
	content := sheets.ValueRange{}
	list := make([][]interface{}, 1)
	list[0] = make([]interface{}, 0)
	list[0] = append(list[0], value.Id)
	list[0] = append(list[0], value.Author)
	list[0] = append(list[0], value.Content)
	for _, tag := range value.Tags {
		list[0] = append(list[0], tag)
	}

	content.Values = list
	req := s.srv.Spreadsheets.Values.Append(s.spreadsheetId, s.getRange(guild), &content)
	req.ValueInputOption("RAW")
	_, err := req.Do()
	if err == nil {
		DatabaseInstance().AddValue(&value, guild)
	} else {
		log.Fatal(err)
	}
	return err == nil
}

func (s *SpreadSheet) Init() {
	var token Token
	if tokenString := os.Getenv("TOKEN"); tokenString != "" {
		err := json.Unmarshal([]byte(tokenString), &token)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("Init spreadsheet token from env")
	} else if jsonFile, err := os.Open(s.tokenFileName); err == nil {
		tokenString, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(tokenString, &token)
		fmt.Println("Init spreadsheet token from json file")
	}

	// Create a JWT configurations object for the Google service account
	conf := &jwt.Config{
		Email:        token.Email,
		PrivateKey:   []byte(token.PrivateKey),
		PrivateKeyID: token.PrivateKeyID,
		TokenURL:     token.TokenURL,
		Scopes:       token.Scopes,
	}

	client := conf.Client(oauth2.NoContext)

	// Create a service object for Google sheets
	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	s.srv = srv

	s.UpdateAllValueFromRemote()
}

func (s *SpreadSheet) UpdateAllValueFromRemote() {
	for k := range s.RawValues {
		delete(s.RawValues, k)
	}
	for k := range s.titleCache {
		delete(s.titleCache, k)
	}

	req := s.srv.Spreadsheets.Get(s.spreadsheetId)
	spreadSheet, err := req.Do()
	if err != nil {
		log.Fatal(err)
	}
	s.spreadSheet = spreadSheet
	s.updateLocalCache()
}

func newSpreadSheet() *SpreadSheet {
	s := SpreadSheet{}
	s.spreadsheetId = "1ONaFSfaRiLkatqb_hMh8d0PEZkREwYcEXpO-AJ1XRWs"
	s.valueRange = "A2:ZZZ"
	s.headerRange = "A:Z"
	s.tokenFileName = "token.json"
	s.RawValues = make(map[string][]RawValue)
	s.titleCache = make(map[string]struct{})

	return &s
}

func (s *SpreadSheet) updateLocalCache() {
	for _, sheet := range s.spreadSheet.Sheets {
		title := sheet.Properties.Title
		fmt.Printf("Try to get value from remote for guild %s\n", title)
		if _, isSheetExist := s.RawValues[title]; !isSheetExist {
			s.RawValues[title] = make([]RawValue, 0)
		}
		s.titleCache[title] = struct{}{}
		req := s.srv.Spreadsheets.Values.Get(s.spreadsheetId, s.getRange(title))
		resp, err := req.Do()
		if err != nil {
			continue
		}
		for _, v := range resp.Values {
			tags := []string{}
			for i := 3; i < len(v); i++ {
				tags = append(tags, fmt.Sprintf("%v", v[i]))
			}
			id := uuid.MustParse(fmt.Sprintf("%v", v[0]))
			author := fmt.Sprintf("%v", v[1])
			content := fmt.Sprintf("%v", v[2])
			raw := RawValue{
				Id:      id,
				Author:  author,
				Content: content,
				Tags:    tags,
			}
			s.RawValues[title] = append(s.RawValues[title], raw)
		}
	}
}

func (s *SpreadSheet) getRange(guild string) string {

	return fmt.Sprintf("%s!%s", guild, s.valueRange)
}

func (s *SpreadSheet) tryAddNewSheet(title string) bool {
	req := sheets.Request{
		AddSheet: &sheets.AddSheetRequest{Properties: &sheets.SheetProperties{Title: title}},
	}
	batchReq := &sheets.BatchUpdateSpreadsheetRequest{Requests: make([]*sheets.Request, 0)}
	batchReq.Requests = append(batchReq.Requests, &req)

	resp, err := s.srv.Spreadsheets.BatchUpdate(s.spreadsheetId, batchReq).Do()
	if err != nil {
		log.Fatal(err)
	}
	for _, reply := range resp.Replies {
		if reply.AddSheet != nil {
			s.titleCache[title] = struct{}{}
			s.addHeader(title)
			return true
		}
	}
	return err == nil
}

func (s *SpreadSheet) addHeader(title string) {
	content := sheets.ValueRange{}
	list := make([][]interface{}, 1)
	list[0] = make([]interface{}, 4)
	list[0][0] = "Id"
	list[0][1] = "Author"
	list[0][2] = "Content"
	list[0][3] = "Tags"
	content.Values = list
	req := s.srv.Spreadsheets.Values.Append(s.spreadsheetId, fmt.Sprintf("%s!%s", title, s.headerRange), &content)
	req.ValueInputOption("RAW")
	req.Do()
}

func (s *SpreadSheet) checkSheetExist(title string) bool {
	_, ok := s.titleCache[title]
	return ok
}
