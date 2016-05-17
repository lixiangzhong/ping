package ping

import (
	"errors"
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	ProtocolICMP = 1
)

var (
	ErrNotReply = errors.New("not echo reply")
)

func doPing(host string) (err error) {
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Fprintln(os.Stderr, "No permission to listen ICMP")
		return
	}
	defer c.Close()
	dst, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff, //why &0xffff ??
			Seq:  r.Int(),
			Data: []byte("R-U-OK"),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		return
	}
	var n int
	if n, err = c.WriteTo(wb, dst); err != nil {
		return
	} else if n != len(wb) {
		return
	}
	rb := make([]byte, 1500)
	if err = c.SetReadDeadline(time.Now().Add(time.Second * 3)); err != nil {
		return
	}
	n, _, err = c.ReadFrom(rb)
	if err != nil {
		return
	}
	rm, err := icmp.ParseMessage(ProtocolICMP, rb[:n])
	if err != nil {
		return
	}
	if rm.Type != ipv4.ICMPTypeEchoReply {
		return ErrNotReply
	}
	return nil
}
func DoPing(host string) error {
	return doPing(host)
}
func Ping(host string) bool {
	return doPing(host) == nil
}
