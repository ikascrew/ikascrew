package video

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

func init() {
}

type Twitter struct {
	*Phrase
	word string
}

type twitterConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	Word              string
}

var app *twitterConfig

func loadConfig() error {

	raw, err := ioutil.ReadFile("twitter.json")
	if err != nil {
		return err
	}
	app = &twitterConfig{}

	err = json.Unmarshal(raw, app)
	if err != nil {
		return err
	}
	return nil
}

func (t *Twitter) search() ([]string, error) {

	config := oauth1.NewConfig(app.ConsumerKey, app.ConsumerSecret)
	token := oauth1.NewToken(app.AccessToken, app.AccessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	param := &twitter.SearchTweetParams{
		Query: t.word,
	}

	param.Count = 100
	/*
		if app.MaxId != 0 {
			param.SinceID = app.MaxId
		}
	*/
	// Search Tweets
	search, _, err := client.Search.Tweets(param)
	if err != nil {
		return nil, err
	}

	//app.MaxId = search.Metadata.MaxID

	if len(search.Statuses) <= 0 {
		return nil, nil
	}

	rtn := make([]string, 0)
	for _, tw := range search.Statuses {
		buf := "@" + tw.User.ScreenName + " " + tw.Text
		rtn = append(rtn, buf)
	}
	fmt.Println(rtn)
	return rtn, nil
}

func NewTwitter() (*Twitter, error) {

	err := loadConfig()
	if err != nil {
		return nil, err
	}

	t := Twitter{}

	t.word = app.Word
	data, err := t.search()
	if err != nil {
		return nil, err
	}

	p, err := NewPhrase(data)
	if err != nil {
		return nil, err
	}

	t.Phrase = p
	return &t, nil
}

func (t *Twitter) initialize() {

	t.Phrase.initialize()
	return
}

func (v *Twitter) Source() string {
	return "_ikascrew_Twitter.mp4"
}
