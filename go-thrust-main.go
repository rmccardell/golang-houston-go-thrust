package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ddliu/go-httpclient"
	"github.com/miketheprogrammer/go-thrust/lib/bindings/window"
	"github.com/miketheprogrammer/go-thrust/lib/commands"
	"github.com/miketheprogrammer/go-thrust/thrust"
	"github.com/miketheprogrammer/go-thrust/tutorials/provisioner"
	"github.com/user/gouis/go-thrust/go-thrust-meetup/models"
)

type authInfo struct {
	consumerKey    string
	consumerSecret string
	redirectURI    string
	groupID        string
	groupURLName   string
}

//Todo: update with your own meetup consumer api values
var meetUpAuth = authInfo{
	consumerKey:    "Your MeetUp Consumer Key Goes Here",
	consumerSecret: "Your MeetUp Secret",
	redirectURI:    "http://localhost:8080/auth",
	groupID:        "18803286",       // <-- can be changed to any valid MeetUp group id
	groupURLName:   "Golang-Houston", // <-- can be changed to any valid MeetUp group url name
}

var authExpiration time.Time
var meetupTokens models.MeetUpToken

func isLoggedIn() bool {

	return (time.Now().Before(authExpiration))

}

func handler(w http.ResponseWriter, r *http.Request) {
	var meetupAuthURL = "https://secure.meetup.com/oauth2/authorize" +
		"?client_id=" + meetUpAuth.consumerKey +
		"&response_type=code" +
		"&redirect_uri=" + meetUpAuth.redirectURI

	http.Redirect(w, r, meetupAuthURL, 301)
}

func authHandler(w http.ResponseWriter, r *http.Request) {

	m, _ := url.ParseQuery(r.URL.RawQuery)
	code := strings.Join(m["code"], "")

	data := url.Values{}
	data.Set("client_id", meetUpAuth.consumerKey)
	data.Add("client_secret", meetUpAuth.consumerSecret)
	data.Add("grant_type", "authorization_code")
	data.Add("redirect_uri", meetUpAuth.redirectURI)
	data.Add("code", code)

	client := &http.Client{}
	r, _ = http.NewRequest("POST", "https://secure.meetup.com/oauth2/access", bytes.NewBufferString(data.Encode())) // <-- URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	body, _ := ioutil.ReadAll(resp.Body)

	err := json.Unmarshal(body, &meetupTokens)

	if err != nil {

		fmt.Println(err)

	} else {

		//we're in like flynn!
		authExpiration = time.Now().Add(time.Second * 3600)
		http.Redirect(w, r, "http://localhost:8080/view", 301)

	}

}

func viewHandler(w http.ResponseWriter, r *http.Request) {

	html, err := ioutil.ReadFile("index.html")

	if err != nil {

		fmt.Println("error loading view..")

	} else {

		view := string(html)

		fmt.Fprint(w, view)
	}

}

//make get request using go-httpclient

// https://github.com/ddliu/go-httpclient
func meetUpAPIGetRequest(uri string, params map[string]string) string {

	params["access_token"] = meetupTokens.AccessToken

	res, err := httpclient.Get(uri, params)

	log.Println(res.StatusCode, err)

	content, _ := ioutil.ReadAll(res.Body)

	log.Println("content", content)

	return string(content)

}

//get golang houston data from meetup api
func getGroupData() string {

	params := map[string]string{}
	params["photo-host"] = "public"
	params["group_urlname"] = meetUpAuth.groupURLName
	params["country"] = "us"
	params["city"] = "houston"
	params["state"] = "tx"
	params["page"] = "20"

	uri := "https://api.meetup.com/2/groups"

	return meetUpAPIGetRequest(uri, params)

}

func getMemberData() string {

	params := map[string]string{}
	params["photo-host"] = "public"
	params["group_id"] = "18803286"
	params["page"] = "100"

	uri := "https://api.meetup.com/2/members"

	return meetUpAPIGetRequest(uri, params)
}

func getEventData() string {

	//https://api.meetup.com/2/events?&photo-host=public&group_id=18803286&page=20

	params := map[string]string{}
	params["photo-host"] = "public"
	params["group_id"] = "18803286"
	params["page"] = "100"

	uri := "https://api.meetup.com/2/events"

	return meetUpAPIGetRequest(uri, params)
}

func main() {

	//redirect logger output to console
	log.SetOutput(os.Stdout)

	//initialize go-httpclient
	httpclient.Defaults(httpclient.Map{
		httpclient.OPT_USERAGENT: "meetup ui client",
		"Accept-Language":        "en-us",
	})

	//define some mux handlers
	http.HandleFunc("/", handler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/view", viewHandler)

	//thrust code starts here
	thrust.InitLogger()
	// Set any Custom Provisioners before Start
	thrust.SetProvisioner(tutorial.NewTutorialProvisioner())

	//must always be called before bindings are created
	thrust.Start()

	thrustWindow := thrust.NewWindow(thrust.WindowOptions{
		RootUrl: "http://localhost:8080/",
	})
	thrustWindow.Show()
	thrustWindow.Maximize()
	thrustWindow.Focus()
	thrustWindow.OpenDevtools()

	//capture remote messages
	_, err := thrustWindow.HandleRemote(func(er commands.EventResult, this *window.Window) {

		//fmt.Println("RemoteMessage Recieved:", er.Message.Payload)

		log.Printf("Remote message Recieved:", er.Message.Payload)

		if er.Message.Payload == "getGroup" {

			//this.SendRemoteMessage("preparing data")
			data := getGroupData()
			log.Printf("data", data)

			msg := "getGroup|" + data
			this.SendRemoteMessage(msg)
		}

		if er.Message.Payload == "getMembers" {

			//this.SendRemoteMessage("preparing data")
			data := getMemberData()
			log.Printf("data", data)

			msg := "getMembers|" + data
			this.SendRemoteMessage(msg)
		}

		if er.Message.Payload == "getEvents" {

			//this.SendRemoteMessage("preparing data")
			data := getEventData()
			log.Printf("data", data)

			msg := "getEvents|" + data
			this.SendRemoteMessage(msg)
		}

	})
	if err != nil {

		log.Printf("remote message error:", err)
	}

	//start go server
	http.ListenAndServe(":8080", nil)
}
