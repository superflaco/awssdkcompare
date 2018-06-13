package fakemediatailor

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/mediatailor"
)

var facepalm_config = MediaTailorConfiguration{AdDecisionServerUrl: "http://7baa3.v.fwmrm.net/ad/g/1?nw=506531\u0026mode=live\u0026prof=96749:global-cocoa\u0026caid=facepalm_tv\u0026csid=ad_plethora\u0026resp=vast2\u0026pvrn=[avail.random]\u0026vprn=[avail.random]\u0026vdty=variable\u0026vrdu=[session.avail_duration_secs];;slid=mid\u0026tpcl=MIDROLL\u0026ptgt=a\u0026cpsq=[avail_num]\u0026mind=[session.avail_duration_secs]\u0026maxd=[session.avail_duration_secs]",
	VideoContentSourceUrl: "http://edge.theplatform-live-test.top.comcast.net/demo/superflaco/superstream/", SlateAdURL: "http://pmd205470tn.download.theplatform.com.edgesuite.net/Tech_Summit_2016__Live/638/205/Free_Clouds_Timelapse_Stock_loop_1824837687_mp4_video_1280x720_1800000_primary_audio_5.mp4"}

var kwality_config = MediaTailorConfiguration{AdDecisionServerUrl: "http://7baa3.v.fwmrm.net/ad/g/1?nw=506531\u0026mode=live\u0026prof=96749:global-cocoa\u0026caid=kwality_video\u0026csid=ad_plethora\u0026resp=vast2\u0026pvrn=[avail.random]\u0026vprn=[avail.random]\u0026vdty=variable\u0026vrdu=[session.avail_duration_secs];;slid=mid\u0026tpcl=MIDROLL\u0026ptgt=a\u0026cpsq=[avail_num]\u0026mind=[session.avail_duration_secs]\u0026maxd=[session.avail_duration_secs]",
	VideoContentSourceUrl: "http://edge.theplatform-live-test.top.comcast.net/demo/superflaco/superstream/"}

func prepDefaultConfig() (aws.Config, error) {

	config, configErr := external.LoadDefaultAWSConfig()
	if nil == configErr {
		config.Region = "us-east-1"
		config.LogLevel = aws.LogDebugWithHTTPBody
	}
	return config, configErr
}

func TestFakeSDKRoundtrip(t *testing.T) {

	config, configErr := prepDefaultConfig()
	if configErr != nil {
		t.Error(configErr)
		t.Fail()
	}

	faketailor := New(config)
	faketailor.AddDebugHandlers()
	tts := "FakeSDKTestStream"

	putReq := faketailor.PutConfigRequest(tts, kwality_config)
	putErr := putReq.Send()
	if putErr != nil {
		t.Log(putErr)
		t.Fail()
	}

	// sleeping to avoid throttling limits
	time.Sleep(time.Second / 4)
	getReq := faketailor.GetConfigRequest(tts)

	sendErr := getReq.Send()
	if sendErr != nil {
		t.Error(sendErr)
		t.Fail()
	}

	mtConfig := getReq.Data.(*MediaTailorConfiguration)
	if mtConfig.Playback() == "" {
		t.Log("failed to find the playback url in the returned media tailor configuration", mtConfig)
		t.Fail()
	}

	t.Log("Playback URL Prefix: ", mtConfig.Playback())
	time.Sleep(time.Second / 4)

	delErr := faketailor.DeleteConfigRequest(tts).Send()
	if delErr != nil {
		t.Log(delErr)
		t.Fail()
	}
}

func TestRealSDKRoundtrip(t *testing.T) {

	config, configErr := prepDefaultConfig()
	if configErr != nil {
		t.Error(configErr)
		t.Fail()
	}

	realtailor := mediatailor.New(config)
	tts := "RealSDKTestStream"

	// As usual, the AWS API wants string pointers for no apparent reason since https://golang.org/pkg/encoding/json/#Marshal provides the following to omit empty values:
	// 'The "omitempty" option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.'
	//
	// In this case, users are required to create an input and then a request which seems a bit cumbersome
	// especially when the put input is basically a subset of  the get input so folks are stuck converting between
	// essentially identical objects
	putInput := &mediatailor.PutPlaybackConfigurationInput{
		Name: &tts,
		VideoContentSourceUrl: &facepalm_config.VideoContentSourceUrl,
		AdDecisionServerUrl:   &facepalm_config.AdDecisionServerUrl,
		SlateAdUrl:            &facepalm_config.SlateAdURL,
	}

	putReq := realtailor.PutPlaybackConfigurationRequest(putInput)
	putResp, putErr := putReq.Send()
	if putErr != nil {
		t.Log(putErr)
		t.Fail()
	}
	// sleeping to avoid throttling limits
	time.Sleep(time.Second / 4)
	// note that the putResp is a *PutPlaybackConfigurationOutput which sure looks a lot like a *GetPlaybackConfigurationOutput yet we are stuck using two different structs
	puthlsconf := putResp.HlsConfiguration

	if puthlsconf == nil || puthlsconf.ManifestEndpointPrefix == nil {
		t.Log("Failed to find a ManifestEndpointPrefix in the HlsConfiguration (PUT)")
		t.Fail()
	}

	putManiPrefix := *puthlsconf.ManifestEndpointPrefix

	if putManiPrefix == "" {
		t.Log("Failed to get a Prefix from PUT call")
		t.Fail()
	}

	t.Log("Got Prefix from PUT call:", putManiPrefix)
	// ok, we found a playback/manifest prefix, now try a get to check we got same thing
	// once again we have to setup an input object which really just has a single Name field
	getInput := &mediatailor.GetPlaybackConfigurationInput{Name: &tts}
	getReq := realtailor.GetPlaybackConfigurationRequest(getInput)
	getResp, getErr := getReq.Send()

	if getErr != nil {
		t.Log(getErr)
		t.Fail()
	}

	// we could easily make a function to handle both Get and Put responses if they were the
	// same struct, at least the HlsConfiguration is consistent between them
	gethlsconf := getResp.HlsConfiguration

	if gethlsconf == nil || gethlsconf.ManifestEndpointPrefix == nil {
		t.Log("Failed to find a ManifestEndpointPrefix in the HlsConfiguration (GET)")
		t.Fail()
	}

	getManiPrefix := *gethlsconf.ManifestEndpointPrefix

	if getManiPrefix == "" {
		t.Log("Failed to get a Prefix from GET call")
		t.Fail()
	}

	t.Log("Got Prefix from GET call: ", getManiPrefix)

	time.Sleep(time.Second / 4)
	deleteInput := &mediatailor.DeletePlaybackConfigurationInput{Name: &tts}
	deleteReq := realtailor.DeletePlaybackConfigurationRequest(deleteInput)

	_, deleteErr := deleteReq.Send()
	if deleteErr != nil {
		t.Log(deleteErr)
		t.Fail()
	}
}
