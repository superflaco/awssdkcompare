package fakemediatailor

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/awserr"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

const opDescribeMTConfig = "GetConfig"

type MediaTailorConfiguration struct {
	AdDecisionServer     string `json:"adDecisionServer,omitempty"`
	OriginPrefix         string `json:"originPrefix,omitempty"`
	SlateURL             string `json:"slateURL,omitempty"`
	CDNContentSegmentURL string `json:"cdnContentSegmentURL,omitempty"`
	CDNAdSegmentURL      string `json:"cdnAdSegmentURL,omitempty"`
	HLSDiscSequence      bool   `json:"hlsDiscSequence,omitempty"`
	DashManifestPrefix   string `json:"dashManifestPrefix,omitempty"`
	PlaybackEndpoint     string `json:"playbackEndpoint,omitempty"`
}

func (mtc MediaTailorConfiguration) String() string {
	pretty, err := json.MarshalIndent(mtc, "", "  ")
	if nil == err {
		return string(pretty)
	} else {
		return err.Error()
	}
}

type MediaTailor struct {
	*aws.Client
}

// Build builds a JSON payload for a JSON RPC request.
func Build(req *aws.Request) {
	req.HTTPRequest.URL.Host = "api." + req.HTTPRequest.URL.Host
	req.HTTPRequest.Header.Add("Content-Type", "application/json; charset=UTF-8")

}

func Unmarshal(req *aws.Request) {
	defer req.HTTPResponse.Body.Close()
	var buf bytes.Buffer

	count, bufErr := buf.ReadFrom(req.HTTPResponse.Body)
	if count > 0 && nil == bufErr {
		data := buf.Bytes()
		switch req.HTTPResponse.StatusCode {
		case 200:
			err := json.Unmarshal(data, req.Data)
			if err != nil {
				req.Error = awserr.New("SerializationError", "failed decoding JSON RPC response", err)
			}
		case 403:
			req.Error = awserr.New("DOH!", string(data), nil)
		}
	}

	return
}

// Used for custom client initialization logic
var initClient func(*MediaTailor)

// Used for custom request initialization logic
var initRequest func(*MediaTailor, *aws.Request)

// Service information constants
const (
	ServiceName = "mediatailor" // Service endpoint prefix API calls made to.
	EndpointsID = ServiceName   // Service ID for Regions and Endpoints metadata.
)

func New(config aws.Config) *MediaTailor {
	var signingName string
	signingName = "mediatailor"
	signingRegion := config.Region

	svc := &MediaTailor{
		Client: aws.NewClient(
			config,
			aws.Metadata{
				ServiceName:   ServiceName,
				SigningName:   signingName,
				SigningRegion: signingRegion,
				APIVersion:    "2017-09-14",
				JSONVersion:   "1.1",
				TargetPrefix:  "MediaTailor_20170914",
			},
		),
	}

	// Handlers
	svc.Handlers.Sign.PushBackNamed(v4.SignRequestHandler)
	svc.Handlers.Build.PushBack(Build)
	svc.Handlers.Unmarshal.PushBack(Unmarshal)
	svc.Handlers.UnmarshalMeta.PushBack(Unmarshal)
	svc.Handlers.UnmarshalError.PushBack(Unmarshal)

	// Run custom client initialization if present
	if initClient != nil {
		initClient(svc)
	}

	return svc
}

// newRequest creates a new request for a MediaTailor operation and runs any
// custom request initialization.
func (c *MediaTailor) newRequest(op *aws.Operation, params, data interface{}) *aws.Request {
	req := c.NewRequest(op, params, data)

	// Run custom request initialization if present
	if initRequest != nil {
		initRequest(c, req)
	}

	return req
}

func (c *MediaTailor) GetConfigRequest(config string) *aws.Request {
	op := &aws.Operation{
		Name:       opDescribeMTConfig,
		HTTPMethod: "GET",
		HTTPPath:   "/v1/config/account/" + config,
	}
	var mtconfig MediaTailorConfiguration
	return c.newRequest(op, nil, &mtconfig)
}
