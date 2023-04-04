package imgs

import (
	"fmt"
	"github.com/yanzay/tbot/v2"
)

func SendImage(m *tbot.Message, client *tbot.Client) {
	_, err := client.SendPhotoFile(m.Chat.ID, `./imgs/default_image.jpg`, tbot.OptCaption("this is image"))
	if err != nil {
		fmt.Println(err)
	}
}

func SaveFile(m *tbot.Message, client *tbot.Client) {
	// here we check if message contains Document
	// you could also check for other types of files:
	// Audio, Photo, Video, etc.

}
