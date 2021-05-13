package rtsp

import (
	"fmt"
	"net"
	"strings"

	"github.com/childback/EasyDarwin/EasyGoLib/utils"
)

type UDPClient struct {
	*Session

	APort        int
	AConn        *net.UDPConn
	AControlPort int
	AControlConn *net.UDPConn
	VPort        int
	VConn        *net.UDPConn
	VControlPort int
	VControlConn *net.UDPConn

	Stoped bool
}

func (s *UDPClient) Stop() {
	if s.Stoped {
		return
	}
	s.Stoped = true
	if s.AConn != nil {
		s.AConn.Close()
		s.AConn = nil
	}
	if s.AControlConn != nil {
		s.AControlConn.Close()
		s.AControlConn = nil
	}
	if s.VConn != nil {
		s.VConn.Close()
		s.VConn = nil
	}
	if s.VControlConn != nil {
		s.VControlConn.Close()
		s.VControlConn = nil
	}
}

func (c *UDPClient) SetupAudio() (err error) {
	var (
		logger = c.logger
	)
	defer func() {
		if err != nil {
			logger.Println(err)
			c.Stop()
		}
	}()
	if err = c.SetupRtpAudio(); err != nil {
		return
	}
	//if err = c.SetupRtcpCAudio(); err != nil {
	//	return
	//}
	if err = c.SetupRtcpSAudio(); err != nil {
		return
	}
	return
}

func (c *UDPClient) SetupRtpAudio() (err error) {
	var (
		logger = c.logger
		addr   *net.UDPAddr
	)
	host := c.Conn.RemoteAddr().String()
	host = host[:strings.LastIndex(host, ":")]
	addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, c.APort))
	if err != nil {
		return
	}
	c.AConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	networkBuffer := utils.Conf().Section("rtsp").Key("network_buffer").MustInt(1048576)
	if err = c.AConn.SetReadBuffer(networkBuffer); err != nil {
		logger.Printf("udp client audio conn set read buffer error, %v", err)
	}
	if err = c.AConn.SetWriteBuffer(networkBuffer); err != nil {
		logger.Printf("udp client audio conn set write buffer error, %v", err)
	}
	return
}

func (c *UDPClient) SetupRtcpCAudio() (err error) {
	var (
		logger = c.logger
		addr   *net.UDPAddr
	)
	host := c.Conn.RemoteAddr().String()
	host = host[:strings.LastIndex(host, ":")]
	if addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, c.AControlPort)); err != nil {
		return
	}
	c.AControlConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	networkBuffer := utils.Conf().Section("rtsp").Key("network_buffer").MustInt(1048576)
	if err = c.AControlConn.SetReadBuffer(networkBuffer); err != nil {
		logger.Printf("udp client audio contrl conn set read buffer error, %v", err)
	}
	if err = c.AControlConn.SetWriteBuffer(networkBuffer); err != nil {
		logger.Printf("udp client audio contrl conn set write buffer error, %v", err)
	}
	return
}

func (c *UDPClient) SetupRtcpSAudio() (err error) {
	var (
		logger = c.logger
		addr   *net.UDPAddr
	)
	if addr, err = net.ResolveUDPAddr("udp", ":0"); err != nil {
		return
	}
	host := addr.String()
	host = host[:strings.LastIndex(host, ":")]
	if addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, c.AControlPort)); err != nil {
		return
	}
	if c.AControlConn, err = net.ListenUDP("udp", addr); err != nil {
		return
	}
	networkBuffer := utils.Conf().Section("rtsp").Key("network_buffer").MustInt(1048576)
	if err = c.AControlConn.SetReadBuffer(networkBuffer); err != nil {
		logger.Printf("udp client audio control conn set read buffer error, %v", err)
	}
	if err = c.AControlConn.SetWriteBuffer(networkBuffer); err != nil {
		logger.Printf("udp client audio control conn set write buffer error, %v", err)
	}
	go func() {
		bufUDP := make([]byte, UDP_BUF_SIZE)
		logger.Printf("udp client start listen audio control port[%d]", c.AControlPort)
		defer logger.Printf("udp client stop listen audio control port[%d]", c.AControlPort)
		for !c.Stoped {
			if n, ra, err := c.AControlConn.ReadFromUDP(bufUDP); err == nil {
				logger.Printf("Package recv from AControlConn.len: %d, %s\n", n, ra.String())
			} else {
				logger.Println("udp server read audio control pack error", err)
				continue
			}
		}
	}()
	return
}

