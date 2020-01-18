package utils

import (
	"fmt"
	"testing"
	"time"
)

func TestPrint(t *testing.T)  {
	fmt.Println("Test Print===")
}

func TestAll(t *testing.T){
	fmt.Println("======start======")
	go func() {
		for{
			fmt.Println(BodyList.BodyList)
			if BodyList.Check("123"){
				body := BodyList.Remove("123")
				fmt.Println(fmt.Sprintf("body :%s",body.Body))
				return
			}
		}
	}()


	go func() {
		q:= AsyncBody{
			Uid:"123",
			Body:[]byte("123"),
		}
		Queue.Add(q)

		<-time.After(1*time.Second)
		q = AsyncBody{
			Uid:"234",
			Body:[]byte("234"),
		}
		Queue.Add(q)

		<-time.After(1*time.Second)
		q = AsyncBody{
			Uid:"345",
			Body:[]byte("345"),
		}
		Queue.Add(q)
	}()

}