package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
