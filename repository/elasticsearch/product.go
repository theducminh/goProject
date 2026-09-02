package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"goproject/domain"

	es "github.com/elastic/go-elasticsearch/v8"
)

type ProductSearchRepository struct {
	client *es.Client
	index  string
}

type productDocument struct {
	ID    uint   `json:"id"`
	Code  string `json:"code"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
	Stock uint   `json:"stock"`
}

func NewProductSearchRepository(client *es.Client, index string) *ProductSearchRepository {
	return &ProductSearchRepository{client: client, index: index}
}

func (search *ProductSearchRepository) Index(ctx context.Context, product *domain.Product) error {
	body, err := json.Marshal(toProductDocument(product))
	if err != nil {
		return err
	}
	response, err := search.client.Index(search.index, bytes.NewReader(body), search.client.Index.WithContext(ctx), search.client.Index.WithDocumentID(strconv.FormatUint(uint64(product.ID), 10)), search.client.Index.WithRefresh("wait_for"))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("elasticsearch index returned status %s", response.Status())
	}
	return nil
}

func (search *ProductSearchRepository) Delete(ctx context.Context, id uint) error {
	response, err := search.client.Delete(search.index, strconv.FormatUint(uint64(id), 10), search.client.Delete.WithContext(ctx), search.client.Delete.WithRefresh("wait_for"))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() && response.StatusCode != 404 {
		return fmt.Errorf("elasticsearch delete returned status %s", response.Status())
	}
	return nil
}

func (search *ProductSearchRepository) Search(ctx context.Context, query string) ([]domain.Product, error) {
	body, err := json.Marshal(map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"name^2", "code"},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	response, err := search.client.Search(search.client.Search.WithContext(ctx), search.client.Search.WithIndex(search.index), search.client.Search.WithBody(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.IsError() {
		return nil, fmt.Errorf("elasticsearch search returned status %s", response.Status())
	}
	var result struct {
		Hits struct {
			Hits []struct {
				Source productDocument `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(io.LimitReader(response.Body, 2<<20)).Decode(&result); err != nil {
		return nil, err
	}
	products := make([]domain.Product, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		products = append(products, domain.Product{ID: hit.Source.ID, Code: hit.Source.Code, Name: hit.Source.Name, Price: hit.Source.Price, Stock: hit.Source.Stock})
	}
	return products, nil
}

func (search *ProductSearchRepository) Ping(ctx context.Context) error {
	response, err := search.client.Ping(search.client.Ping.WithContext(ctx))
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.IsError() {
		return fmt.Errorf("elasticsearch ping returned status %s", response.Status())
	}
	return nil
}

func toProductDocument(product *domain.Product) productDocument {
	return productDocument{ID: product.ID, Code: product.Code, Name: product.Name, Price: product.Price, Stock: product.Stock}
}
