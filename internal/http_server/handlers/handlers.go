package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog"
	"net/http"
	"sf-api-gateway/internal/config"
	"sf-api-gateway/internal/constants"
	"sf-api-gateway/internal/service/censor"
	"sf-api-gateway/internal/service/comments"
	"sf-api-gateway/internal/service/news"
)

type Handler struct {
	cfg      *config.Config
	lgr      zerolog.Logger
	censor   *censor.Censor
	news     *news.News
	comments *comments.Comments
}

func NewHandler(cfg *config.Config, lgr zerolog.Logger, censor *censor.Censor, news *news.News, comments *comments.Comments) *Handler {
	return &Handler{
		cfg:      cfg,
		lgr:      lgr,
		censor:   censor,
		news:     news,
		comments: comments,
	}
}

func (h *Handler) GetLastNews(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(GetLastNewsReq)
	err := decoder.Decode(request)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect request: %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := h.lgr.With().
		Str("handler", "GetLastNews").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			//Time("date_gte", request.DateGte).
			//Time("date_lte", request.DateLte).
			//Uint64("limit", request.Limit).
			//Uint64("offset", request.Offset).
			Uint64("count", request.Count)).
		Logger()

	getNewsResp, err := h.news.GetNews(ctx, request.Count)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")

	newsEntities := make([]NewsEntity, 0, 10)
	for _, n := range getNewsResp.News {
		newsEntities = append(newsEntities, NewsEntity{
			Id:      n.Id,
			Title:   n.Title,
			Link:    n.Link,
			Desc:    n.Desc,
			PubDate: n.PubDate,
		})
	}

	resp, _ := json.Marshal(GetLastNewsResp{News: newsEntities})
	fmt.Fprintf(w, string(resp))
}

func (h *Handler) GetNewsDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(GetNewsDetailReq)
	err := decoder.Decode(request)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect request: %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := h.lgr.With().
		Str("handler", "GetNewsDetail").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("new_id", request.NewId)).
		Logger()

	getNewResp, err := h.news.GetNew(ctx, request.NewId)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	getCommentsResp, err := h.comments.GetComments(ctx, request.NewId)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	comments := make([]CommentsEntity, 0, 10)
	for _, c := range getCommentsResp.Comments {
		comments = append(comments, CommentsEntity{
			Id:       c.Id,
			ParentId: c.ParentId,
			Text:     c.Text,
		})
	}

	lgr.Debug().Msg("executed")

	resp, _ := json.Marshal(GetNewsDetailResp{
		New: NewsEntity{
			Id:      getNewResp.Id,
			Title:   getNewResp.Title,
			Link:    getNewResp.Link,
			Desc:    getNewResp.Desc,
			PubDate: getNewResp.PubDate,
		},
		Comments: comments,
	})
	fmt.Fprintf(w, string(resp))
}

func (h *Handler) AddNewsComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	request := new(AddNewsCommentReq)
	err := decoder.Decode(request)
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("incorrect request: %s", err.Error())})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	ctx := r.Context()
	requestId, _ := ctx.Value(constants.RequestIdKey).(string)

	lgr := h.lgr.With().
		Str("handler", "AddNewsComment").
		Str(constants.RequestIdKey, requestId).
		Dict("request", zerolog.Dict().
			Uint64("NewId", request.NewId).
			Interface("ParentId", request.ParentId).
			Str("Text", request.Text)).
		Logger()

	checkResp, err := h.censor.Check(ctx, &censor.CheckReq{Text: request.Text})
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}
	if !checkResp.Status {
		resp, _ := json.Marshal(ErrorResp{Error: "text: did not pass the censorship"})
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, string(resp))
		return
	}

	addCommentResp, err := h.comments.AddComment(ctx, &comments.AddCommentReq{
		NewId:    request.NewId,
		ParentId: request.ParentId,
		Text:     request.Text,
	})
	if err != nil {
		resp, _ := json.Marshal(ErrorResp{Error: fmt.Sprintf("internal error: %s", err.Error())})
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, string(resp))
		return
	}

	lgr.Debug().Msg("executed")

	resp, _ := json.Marshal(AddNewsCommentResp{Status: addCommentResp.Status})
	fmt.Fprintf(w, string(resp))
}
