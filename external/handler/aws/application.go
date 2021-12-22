package gcp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/akubi0w1/chatbot-202201/code"
	"github.com/akubi0w1/chatbot-202201/config"
	"github.com/akubi0w1/chatbot-202201/external/handler"
	cslack "github.com/akubi0w1/chatbot-202201/external/slack"
	"github.com/aws/aws-lambda-go/events"
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

func (app *Application) Handle(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch r.Path {
	case "/event":
		return app.handleEvent(ctx, r)

	case "/interactive":
		return app.handleAction(ctx, r)

	case "/watch":
		app.handleWatch(ctx)
		return makeLambdaResponse(http.StatusOK, ""), nil
	}

	return makeLambdaResponse(http.StatusNotFound, code.Errorf(code.NotFound, "request path is invalid: %v", r.Path).Error()), nil
}

func (a *Application) handleAction(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body, err := url.QueryUnescape(request.Body[8:])
	if err != nil {
		err := code.Errorf(code.InvalidQuery, "failed to encode body: %v", err)
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), err
	}
	var payload slack.InteractionCallback
	err = json.Unmarshal([]byte(body), &payload)
	if err != nil {
		err := code.Errorf(code.JSON, "failed to parse action response JSON: %v", err)
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), err
	}

	if err = a.action.Handle(payload); err != nil {
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), nil
	}

	return makeLambdaResponse(http.StatusOK, ""), nil
}

func (a *Application) handleEvent(ctx context.Context, r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := []byte(r.Body)
	header := map[string][]string{}
	for k, v := range r.Headers {
		header[k] = strings.Split(v, ",")
	}

	sv, err := slack.NewSecretsVerifier(header, config.SlackSigningSecret())
	if err != nil {
		err := code.Errorf(code.Slack, "failed to new secret verifier: %v", err)
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), err
	}
	if _, err := sv.Write(body); err != nil {
		err := code.Errorf(code.Slack, "failed to write body: %v", err)
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), err
	}
	if err := sv.Ensure(); err != nil {
		err := code.Errorf(code.Slack, "failed to ensure: %v", err)
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), err
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		err := code.Errorf(code.Slack, "failed to parse event: %v", err)
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), err
	}

	// challenge response
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var chRes *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &chRes)
		if err != nil {
			err = code.Errorf(code.JSON, "failed to unmarshal json: %v", err)
			return makeLambdaResponse(code.GetStatusCode(err), err.Error()), err
		}
		res := makeLambdaResponse(http.StatusOK, chRes.Challenge)
		res.Headers = map[string]string{
			"Content-Type": "text/plain",
		}
		return res, nil
	}

	// application
	if err = a.event.Handle(eventsAPIEvent); err != nil {
		return makeLambdaResponse(code.GetStatusCode(err), err.Error()), nil
	}

	return makeLambdaResponse(http.StatusOK, ""), nil
}

func (a *Application) handleWatch(ctx context.Context) {
	now := time.Now().Local()

	a.watch.Handle(now)
}

func (a *Application) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func makeLambdaResponse(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       body,
	}
}
