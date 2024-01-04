package gcal

import (
	"context"
	"log"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/ddeville/kattungar-notify/apns"
	"github.com/ddeville/kattungar-notify/store"
)

type CalendarConfig struct {
	GoogleCredentialsPath string
	GoogleRefreshToken    string
	CalendarId            string
	ApnsClient            *apns.ApnsClient
	Store                 *store.Store
}

type CalendarClient struct {
	svc        *calendar.Service
	calendarId string
}

func NewClient(cfg CalendarConfig) (*CalendarClient, error) {
	cfgData, err := os.ReadFile(cfg.GoogleCredentialsPath)
	if err != nil {
		return nil, err
	}

	googleCfg, err := google.ConfigFromJSON(cfgData, calendar.CalendarReadonlyScope, calendar.CalendarEventsReadonlyScope)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{RefreshToken: cfg.GoogleRefreshToken}
	client := googleCfg.Client(context.Background(), token)

	svc, err := calendar.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return &CalendarClient{svc, cfg.CalendarId}, nil
}

const tickerDuration = 5 * time.Minute
const eventSpread = 4 * time.Minute

func (c *CalendarClient) Run(ctx context.Context) {
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	c.checkEvents()

	for {
		select {
		case <-ticker.C:
			c.checkEvents()
		case <-ctx.Done():
			log.Println("Calendar runloop was canceled")
			return
		}
	}
}

func (c *CalendarClient) checkEvents() {
	events, err := c.svc.Events.List(c.calendarId).
		TimeMin(time.Now().Add(-eventSpread).Format(time.RFC3339)).
		TimeMax(time.Now().Add(+eventSpread).Format(time.RFC3339)).
		Do()
	if err != nil {
		log.Printf("Unable to retrieve events: %v\n", err)
		return
	}

	for _, event := range events.Items {
		// TODO(damien): Check whether we should post a notification for this event
		c.postNotification(event)
	}
}

func (c *CalendarClient) postNotification(event *calendar.Event) {
	// TODO(damien): Post notification
	println(event.Summary)

	// TODO(damien): Record notification in store
}
