package sip

import (
	"fmt"
	"github.com/marv2097/siprocket"
	"log"
	"sip-parser/pkg/utils"
	"strings"
	"time"
)

type Message struct {
	Timestamp time.Duration
	pct       siprocket.SipMsg
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
	CallID    string           // 唯一标识会话的 Call-ID
	Messages  []Message        // 该会话中所有的请求/响应消息
	CreatedAt int64            // 会话的创建时间（通常是第一个消息的时间）
	EndedAt   int64            // 会话结束时间（通常是最后一个消息的时间）
	Duration  int64            // 会话持续时长
	Stage     string           //会话阶段
	Status    SipSessionStatus // 会话的状态（如进行中、已结束等）
}

// 基于传入的消息计算当前会话的状态
func (s *SipSession) CalcStatus(simMsg SipMessage) (Message, error) {
	method, startLine := utils.GetRequestLine(string(simMsg.pct.Req.Src))
	if startLine == "" {
		return Message{}, fmt.Errorf("invalid request line")
	}

	callId := string(simMsg.pct.CallId.Value)
	cSeq := string(simMsg.pct.Cseq.Src)
	toAddr := string(simMsg.pct.To.Src)
	fromAddr := string(simMsg.pct.From.Src)

	// 创建一个 SIPInfo 实例
	msg := Message{
		Timestamp: simMsg.Timestamp,
		Method:    method,
		CallID:    callId,
		StartLine: startLine,
		CSeq:      cSeq,
		ToAddr:    toAddr,
		FromAddr:  fromAddr,
	}

	if method == "INVITE" {
		s.Stage = "INVITE" //INVITE阶段
	} else if method == "CANCEL" {
		s.Stage = "CANCEL"   //CANCEL 取消会话
		s.Status = CANCELLED //取消会话
	} else if method == "BYE" {
		s.Stage = "BYE" //BYE 结束会话

		//s.Status = BYE  //结束会话
	} else {
		if strings.Contains(startLine, "SIP/2.0 100") { //握手阶段
			s.Stage = "Trying"
		} else if strings.Contains(startLine, "SIP/2.0 503") {
			s.Stage = "Service Unavailable"
			s.Status = REJECTED //返回503 这是服务不可用
		} else if strings.Contains(startLine, "SIP/2.0 180") {
			s.Stage = "Ringing"
			//s.Status = REJECTED //返回503 这是服务不可用
		} else if strings.Contains(startLine, "SIP/2.0 200 OK") { //对端返回200响应 代表收到
			if strings.Contains(cSeq, "INVITE") { //这是代表对端收到我方INVITE 请求
				s.Stage = "INVITE OK"
			}
			if strings.Contains(cSeq, "BYE") { //这是代表对端收到我方INVITE 请求
				s.Stage = "BYE OK"
				s.Status = COMPLETED
			}
		}
	}

	//fmt.Println(msg.String())

	return msg, nil
}

// AddMessage 添加一条消息到会话中
func (s *SipSession) AddMessage(simMsg SipMessage) {
	msg, err := s.CalcStatus(simMsg)
	if err != nil {
		return
	}
	s.Messages = append(s.Messages, msg)
	if len(s.Messages) == 1 {
		s.CreatedAt = msg.Timestamp.Microseconds()
	}
	s.EndedAt = msg.Timestamp.Microseconds()
	s.Duration = s.EndedAt - s.CreatedAt

}

// NewSipSession 创建一个新的 SIP 会话
func NewSipSession(callID string) *SipSession {
	return &SipSession{
		CallID: callID,
		Status: UNKNOWN, // 默认状态为进行中
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

func (m *SipSessionManager) GetAllSessions() {
	for call_id, session := range m.Sessions {
		log.Println(call_id, len(session.Messages))
	}
}
