package mw

import (
	"context"
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var m = protojson.MarshalOptions{
	EmitUnpopulated: true,
	UseProtoNames:   true,
}

func init() {
	hlog.Info("using protojson")
	// TODO: fix int64 problem to use protojson back.
	//render.ResetJSONMarshal(marshal)
}

func ProtoJsonMiddleware() app.HandlerFunc {
	return func(ctx context.Context, rc *app.RequestContext) {
		rc.Next(ctx)
	}
}

func marshal(v interface{}) ([]byte, error) {
	switch v := v.(type) {
	case proto.Message:
		return m.Marshal(v)
	default:
		return json.Marshal(v)
	}
}
