package handlers

type ErrorResp struct {
	Error string `json:"error"`
}

type GetLastNewsReq struct {
	//DateGte time.Time `json:"date_gte"`
	//DateLte time.Time `json:"date_lte"`
	//Limit   uint64    `json:"limit"`
	//Offset  uint64    `json:"offset"`
	Count uint64 `json:"count"`
}

type NewsEntity struct {
	Id      uint64 `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"description"`
	Link    string `json:"link"`
	PubDate string `json:"pub_date"`
}

type GetLastNewsResp struct {
	News []NewsEntity `json:"news"`
}

type AddNewsCommentReq struct {
	NewId    uint64  `json:"new_id"`
	ParentId *uint64 `json:"parent_id"`
	Text     string  `json:"text"`
}

type AddNewsCommentResp struct {
	Status string `json:"status"`
}

type GetNewsDetailReq struct {
	NewId uint64 `json:"new_id"`
}

type CommentsEntity struct {
	Id       uint64  `json:"id"`
	ParentId *uint64 `json:"parent_id"`
	Text     string  `json:"text"`
}

type GetNewsDetailResp struct {
	New      NewsEntity       `json:"new"`
	Comments []CommentsEntity `json:"comments"`
}
