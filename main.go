package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please provide path to the JSON export as argument.")
	}

	p := expandTilde(os.Args[1])
	downloadDir := filepath.Join(filepath.Dir(p), "attachments")
	os.Mkdir(downloadDir, 0777)

	jsonContent, err := ioutil.ReadFile(p)
	if err != nil {
		log.Fatal(err)
	}

	var trelloExport TrelloJSON
	if err := json.Unmarshal(jsonContent, &trelloExport); err != nil {
		panic(err)
	}

	client := grab.NewClient()

	for _, card := range trelloExport.Cards {
		for i, attach := range card.Attachments {
			if !attach.IsUpload {
				continue
			}
			log.Println("downloading", attach.ID, "of", card.ID)
			ddir := filepath.Join(downloadDir, card.ID)
			os.Mkdir(ddir, 0777)
			req, _ := grab.NewRequest(ddir, attach.URL)
			// add attachment number to the filename to resolve conflicts, yes it is lazy
			patchedFilename := fmt.Sprintf("%d-%s", i, filepath.Base(attach.URL))
			req.Filename = filepath.Join(ddir, patchedFilename)
			resp := client.Do(req)
			log.Println("finished download")
			if err := resp.Err(); err != nil {
				log.Println("Could not download", attach.URL)
				log.Fatalln(err)
			} else {
				// log.Println("Download saved to", resp.Filename)
			}
		}
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func expandTilde(path string) string {
	if strings.HasPrefix(path, "~/") {
		usr, err := user.Current()
		if err == nil {
			return filepath.Join(usr.HomeDir, path[2:])
		}
	}
	return path
}

type TrelloJSON struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Desc           string      `json:"desc"`
	DescData       interface{} `json:"descData"`
	Closed         bool        `json:"closed"`
	IDOrganization interface{} `json:"idOrganization"`
	Cards          []struct {
		ID                string        `json:"id"`
		Address           interface{}   `json:"address"`
		CheckItemStates   interface{}   `json:"checkItemStates"`
		Closed            bool          `json:"closed"`
		Coordinates       interface{}   `json:"coordinates"`
		CreationMethod    interface{}   `json:"creationMethod"`
		DateLastActivity  time.Time     `json:"dateLastActivity"`
		Desc              string        `json:"desc"`
		DescData          interface{}   `json:"descData"`
		DueReminder       interface{}   `json:"dueReminder"`
		IDBoard           string        `json:"idBoard"`
		IDLabels          []interface{} `json:"idLabels"`
		IDList            string        `json:"idList"`
		IDMembersVoted    []interface{} `json:"idMembersVoted"`
		IDShort           int           `json:"idShort"`
		IDAttachmentCover string        `json:"idAttachmentCover"`
		Limits            struct {
			Attachments struct {
				PerCard struct {
					Status    string `json:"status"`
					DisableAt int    `json:"disableAt"`
					WarnAt    int    `json:"warnAt"`
				} `json:"perCard"`
			} `json:"attachments"`
			Checklists struct {
				PerCard struct {
					Status    string `json:"status"`
					DisableAt int    `json:"disableAt"`
					WarnAt    int    `json:"warnAt"`
				} `json:"perCard"`
			} `json:"checklists"`
			Stickers struct {
				PerCard struct {
					Status    string `json:"status"`
					DisableAt int    `json:"disableAt"`
					WarnAt    int    `json:"warnAt"`
				} `json:"perCard"`
			} `json:"stickers"`
		} `json:"limits"`
		LocationName          interface{} `json:"locationName"`
		ManualCoverAttachment bool        `json:"manualCoverAttachment"`
		Name                  string      `json:"name"`
		Pos                   float64     `json:"pos"`
		ShortLink             string      `json:"shortLink"`
		Badges                struct {
			AttachmentsByType struct {
				Trello struct {
					Board int `json:"board"`
					Card  int `json:"card"`
				} `json:"trello"`
			} `json:"attachmentsByType"`
			Location           bool        `json:"location"`
			Votes              int         `json:"votes"`
			ViewingMemberVoted bool        `json:"viewingMemberVoted"`
			Subscribed         bool        `json:"subscribed"`
			Fogbugz            string      `json:"fogbugz"`
			CheckItems         int         `json:"checkItems"`
			CheckItemsChecked  int         `json:"checkItemsChecked"`
			Comments           int         `json:"comments"`
			Attachments        int         `json:"attachments"`
			Description        bool        `json:"description"`
			Due                interface{} `json:"due"`
			DueComplete        bool        `json:"dueComplete"`
		} `json:"badges"`
		DueComplete  bool          `json:"dueComplete"`
		Due          interface{}   `json:"due"`
		Email        string        `json:"email"`
		IDChecklists []interface{} `json:"idChecklists"`
		IDMembers    []interface{} `json:"idMembers"`
		Labels       []interface{} `json:"labels"`
		ShortURL     string        `json:"shortUrl"`
		Subscribed   bool          `json:"subscribed"`
		URL          string        `json:"url"`
		Attachments  []struct {
			Bytes     int         `json:"bytes"`
			Date      time.Time   `json:"date"`
			EdgeColor string      `json:"edgeColor"`
			IDMember  string      `json:"idMember"`
			IsUpload  bool        `json:"isUpload"`
			MimeType  interface{} `json:"mimeType"`
			Name      string      `json:"name"`
			Previews  []struct {
				URL    string `json:"url"`
				Bytes  int    `json:"bytes"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
				ID     string `json:"_id"`
				Scaled bool   `json:"scaled"`
			} `json:"previews"`
			URL string `json:"url"`
			Pos int    `json:"pos"`
			ID  string `json:"id"`
		} `json:"attachments"`
		PluginData       []interface{} `json:"pluginData"`
		CustomFieldItems []interface{} `json:"customFieldItems"`
	} `json:"cards"`
}
