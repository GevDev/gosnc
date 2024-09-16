package tableapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type GlideRecord struct {
	TableName string
	SysId     string
	Rows      map[string]any
	//TODO: sys_tags which I can marshall
	//sys_created_on
	//sys_created_by
	//sys_updated_on
	//sys_updated_by
}

type TableAPI struct {
	InstanceURL string
	AuthToken   string
	BasePath    string
	HttpClient  *http.Client
}

func NewTableAPI(instanceUrl string, authToken string, basePath string, httpClient *http.Client) *TableAPI {
	return &TableAPI{
		InstanceURL: instanceUrl,
		AuthToken:   authToken,
		BasePath:    basePath,
		HttpClient:  httpClient,
	}
}

func (tableApi *TableAPI) GetEmptyGlideRecord() GlideRecord {
	return GlideRecord{}
}

func (tableApi *TableAPI) GetRecord(gr *GlideRecord, queryParams map[string]string) (*GlideRecord, error) {
	return &GlideRecord{
		gr.TableName,
		gr.SysId,
		gr.Rows,
		//make([]map[string]string, 1),
	}, nil
}

func (tableApi *TableAPI) GetRecords(tableName string, queryParams map[string]string) ([]GlideRecord, error) {
	return []GlideRecord{
		GlideRecord{
			tableName,
			"sysId",
			make(map[string]any, 1),
		},
	}, nil
}

func (tableApi *TableAPI) DeleteRecord(tableName string, sysId string) (bool, error) {
	return true, nil
}

func (tableApi *TableAPI) CreateRecord(gr *GlideRecord) (sysId string, err error) {

	jsonFields, err := json.Marshal(gr.Rows)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", tableApi.InstanceURL+tableApi.BasePath+gr.TableName, bytes.NewBuffer(jsonFields))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", tableApi.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := tableApi.HttpClient.Do(req)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	// Read and parse the JSON response
	var result GlideRecord // Or use a custom struct
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	err = json.Unmarshal(body, &result.Rows)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		return
	}

	processedFields := make(map[string]string)
	processRecordFields("", result.Rows["result"].(map[string]any), &processedFields)

	for k, v := range processedFields {
		fmt.Printf("Key: %s, Value: %s\n", k, v)
	}

	// Output the parsed response
	fmt.Println("Response:", result.Rows["result"].(map[string]any)["sys_id"])
	//sysId = result.Rows["result"].(map[string]any)["sys_id"].(string) // TODO: Need proper error checks and also this is disgusting, I wonder if there is a better way to do this
	sysId = processedFields["sys_id"]
	return sysId, nil
}

func (tableApi *TableAPI) UpdateRecord(gr *GlideRecord) (*GlideRecord, error) {
	return &GlideRecord{
		gr.TableName,
		gr.SysId,
		gr.Rows,
		//make([]map[string]string, 1),
	}, nil
}

// Recursively process a map and print key/value pairs
func processRecordFields(prefix string, m map[string]any, processedFields *map[string]string) {
	for key, value := range m {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key // Concatenate keys to show nested structure
		}

		switch v := value.(type) {
		case string:
			fmt.Printf("Key: %s, Value: %s\n", fullKey, v)
			(*processedFields)[fullKey] = v
		case float64:
			(*processedFields)[fullKey] = strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			(*processedFields)[fullKey] = strconv.FormatBool(v)
		case []interface{}:
			fmt.Printf("Key: %s, Value: %v (slice)\nTBH This isn't supposed to happen, idk how we got a slice from ServiceNow\n", fullKey, v)
		case map[string]interface{}:
			processRecordFields(fullKey, v, processedFields) // Recursively process nested maps
		default:
			fmt.Printf("Key: %s, Value: %v (unknown type)\n", fullKey, v)
		}
	}
}
