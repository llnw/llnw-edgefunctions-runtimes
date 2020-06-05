package main

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/llnw/llnw-edgefunctions-runtimes/go/edgefunction"
	"github.com/llnw/llnw-edgefunctions-runtimes/go/edgefunction/events"
)

// Host passed through env variable
type Host struct {
	Name   string
	Weight float64
}

func handler(request *events.EPInvokeRequest) (*events.EPInvokeResponse, error) {
	// Grab envvar from queries or set to LLNW_HOSTS
	envvar, ok := request.Queries["envvar"]
	if !ok {
		envvar = "LLNW_HOSTS"

	}
	// Look a for the set of hosts in the environment variables
	hostSet, exists := os.LookupEnv(envvar)
	if !exists {
		return nil, errors.New("environment variable not found")
	}

	// Grab path from queries
	path, ok := request.Queries["path"]
	if !ok {
		path = ""
	}

	hostBytes := []byte(hostSet)

	// convert environment variable into list of hosts
	var hosts []Host
	err := json.Unmarshal(hostBytes, &hosts)

	if err != nil {
		return nil, err
	}

	// Grab tag from queries
	tag, ok := request.Queries["tag"]
	// If tag is not provided, set the seed to the time
	if !ok {
		rand.Seed(time.Now().UTC().UnixNano())
	} else {
		// convert tag to int and set seed to the tag
		tagInt, err := strconv.ParseInt(tag, 10, 64)
		if err != nil {
			return nil, errors.New("Tag must be a valid integer")
		}
		rand.Seed(tagInt)
	}

	// generate random number in range [0.0, 1.0)
	val := rand.Float64()

	distributeWeighting(hosts)

	redirect := filterHosts(hosts, val) + "/" + path

	resp := &events.EPInvokeResponse{
		StatusCode: http.StatusOK,
		Body:       redirect,
	}

	return resp, nil

}

// return first host with cumulative weight greater than the random value
func filterHosts(hostList []Host, randomVal float64) string {
	for i := range hostList {
		if hostList[i].Weight > randomVal {
			return hostList[i].Name
		}
	}
	return hostList[len(hostList)].Name
}

// Set hosts weight to percentage of total weight,
func distributeWeighting(hostList []Host) {
	defaultWeight := 0.5
	cumulativeWeight := 0.0

	// get total weight of all hosts
	totalWeight := 0.0
	for i, h := range hostList {
		// If weight is 0, set it to default weight
		if h.Weight <= 0 {
			hostList[i].Weight = defaultWeight
		}
		totalWeight += hostList[i].Weight
	}

	// sort hostList from lowest to highest weight
	sort.Slice(hostList, func(i, j int) bool {
		return (hostList)[i].Weight < (hostList)[j].Weight
	})

	// normalize weights to sum to 1
	for i := range hostList {
		hostList[i].Weight = hostList[i].Weight / totalWeight
	}

	// set host weights to cumulative weights
	for i := range hostList {
		hostList[i].Weight += cumulativeWeight
		cumulativeWeight = hostList[i].Weight

	}

}

func main() {
	edgefunction.Start(handler)
}
