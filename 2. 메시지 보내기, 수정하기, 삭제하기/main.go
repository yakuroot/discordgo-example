package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token   = "여기에 봇 토큰을 입력합니다."
	Session *discordgo.Session
)

func init() {
	var err error
	Session, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalf("클라이언트 생성 오류: %v", err)
	}

	err = Session.Open()
	if err != nil {
		log.Fatalf("세션 오픈 오류: %v", err)
	}

	log.Printf("%s (%s)에 로그인 됨", Session.State.User.String(), Session.State.User.ID)
}

func main() {
	Session.AddHandler(messageCreate)

	defer Session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("봇 종료됨")
}

func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if session.State.User.ID == message.Author.ID {
		return
	}

	if message.Content == "!dm" {
		// 유저에게 DM(1:1 메시지)을 보내려면, 먼저 UserChannelCreate 함수를 통해 유저 채널 값을 얻어야 합니다.
		// UserChannelCreate 함수는 *discordgo.Channel과 error 형식을 반환합니다.
		userChannel, _ := session.UserChannelCreate(message.Author.ID)

		// 유저 채널 값이 위에 선언한 변수 userChannel에 들어갔습니다.
		// ChannelMessageSend의 첫 번째 매개 변수에는 채널 아이디를 넣어야 하므로, userChannel.ID를 넣어줍니다.
		session.ChannelMessageSend(userChannel.ID, "hello, DM!")
	}

	if message.Content == "!hello" {
		// ChannelMessageSendReply 함수는 디스코드의 '답장'기능을 봇이 사용할 수 있도록 하는 함수입니다.
		// 첫 번째 매개 변수로는 메시지를 보낼 채널의 ID, 두 번째 매개 변수로는 보낼 메시지 내용을 입력해 줍니다.
		// 그리고 세 번째 매개 변수로는 *discordgo.MessageReference 형식을 넣어주는데요.
		// *discordgo.MessageReference 형식에는 MessageID, ChannelID, GuildID, 이렇게 3가지 항목이 있습니다.
		// 지난 예제에서 매개 변수로 받아 온 message에는 수신된 메시지의 정보가 들어있다고 말씀 드렸던 것, 기억 나시나요?
		// message.ID에는 메시지의 아이디, message.ChannelID에는 메시지가 수신된 채널의 아이디, message.GuildID에는 메시지가 수신된 서버의 아이디가 들어가 있습니다.
		// 이 항목들을 차례대로 *discordgo.MessageReference 형식에 입력해 주면 됩니다.

		// message.Author.String()에는 메시지를 보낸 사람의 닉네임#0000(태그)이 들어있습니다.
		session.ChannelMessageSendReply(message.ChannelID, "hello, "+message.Author.String(), message.Reference())
	}

	if message.Content == "!embed" {
		// 임베드를 보내기 위해서는 ChannelMessageSendEmbed 함수를 사용합니다.
		// 첫 번째 매개 변수로는 메시지를 보낼 채널의 ID, 두 번째 매개 변수로는 *discordgo.MessageEmbed 형식을 넣어 주어야 합니다.
		// *discordgo.MessageEmbed 형식에 대해서는 아래 코드 각 줄마다 주석을 통해 설명하겠습니다.
		session.ChannelMessageSendEmbed(message.ChannelID, &discordgo.MessageEmbed{
			// 임베드의 색상을 작성합니다. 색상 코드를 넣어주면 되는데, 쉽게 말하자면 #000000에서 #을 0x로만 바꾸어 넣으면 됩니다.
			// 이 항목이 없으면 임베드의 색상은 검정색으로 자동 지정됩니다.
			Color: 0xffcc66,

			// 임베드의 제목을 작성합니다.
			// 이 항목이 없으면 임베드의 제목이 없는 것으로 나타납니다.
			Title: "임베드 제목",

			// 임베드의 내용을 결정합니다.
			// 이 항목이 없으면 임베드의 내용이 없는 것으로 나타납니다.
			Description: "임베드 내용",

			// 임베드의 작성자 항목을 편집합니다.
			Author: &discordgo.MessageEmbedAuthor{
				// 작성자 항목을 클릭하면 나오는 링크를 써넣습니다. (선택 항목)
				URL: "https://github.com/Neoration/discordgo-example",
				// 작성자 항목의 이름을 써넣습니다. (필수 항목)
				Name: message.Author.String(),
				// 작성자 항목의 아이콘 URL을 써넣습니다. (선택 항목)
				IconURL: message.Author.AvatarURL(""),
			},

			// 임베드의 필드 항목을 편집합니다.
			// 필드는 *discordgo.MessageEmbedField의 집합체(배열; 슬라이스) 형태로 입력할 수 있습니다.
			Fields: []*discordgo.MessageEmbedField{
				{
					// 필드에서 Name, Value는 필수 작성 항목이며, Inline은 선택 작성 항목입니다.
					// Inline을 작성하지 않으면 기본 값은 false로 지정됩니다.
					Name:   "첫 번째 필드 이름",
					Value:  "첫 번째 필드 값",
					Inline: true,
				},
				{
					Name:   "두 번째 필드 이름",
					Value:  "두 번째 필드 값",
					Inline: true,
				},
				{
					Name:   "세 번째 필드 이름",
					Value:  "세 번째 필드 값",
					Inline: false,
				},
			},

			// 임베드에 사진을 넣습니다.
			Image: &discordgo.MessageEmbedImage{
				// URL에는 사진 링크, Width에는 사진의 폭, Height에는 사진의 높이를 작성합니다.
				// Width, Height는 생략할 수 있습니다.
				URL:    "https://github.com/Neoration/discordgo-example/raw/main/golang.jpg",
				Width:  600,
				Height: 300,
			},

			// 임베드에 썸네일을 넣습니다. (임베드 오른쪽 상단의 사진)
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				// URL에는 사진 링크, Width에는 사진의 폭, Height에는 사진의 높이를 작성합니다.
				// Width, Height는 생략할 수 있습니다.
				// message.Author.AvatarURL() 함수를 통해 메시지를 보낸 사람의 프로필 사진 링크를 얻을 수 있습니다.
				URL:    message.Author.AvatarURL(""),
				Width:  100,
				Height: 100,
			},

			// 임베드의 푸터를 설정합니다. (맨 아래에 나오는 작은 글씨)
			Footer: &discordgo.MessageEmbedFooter{
				// Text에는 푸터에 들어갈 내용, 그리고 IconURL에는 푸터 아이콘에 들어갈 사진 링크를 작성합니다.
				// IconURL은 선택 항목이지만, Text는 무조건 작성해야 하는 필수 항목입니다.
				Text:    message.Author.String(),
				IconURL: message.Author.AvatarURL(""),
			},

			// 푸터의 Text란이 끝나면 해당 메시지가 보내진 시간을 알려주는 TImestamp를 설정할 수 있습니다.
			// DiscordGo에서는 아래와 동일하게 작성하여야 Timestamp가 제대로 작동힙니다.
			Timestamp: time.Now().Format(time.RFC3339),
		})
	}

	if message.Content == "!countdown" {
		// ChannelMessageSend 함수는 메시지를 보내는 역할 뿐만이 아닌, *Message 형식과 error 형식을 내보내는 역할도 합니다.
		// 즉, 아래와 같이 변수로 ChannelMessageSend 함수를 받아주게 되면
		// 아래 msg 변수 내에는 봇이 전송한 메시지의 정보가 담기게 됩니다.
		msg, _ := session.ChannelMessageSend(message.ChannelID, "3초 뒤 메시지가 삭제됩니다...")

		for i := 2; i > 0; i-- {
			// time 패키지와 strconv 패키지에 대한 설명은 간략히 하고 넘어가겠습니다.
			// 바로 아래 코드는 1초 간 코드 실행을 지연시키는 코드입니다.
			// 그리고 strconv.Itoa()는 숫자를 문자열로 바꿔주는 함수이고요.
			time.Sleep(time.Second * time.Duration(1))

			// ChannelMessageEdit 함수는 특정 채널에 있는 특정 메시지를 수정하는 역할을 합니다.
			// 첫 번째 매개 변수로는 수정할 메시지가 있는 채널의 ID, 두 번째 매개 변수로는 수정할 메시지의 ID, 세 번째 매개 변수로는 수정할 메시지 내용을 넣어주면 됩니다.
			// 단, 여러분이 메시지를 수정할 때 본인의 메시지만 수정할 수 있는 것처럼 봇도 봇 본인이 보낸 메시지만 수정할 수 밖에 없습니다.
			// 아까 봇이 보낸 메시지의 정보를 담아 놓은(L152) msg 변수가 있었습니다.
			// msg.ChannelID는 봇이 보낸 메시지가 어느 채널로 보내졌는지, msg.ID는 봇이 보낸 메시지의 ID를 담고 있겠죠.
			// 이 둘을 ChannelMessageEdit 함수에 차례대로 넘겨줍니다.
			session.ChannelMessageEdit(msg.ChannelID, msg.ID, strconv.Itoa(i)+"초 뒤 메시지가 삭제됩니다...")

			// + 만약 수정하고자 하는 메시지가 텍스트 메시지가 아니라 임베드 형식 메시지라면,
			// session.ChannelMessageEditEmbed() 함수를 사용하여야 합니다.
			// 매개 변수는 수정할 메시지 내용을 *discordgo.MessageEmbed 형식으로만 바꿔주는 것 빼고는 위 ChannelMessageEdit 함수와 동일합니다.
		}
		time.Sleep(time.Second * time.Duration(1))

		// ChannelMessageDelete 함수는 봇이 메시지를 삭제하도록 하는 역할을 합니다.
		// 첫 번째 매개 변수로는 삭제할 메시지가 있는 채널의 ID, 두 번째 매개 변수로는 삭제할 메시지의 ID를 넣어주면 됩니다.
		// 단, 봇에게 메시지를 삭제할 수 있는 권한이 없는 경우 ChannelMessageDelete 함수는 실행되지 않으니,
		// 만약 매개 변수를 제대로 작성했는데도 불구하고 메시지가 삭제되지 않는 경우 봇의 권한을 먼저 확인해 주세요.
		session.ChannelMessageDelete(msg.ChannelID, msg.ID)
	}
}
