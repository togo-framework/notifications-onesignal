// Package onesignal is a OneSignal push channel for togo notifications.
// Install: `togo install togo-framework/notifications-onesignal`.
package onesignal

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/togo-framework/notifications"
	"github.com/togo-framework/togo"
)

const endpoint = "https://onesignal.com/api/v1/notifications"

func init() {
	notifications.RegisterChannel("onesignal", func(k *togo.Kernel) notifications.Channel {
		return &channel{
			appID:  os.Getenv("ONESIGNAL_APP_ID"),
			key:    os.Getenv("ONESIGNAL_API_KEY"),
			client: &http.Client{Timeout: 15 * time.Second},
		}
	})
}

type channel struct {
	appID, key string
	client     *http.Client
}

func (c *channel) Send(ctx context.Context, to notifications.Notifiable, n notifications.Notification) error {
	pn, ok := n.(notifications.PushNotification)
	if !ok {
		return nil
	}
	tokens := to.RoutePushTokens()
	if len(tokens) == 0 || c.appID == "" || c.key == "" {
		return nil
	}
	msg := pn.ToPush(to)
	body := map[string]any{
		"app_id":             c.appID,
		"include_player_ids": tokens,
		"headings":           map[string]string{"en": msg.Title},
		"contents":           map[string]string{"en": msg.Body},
	}
	if len(msg.Data) > 0 {
		body["data"] = msg.Data
	}
	buf, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic "+c.key)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("onesignal: status %d: %s", resp.StatusCode, string(b))
	}
	return nil
}
