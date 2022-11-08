package gstreamer

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"mediaservice/apputils"
	"testing"
	"time"
)

const (
	//minport int = 10000
	streamfile string = "female_5min_8000_mono.wav"
)


func mreceiver(port int, codec,callid string, ipv6 bool) {

	log.Println("Create NEW Receiver:")
	var pipeline string
	var err error
	var gst *Pipeline

	bind := "0.0.0.0"
	if ipv6 == true {
		bind = "::"
	}
	switch codec {
	case "alaw":
		pipeline = fmt.Sprintf("udpsrc port=%d address=\"%s\" caps=\"application/x-rtp,payload=8,clock-rate=8000\" ! rtppcmadepay ! alawdec ! audioconvert ! audioresample ! audiorate ! wavenc ! filesink name=dstfile location=%s.wav", port, bind, callid)
	case "ulaw":
		pipeline = fmt.Sprintf("udpsrc port=%d address=\"%s\" caps=\"application/x-rtp,payload=0,clock-rate=8000\" ! rtppcmudepay ! mulawdec ! audioconvert ! audioresample ! audiorate ! wavenc ! filesink name=dstfile location=%s.wav", port, bind, callid)
	case "opus":
		pipeline = fmt.Sprintf("udpsrc port=%d address=\"%s\" caps=\"application/x-rtp,media=audio,clock-rate=48000\" ! queue !  rtpopusdepay ! opusdec plc=true ! audioconvert ! audioresample  ! audio/x-raw, rate=8000 ! audiorate ! wavenc ! filesink name=dstfile location=%s.wav", port, bind, callid)
	case "g722":
		pipeline = fmt.Sprintf("udpsrc port=%d address=\"%s\" caps=\"application/x-rtp,media=audio,encoding-name=G722, payload=9, encoding-params=1, clock-rate=8000\"  ! rtpg722depay ! avdec_g722 ! audioconvert ! audioresample ! audiorate !  wavenc ! filesink name=dstfile location=%s.wav", port, bind, callid)
	}
	fmt.Println("gst-launch-1.0",pipeline)
	gst, err = New(pipeline)
	if err != nil {
		log.Println("pipeline create error")
	}
	gst.Start()
	log.Println("receiver started")
}

func mstreamer(codec, rhost string, lport, rport int) {
	log.Println("Create NEW streamer:")
	var gst *Pipeline
	var pipeline string
	var err error
	filename := streamfile

	switch codec {
	case "alaw":
		pipeline = fmt.Sprintf("filesrc name=srcfile location=%s ! wavparse ! audioconvert ! audioresample ! alawenc ! rtppcmapay mtu=128 ! udpsink name=usink bind-port=%d host=%s port=%d", filename, lport, rhost, rport)
	case "ulaw":
		pipeline = fmt.Sprintf("filesrc name=srcfile location=%s ! wavparse ! audioconvert ! audioresample ! mulawenc ! rtppcmupay mtu=128 ! udpsink name=usink bind-port=%d host=%s port=%d", filename, lport, rhost, rport)
	case "opus":
		pipeline = fmt.Sprintf("filesrc name=srcfile location=%s ! wavparse ! audioconvert ! audioresample ! opusenc ! rtpopuspay mtu=128 ! udpsink name=usink bind-port=%d host=%s port=%d", filename, lport, rhost, rport)
	case "g722":
		pipeline = fmt.Sprintf("filesrc name=srcfile location=%s ! wavparse ! audioconvert ! audioresample ! avenc_g722 ! rtpg722pay mtu=128 ! udpsink name=usink bind-port=%d host=%s port=%d", filename, lport, rhost, rport)
	}
	gst, err = New(pipeline)
	if err != nil {
		log.Println("pipeline create error")
	}
	fmt.Println("gst-launch-1.0",pipeline)
	gst.Start()
	log.Println("streamer started")
}

func unitm(lhost string, lport int){
	callide := uuid.New().String()
	callidd := uuid.New().String()
	go mstreamer("alaw", lhost,lport,lport+1)
	go mreceiver(lport,"alaw",callide,true)
	go mstreamer("alaw", lhost,lport+1,lport)
	go mreceiver(lport+1,"alaw",callidd,true)
}

func TestGstreamer(t *testing.T) {
	localV6Addr := apputils.GetLocalIpv6().String()

	unitm(localV6Addr,10003)
	time.Sleep(10  *time.Second)

}