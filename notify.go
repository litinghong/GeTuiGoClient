package GeTuiGo

type INotify interface {
	GetNotify() interface{}
}

// 消息应用类型
const (
	TypeNotification  = "notification"
	TypeLink          = "link"
	TypeNotypopload   = "notypopload"
	TypeStartActivity = "startactivity"
	TypeTransmission  = "transmission"
)

// 点开通知打开应用模板
type TmplNotification struct {
	TransmissionType    bool   `json:"transmission_type"`    // 收到消息是否立即启动应用，true为立即启动，false则广播等待启动，默认是否
	TransmissionContent string `json:"transmission_content"` // 透传内容
	DurationBegin       string `json:"duration_begin"`       // 设定展示开始时间，格式为yyyy-MM-dd HH:mm:ss
	DurationEnd         string `json:"duration_end"`         // 设定展示结束时间，格式为yyyy-MM-dd HH:mm:ss
	Style               IStyle `json:"style"`                // 通知栏消息布局样式，见底下Style说明
}

func (t TmplNotification) GetNotify() interface{} {
	return t
}

// 点开通知打开网页模板
type TmplLink struct {
	Url           string `json:"url"`            // 必传: 打开网址
	DurationBegin string `json:"duration_begin"` // 设定展示开始时间，格式为yyyy-MM-dd HH:mm:ss
	DurationEnd   string `json:"duration_end"`   // 设定展示结束时间，格式为yyyy-MM-dd HH:mm:ss
	Style         IStyle `json:"style"`          // 通知栏消息布局样式，见底下Style说明
}

func (t TmplLink) GetNotify() interface{} {
	return t
}

// 点击通知弹窗下载模板
type TmplNotifyPopLoad struct {
	NotifyIcon    string `json:"notyicon"`       //	必传: 通知栏图标
	NotifyTitle   string `json:"notytitle"`      //	必传: 通知标题
	NotifyContent string `json:"notycontent"`    //	必传: 通知内容
	PopTitle      string `json:"poptitle"`       //	必传: 弹出框标题
	PopContent    string `json:"popcontent"`     //	必传: 弹出框内容
	PopImage      string `json:"popimage"`       //	必传: 弹出框图标
	PopButton1    string `json:"popbutton_1"`    //	必传: 弹出框左边按钮名称
	PopButton2    string `json:"popbutton_2"`    //	必传: 弹出框右边按钮名称
	LoadIcon      string `json:"loadicon"`       //	现在图标
	LoadTitle     string `json:"loadtitle"`      //	下载标题
	LoadUrl       string `json:"loadurl"`        //	必传:下载文件地址
	IsAutoInstall bool   `json:"is_autoinstall"` //	是否自动安装，默认值false
	IsActive      bool   `json:"is_actived"`     //	安装完成后是否自动启动应用程序，默认值false
	AndroidMark   string `json:"androidmark"`    //	安卓标识
	SymbianMark   string `json:"symbianmark"`    //	塞班标识
	IphoneMark    string `json:"iphonemark"`     //	苹果标志
	DurationBegin string `json:"duration_begin"` //	设定展示开始时间，格式为yyyy-MM-dd HH:mm:ss
	DurationEnd   string `json:"duration_end"`   //	设定展示结束时间，格式为yyyy-MM-dd HH:mm:ss
}

func (t TmplNotifyPopLoad) GetNotify() interface{} {
	return t
}

// 点开通知打开应用内特定页面模板
type TmplStartActivity struct {
	TransmissionType    bool   `json:"transmission_type"`    // 收到消息是否立即启动应用，true为立即启动，false则广播等待启动，默认否
	TransmissionContent string `json:"transmission_content"` // 透传内容
	DurationBegin       string `json:"duration_begin"`       // 设定展示开始时间，格式为yyyy-MM-dd HH:mm:ss
	DurationEnd         string `json:"duration_end"`         // 设定展示结束时间，格式为yyyy-MM-dd HH:mm:ss
	// 必传: 应用内页面intent 【Android】长度小于1000字节，
	//  intent参数（以intent:开头;end结尾）
	//  示例：intent:#Intent;component=你的包名/你要打开的 activity 全路径;S.parm1=value1;S.parm2=value2;end
	Intent IStyle `json:"intent"`
}

func (t TmplStartActivity) GetNotify() interface{} {
	return t
}

// 透传消息模板
type TmplTransmission struct {
	TransmissionType    bool   `json:"transmission_type"`    // 收到消息是否立即启动应用，true为立即启动，false则广播等待启动，默认是否
	TransmissionContent string `json:"transmission_content"` // 必传:透传内容
	DurationBegin       string `json:"duration_begin"`       // 设定展示开始时间，格式为yyyy-MM-dd HH:mm:ss
	DurationEnd         string `json:"duration_end"`         // 设定展示结束时间，格式为yyyy-MM-dd HH:mm:ss
}

func (t TmplTransmission) GetNotify() interface{} {
	return t
}
