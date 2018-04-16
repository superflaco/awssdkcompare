package fakemediatailor

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"testing"
)

func TestMediaTailor_GetConfigRequest(t *testing.T) {
	config, configErr := external.LoadDefaultAWSConfig()
	if nil == configErr {
		config.Region = "us-east-1"
		//config.LogLevel = aws.LogDebug
		fmt := New(config)
		//fmt.AddDebugHandlers()
		req := fmt.GetConfigRequest("SuperStream")
		//t.Log(" URL:", req.HTTPRequest.URL.String())
		sendErr := req.Send()
		if nil == sendErr {
			mtConfig := req.Data.(*MediaTailorConfiguration)
			if "" != mtConfig.DashManifestPrefix {
				t.Log("Playback URL Prefix: ", mtConfig.DashManifestPrefix)
			} else {
				t.Log("failed to find the playback url in the returned media tailor configuration", mtConfig)
				t.Fail()
			}
		} else {
			t.Error(sendErr)
			t.Fail()
		}
	} else {
		t.Error(configErr)
		t.Fail()
	}
}