func (c *UDPClient) SetupVideo() (err error) {
	var (
		logger = c.logger
	)
	defer func() {
		if err != nil {
			logger.Println(err)
			c.Stop()
		}
	}()
	if err = c.SetupRtpVideo(); err != nil {
		return
	}
	//if err = c.SetupRtcpCVideo(); err != nil {
	//	return
	//}
	if err = c.SetupRtcpSVideo(); err != nil {
		return
	}
	return
}

func (c *UDPClient) SetupRtpVideo() (err error) {
	var (
		logger = c.logger
		addr   *net.UDPAddr
	)
	host := c.Conn.RemoteAddr().String()
	host = host[:strings.LastIndex(host, ":")]
	addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, c.VPort))
	if err != nil {
		return
	}
	c.VConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	networkBuffer := utils.Conf().Section("rtsp").Key("network_buffer").MustInt(1048576)
	if err = c.VConn.SetReadBuffer(networkBuffer); err != nil {
		logger.Printf("udp client video conn set read buffer error, %v", err)
	}
	if err = c.VConn.SetWriteBuffer(networkBuffer); err != nil {
		logger.Printf("udp client video conn set write buffer error, %v", err)
	}
	return
}

func (c *UDPClient) SetupRtcpCVideo() (err error) {
	var (
		logger = c.logger
		addr   *net.UDPAddr
	)
	host := c.Conn.RemoteAddr().String()
	host = host[:strings.LastIndex(host, ":")]
	if addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, c.VControlPort)); err != nil {
		return
	}
	c.VControlConn, err = net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	networkBuffer := utils.Conf().Section("rtsp").Key("network_buffer").MustInt(1048576)
	if err = c.VControlConn.SetReadBuffer(networkBuffer); err != nil {
		logger.Printf("udp client video contrl conn set read buffer error, %v", err)
	}
	if err = c.VControlConn.SetWriteBuffer(networkBuffer); err != nil {
		logger.Printf("udp client video contrl conn set write buffer error, %v", err)
	}
	return
}

func (c *UDPClient) SetupRtcpSVideo() (err error) {
	var (
		logger = c.logger
		addr   *net.UDPAddr
	)
	if addr, err = net.ResolveUDPAddr("udp", ":0"); err != nil {
		return
	}
	host := addr.String()
	host = host[:strings.LastIndex(host, ":")]
	if addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, c.VControlPort)); err != nil {
		return
	}
	if c.VControlConn, err = net.ListenUDP("udp", addr); err != nil {
		return
	}
	networkBuffer := utils.Conf().Section("rtsp").Key("network_buffer").MustInt(1048576)
	if err = c.VControlConn.SetReadBuffer(networkBuffer); err != nil {
		logger.Printf("udp client video control conn set read buffer error, %v", err)
	}
	if err = c.VControlConn.SetWriteBuffer(networkBuffer); err != nil {
		logger.Printf("udp client video control conn set write buffer error, %v", err)
	}
	go func() {
		bufUDP := make([]byte, UDP_BUF_SIZE)
		logger.Printf("udp client start listen video control port[%d]", c.VControlPort)
		defer logger.Printf("udp client stop listen video control port[%d]", c.VControlPort)
		for !c.Stoped {
			if n, ra, err := c.VControlConn.ReadFromUDP(bufUDP); err == nil {
				logger.Printf("Package recv from VControlConn.len: %d, %s\n", n, ra.String())
			} else {
				logger.Println("udp server read video control pack error", err)
				continue
			}
		}
	}()
	return
}

func (c *UDPClient) SendRTP(pack *RTPPack) (err error) {
	if pack == nil {
		err = fmt.Errorf("udp client send rtp got nil pack")
		return
	}
	var conn *net.UDPConn
	switch pack.Type {
	case RTP_TYPE_AUDIO:
		conn = c.AConn
	case RTP_TYPE_AUDIOCONTROL:
		conn = c.AControlConn
	case RTP_TYPE_VIDEO:
		conn = c.VConn
	case RTP_TYPE_VIDEOCONTROL:
		conn = c.VControlConn
	default:
		err = fmt.Errorf("udp client send rtp got unkown pack type[%v]", pack.Type)
		return
	}
	if conn == nil {
		err = fmt.Errorf("udp client send rtp pack type[%v] failed, conn not found", pack.Type)
		return
	}
	var n int
	if n, err = conn.Write(pack.Buffer.Bytes()); err != nil {
		err = fmt.Errorf("udp client write bytes error, %v", err)
		return
	}
	// logger.Printf("udp client write [%d/%d]", n, pack.Buffer.Len())
	c.Session.OutBytes += n
	return
}
