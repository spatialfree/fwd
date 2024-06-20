package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func rec(path string, duration time.Duration) {
	cmd := exec.Command("wf-recorder", "-f", path)
	cmd.Start()
	time.Sleep(duration)
	cmd.Process.Signal(os.Interrupt)
}

func main() {
	duration := 5 * time.Second     // 25 minutes
	intermission := 1 * time.Second // 5 minutes
	for i := 0; i < 4; i++ {
		// print system timestamp HH:MM
		fmt.Print(time.Now().Format("15:04"))
		fmt.Printf(" | %d/4 5s\n", i+1)
		rec(fmt.Sprintf("./%d.mkv", i+1), duration)
		exec.Command("notify-send", "fwd", "intermission").Start()
		time.Sleep(intermission)
	}

	// compile
	// cmd = [
	// 	'ffmpeg', '-y',
	// 	'-safe', '0',
	// 	'-f', 'concat',
	// 	'-i', 'list.txt',
	// 	'-c', 'copy',
	// 	'temp/compiled.mp4'
	// ]
}
