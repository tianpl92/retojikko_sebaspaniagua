package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DatosGovService fetches public proposals from the external datos.gov.co API.
type DatosGovService struct {
	baseURL    string
	httpClient *http.Client
}

func NewDatosGovService(baseURL string) *DatosGovService {
	return &DatosGovService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchProposals queries the datos.gov.co SODA API and returns raw JSON records.
func (s *DatosGovService) FetchProposals(query, fase, entidad string, limit, offset int) ([]map[string]interface{}, error) {
	params := url.Values{}

	// Build SoQL filters
	var whereClauses []string
	if query != "" {
		whereClauses = append(whereClauses,
			fmt.Sprintf("lower(nombre_del_procedimiento) like '%%%s%%'", strings.ToLower(url.QueryEscape(query))))
	}
	if fase != "" {
		whereClauses = append(whereClauses,
			fmt.Sprintf("fase = '%s'", fase))
	}
	if entidad != "" {
		whereClauses = append(whereClauses,
			fmt.Sprintf("lower(entidad) like '%%%s%%'", strings.ToLower(url.QueryEscape(entidad))))
	}
	if len(whereClauses) > 0 {
		params.Set("$where", strings.Join(whereClauses, " AND "))
	}

	if limit > 0 {
		params.Set("$limit", fmt.Sprintf("%d", limit))
	}
	if offset > 0 {
		params.Set("$offset", fmt.Sprintf("%d", offset))
	}
	params.Set("$order", "fecha_de_publicacion_del DESC")

	fullURL := s.baseURL + "?" + params.Encode()

	req, err := http.NewRequest(http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("external API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	var records []map[string]interface{}
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return records, nil
}
