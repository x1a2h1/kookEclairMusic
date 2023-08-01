package kook

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/x1a2h1/kookvoice"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type VoiceInstance struct {
	Token         string
	ChannelId     string
	wsConnect     *websocket.Conn
	streamProcess *os.Process
	sourceProcess *os.Process
}

func (i *VoiceInstance) Init() error {
	makeFifoCmd := exec.Command(
		"mkfifo",
		"streampipe"+i.ChannelId,
	)
	err := makeFifoCmd.Run()
	if err != nil {
		return err
	}

	keepFifoOpenCmd := exec.Command(
		"bash",
		"-c",
		"exec 7<>streampipe"+i.ChannelId,
	)
	err = keepFifoOpenCmd.Run()
	if err != nil {
		return err
	}

	silentSourceCmd := exec.Command(
		"bash",
		"-c",
		"ffmpeg -f lavfi -i anullsrc -f wav -c:a pcm_s16le -b:a 1411200 -ar 44.1k -ac 2 pipe:1 > streampipe"+i.ChannelId,
	)
	silentSourceCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err = silentSourceCmd.Start()
	if err != nil {
		return err
	}
	i.sourceProcess = silentSourceCmd.Process

	gatewayUrl := kookvoice.GetGatewayUrl(i.Token, i.ChannelId)
	connect, rtpUrl := kookvoice.InitWebsocketClient(gatewayUrl)

	go kookvoice.KeepWebsocketClientAlive(connect)
	go kookvoice.KeepRecieveMessage(connect)

	i.wsConnect = connect

	streamCmd := exec.Command(
		"ffmpeg",
		"-re",
		"-loglevel",
		"level+info",
		"-nostats",
		"-i",
		"streampipe"+i.ChannelId,
		"-map",
		"0:a:0",
		"-acodec",
		"libopus",
		"-ab",
		"128k",
		"-filter:a",
		"volume=0.8",
		"-ac",
		"2",
		"-ar",
		"48000",
		"-f",
		"tee",
		fmt.Sprintf("[select=a:f=rtp:ssrc=1357:payload_type=100]%v", rtpUrl),
	)
	err = streamCmd.Start()
	if err != nil {
		return err
	}
	i.streamProcess = streamCmd.Process

	return nil
}
func (i *VoiceInstance) PlayMusic(input string) error {
	time.Sleep(500 * time.Millisecond)
	if err := syscall.Kill(-i.sourceProcess.Pid, syscall.SIGKILL); err != nil {
		return errors.New(fmt.Sprintf("无法终止源进程, err: %v", err))
	}
	//if i.sourceProcess != nil {
	//	if err := syscall.Kill(-i.sourceProcess.Pid, syscall.SIGKILL); err != nil {
	//		return errors.New(fmt.Sprintf("终止音频源出错:%v", err))
	//	}
	//}

	musicSourceCmd := exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("ffmpeg -re -i %v -f s16le -c:a pcm_s16le -b:a 1411200 -ar 44.1k -ac 2 pipe:1 > streampipe"+i.ChannelId, input),
	)
	musicSourceCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err := musicSourceCmd.Start()
	if err != nil {
		return errors.New(fmt.Sprintf("启动音乐进程失败, err: %v", err))
	}
	i.sourceProcess = musicSourceCmd.Process

	err = musicSourceCmd.Wait()
	if err != nil {
		return errors.New(fmt.Sprintf("等待音乐处理失败, err: %v", err))
	}

	silentSourceCmd := exec.Command(
		"bash",
		"-c",
		"ffmpeg -f lavfi -i anullsrc -f wav -c:a pcm_s16le -b:a 1411200 -ar 44.1k -ac 2 pipe:1 > streampipe"+i.ChannelId,
	)
	silentSourceCmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	err = silentSourceCmd.Start()
	if err != nil {
		return errors.New(fmt.Sprintf("启动静默源进程失败, err: %v", err))
	}
	i.sourceProcess = silentSourceCmd.Process

	return nil
}

func (i *VoiceInstance) Close() error {
	if _, err := os.Stat("streampipe" + i.ChannelId); errors.Is(err, os.ErrNotExist) {
		fmt.Println("匿名管道不存在")
	}
	err := os.Remove("streampipe" + i.ChannelId)
	if err != nil {
		return err
	}
	return nil
}
