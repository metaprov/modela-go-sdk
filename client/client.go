package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	api "github.com/metaprov/modeldapi/services/predictiond/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// SDK is an instance of the Agones SDK
type PredictorClient struct {
	client api.PredictionServerClient
	ctx    context.Context
	host   string
	port   int32
}

func NewPredictorClient(host string, port int32) (*PredictorClient, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	s := &PredictorClient{
		ctx:  context.Background(),
		host: host,
		port: port,
	}
	// block for at least 30 seconds
	ctx, cancel := context.WithTimeout(s.ctx, 30*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return s, errors.Wrapf(err, "could not connect to %s", addr)
	}
	s.client = api.NewPredictionServerClient(conn)
	return s, errors.Wrap(err, "could not set up health check")
}

func (r *PredictorClient) Ready() (bool, error) {
	req := &api.ServerReadyRequest{}
	res, err := r.client.ServerReady(r.ctx, req)
	if err != nil {
		return false, errors.Wrap(err, "failed not send Ready message")
	}
	return res.Ready, nil

}

func (r *PredictorClient) Alive() (bool, error) {
	req := &api.ServerLiveRequest{}
	res, err := r.client.ServerLive(r.ctx, req)
	if err != nil {
		return false, errors.Wrap(err, "failed not send Ready message")
	}
	return res.Live, nil
}

func (r *PredictorClient) ModelReady(name string, version string) (bool, error) {
	req := &api.ModelReadyRequest{
		Name:    name,
		Version: version,
	}
	res, err := r.client.ModelReady(r.ctx, req)
	if err != nil {
		return false, errors.Wrap(err, "failed not send ready message")
	}
	return res.Ready, nil
}

func (r *PredictorClient) Predict(colsJson string, dataJson string, full bool) (string, error) {
	req := &api.PredictRequest{
		Name:     "",
		Validate: false,
		Explain:  false,
		Format:   0,
		Payload:  "",
	}

	result, err := r.client.Predict(r.ctx, req)
	if err != nil {
		return "", errors.Wrap(err, "failed prediction")
	}
	res, err := json.Marshal(result.Items)
	return string(res), err
}

// SDK is an instance of the Agones SDK
type PredictionRequest struct {
	client   api.PredictionServerClient
	ctx      context.Context
	explain  bool // return an explenation for each prediction
	validate bool // using the predictor schema, validate the prediction
	model    string
	payload  string
}

func NewPrediction() *PredictionRequest {
	return &PredictionRequest{ctx: context.Background()}
}

type PredictionResult struct {
	raw *api.PredictResponse
}

// return raw json
func (res *PredictionResult) AsJson() (string, error) {
	marshaler := jsonpb.Marshaler{}
	resultPayload, err := marshaler.MarshalToString(res.raw)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal result")
	}
	return resultPayload, nil

}

type PredictionResultLineItem struct {
	raw api.PredictResultLineItem
}

func (req *PredictionRequest) Explain() *PredictionRequest {
	req.explain = true
	return req
}

func (req *PredictionRequest) Validate() *PredictionRequest {
	req.validate = true
	return req
}

func (req *PredictionRequest) WithJson(payload string) *PredictionRequest {
	req.payload = payload
	return req
}

// send the request to the predictor
func (req *PredictionRequest) SendInsecure(host string, port int) (*PredictionResult, error) {
	// set default context
	// block for at least 30 seconds
	addr := fmt.Sprintf("%s:%d", host, port)
	ctx, cancel := context.WithTimeout(req.ctx, 30*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to %s", addr)
	}
	client := api.NewPredictionServerClient(conn)
	grpcReq := &api.PredictRequest{
		Name:     "",
		Validate: req.validate,
		Explain:  req.explain,
		Format:   api.PredictFormat_PREDICT_FORMAT_JSON,
		Payload:  req.payload,
	}

	result, err := client.Predict(req.ctx, grpcReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed call prediction")
	}

	return &PredictionResult{raw: result}, nil

}

func (r *PredictionRequest) GetPredictor(ctx context.Context, in *api.GetPredictorRequest) (*api.GetPredictorResponse, error) {
	return &api.GetPredictorResponse{}, nil

}

func (r *PredictionRequest) GetModel(ctx context.Context, in *api.GetModelRequest) (*api.GetModelRequest, error) {
	return &api.GetModelRequest{}, nil
}
