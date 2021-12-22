package gcp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/akubi0w1/chatbot-202201/code"
	"github.com/akubi0w1/chatbot-202201/config"
	"github.com/akubi0w1/chatbot-202201/external/handler"
	"github.com/akubi0w1/chatbot-202201/external/response"
	cslack "github.com/akubi0w1/chatbot-202201/external/slack"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

type Application struct {
	slack  *cslack.SlackClient
	event  *handler.Event
	action *handler.Action
	watch  *handler.Watch
}

func NewApplication() *Application {
	sl := cslack.NewSlackClient()

	eventHandler := handler.NewEvent(sl)
	actionHandler := handler.NewAction(sl)
	watchHandler := handler.NewWatch(sl)
	return &Application{
		slack:  sl,
		event:  eventHandler,
		action: actionHandler,
		watch:  watchHandler,
	}
}

func (a *Application) Routing() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/health", a.healthCheck)

	mux.Post("/event", a.handleEvent)
	mux.Post("/interactive", a.handleAction)
	mux.Post("/watch", a.handleWatch)

	return mux
}

// handle Action
func (a *Application) handleAction(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	r.ParseForm()
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		response.Error(w, r, code.Errorf(code.JSON, "failed to parse action response JSON: %v", err))
		return
	}

	if err = a.action.Handle(payload); err != nil {
		response.Error(w, r, err)
		return
	}
	response.Success(w, r, "")
}

func (a *Application) handleEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.Error(w, r, code.Errorf(code.Internal, "failed to read body: %v", err))
		return
	}

	sv, err := slack.NewSecretsVerifier(r.Header, config.SlackSigningSecret())
	if err != nil {
		response.Error(w, r, code.Errorf(code.Slack, "failed to new secret verifier: %v", err))
		return
	}
	if _, err := sv.Write(body); err != nil {
		response.Error(w, r, code.Errorf(code.Slack, "failed to write body: %v", err))
		return
	}
	if err := sv.Ensure(); err != nil {
		response.Error(w, r, code.Errorf(code.Slack, "failed to ensure: %v", err))
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		response.Error(w, r, code.Errorf(code.Slack, "failed to parse event: %v", err))
		return
	}

	// challenge response
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var chRes *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &chRes)
		if err != nil {
			response.Error(w, r, code.Errorf(code.JSON, "failed to unmarshal json: %v", err))
			return
		}
		render.Status(r, http.StatusOK)
		render.PlainText(w, r, chRes.Challenge)
		return
	}

	// application
	if err = a.event.Handle(eventsAPIEvent); err != nil {
		response.Error(w, r, err)
		return
	}
	response.Success(w, r, nil)
}

func (a *Application) handleWatch(w http.ResponseWriter, r *http.Request) {
	now := time.Now().Local()
	if err := a.watch.Handle(now); err != nil {
		response.Error(w, r, err)
		return
	}

	response.Success(w, r, nil)
}

func (a *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}
