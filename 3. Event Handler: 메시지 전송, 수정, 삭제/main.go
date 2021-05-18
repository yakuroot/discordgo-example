package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	Token = "여기에 봇 토큰을 입력합니다."
	// 이벤트 핸들러를 이해하기 위해 이번 예제는 "로그 전송하기"를 만들어 보려고 합니다.
	// 이벤트가 감지되었을 때, 메시지를 보낼 채널 아이디를 입력해 주세요.
	ChannelID = "여기에 채널 아이디를 입력합니다."
	Session   *discordgo.Session
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
	// 수정되기 전, 또는 삭제되기 전 메시지의 정보를 불러오려면 메시지를 캐시하여야 합니다.
	// (*메시지 캐시: 메시지 정보를 클라이언트에 임시로 보관해 두는 것 - 이를 통해 수정, 삭제되기 전 메시지 정보를 불러올 수 있음)
	// DiscordGo에서는 기본적으로 캐시되는 메시지가 없기 때문에 직접 캐시 할 메시지 수를 정해줘야 하는데요.
	// 캐시 할 메시지 수는 Session.State.MaxMessageCount를 통해 정할 수 있습니다.
	// 이번 예제에서는 메시지 100개를 캐시하도록 하겠습니다.
	Session.State.MaxMessageCount = 100

	Session.AddHandler(messageUpdateHandler)
	Session.AddHandler(messageDeleteHandler)

	defer Session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("봇 종료됨")
}

// 메시지가 수정된 것을 비롯하여 메시지의 상태가 변경(예: 메시지 고정 등)됐을 때 이를 감지하려면
// *discordgo.Session 및 *discordgo.MessageUpdate 형식을 매개 변수로 받는 함수가 필요합니다.
func messageUpdateHandler(session *discordgo.Session, message *discordgo.MessageUpdate) {
	// 상태가 변경되기 이전의 메시지 정보를 message.BeforeUpdate를 통해 가져올 수 있습니다.
	// 물론 캐시되지 않은 메시지의 경우 이 값은 비어있게 됩니다.
	oldMessage := message.BeforeUpdate
	// 상태가 변경된 후의 메시지 정보는 그냥 message에 들어있습니다.
	// 이번 예제에서는 이해를 돕기 위해 따로 변수를 선언하겠습니다.
	newMessage := message

	// 상태가 변경되기 이전의 메시지 정보가 비어있다면,
	// 아래 메시지를 보내고 함수를 종료합니다.
	if oldMessage == nil {
		session.ChannelMessageSend(ChannelID, "메시지의 상태가 변경되었습니다.")
		return
	}

	// 그러나 이전 메시지 정보가 비어있지 않고, 메시지의 내용이 서로 다르다면
	// 이전 메시지와 새로운 메시지의 내용을 보여주는 메시지를 보냅니다.
	if newMessage.Content != oldMessage.Content {
		session.ChannelMessageSend(ChannelID, "메시지가 수정되었습니다.\n이전 메시지: "+oldMessage.Content+"\n바뀐 메시지: "+newMessage.Content)
	}
}

// 메시지가 삭제된 것을 감지하려면
// *discordgo.Session 및 *discordgo.MessageDelete 형식을 매개 변수로 받는 함수가 필요합니다.
func messageDeleteHandler(session *discordgo.Session, message *discordgo.MessageDelete) {
	// 삭제되기 전 메시지 정보를 message.BeforeDelete를 통해 가져올 수 있습니다.
	// 물론 캐시되지 않은 메시지의 경우 이 값은 비어있게 됩니다.
	deletedMessage := message.BeforeDelete

	// 삭제되기 전 메시지 정보가 비어있다면,
	// 아래 메시지를 보내고 함수를 종료합니다.
	if deletedMessage == nil {
		session.ChannelMessageSend(ChannelID, "메시지가 삭제되었습니다.")
		return
	}

	// 그러나 삭제되기 전 메시지 정보가 비어있지 않다면,
	// 어떤 내용의 메시지가 삭제되었는지 알려주는 메시지를 보냅니다.
	session.ChannelMessageSend(ChannelID, "메시지가 삭제되었습니다.\n삭제된 메시지: "+deletedMessage.Content)
}
