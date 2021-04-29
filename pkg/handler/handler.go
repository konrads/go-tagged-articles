package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/konrads/go-tagged-articles/pkg/db"
	"github.com/konrads/go-tagged-articles/pkg/model"
	"github.com/konrads/go-tagged-articles/pkg/utils"
)

type Handler struct {
	db *db.DB
}

// Handler for REST methods, including param handling, body generation, response handling
func NewHandler(db *db.DB) Handler {
	return Handler{db}
}

func (h *Handler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	a, err := (*h.db).GetArticle(id)
	if err == nil {
		c.JSON(http.StatusOK, a)
	} else if a == nil {
		c.Status(http.StatusNotFound)
	} else {
		c.JSON(http.StatusInternalServerError, err)
	}
}

func (h *Handler) PostArticle(c *gin.Context) {
	defer c.Request.Body.Close()
	var areq model.ArticleReq
	err := c.ShouldBind(&areq)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprint(err))
	} else {
		uuid := uuid.New().String()
		a := model.Article{ID: uuid, ArticleReq: &areq}
		saved, err := (*h.db).SaveArticle(&a)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else if saved != 1 {
			c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to persist article %v in db", a))
		} else {
			c.JSON(http.StatusOK, a)
		}
	}
}

func (h *Handler) GetTagInfos(c *gin.Context) {
	tag := c.Param("tag")
	date := c.Param("date")
	parsedDate, err := utils.ToUrlParamDate(date)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	} else {
		as, err := (*h.db).GetArticles(tag, *parsedDate)
		if err == nil {
			allArticleIds := []string{}
			allTags := []string{}
			for _, a := range as {
				allArticleIds = append(allArticleIds, a.ID)
				if a.Tags != nil {
					for _, t := range a.Tags {
						allTags = append(allTags, t)
					}
				}
			}
			tagInfo := model.TagInfo{
				Tag:         tag,
				Count:       len(allTags),
				Articles:    allArticleIds,
				RelatedTags: utils.UniqueStrings(utils.FilterStrings(allTags, tag)),
			}

			c.JSON(http.StatusOK, tagInfo)
		} else if as == nil {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, err)
		}
	}
}

// FIXME: add failure scenarios
