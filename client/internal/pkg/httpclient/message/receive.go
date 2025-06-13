package message

import (
	"fmt"
	"net/http"
	"push/client/internal/pkg/httpclient"
)

func MessageReceive(msgId uint64) error {
	if msgId == 0 {
		return nil
	}
	resp, err := httpclient.DoJSONRequest("POST", "http://localhost:8800/api/v1/messages/"+fmt.Sprint(msgId)+"/receive", nil)
	if err != nil || resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to receive message: %w", err)
	}

	return nil
}
