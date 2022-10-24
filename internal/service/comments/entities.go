package comments

import (
	"fmt"
	"strings"
	"time"
)

type JsonTime time.Time

func (t JsonTime) MarshalJSON() (b []byte, err error) {
	tm := time.Time(t)

	return []byte(fmt.Sprintf(`"%s"`, tm.Format(time.RFC3339))), nil
}

func (t JsonTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)

	tm, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	t = JsonTime(tm)

	return nil
}

type AddCommentReq struct {
	NewId    uint64  `json:"new_id"`
	ParentId *uint64 `json:"parent_id"`
	Text     string  `json:"text"`
}

type AddCommentResp struct {
	Status string `json:"status"`
}

type GetCommentsReq struct {
	NewId uint64 `json:"new_id"`
}

type CommentEntity struct {
	Id       uint64   `json:"id"`
	ParentId *uint64  `json:"parent_id"`
	Text     string   `json:"text"`
	PubDate  JsonTime `json:"pub_date"`
}

type GetCommentsResp struct {
	Comments []CommentEntity `json:"comments"`
}
