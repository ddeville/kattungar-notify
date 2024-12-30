package gcal

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"

	"github.com/ddeville/kattungar-notify/internal/apns"
	"github.com/ddeville/kattungar-notify/internal/store"
	"github.com/ddeville/kattungar-notify/internal/types"
)

type CalendarConfig struct {
	GoogleCredentialsPath string
	GoogleRefreshToken    string
	CalendarId            string
	ApnsClient            *apns.ApnsClient
	Store                 *store.Store
}

type CalendarClient struct {
	svc                 *calendar.Service
	apns                *apns.ApnsClient
	store               *store.Store
	calendarId          string
	notifiedEventsCache LRUCache
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

	return &CalendarClient{svc, cfg.ApnsClient, cfg.Store, cfg.CalendarId, NewLRUCache(512)}, nil
}

const tickerDuration = 5 * time.Minute
const eventSpread = 10 * time.Minute

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
	log.Println("Checking calendar events")

	events, err := c.svc.Events.List(c.calendarId).
		TimeMin(time.Now().Add(-eventSpread).Format(time.RFC3339)).
		TimeMax(time.Now().Add(+eventSpread).Format(time.RFC3339)).
		Do()
	if err != nil {
		log.Printf("Unable to retrieve events: %v\n", err)
		return
	}

	for _, event := range events.Items {
		// Check whether this is a recurring event, in which case fetch the actual instance
		if len(event.Recurrence) != 0 {
			instances, err := c.svc.Events.Instances(c.calendarId, event.Id).
				TimeMin(time.Now().Add(-eventSpread).Format(time.RFC3339)).
				TimeMax(time.Now().Add(eventSpread).Format(time.RFC3339)).
				Do()
			if err != nil {
				log.Printf("Unable to retrieve recurring instances for event: %v: %v\n", event.Summary, err)
				continue
			}

			if len(instances.Items) == 0 {
				log.Printf("Recurring event doesn't have instances in the time span: %v\n", event.Summary)
				continue
			}

			event = instances.Items[0]
		}

		// Look for all-day events and filter them out (they have `Date` set but not `DateTime`)
		if event.Start != nil && event.Start.DateTime == "" {
			log.Printf("Skipping all-day event: %v\n", event.Summary)
			continue
		}

		// Make sure that the event includes the tag in its description
		if !strings.Contains(event.Description, "#kattungar-notify") {
			log.Printf("Skipping event that doesn't have #kattungar-notify tag: %v\n", event.Summary)
			continue
		}

		if c.notifiedEventsCache.Contains(event.Id) || c.store.HasNotifiedCalendarEvent(event.Id) {
			log.Printf("Skipping event that already triggered a notification: %v\n", event.Summary)
			continue
		}

		c.postNotification(event)

		c.store.AddCalendarEvent(event.Id, true)
		c.notifiedEventsCache.Add(event.Id)
	}
}

func (c *CalendarClient) postNotification(event *calendar.Event) {
	devices, err := c.store.ListDevices()
	if err != nil {
		log.Printf("Failed to retrieve devices, not sending notification: %v\n", err)
	}

	for _, device := range devices {
		log.Printf("Sending notification for event \"%v\" to device %v\n", event.Summary, device.Name)
		notification := types.Notification{
			DeviceKey: device.Name,
			Body:      fmt.Sprintf("🗓️ %v", event.Summary),
		}

		c.apns.Notify(&device, notification)
		c.store.RecordNotification(notification)
	}
}
