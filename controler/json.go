package controler

import "encoding/json"

func makeJSONResponse(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func jsonErrMessage(s string) string {
	return makeJSONResponse(jsonRep{s})
}

type jsonRep struct {
	Message string
}
