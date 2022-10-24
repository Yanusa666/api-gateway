package news

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"sf-api-gateway/internal/constants"
)

func (n *News) GetNews(ctx context.Context, count uint64) (result *GetNewsResp, err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := n.lgr.With().
		Str("api", "GetNews").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("count", count)).
		Logger()

	//reqBody := bytes.NewReader()
	uri := n.url.String() + fmt.Sprintf("/news/list/%d", count)

	httpReq, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		lgr.Error().Err(err).Msg("http.NewRequest failed")
		return nil, err
	}
	httpReq.Header.Set(constants.RequestIdKey, requestId)

	resp, err := n.client.Do(httpReq)
	if err != nil {
		lgr.Error().Err(err).Msg("client.Do failed")
		return nil, err
	}
	defer resp.Body.Close()

	result = new(GetNewsResp)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		lgr.Error().Err(err).Msg("decode response failed")
		return nil, err
	}

	return result, nil
}

func (n *News) GetNew(ctx context.Context, newId uint64) (result *GetNewResp, err error) {
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := n.lgr.With().
		Str("api", "GetNew").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("new_id", newId)).
		Logger()

	//reqBody := bytes.NewReader()
	uri := n.url.String() + fmt.Sprintf("/news/get/%d", newId)

	httpReq, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		lgr.Error().Err(err).Msg("http.NewRequest failed")
		return nil, err
	}
	httpReq.Header.Set(constants.RequestIdKey, requestId)

	resp, err := n.client.Do(httpReq)
	if err != nil {
		lgr.Error().Err(err).Msg("client.Do failed")
		return nil, err
	}
	defer resp.Body.Close()

	result = new(GetNewResp)
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		lgr.Error().Err(err).Msg("decode response failed")
		return nil, err
	}

	return result, nil
}
