package main

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"io"
	"log"
	"io/ioutil"
)

type Literal interface {
	io.Reader
	Len() int
}

type ResponseMessage struct {
	From    string
	Subject string
	Content string
}

func AcceptAllMail(addr, user, pass string) ([]Literal, error) {
	client, err := client.DialTLS(addr, nil)
	if err != nil {
		return nil, err
	}
	defer client.Logout()
	if err := client.Login(user, pass); err != nil {
		log.Fatal(err)
	}
	// 收件箱
	mbox, err := client.Select("INBOX", true)
	if err != nil {
		return nil, err
	}

	if mbox.Messages == 0 {
		return nil, nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(1), mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- client.Fetch(seqset, []string{"BODY[]"}, messages)
	}()
	// 返回体
	response := []Literal{}
	// 收件箱的所有邮件
	for msg := range messages {
		r := msg.GetBody("BODY[]")
		if r == nil {
			return nil, fmt.Errorf("没有邮件内容")
		}
		response = append(response, r)
	}
	return response, nil
}

func main() {
	//request, err := AcceptAllMail("imap.qq.com:993", "1436863821@qq.com", "mrosvtjgojhdgddj")
	if err != nil {
		log.Fatal(err)
	}

	// 返回体
	response := []*ResponseMessage{}
	for _, r := range request {
		// Create a new mail reader
		mr, err := mail.CreateReader(r)
		if err != nil {
			log.Fatal(err)
		}

		var existEntity = new(ResponseMessage)
		header := mr.Header
		if from, err := header.AddressList("From"); err == nil {
			for _, value := range from {
				existEntity.From = value.Address
			}
		}

		if subject, err := header.Subject(); err == nil {
			existEntity.Subject = subject
		}

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatal(err)
			}
			switch  p.Header.(type) {
			case mail.TextHeader:
				// This is the message's text (can be plain-text or HTML)
				b, err := ioutil.ReadAll(p.Body)
				if err != nil{
					fmt.Println(err)
				}
				existEntity.Content = string(b)
			}
		}
		response = append(response, existEntity)
	}

	// 测试输出
	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>><<><><><><><><><><><>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>><<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	for _,value := range response{
		fmt.Println("发件人：",value.From)
		fmt.Println("主题：",value.Subject)
		fmt.Println("发件内容：",value.Content)
		fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>><>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>><>>>>>>>>>>>>>>><<<<<<<<<<<>>>>>>>>><<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	}
}


