package client

import (
	"context"
	"fmt"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	api "github.com/metaprov/modelaapi/services/grpcinferenceservice/v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// SDK is an instance of the Agones SDK
type PredictorClient struct {
	client api.GRPCInferenceServiceClient
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
	s.client = api.NewGRPCInferenceServiceClient(conn)
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

func (r *PredictorClient) Shutdown() error {
	req := &api.ServerShutdownRequest{}
	_, err := r.client.Shutdown(r.ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed shutdown server")
	}
	return nil

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

// SDK is an instance of the Agones SDK
type PredictionRequest struct {
	client   api.GRPCInferenceServiceClient
	ctx      context.Context
	explain  bool // return an explenation for each prediction
	validate bool // using the predictor schema, validate the prediction
	model    string
	format   string
	payload  string
	labeled  bool
	metrics  []string
}

func NewPrediction() *PredictionRequest {
	return &PredictionRequest{
		ctx:     context.Background(),
		metrics: make([]string, 0),
	}
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

//TestResults answers the result of executing the test
func (res *PredictionResult) TestResults() (map[string]float32, error) {
	return res.raw.Scores, nil
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
	req.format = "json"
	return req
}

func (req *PredictionRequest) WithCsv(payload string) *PredictionRequest {
	req.payload = payload
	req.format = "csv"
	return req
}

func (req *PredictionRequest) WithLabeled() *PredictionRequest {
	req.labeled = true
	return req
}

func (req *PredictionRequest) WithMetrics(metrics []string) *PredictionRequest {
	req.metrics = metrics
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
	defer conn.Close()
	client := api.NewGRPCInferenceServiceClient(conn)
	grpcReq := &api.PredictRequest{
		Name:     "",
		Validate: req.validate,
		Explain:  req.explain,
		Format:   req.format,
		Payload:  req.payload,
		Labeled:  req.labeled,
		Metrics:  req.metrics,
	}

	result, err := client.Predict(req.ctx, grpcReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed call prediction")
	}

	return &PredictionResult{raw: result}, nil

}

func (req *PredictionRequest) GetPredictor(host string, port int) (*api.GetPredictorResponse, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	ctx, cancel := context.WithTimeout(req.ctx, 30*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to %s", addr)
	}
	defer conn.Close()
	client := api.NewGRPCInferenceServiceClient(conn)
	grpcReq := &api.GetPredictorRequest{}

	result, err := client.GetPredictor(req.ctx, grpcReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed call prediction")
	}

	return result, nil

}

func (r *PredictionRequest) GetModel(predictorName string, name string, host string, port int, ctx context.Context) (*api.GetModelResponse, error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "could not connect to %s", addr)
	}
	defer conn.Close()
	client := api.NewGRPCInferenceServiceClient(conn)
	grpcReq := &api.GetModelRequest{
		Name: name,
	}

	result, err := client.GetModel(ctx, grpcReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed GetModel")
	}
	return result, nil
}
