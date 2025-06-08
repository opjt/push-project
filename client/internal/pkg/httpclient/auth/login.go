package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"push/client/internal/pkg/httpclient"
	commonDto "push/common/dto"
	"push/linker/dto"
)

func AuthLogin(req dto.AuthLoginReq) (dto.AuthLoginRes, error) {
	resp, err := httpclient.DoJSONRequest("POST", "http://localhost:8800/api/v1/auth/login", req)
	if err != nil {
		return dto.AuthLoginRes{}, fmt.Errorf("failed to login: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp commonDto.CommonResponse[any]
		if err := json.Unmarshal(resp.Body, &errResp); err != nil {
			return dto.AuthLoginRes{}, fmt.Errorf("login failed with unknown error: %s", string(resp.Body))
		}
		return dto.AuthLoginRes{}, fmt.Errorf("login failed: %s", errResp.Error)
	}

	var commonResp commonDto.CommonResponse[dto.AuthLoginRes]
	if err := json.Unmarshal(resp.Body, &commonResp); err != nil {
		return dto.AuthLoginRes{}, fmt.Errorf("failed to parse login response: %w", err)
	}

	return commonResp.Data, nil
}
