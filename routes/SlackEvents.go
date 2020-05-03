package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"slack-bot/bot"

	"github.com/gin-gonic/gin"
	"github.com/slack-go/slack/slackevents"
)

func SlackEventHandler(c *gin.Context) {

	buf := new(bytes.Buffer)
	buf.ReadFrom(c.Request.Body)
	body := buf.String()

	// Verify the request came from slack
	eventsAPIEvent, e := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: os.Getenv("VERIFICATION_TOKEN")}))
	if e != nil {
		fmt.Println(e.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
		return
	}

	// Verify event URL when setting up bot
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, r.Challenge)
		return
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			bot.HandleMentionEvent(ev)
		}
	}

	c.JSON(200, "OK")
}
