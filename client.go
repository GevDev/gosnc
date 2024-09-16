package gosnc

import (
	"github.com/gevdev/gosnc/tableapi"
	"net/http"
)

type TableAPI[glideRecord any] interface {
	GetEmptyGlideRecord() glideRecord
	GetRecord(gr *glideRecord, queryParams map[string]string) (*glideRecord, error)
	GetRecords(tableName string, queryParams map[string]string) ([]glideRecord, error)
	DeleteRecord(tableName string, sysId string) (bool, error)
	CreateRecord(gr *glideRecord) (string, error)
	UpdateRecord(gr *glideRecord) (*glideRecord, error)
}

type NowClient struct {
	InstanceURL      string
	TokenHeaderValue string
	TableAPI         *TableAPI[tableapi.GlideRecord]
}

func NewNowClient(instanceUrl string, tokenHeaderValue string) (*NowClient, error) {

	httpClient := &http.Client{}
	var tableApi TableAPI[tableapi.GlideRecord] = tableapi.NewTableAPI(instanceUrl, tokenHeaderValue, "/api/now/table/", httpClient)

	return &NowClient{
		InstanceURL:      instanceUrl,
		TokenHeaderValue: tokenHeaderValue,
		TableAPI:         &tableApi,
	}, nil
}
