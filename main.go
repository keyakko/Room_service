package main

// base libraries
import (
  "log"
  "os"
  "fmt"
  "strings"
  "time"
)

// slackAPI lib
import (
  "github.com/nlopes/slack"
)

// gobot libraries
import (
  //"gobot.io/x/gobot"
  "gobot.io/x/gobot/drivers/gpio"
  "gobot.io/x/gobot/platforms/raspi"
)


func run (api *slack.Client) (int) {

  // survo settings
  adaptor := raspi.NewAdaptor()
  servo := gpio.NewServoDriver(adaptor, "12")
/*
  work := func () {
    servo.Move(uint8(27))
    servo.Move(uint8(32))
    servo.Move(uint8(27))
  }

  robot := gobot.NewRobot(
    "servo",
    []gobot.Connection{adaptor},
    []gobot.Device{servo},
    work,
  )
*/

  rtm := api.NewRTM()
  go rtm.ManageConnection()

  for {
    select {
    case msg := <-rtm.IncomingEvents:
        switch ev := msg.Data.(type) {
        case *slack.HelloEvent:
          log.Print("Hello Event")

        case *slack.MessageEvent:
          log.Printf("text: %+v\n", ev.Msg.Text)
          if strings.Index(ev.Msg.Text, "開けて") != -1 {
            rtm.SendMessage(rtm.NewOutgoingMessage("まかせて！", ev.Channel))
            //robot.Start()
            servo.Move(uint8(20))
            servo.Move(uint8(22))
	    time.Sleep(1 * time.Second)
            servo.Move(uint8(20))
          } else {
            //rtm.SendMessage(rtm.NewOutgoingMessage("理解できないメッセージだから無視するよ。", ev.Channel))
          }


        case *slack.InvalidAuthEvent:
          log.Print("Invalid credentials")
          return 1

        }
    }
  }
}

func main () {

  token := os.Getenv("SLACKBOT_TOKEN")

  if token == "" {
    fmt.Println("error: token is empty")
    os.Exit(1)
  }

  api := slack.New(token)
  os.Exit(run(api))

}
