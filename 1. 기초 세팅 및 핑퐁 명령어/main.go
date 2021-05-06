package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

var (
	Token   = "Nzc1MDA1OTYxMTYwMjI4ODg0.X6gCjA.gQGSiGOV7Yfntf1GLtwkZs-kq44"
	Session *discordgo.Session
)

func init() {
	var err error
	// 봇 클라이언트를 생성합니다.
	Session, err = discordgo.New("Bot " + Token)
	// 에러가 있는 경우 프로그램을 종료합니다.
	if err != nil {
		log.Fatalf("클라이언트 생성 오류: %v", err)
	}

	// 클라이언트를 디스코드와 연결합니다.
	err = Session.Open()
	// 에러가 있는 경우 프로그램을 종료합니다.
	if err != nil {
		log.Fatalf("세션 오픈 오류: %v", err)
	}

	// 위 모든 절차가 완료된 경우 로그를 찍습니다.
	log.Printf("%s (%s)에 로그인 됨", Session.State.User.String(), Session.State.User.ID)
}

func main() {
	// 세션의 이벤트가 일어났을 경우 AddHandler의 괄호 안에 있는 함수를 호출합니다.
	// (위 Ready 이벤트처럼 익명 함수를 사용할 수도 있습니다.)
	// 디스코드 봇의 대표적인 이벤트로는
	// "메시지 수신(Message Create)", "새로운 채널 생성(Channel Create)" 등이 있습니다.
	// 이번에는 핑퐁 명령어 제작을 위하여 "메시지 수신" 이벤트를 활용해 보도록 하겠습니다.
	Session.AddHandler(messageCreate)

	// ^C 등으로 프로그램을 종료하였을 때,
	// 프로그램이 꺼지기 전 세션을 닫습니다.
	defer Session.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("봇 종료됨")
}

// 이벤트 수신을 위하여 첫 매개변수에는 *discordgo.Session 형식을 받고,
// 두번째 매개변수에는 *discordgo.이벤트 형식을 받아줍니다.
// 위에서 말한대로, 메시지 수신을 위하여 MessageCreate 이벤트를 사용하도록 하겠습니다.
func messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	// 여기서는 메시지 수신 시 어느 대상에게 반응할 지를 거르는 필터 역할을 한다고 생각하면 편합니다.
	// 매개 변수로 받아온 session에는 우리가 방금 실행한 디스코드 봇의 정보가 담겨져 있으며,
	// session.State.User.ID에는 봇의 아이디가 들어가 있습니다.
	// 그리고 또 매개 변수로 받아온 message에는 수신된 메시지의 내용, 작성자, 메시지가 온 채널 등의 정보가 담겨 있습니다.
	// 그 중 message.Author.ID는 메시지를 작성한 사람의 아이디가 들어가 있습니다.
	// 즉, 아래 if문은 봇의 아이디와 메시지를 작성한 사람의 아이디가 동일한 경우
	// 해당 메시지에 반응하지 않도록 하는 필터 역할을 하고 있습니다.
	if session.State.User.ID == message.Author.ID {
		return
	}

	// message에는 수신된 메시지의 내용, 작성자, 메시지가 온 채널 등의 정보가 담겨 있다고 위에 언급했었습니다.
	// message.Content는 수신된 메시지의 내용이 담겨 있습니다.
	// 즉, 아래 if문은 수신된 메시지의 내용이 "!ping"일 경우
	// if문 내부 블록에 있는 코드를 실행한다고 볼 수 있습니다.
	if message.Content == "!ping" {
		// ChannelMessageSend는 함수 이름 그대로 해당 채널에 봇이 메시지를 보내는 함수입니다.
		// ChannelMessageSend에는 첫번째 매개 변수로 메시지를 보낼 채널의 ID,
		// 그리고 두번째 매개 변수로 보낼 메시지를 입력해 주면 됩니다.
		// 아까 message에 메시지가 수신된 채널의 정보도 담겨 있다고 말씀 드렸는데, message.ChannelID에 그 정보가 담겨있습니다.
		// 즉, 아래 코드는 메시지가 수신된 채널에 "pong!"이라는 메시지를 보내는 코드입니다.
		session.ChannelMessageSend(message.ChannelID, "pong!")
	}

	if message.Content == "!pong" {
		session.ChannelMessageSend(message.ChannelID, "ping!")
	}

	// 그렇다면 아래 코드를 수정하여
	// 메시지를 보낸 사람이 "자기 자신"인 경우,
	// 봇이 "hello, (유저 닉네임#태그)"를 보내도록 해보세요.

	// 메시지를 보낸 사람을 걸러내는 방법은 이미 67줄의 코드에서 설명하였습니다.
	if message.Content == "" {
		// 메시지를 보낸 사람의 '유저 닉네임#태그'를 불러오려는 경우
		// message.Author.String()을 활용할 수 있습니다,
		session.ChannelMessageSend(message.ChannelID, "")
	}
}
