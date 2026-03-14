package httpx

import (
	"encoding/json"
	"net/http"
)

// Json 响应 JSON 数据
func Json(w http.ResponseWriter, code int32, desc string, data any) error {
	bytes, err := json.Marshal(struct {
		Code int32  `json:"code"`
		Desc string `json:"desc"`
		Data any    `json:"data"`
	}{
		Code: code,
		Desc: desc,
		Data: data,
	})
	if err != nil {
		return err
	}

	if _, err = w.Write(bytes); err != nil {
		return err
	}
	return nil
}
