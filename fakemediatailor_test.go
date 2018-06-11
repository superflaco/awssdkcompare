package fakemediatailor

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"testing"
	"time"
)

var facepalm_config = MediaTailorConfiguration{AdDecisionServerUrl: "http://7baa3.v.fwmrm.net/ad/g/1?nw=506531\u0026mode=live\u0026prof=96749:global-cocoa\u0026caid=facepalm_tv\u0026csid=ad_plethora\u0026resp=vast2\u0026pvrn=[avail.random]\u0026vprn=[avail.random]\u0026vdty=variable\u0026vrdu=[session.avail_duration_secs];;slid=mid\u0026tpcl=MIDROLL\u0026ptgt=a\u0026cpsq=[avail_num]\u0026mind=[session.avail_duration_secs]\u0026maxd=[session.avail_duration_secs]",
	VideoContentSourceUrl: "http://edge.theplatform-live-test.top.comcast.net/demo/superflaco/superstream/"}

var kwality_config = MediaTailorConfiguration{AdDecisionServerUrl: "http://7baa3.v.fwmrm.net/ad/g/1?nw=506531\u0026mode=live\u0026prof=96749:global-cocoa\u0026caid=kwality_video\u0026csid=ad_plethora\u0026resp=vast2\u0026pvrn=[avail.random]\u0026vprn=[avail.random]\u0026vdty=variable\u0026vrdu=[session.avail_duration_secs];;slid=mid\u0026tpcl=MIDROLL\u0026ptgt=a\u0026cpsq=[avail_num]\u0026mind=[session.avail_duration_secs]\u0026maxd=[session.avail_duration_secs]",
	VideoContentSourceUrl: "http://edge.theplatform-live-test.top.comcast.net/demo/superflaco/superstream/"}

func TestMediaTailor_GetConfigRequest(t *testing.T) {
	config, configErr := external.LoadDefaultAWSConfig()
	if nil == configErr {
		config.Region = MT_DEFAULT_REGION
		config.LogLevel = aws.LogDebugWithHTTPBody
		fmt := New(config)
		//fmt.AddDebugHandlers()
		req := fmt.GetConfigRequest("SuperStream")
		//t.Log(" URL:", req.HTTPRequest.URL.String())
		sendErr := req.Send()
		if nil == sendErr {
			mtConfig := req.Data.(*MediaTailorConfiguration)
			if "" != mtConfig.Playback() {
				t.Log("Playback URL Prefix: ", mtConfig.Playback())
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

func TestMediaTailor_RoundTripConfigRequest(t *testing.T) {
	config, configErr := external.LoadDefaultAWSConfig()
	if nil == configErr {
		config.Region = "us-east-1"
		config.LogLevel = aws.LogDebugWithHTTPBody
		fmt := New(config)
		fmt.AddDebugHandlers()
		tts := "TemporaryTestStream"
		putReq := fmt.PutConfigRequest(tts, kwality_config)
		//t.Log(" URL:", req.HTTPRequest.URL.String())
		var putErr = putReq.Send()
		if nil == putErr {
			time.Sleep(time.Second / 4)
			getReq := fmt.GetConfigRequest(tts)
			sendErr := getReq.Send()
			if nil == sendErr {
				mtConfig := getReq.Data.(*MediaTailorConfiguration)
				if "" != mtConfig.Playback() {
					t.Log("Playback URL Prefix: ", mtConfig.Playback())
					time.Sleep(time.Second / 4)
					delErr := fmt.DeleteConfigRequest(tts).Send()
					if nil != delErr {
						t.Log(delErr)
						t.Fail()
					}
				} else {
					t.Log("failed to find the playback url in the returned media tailor configuration", mtConfig)
					t.Fail()
				}
			} else {
				t.Error(sendErr)
				t.Fail()
			}
		} else {
			t.Log(putErr)
			t.Fail()
		}
	} else {
		t.Error(configErr)
		t.Fail()
	}

}
