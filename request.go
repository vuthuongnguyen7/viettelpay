package viettelpay

type Envelope interface {
	SetData([]byte)
	SetOrderID(string)
	SetPassword(string)
	SetServiceCode(string)
	SetUsername(string)
}

type Request interface {
	Command() string
	Data() interface{}
	Envelope() Envelope
}

func NewRequest(cmd string, data interface{}, env Envelope) Request {
	return &request{
		cmd:      cmd,
		data:     data,
		envelope: env,
	}
}

type request struct {
	cmd      string
	data     interface{}
	envelope Envelope
}

func (r request) Command() string {
	return r.cmd
}
func (r request) Data() interface{} {
	return r.data
}
func (r request) Envelope() Envelope {
	if r.envelope == nil {
		return &EnvelopeBase{}
	}
	return r.envelope
}
