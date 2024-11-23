package sip

import (
	"fmt"
	"sip-parser/pkg/siprocket"
	"sip-parser/pkg/utils"
	"strings"
	"time"
)

type Message struct {
	Timestamp time.Time
	//pct       siprocket.SipMsg
	Method    string
	CallID    string
	StartLine string
	CSeq      string
	ToAddr    string
	FromAddr  string
}

// 实现 String() 方法来模拟 Python 的 repr
func (sip Message) String() string {
	return fmt.Sprintf(" call_id=%s, method=%s, c_seq=%s start_line=%s", sip.CallID, sip.Method, sip.CSeq, sip.StartLine)
}

// SipSessionStatus 表示会话状态的枚举类型
type SipSessionStatus int

const (
	REJECTED  SipSessionStatus = iota // 会话拒绝
	CANCELLED                         // 会话取消
	COMPLETED                         // 会话完成
	INCALL                            // 打入会话
	CALLSETUP                         //会话建立
	UNKNOWN                           //未知
)

// String 方法使得枚举可以以字符串形式表示
func (s SipSessionStatus) String() string {
	switch s {
	case REJECTED:
		return "REJECTED"
	case CANCELLED:
		return "CANCELLED"
	case COMPLETED:
		return "COMPLETED"
	case INCALL:
		return "IN_CALL"
	case CALLSETUP:
		return "CALL_SETUP"
	case UNKNOWN:
		return "UNKNOWN"
	default:
		return "Unknown"
	}
}

// SipSession 表示一个完整的 SIP 会话
type SipSession struct {
	CallID        string     // 唯一标识会话的 Call-ID
	ANI           string     //ANI to_addr
	DNIS          string     //DNIS from_addr
	Via           string     //via 数据
	CallBound     bool       //呼叫方向 true 呼出 false 呼入
	Messages      []*Message // 该会话中所有的请求/响应消息
	CreatedAt     int64      // 会话的创建时间（通常是第一个消息的时间）
	EndedAt       int64      // 会话结束时间（通常是最后一个消息的时间）
	InviteTime    int64      //发起通话时间 Milliseconds
	RingTime      int64      //响铃时间 Milliseconds
	AnswerTime    int64      //应答时间 Milliseconds
	HangUpTime    int64      //挂起时间 Milliseconds
	Duration      int64      // 会话持续时长 Milliseconds
	IsFirstInvite bool
	IsFirst200    bool

	Stage  string           //会话阶段
	Status SipSessionStatus // 会话的状态（如进行中、已结束等）
}

func (s SipSession) String() string {
	return fmt.Sprintf(" call_id=%s, length=%d, dur=%d stage=%s status=%s invite=%d dur=%d", s.CallID, len(s.Messages), s.Duration, s.Stage, s.Status, s.InviteTime, s.Duration)
}

