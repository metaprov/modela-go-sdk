package integration

import (
	"github.com/metaprov/mdgoclient/pkg/client"
	"testing"
)

func Test_one(t *testing.T) {
	cl, _ := client.NewPredictorClient("127.0.0.1", 8080)
	cl.Alive()

}
