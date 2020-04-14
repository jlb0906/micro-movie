package handler

import (
	"context"
	"encoding/json"
	aria2 "github.com/jlb0906/micro-movie/aria2-srv/proto/aria2"
	"github.com/jlb0906/micro-movie/basic/common"
	"github.com/micro/go-micro/v2/client"
	"net/http"
	"time"
	//aria2 "path/to/service/proto/aria2"
)

func Aria2Call(w http.ResponseWriter, r *http.Request) {
	// decode the incoming request as json
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// call the backend service
	aria2Client := aria2.NewAria2Service(common.Aria2Srv, client.DefaultClient)
	rsp, err := aria2Client.AddURI(context.TODO(), &aria2.AddURIReq{
		Uri: request["uri"].(string),
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// we want to augment the response
	response := map[string]interface{}{
		"gid": rsp.Gid,
		"ref": time.Now().UnixNano(),
	}

	// encode and write the response as json
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
