package comments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"sf-api-gateway/internal/constants"
)

func (c *Comments) GetComments(ctx context.Context, newId uint64) (result *GetCommentsResp, err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := c.lgr.With().
		Str("api", "GetNews").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("new_id", newId)).
		Logger()

	//reqBody := bytes.NewReader()
	uri := c.url.String() + fmt.Sprintf("/comments/%d", newId)

	httpReq, err := http.NewRequest(http.MethodGet, uri, nil)
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

	result = new(GetCommentsResp)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		lgr.Error().Err(err).Msg("decode response failed")
		return nil, err
	}

	return result, nil
}

func (c *Comments) AddComment(ctx context.Context, req *AddCommentReq) (result *AddCommentResp, err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := c.lgr.With().
		Str("api", "AddComment").
		Str(constants.RequestIdKey, requestId).
		Interface("request", req).
		Logger()

	reqBytes, _ := json.Marshal(req)
	reqBody := bytes.NewReader(reqBytes)
	uri := c.url.String() + "/comments"

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
		err = fmt.Errorf("incorrect response code: %d", resp.StatusCode)
		lgr.Error().Err(err).Msg("client.Do failed")
		return nil, err
	}

	result = new(AddCommentResp)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		lgr.Error().Err(err).Msg("decode response failed")
		return nil, err
	}

	return result, nil
}
