package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/konrads/go-tagged-articles/pkg/db"
	"github.com/konrads/go-tagged-articles/pkg/handler"
	"github.com/konrads/go-tagged-articles/pkg/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedDb struct {
	mock.Mock
}

func (m *MockedDb) GetArticle(id string) (*model.Article, error) {
	args := m.Called(id)
	res := args.Get(0).(*model.Article)
	err := args.Error(1)
	return res, err
}

func (m *MockedDb) SaveArticle(article *model.Article) (int, error) {
	args := m.Called(article)
	res := args.Int(0)
	err := args.Error(1)
	return res, err
}

func (m *MockedDb) GetArticles(tag string, date model.Date) ([]*model.Article, error) {
	args := m.Called(tag, date)
	res := args.Get(0).([]*model.Article)
	err := args.Error(1)
	return res, err
}

func (m *MockedDb) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetArticle(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	var mockDb = &MockedDb{}
	var asDb db.DB = mockDb

	res := model.Article{
		ID: "1",
		ArticleReq: &model.ArticleReq{
			Title: "t1",
			Body:  "b1",
			Date:  model.Date(time.Date(2017, 7, 2, 0, 0, 0, 0, time.UTC)),
			Tags:  []string{"tt1"},
		},
	}
	mockDb.On("GetArticle", "1").Return(&res, nil)

	handler := handler.NewHandler(&asDb)
	handler.GetArticle(c)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"id":"1","title":"t1","date":"2017-07-02","body":"b1","tags":["tt1"]}`, w.Body.String())
}

func TestPostArticle(t *testing.T) {
	artReq := &model.ArticleReq{
		Title: "t1",
		Body:  "b1",
		Date:  model.Date(time.Date(2017, 7, 2, 0, 0, 0, 0, time.UTC)),
		Tags:  []string{"tt1"},
	}
	artReqBody, _ := json.Marshal(artReq)
	res := model.Article{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("POST", "/articles", bytes.NewReader(artReqBody))
	c.Request.Header.Set("Content-Type", "application/json")
	var mockDb = &MockedDb{}
	var asDb db.DB = mockDb

	mockDb.On("SaveArticle", mock.Anything).Return(1, nil)

	handler := handler.NewHandler(&asDb)
	handler.PostArticle(c)
	// adding extra unmarshalling to deal with unique uuid for the Article
	json.Unmarshal([]byte(w.Body.String()), &res)
	assert.Equal(t, 200, w.Code)
	assert.NotNil(t, res.ID)
	assert.Equal(t, artReq, res.ArticleReq)
}

func TestGetTagInfos(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "tag", Value: "tt1"}, {Key: "date", Value: "20170702"}}
	var mockDb = &MockedDb{}
	var asDb db.DB = mockDb

	date := model.Date(time.Date(2017, 7, 2, 0, 0, 0, 0, time.UTC))

	dbRes := []*model.Article{
		{
			ID: "1",
			ArticleReq: &model.ArticleReq{
				Title: "t1",
				Body:  "b1",
				Date:  date,
				Tags:  []string{"tt1", "tt2"},
			},
		},
		{
			ID: "1",
			ArticleReq: &model.ArticleReq{
				Title: "t2",
				Body:  "b2",
				Date:  date,
				Tags:  []string{"tt1", "tt3"},
			},
		},
	}
	mockDb.On("GetArticles", "tt1", date).Return(dbRes, nil)

	handler := handler.NewHandler(&asDb)
	handler.GetTagInfos(c)
	json.Unmarshal([]byte(w.Body.String()), &dbRes)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"tag":"tt1","count":4,"articles":["1","1"],"related_tags":["tt2","tt3"]}`, w.Body.String())
}
