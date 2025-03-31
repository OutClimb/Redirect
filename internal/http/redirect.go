package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/OutClimb/Redirect/internal/app"
	"github.com/gin-gonic/gin"
)

type redirectPublic struct {
	Id       uint   `json:"id"`
	FromPath string `json:"fromPath"`
	ToUrl    string `json:"toUrl"`
	StartsOn int64  `json:"startsOn"`
	StopsOn  int64  `json:"stopsOn"`
}

func (r *redirectPublic) Publicize(redirect *app.RedirectInternal) {
	r.Id = redirect.ID
	r.FromPath = redirect.FromPath
	r.ToUrl = redirect.ToUrl

	r.StartsOn = 0
	if redirect.StartsOn != nil {
		r.StartsOn = redirect.StartsOn.UnixMilli()
	}

	r.StopsOn = 0
	if redirect.StopsOn != nil {
		r.StopsOn = redirect.StopsOn.UnixMilli()
	}
}

func (h *httpLayer) createRedirect(c *gin.Context) {
	// Get the body data
	bodyAsByteArray, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve request body"})
		return
	}

	// Parse the body data
	body := redirectPublic{}
	err = json.Unmarshal(bodyAsByteArray, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse request body"})
		return
	}

	if redirect, err := h.app.CreateRedirect(body.FromPath, body.ToUrl, body.StartsOn, body.StopsOn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create redirect"})
	} else {
		redirectPublic := redirectPublic{}
		redirectPublic.Publicize(redirect)

		c.JSON(http.StatusOK, redirectPublic)
	}
}

func (h *httpLayer) deleteRedirect(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = h.app.DeleteRedirect(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Redirect not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (h *httpLayer) findRedirect(c *gin.Context) {
	if redirect, err := h.app.FindRedirect(c.Request.URL.Path); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "https://outclimb.gay")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirect.ToUrl)
	}
}

func (h *httpLayer) getRedirect(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	internalRedirect, error := h.app.GetRedirect(uint(id))
	if error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Redirect not found"})
		return
	}

	redirect := redirectPublic{}
	redirect.Publicize(internalRedirect)

	c.JSON(http.StatusOK, redirect)
}

func (h *httpLayer) getRedirects(c *gin.Context) {
	if internalRedirects, err := h.app.GetAllRedirects(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve redirects"})
	} else {
		redirects := make([]redirectPublic, len(*internalRedirects))
		for i, redirect := range *internalRedirects {
			redirects[i].Publicize(&redirect)
		}

		c.JSON(http.StatusOK, redirects)
	}
}

func (h *httpLayer) updateRedirect(c *gin.Context) {
	// Get the id from the URL
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get the body data
	bodyAsByteArray, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to retrieve request body"})
		return
	}

	// Parse the body data
	body := redirectPublic{}
	err = json.Unmarshal(bodyAsByteArray, &body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse request body"})
		return
	}

	if redirect, err := h.app.UpdateRedirect(uint(id), body.FromPath, body.ToUrl, body.StartsOn, body.StopsOn); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to update redirect"})
	} else {
		redirectPublic := redirectPublic{}
		redirectPublic.Publicize(redirect)

		c.JSON(http.StatusOK, redirectPublic)
	}
}