// 基于传入的消息计算当前会话的状态
func (s *SipSession) CalcStatus(simMsg *siprocket.SipMsg) (*Message, error) {
	method, startLine := utils.GetRequestLine(string(simMsg.Req.Src))
	if startLine == "" {
		return &Message{}, fmt.Errorf("invalid request line")
	}

	callId := string(simMsg.CallId.Value)
	cSeq := string(simMsg.Cseq.Src)
	toAddr := string(simMsg.To.Src)
	fromAddr := string(simMsg.From.Src)

	//fmt.Printf("method(%s) startLine(%s)\n", method, startLine)

	// 创建一个 SIPInfo 实例
	msg := &Message{
		Timestamp: simMsg.Timestamp,
		Method:    method,
		CallID:    callId,
		StartLine: startLine,
		CSeq:      cSeq,
		ToAddr:    toAddr,
		FromAddr:  fromAddr,
	}

	if method == "INVITE" {
		//fmt.Println(method)
		if s.IsFirstInvite { //第一次收到200响应 设置应答时间
			s.InviteTime = msg.Timestamp.UnixMilli() //设置发起时间
			s.Stage = "INVITE"                       //INVITE阶段
			s.IsFirstInvite = false

			s.ANI = msg.ToAddr
			s.DNIS = msg.FromAddr

			if len(simMsg.Via) > 0 {
				via := simMsg.Via[0]
				s.Via = string(via.Src)
			}

			if utils.IsOutbound(s.DNIS) {
				s.CallBound = true
			} else {
				s.CallBound = false
			}

		}
	} else if method == "CANCEL" {
		s.Stage = "CANCEL"   //CANCEL 取消会话
		s.Status = CANCELLED //取消会话
	} else if method == "BYE" {
		if s.IsFirstInvite { //收到BYE但是没收到invite请求 异常请求
			s.Status = UNKNOWN
		} else {
			s.Stage = "BYE OK"
			s.Status = COMPLETED
			s.HangUpTime = msg.Timestamp.UnixMilli() //挂起时间
			s.Duration = s.HangUpTime - s.InviteTime //计算通话时间
		}

	} else if method == "PRACK" {
		s.Stage = "PRACK" //PRACK
	} else if method == "UPDATE" {
		s.Stage = "UPDATE" //UPDATE
	} else if method == "ACK" {
		s.Stage = "ACK" //ACK
	} else {
		if strings.Contains(startLine, "SIP/2.0 100") { //握手阶段
			s.Stage = "Trying"
			s.Status = CALLSETUP
		} else if strings.Contains(startLine, "SIP/2.0 503") {
			s.Stage = "Service Unavailable"
			s.Status = REJECTED //返回503 这是服务不可用
		} else if strings.Contains(startLine, "SIP/2.0 180") {
			s.Stage = "Ringing 180"
			s.RingTime = msg.Timestamp.UnixMilli()
		} else if strings.Contains(startLine, "SIP/2.0 183") {
			s.Stage = "Ringing 183"
			s.RingTime = msg.Timestamp.UnixMilli()

		} else if strings.Contains(startLine, "SIP/2.0 487") {
			s.Stage = "Request Terminated"
			s.Status = REJECTED
		} else if strings.Contains(startLine, "SIP/2.0 408") {
			s.Stage = "Request Timeout"
			s.Status = REJECTED
		} else if strings.Contains(startLine, "SIP/2.0 200 OK") { //对端返回200响应 代表收到
			if s.IsFirst200 { //第一次收到200响应 设置应答时间
				s.AnswerTime = msg.Timestamp.UnixMilli()
				s.IsFirst200 = false
			}
			if strings.Contains(cSeq, "INVITE") { //这是代表对端收到我方INVITE 请求
				s.Stage = "INVITE OK"
			}
			if strings.Contains(cSeq, "PRACK") { //这是代表对端收到我方INVITE 请求
				s.Stage = "PRACK OK"
			}
			if strings.Contains(cSeq, "BYE") { //这是代表对端收到我方INVITE 请求
				if s.IsFirstInvite { //收到BYE但是没收到invite请求 异常请求
					s.Status = UNKNOWN
				} else {
					s.Stage = "BYE OK"
					s.Status = COMPLETED
					s.HangUpTime = msg.Timestamp.UnixMilli() //挂起时间
					s.Duration = s.HangUpTime - s.InviteTime //计算通话时间
				}

			}
		}
	}

	//fmt.Println(msg.String())

	return msg, nil
}

// AddMessage 添加一条消息到会话中
func (s *SipSession) AddMessage(simMsg *siprocket.SipMsg) {
	msg, err := s.CalcStatus(simMsg)
	if err != nil {
		return
	}
	s.Messages = append(s.Messages, msg)
	if len(s.Messages) == 1 {
		s.CreatedAt = msg.Timestamp.UnixMilli()
	}
	s.EndedAt = msg.Timestamp.UnixMilli()
	//s.Duration = s.EndedAt - s.CreatedAt

}

// NewSipSession 创建一个新的 SIP 会话
func NewSipSession(callID string) *SipSession {
	return &SipSession{
		CallID:        callID,
		Status:        UNKNOWN, // 默认状态为进行中
		IsFirst200:    true,
		IsFirstInvite: true,
	}
}

// SipSessionManager 用于管理多个 SIP 会话
type SipSessionManager struct {
	Sessions map[string]*SipSession // 使用 Call-ID 作为键存储会话
}

// NewSipSessionManager 创建一个新的 SIP 会话管理器
func NewSipSessionManager() *SipSessionManager {
	return &SipSessionManager{
		Sessions: make(map[string]*SipSession), // 初始化会话映射
	}
}

// AddSession 向管理器中添加一个新的 SIP 会话
func (manager *SipSessionManager) AddSession(session *SipSession) {
	manager.Sessions[session.CallID] = session
}

// GetSession 获取一个已存在的 SIP 会话
func (manager *SipSessionManager) GetSession(callID string) (*SipSession, bool) {
	session, exists := manager.Sessions[callID]
	return session, exists
}
