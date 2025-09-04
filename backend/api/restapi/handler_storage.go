package restapi

import (
	"net/http"
	"net/url"
	"wedding-app/domain/service"

	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	storage service.StorageService
}

func NewBasicHandler(ss service.StorageService) *StorageHandler {
	return &StorageHandler{storage: ss}
}

func (h *StorageHandler) getUserSvgs(c *gin.Context) {
	svgs, err := h.storage.GetUserSvgs(c)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, svgs)
}

func (h *StorageHandler) getBucketUrls(c *gin.Context) {
	bucketName := c.Query(QueryBucket)
	suffix := c.Query(QuerySuffix)

	if bucketName == "" || suffix == "" {
		c.Error(NewAPIError(http.StatusBadRequest, "missing bucket or suffix query parameter", nil))
		return
	}

	urls, err := h.storage.GetBucketUrls(c, bucketName, suffix)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.JSON(http.StatusOK, urls)
}

func (h *StorageHandler) getBucketItemData(c *gin.Context) {
	bucketName := c.Query(QueryBucket)
	itemName := c.Query(QueryName)

	if bucketName == "" || itemName == "" {
		c.Error(NewAPIError(http.StatusBadRequest, "missing bucket or name query parameter", nil))
		return
	}

	item, err := h.storage.GetBucketData(c, bucketName, itemName)
	if err != nil {
		c.Error(NewInternalAPIError(err))
		return
	}

	c.Data(http.StatusOK, item.CType, item.Data)
}

const (
	BucketUrlsEndpoint = "/bucket-urls"
	BucketDataEndpoint = "/bucket-data"
	QueryBucket        = "bucket"
	QueryName          = "name"
	QuerySuffix        = "suffix"
)

func MakeBucketDataURL(bucketName, itemName string) string {
	values := url.Values{}
	values.Set(QueryBucket, bucketName)
	values.Set(QueryName, itemName)
	return BucketDataEndpoint + "?" + values.Encode()
}
