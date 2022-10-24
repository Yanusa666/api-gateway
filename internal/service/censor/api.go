package censor

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"sf-api-gateway/internal/constants"
)

func (c *Censor) Check(ctx context.Context, req *CheckReq) (result *CheckResp, err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := c.lgr.With().
		Str("api", "Check").
		Str(constants.RequestIdKey, requestId).
		Interface("request", req).
		Logger()

	reqBytes, _ := json.Marshal(req)
	reqBody := bytes.NewReader(reqBytes)
	uri := c.url.String() + "/check"

	httpReq, err := http.NewRequest(http.MethodPost, uri, reqBody)
	if err != nil {
		lgr.Error().Err(err).Msg("http.NewRequest failed")
		return nil, err
	}
	httpReq.Header.Set(constants.RequestIdKey, requestId)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		lgr.Error().Err(err).Msg("client.Do failed")
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		lgr.Error().Msgf("incorrect response code: %d", resp.StatusCode)
		return nil, err
	}

	result = new(CheckResp)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		lgr.Error().Err(err).Msg("decode response failed")
		return nil, err
	}

	return result, nil
}
