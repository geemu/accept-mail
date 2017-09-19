package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
)


func AcceptAllMail(addr,user,pass string) {
	client, err := client.DialTLS(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Logout()
	if err := client.Login(user, pass); err != nil {
		log.Fatal(err)
	}
	// 收件箱
	mbox, err := client.Select("INBOX", true)
	if err != nil {
		log.Fatal(err)
	}

	if mbox.Messages == 0 {
		log.Fatal("信箱中没有信息")
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, []string{"BODY[]"}, messages)
	}()

	// 收件箱的所有邮件
	for msg := range messages {
		fmt.Println(msg.Body)
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println()
		fmt.Println()
	}
}

func main()  {
	//	Addr = "imap.qq.com:993"
	//	User = "1436863821@qq.com"
	//	Pass = "mrosvtjgojhdgddj"
	AcceptAllMail("imap.qq.com:993","1436863821@qq.com","mrosvtjgojhdgddj")
}
