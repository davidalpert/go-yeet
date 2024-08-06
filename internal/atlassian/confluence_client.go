package atlassian

import (
	"context"
	"fmt"
	"github.com/ctreminiom/go-atlassian/confluence/v2"
	"github.com/ctreminiom/go-atlassian/pkg/infra/models"
	"net/url"
	"strings"
)

type Confluence struct {
	*confluence.Client
}

func New(cfg *CloudConfig) (*Confluence, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg is required")
	}

	if err := cfg.Validate("atlassian_config"); len(err) > 0 {
		return nil, fmt.Errorf(strings.Join(err, ", "))
	}

	instance, err := confluence.New(nil, cfg.Host)
	if err != nil {
		return nil, err
	}

	instance.Auth.SetBasicAuth(cfg.Email, cfg.APIToken)
	instance.Auth.SetUserAgent("curl/7.54.0")

	return &Confluence{
		Client: instance,
	}, nil
}

func (c *Confluence) GetPages() (interface{}, error) {
	options := &models.PageOptionsScheme{
		PageIDs:    nil,
		SpaceIDs:   nil,
		Sort:       "created-date",
		Status:     []string{"current"},
		Title:      "",
		BodyFormat: "atlas_doc_format",
	}

	results := make([]*models.PageScheme, 0)
	var cursor string
	for {
		chunk, response, err := c.Client.Page.Gets(context.Background(), options, cursor, 20)
		if err != nil {

			if response != nil {
				return nil, fmt.Errorf("%d %s: %#v", response.Code, response.Status, err)
			}

			return nil, err
		}

		results = append(results, chunk.Results...)

		if chunk.Links != nil && chunk.Links.Next == "" {
			break
		}

		values, err := url.ParseQuery(chunk.Links.Next)
		if err != nil {
			return nil, fmt.Errorf("error parsing next url: %#v", err)
		}

		_, containsCursor := values["cursor"]
		if containsCursor {
			cursor = values["cursor"][0]
		}
	}
	return results, nil
}

func (c *Confluence) GetChildPagesOf(pageID int) ([]*models.ChildPageScheme, error) {
	results := make([]*models.ChildPageScheme, 0)
	var cursor string
	for {
		chunk, response, err := c.Client.Page.GetsByParent(context.Background(), pageID, cursor, 20)
		if err != nil {
			if response != nil {
				return nil, fmt.Errorf("children of page (%d): %d %s: %#v", pageID, response.Code, response.Status, err)
			}

			return nil, fmt.Errorf("children of page (%d): %#v", pageID, err)
		}

		results = append(results, chunk.Results...)

		if chunk.Links != nil && chunk.Links.Next == "" {
			break
		}

		values, err := url.ParseQuery(chunk.Links.Next)
		if err != nil {
			return nil, fmt.Errorf("error parsing next url: %#v", err)
		}

		_, containsCursor := values["cursor"]
		if containsCursor {
			cursor = values["cursor"][0]
		}
	}
	return results, nil
}

func (c *Confluence) GetPageByID(pageID int, includeChildren bool) (*PageWithChildren, error) {
	page, response, err := c.Client.Page.Get(context.Background(), pageID, "atlas_doc_format", false, 3)
	if err != nil {
		if response != nil {
			return nil, fmt.Errorf("getting requested page (%d): %d %s: %#v", pageID, response.Code, response.Status, err)
		}
		return nil, fmt.Errorf("getting requested page (%d): %#v", pageID, err)
	}

	result := PageWithChildren{Page: page}

	if includeChildren {
		children, err2 := c.GetChildPagesOf(pageID)
		if err2 != nil {
			return nil, err2
		}
		result.Children = children
	}

	return &result, nil
}

type PageWithChildren struct {
	Page     *models.PageScheme
	Children []*models.ChildPageScheme
}

//type getChunkByCursor func(cursor string) (*models.PageChunkScheme, *models.ResponseScheme, error)
//
//func (c *Confluence) collectOverAllChunks(getChunk getChunkByCursor) ([]*models.PageScheme, error) {
//}
