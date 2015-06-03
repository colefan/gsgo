package netqos

import (
	"time"
	"testing"
)
func TestQos(t *testing.T){
	q := NewServerQos()
	q.Stat()
	
	for i:=0;i<100;i++{
		q.StatAccpetConns()
		go func(){
			t2 := time.NewTicker(2*time.Second)
			for{
					select {
					case <-t2.C:
					
					q.StatReadMsgs()
					
					}
				}	
		}()
	}
	
	time.Sleep(time.Minute*1)
	
	
	
}