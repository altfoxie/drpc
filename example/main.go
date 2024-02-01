package main

import (
	"time"

	"github.com/altfoxie/drpc"
)

func main() {
	client := drpc.New("975346661540909056")

	err := client.SetActivity(drpc.Activity{
		Details: "Details",
		State:   "State",
		Timestamps: &drpc.Timestamps{
			// Start: time.Now().Add(-time.Minute * 5),
			End: time.Now().Add(time.Minute * 5),
		},
		Assets: &drpc.Assets{
			LargeImage: "music",
			LargeText:  "Large Image Text",
		},
		Buttons: []drpc.Button{
			{
				Label: "Google",
				URL:   "https://youtube.com",
			},
			{
				Label: "Discord",
				URL:   "https://discord.com",
			},
		},
		Party: &drpc.Party{
			ID:   "12345",
			Size: [2]int{1, 6},
		},
	})
	if err != nil {
		panic(err)
	}

	// Discord keeps activity as long as the client is connected.
	time.Sleep(time.Second * 30)
}
