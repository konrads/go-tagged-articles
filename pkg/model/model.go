package model

import (
	"fmt"
	"strings"
	"time"
)

// from: https://stackoverflow.com/questions/45303326/how-to-parse-non-standard-time-format-from-json
type Date time.Time

func (j *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = Date(t)
	return nil
}

func (j Date) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))), nil
	// return json.Marshal(time.Time(j))
}

func (j Date) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

type ArticleReq struct {
	Title string   `json:"title" binding:"required"`
	Date  Date     `json:"date" binding:"required"`
	Body  string   `json:"body" binding:"required"`
	Tags  []string `json:"tags"`
}

type Article struct {
	ID string `json:"id" binding:"required"`
	*ArticleReq
}

type TagInfo struct {
	Tag         string   `json:"tag"`
	Count       int      `json:"count"`
	Articles    []string `json:"articles"`
	RelatedTags []string `json:"related_tags"`
}
