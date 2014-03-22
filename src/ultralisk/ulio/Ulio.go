package ulio

const (
	DEFAULT_BUFFER_SIZE = 16 * 1024
)

// 消息的IObuffer
// 消息buffer的处理,使用类似栈的方式进行处理
type IOBuffer interface {
	Init(inputMsgLenLimit, outputMsgLenLimit uint32) error // 初始化
	// 输入数据处理
	PushInputData([]byte) error    // 压入输入数据
	PopInputData() ([]byte, error) // 弹出输入数据
	// 输出数据处理
	PushOutputData([]byte) error    // 压入输出数据
	PopOutputData() ([]byte, error) // 弹出输出数据
}
