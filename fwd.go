package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	Duration     = 10 * time.Minute // 25 minutes
	Intermission = 5 * time.Minute
)

func rec(path string, duration time.Duration) {
	cmd := exec.Command("wf-recorder", "-f", path)
	cmd.Start()
	time.Sleep(duration)
	cmd.Process.Signal(os.Interrupt)
}

func main() {
	// remove old clips
	os.Remove("clips/1.mkv")
	os.Remove("clips/2.mkv")
	os.Remove("clips/3.mkv")
	os.Remove("clips/4.mkv")
	os.Remove("clips/5_compiled.mkv")
	os.Remove("clips/6_review.mkv")
	os.Remove("clips/7_audio.wav")
	os.Remove("clips/8_analysis.mkv")
	os.Remove("clips/9_analysis.mp4")

	// core loop
	for i := 0; i < 4; i++ {
		// print system timestamp HH:MM
		// fmt.Print(time.Now().Format("15:04"))
		// fmt.Printf(" | %d/4 5s\n", i+1)
		rec(fmt.Sprintf("./clips/%d.mkv", i+1), Duration)
		exec.Command("notify-send", "fwd", fmt.Sprintf("take 5, resume at %s", time.Now().Add(Intermission).Format("15:04"))).Run()
		time.Sleep(Intermission)
	}

	// compile video
	exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", "clips.txt", "-c", "copy", "clips/5_compiled.mkv").Run()

	// 10x visual review
	exec.Command("ffmpeg", "-i", "clips/5_compiled.mkv", "-vf", "setpts=0.1*PTS", "-r", "60", "clips/6_review.mkv").Run()
	exec.Command("mpv", "clips/6_review.mkv").Run()

	// voice over analysis
	audio_cmd := exec.Command("ffmpeg", "-f", "pulse", "-i", "default", "-acodec", "pcm_s16le", "clips/7_audio.wav")
	from_timestamp := time.Now()
	// cmd line input to stop audio recording
	_ = audio_cmd.Start()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("press [enter] to end audio recording...\n")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSuffix(content, "\n")
	fmt.Printf("%s preparing...\n", content)
	to_timestamp := time.Now()
	audio_cmd.Process.Signal(os.Interrupt)

	// scale video to match
	audio_duration := to_timestamp.Sub(from_timestamp)
	total_time := Duration * 4
	setpts := audio_duration.Seconds() / total_time.Seconds()
	fmt.Printf("setpts=%f*PTS\n", setpts)
	exec.Command("ffmpeg", "-i", "clips/5_compiled.mkv", "-vf", fmt.Sprintf("setpts=%f*PTS", setpts), "-r", "60", "clips/8_analysis.mkv").Run()

	exec.Command("ffmpeg", "-i", "clips/8_analysis.mkv", "-i", "clips/7_audio.wav", "-c:v", "copy", "-c:a", "aac", "-strict", "experimental", "clips/9_analysis.mp4").Run()

	fmt.Println("video ready at clips/9_analysis.mp4\n	now playing...")
	exec.Command("mpv", "clips/9_analysis.mp4").Run()
}
