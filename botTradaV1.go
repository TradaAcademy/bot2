package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Step int

const (
	stepToChooseCourse Step = iota
	stepToEnterDifferentCousre
	//stepToAskName
	//stepToAskCourse
	// stepToAskPhoneNum
	// stepToConfirmDisplayYear
	// stepToAskYearOfBirth

	// stepToDisplay
	stepDone
)

//"721136312"

type Display struct {
	name     string
	course   string
	external string
	yearob   string
	dif      int
}

type User struct {
	name              string
	course            string
	phonenum          string
	yearob            string
	lastSigninRequest time.Time
	registrationStep  Step
}

// func getFullName(m *tb.Message) string {
// 	if user, ok := userMap[m.Sender.ID]; ok {
// 		if user.name != "" {
// 			return user.name
// 		}
// 	}
// 	return fmt.Sprintf("%s %s", m.Sender.FirstName, m.Sender.LastName)
// }
var userMap = make(map[int]*User)
var u User
var ud Display

// var courseMap = map[int]string{
// 	111: "Dapp",
// 	112: "Machine Learning",
// 	113: "Web Development",
// }

func main() {
	b, err := tb.NewBot(tb.Settings{

		Token: os.Getenv("BOT_TOKEN"),
		//Token:  "702307311:AAGVv4Xbp0pitk2IstBh9jxE5WYqmpXHbt8",
		Poller: &tb.LongPoller{Timeout: 2 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle(tb.OnText, func(m *tb.Message) {

		checkStep(b, m)
	})

	// b.Handle("/me", func(m *tb.Message) {
	// 	b.Send(m.Sender, "Tính năng đang được phát triển, vui lòng chọn tính năng khác trong mục /help")
	// })

	b.Handle("/start", func(m *tb.Message) {

		b.Send(m.Sender, "Cảm ơn bạn đã liên lạc với Trada Tech. Đây là bot trả lời tự động hỗ trợ các lệnh sau:\n/info - thông tin về dịch vụ và khoá học\n/register - đăng kí khoá học, dịch vụ, hoặc đặt câu hỏi cho Trada Tech.\n/cancel - thông báo muốn huỷ đăng kí.\n/help - hiện hướng dẫn.")

	})
	b.Handle("/help", func(m *tb.Message) {

		b.Send(m.Sender, "/info - thông tin về dịch vụ và khoá học\n/register - đăng kí khoá học, dịch vụ, hoặc đặt câu hỏi cho Trada Tech.\n/cancel - thông báo muốn huỷ đăng kí.\n/help - hiện hướng dẫn.")

	})
	b.Handle("/info", func(m *tb.Message) {

		b.Send(m.Sender, "Thông tin về dịch vụ và khoá học:\n1. Đào tạo Ethereum DApp Developer\n2. Đào tạo theo yêu cầu riêng của công ty bạn\n3. Tư vấn về phát triển và ứng dụng blockchain")

	})

	b.Handle("/register", func(m *tb.Message) {

		next(b, m)

	})

	b.Handle("/cancel", func(m *tb.Message) {

		if ud.course == "" {
			sendAndHideKeyboard(b, m, "Bạn vẫn chưa đăng kí khoá học nào, không thể huỷ khoá học!")
		} else {
			sendAndHideKeyboard(b, m, "Yêu cầu huỷ của bạn đã được thông báo cho nhân viên Trada Tech. Chúng tôi sẽ liên lạc lại với bạn để xác nhận.")
			sendCancelRequest(b, m)
		}

	})

	b.Start()
}

func containAny(array []string, item string) bool {
	for _, element := range array {
		if strings.EqualFold(element, item) {
			return true
		}
	}

	return false
}

func isCourse1(text string) bool {
	values := []string{"DApp"}
	//, "course 111", "111", "course 1", "first course"
	return containAny(values, strings.TrimSpace(text))
}

func isCourse2(text string) bool {
	values := []string{"Khác", "Khac"}
	//, "course 112", "112", "course 2", "second course"
	return containAny(values, strings.TrimSpace(text))
}

// func isCourse3(text string) bool {
// 	values := []string{"Web Development", "course 113", "113", "course 3", "third course"}
// 	return containAny(values, strings.TrimSpace(text))
// }

// func isYes(text string) bool {
// 	values := []string{"có", "yes", "sure", "certainly", "ok", "okay", "fine", "indeed",
// 		"definitely", "of course", "affirmative", "obviously", "absolutely",
// 		"indubitably", "undoubtedly"}
// 	return containAny(values, strings.TrimSpace(text))
// }

// func isNo(text string) bool {
// 	values := []string{"no", "không", "never", "by no means", "no way"}
// 	return containAny(values, strings.TrimSpace(text))
// }

func listCourse(b *tb.Bot, m *tb.Message) {
	//sendfAndHideKeyboard(b, m, "Các khoá học hiện có:\nDApp - Khoá học Ethereum DApp Development\nKhác - Các yêu cầu khác")

	sendCourseChoices(b, m, "Các khoá học hiện có:\nDApp - Khoá học Ethereum DApp Development\nKhác - Các yêu cầu khác\nVui lòng chọn khoá học:")
}

func confirmDisplayYear(b *tb.Bot, m *tb.Message) {
	sendYesNo(b, m, "Bạn có muốn nhập và hiển thị công khai năm sinh không?")
}

func sendCourseChoices(b *tb.Bot, m *tb.Message, text string) (*tb.Message, error) {

	c1Btn := tb.ReplyButton{Text: "DApp"}
	c2Btn := tb.ReplyButton{Text: "Khác"}
	//c3Btn := tb.ReplyButton{Text: "Web Development"}
	replyChoice := [][]tb.ReplyButton{
		[]tb.ReplyButton{c1Btn, c2Btn},
		// []tb.ReplyButton{c2Btn},
		//[]tb.ReplyButton{c3Btn},
	}

	return b.Reply(m,
		text,
		&tb.ReplyMarkup{
			ReplyKeyboard:       replyChoice,
			ResizeReplyKeyboard: true,
			OneTimeKeyboard:     true,
			ReplyKeyboardRemove: true,
		})
}

func sendYesNo(b *tb.Bot, m *tb.Message, text string) (*tb.Message, error) {
	yesBtn := tb.ReplyButton{Text: "Yes"}
	noBtn := tb.ReplyButton{Text: "No"}
	replyYesNo := [][]tb.ReplyButton{
		[]tb.ReplyButton{yesBtn, noBtn},
	}
	return b.Send(m.Sender,
		text,
		&tb.ReplyMarkup{
			ReplyKeyboard:       replyYesNo,
			ResizeReplyKeyboard: true,
			OneTimeKeyboard:     true,
		})
}

func sendMessageToGroup(b *tb.Bot, m *tb.Message) {
	var text string
	username := m.Sender.Username
	//sendfAndHideKeyboard(b, m, "Hello @%s!", m.Sender.Username)

	if username == "" {
		if ud.dif == 0 {
			text = "Có đăng kí mới!%0ATừ một người không public username%0AKhoá học: " + ud.course + "%0AThông tin thêm: "
		} else {
			text = "Có đăng kí mới!%0ATừ một người không public username%0AKhoá học: " + ud.course + "%0AThông tin thêm: " + ud.external
		}
		//text = "Hey boss, someone without an username has just registered!%0AHere's the info:%0ATelegram ID: " + strconv.Itoa(m.Sender.ID) + "%0AName: " + ud.name + "%0ACourse: " + ud.course + "%0APhone number: " + ud.phonenum + "%0AYear of birth: " + ud.yearob

	} else {
		//text = "Hey boss, @" + username + " has just registered!%0AHere's the info:%0ATelegram ID: " + strconv.Itoa(m.Sender.ID) + "%0ADisplay name: " + ud.name + "%0ACourse: " + ud.course + "%0APhone number: " + ud.phonenum + "%0AYear of birth: " + ud.yearob

		if ud.dif == 0 {
			text = "Có đăng kí mới!%0ATừ: @" + username + "%0AKhoá học: " + ud.course + "%0AThông tin thêm:"
		} else {
			text = "Có đăng kí mới!%0ATừ: @" + username + "%0AKhoá học: " + ud.course + "%0AThông tin thêm: " + ud.external
		}
	}

	_, err := http.Get("https://api.telegram.org/bot" + os.Getenv("BOT_TOKEN") + "/sendMessage?chat_id=" + os.Getenv("_ID") + "&text=" + text)
	if err != nil {
		fmt.Print("error: %s", err)
	}

}

func sendCancelRequest(b *tb.Bot, m *tb.Message) {
	var text string
	username := m.Sender.Username
	//sendfAndHideKeyboard(b, m, "Hello @%s!", m.Sender.Username)

	if username == "" {
		//text = "Hey boss, someone without an username has just registered!%0AHere's the info:%0ATelegram ID: " + strconv.Itoa(m.Sender.ID) + "%0AName: " + ud.name + "%0ACourse: " + ud.course + "%0APhone number: " + ud.phonenum + "%0AYear of birth: " + ud.yearob
		text = "Người dùng private vừa huỷ khoá học mới đăng kí."
	} else {
		//text = "Hey boss, @" + username + " has just registered!%0AHere's the info:%0ATelegram ID: " + strconv.Itoa(m.Sender.ID) + "%0ADisplay name: " + ud.name + "%0ACourse: " + ud.course + "%0APhone number: " + ud.phonenum + "%0AYear of birth: " + ud.yearob
		text = "@" + username + " vừa huỷ khoá học mới đăng kí."

	}

	_, err := http.Get("https://api.telegram.org/bot" + os.Getenv("BOT_TOKEN") + "/sendMessage?chat_id=" + os.Getenv("_ID") + "&text=" + text)
	if err != nil {
		fmt.Print("error: %s", err)
	}
}

func askDisplayName(b *tb.Bot, m *tb.Message) {

	sendAndHideKeyboard(b, m, "Tên của bạn là: ")

}

func askPhoneNumber(b *tb.Bot, m *tb.Message) {
	sendAndHideKeyboard(b, m, "Vui lòng nhập số điện thoại của bạn: ")
}

func reEnterPhonenum(b *tb.Bot, m *tb.Message) {
	sendAndHideKeyboard(b, m, "Bạn vừa nhập KHÔNG ĐÚNG định dạng của số điện thoại, xin vui lòng nhập lại theo định dạng ĐÚNG!")
}

func reEnterBirthYear(b *tb.Bot, m *tb.Message) {
	sendAndHideKeyboard(b, m, "Bạn vừa nhập KHÔNG ĐÚNG năm, xin vui lòng nhập lại!")
}

func isValidPhoneNum(text string) bool {
	if len(text) != 10 {
		return false
	}
	for _, e := range text {
		if e < 48 || e > 57 {
			return false
		}
	}
	return true
}

func askYearOfBirth(b *tb.Bot, m *tb.Message) {
	sendAndHideKeyboard(b, m, "Vui lòng nhập năm sinh của bạn: ")
}

func isValidYear(text string) bool {
	for _, e := range text {
		if e < 48 || e > 57 {
			return false
		}
	}
	y, _ := strconv.Atoi(text)
	if y > time.Now().Year()-10 || y < time.Now().Year()-100 {
		return false
	}
	return true

}

func startRegistration(b *tb.Bot, m *tb.Message) {
	newUser := User{registrationStep: stepToChooseCourse}
	userMap[m.Sender.ID] = &newUser
	listCourse(b, m)
	//askDisplayName(b, m)

}

func sendAndHideKeyboard(b *tb.Bot, m *tb.Message, text string) (*tb.Message, error) {
	return b.Send(m.Sender, text, &tb.ReplyMarkup{ReplyKeyboardRemove: true})
}
func sendfAndHideKeyboard(b *tb.Bot, m *tb.Message, text string, a ...interface{}) (*tb.Message, error) {
	return sendAndHideKeyboard(b, m, fmt.Sprintf(text, a...))
}

// func disPlayInformation(b *tb.Bot, m *tb.Message) {
// 	u := userMap[m.Sender.ID]
// 	if ud.yearob != "" {

// 		_, err := sendfAndHideKeyboard(b, m, "So, here is your information:\nName: %s\nCourse: %s\nPhone number: %s\nYear of birth: %s\nThank you for signing in our course!",
// 			u.name,
// 			u.course,
// 			u.phonenum,
// 			u.yearob,
// 		)
// 		if err == nil {
// 			u.lastSigninRequest = time.Now()
// 		}
// 	} else {
// 		_, err := sendfAndHideKeyboard(b, m, "So, here is your information:\nName: %s\nCourse: %s\nPhone number: %s\nYear of birth: N/A\nThank you for signing in our course!",
// 			u.name,
// 			u.course,
// 			u.phonenum,
// 		)
// 		if err == nil {
// 			u.lastSigninRequest = time.Now()
// 		}
// 	}

// 	sendAndHideKeyboard(b, m, "\nType anything to finish registration ...")

// }

func sayGoodBye(b *tb.Bot, m *tb.Message) {

	sendAndHideKeyboard(b, m, "\nCảm ơn bạn đã đăng kí, thông tin của bạn đã được gửi cho nhân viên Trada Tech xử lý, chúng tôi sẽ liên lạc lại để xác nhận thông tin.")
}

func differentCourse(b *tb.Bot, m *tb.Message) {
	sendAndHideKeyboard(b, m, "Vui lòng cung cấp thêm chi tiết về khoá học bạn muốn đăng kí")
}

func awaitCommand(b *tb.Bot, m *tb.Message) {
	b.Send(m.Sender, "\nPhiên đăng kí đã kết thúc, gõ /help để hiện menu trợ giúp.")
}

func checkStep(b *tb.Bot, m *tb.Message) {

	if u, ok := userMap[m.Sender.ID]; ok {
		switch u.registrationStep {
		case stepToChooseCourse:
			if isCourse1(m.Text) {
				u.registrationStep = stepDone
				ud.dif = 0
				u.course = strings.Title(strings.TrimSpace(m.Text))
				ud.course = u.course

				sayGoodBye(b, m)
				sendMessageToGroup(b, m)
				next(b, m)
			}
			if isCourse2(m.Text) {
				u.registrationStep = stepToEnterDifferentCousre
				u.course = strings.Title(strings.TrimSpace(m.Text))
				ud.course = u.course
				ud.dif = 1
				//sendAndHideKeyboard(b,m,"Cảm ơn bạn. Nhân viên của Trada Tech sẽ liên lạc lại với bạn để hỏi thêm chi tiết.")
				next(b, m)
			} else if !isCourse1(m.Text) && !isCourse2(m.Text) {
				sendAndHideKeyboard(b, m, "Vui lòng dùng 2 nút có sẵn để trả lời.")

				next(b, m)
			}
			/*if isCourse3(m.Text) {
				u.registrationStep = stepToAskName
				//u.course = strings.Title(strings.TrimSpace(m.Text))
				u.course = courseMap[113]
				ud.course = u.course
				next(b, m)
			}*/
		case stepToEnterDifferentCousre:
			u.registrationStep = stepDone
			//u.course = strings.Title(strings.TrimSpace(m.Text))
			ud.external = strings.Title(strings.TrimSpace(m.Text))
			sendAndHideKeyboard(b, m, "Nhân viên của Trada Tech sẽ liên lạc lại với bạn để hỏi thêm chi tiết, cảm ơn bạn đã đăng kí.")
			sendMessageToGroup(b, m)
			//	sayGoodBye(b, m)
			next(b, m)
			//removeRegisteredUser(m)

		//case stepDone:
		//u.registrationStep = stepToChooseCourse
		/*case stepToAskName:
			u.name = strings.Title(strings.TrimSpace(m.Text))
			ud.name = u.name
			u.registrationStep = stepToAskPhoneNum
			next(b, m)

		// case stepToAskCourse:

		// 	u.course = strings.Title(strings.TrimSpace(m.Text))
		// 	ud.course = u.course
		// 	u.registrationStep = stepToDisplay
		// 	next(b, m)
		case stepToAskPhoneNum:

			if isValidPhoneNum(m.Text) {
				u.phonenum = strings.Title(strings.TrimSpace(m.Text))
				ud.phonenum = u.phonenum
				u.registrationStep = stepToConfirmDisplayYear
				next(b, m)

			} else {
				//u.registrationStep = stepToAskPhoneNum
				reEnterPhonenum(b, m)
				next(b, m)
			}

		case stepToConfirmDisplayYear:
			if isYes(m.Text) {
				u.registrationStep = stepToAskYearOfBirth
				next(b, m)
			}
			if isNo(m.Text) {
				u.registrationStep = stepToDisplay
				next(b, m)
			}
		case stepToAskYearOfBirth:

			if isValidYear(m.Text) {
				u.yearob = strings.Title(strings.TrimSpace(m.Text))
				ud.yearob = u.yearob
				u.registrationStep = stepDone
				next(b, m)
			} else {
				//u.registrationStep = stepToAskYearOfBirth
				reEnterBirthYear(b, m)
				next(b, m)
			}

		/*case stepToDisplay:
		u.registrationStep = stepDone
		sendMessageToGroup(b, m)
		sayGoodBye(b, m)
		removeRegisteredUser(m)
		*/

		default:
			//u.registrationStep = stepToChooseCourse
			awaitCommand(b, m)

		}

	} else {

		awaitCommand(b, m)
		//sayGoodBye(b, m)
		//sendMessageToGroup()
	}
}

func next(b *tb.Bot, m *tb.Message) {

	if user, ok := userMap[m.Sender.ID]; ok {
		funcArray := []func(*tb.Bot, *tb.Message){
			listCourse,
			differentCourse,
			removeRegisteredUser,
			//askDisplayName,
			//askPhoneNumber,
			//confirmDisplayYear,
			//askYearOfBirth,
			//disPlayInformation,
		}
		funcArray[user.registrationStep](b, m)
	} else {

		startRegistration(b, m)
	}
}

func removeRegisteredUser(b *tb.Bot, m *tb.Message) {
	b.Send(m.Sender, "")
	delete(userMap, m.Sender.ID)
}
