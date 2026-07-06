package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// SECOPRecord represents a public call from SECOP II.
type SECOPRecord struct {
	Entidad                 string `json:"entidad"`
	NitEntidad              string `json:"nit_entidad"`
	NombreProcedimiento     string `json:"nombre_del_procedimiento"`
	DescripcionProcedimiento string `json:"descripci_n_del_procedimiento"`
	Fase                    string `json:"fase"`
	FechaPublicacion        string `json:"fecha_de_publicacion_del"`
	FechaUltimaPublicacion  string `json:"fecha_ultima_publicacion"`
	ModalidadContratacion   string `json:"modalidad_de_contratacion"`
	PrecioBase              string `json:"precio_base"`
	EstadoProcedimiento     string `json:"estado_del_procedimiento"`
	URLProceso              string `json:"url_proceso"`
	ValorTotalAdjudicacion  string `json:"valor_total_adjudicacion"`
	NombreProveedor         string `json:"nombre_proveedor"`
}

// FetchPublicCalls performs an HTTP GET to the given base URL with the provided query string,
// decodes the JSON response, and returns a slice of SECOPRecord.
func FetchPublicCalls(baseURL, queryString string) ([]SECOPRecord, error) {
	fullURL := baseURL + "?" + queryString

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("http GET failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	var records []SECOPRecord
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	return records, nil
}

// BuildQuery constructs a URL-encoded query string from a map of SODA API parameters.
// Typical keys: $limit, $where, $order, $select.
func BuildQuery(params map[string]string) string {
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values.Encode()
}

func main() {
	baseURL := "https://datos.gov.co/resource/p6dx-8zbt.json"

	// Build query parameters
	params := map[string]string{
		"$limit": "5",
		"$where": "nombre_del_procedimiento like '%convocatoria%'",
		"$order": "fecha_de_publicacion_del DESC",
		"$select": "entidad,nit_entidad,nombre_del_procedimiento,descripci_n_del_procedimiento,fase,fecha_de_publicacion_del,fecha_ultima_publicacion,modalidad_de_contratacion,precio_base,estado_del_procedimiento,url_proceso,valor_total_adjudicacion,nombre_proveedor",
	}

	queryString := BuildQuery(params)

	records, err := FetchPublicCalls(baseURL, queryString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Total records fetched: %d\n\n", len(records))
	for i, r := range records {
		fmt.Printf("Record %d:\n", i+1)
		fmt.Printf("  Entidad:                  %s\n", r.Entidad)
		fmt.Printf("  NIT Entidad:              %s\n", r.NitEntidad)
		fmt.Printf("  Nombre Procedimiento:     %s\n", r.NombreProcedimiento)
		fmt.Printf("  Descripción Procedimiento:%s\n", r.DescripcionProcedimiento)
		fmt.Printf("  Fase:                     %s\n", r.Fase)
		fmt.Printf("  Fecha Publicación:        %s\n", r.FechaPublicacion)
		fmt.Printf("  Fecha Última Publicación: %s\n", r.FechaUltimaPublicacion)
		fmt.Printf("  Modalidad Contratación:   %s\n", r.ModalidadContratacion)
		fmt.Printf("  Precio Base:              %s\n", r.PrecioBase)
		fmt.Printf("  Estado Procedimiento:     %s\n", r.EstadoProcedimiento)
		fmt.Printf("  URL Proceso:              %s\n", r.URLProceso)
		fmt.Printf("  Valor Total Adjudicación: %s\n", r.ValorTotalAdjudicacion)
		fmt.Printf("  Nombre Proveedor:         %s\n", r.NombreProveedor)
		fmt.Println()
	}
}
