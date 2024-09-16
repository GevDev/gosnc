package gosnc

type NowRequest struct {
	Client *NowClient
}

func (nr *NowRequest) get(endpoint string, params map[string]string, headers map[string]string) {

}
