package GeTuiGo

const (
	ResultOk                 = "ok"                   // 成功
	ResultNoMsg              = "no_msg"               // 没有消息体
	ResultAliasError         = "alias_error"          // 找不到别名
	ResultBlackIp            = "black_ip"             // 黑名单ip
	ResultSignError          = "sign_error"           // 鉴权失败
	ResultPushNumOverLimit   = "pushnum_overlimit"    // 推送次数超限
	ResultNoAppid            = "no_appid"             // 找不到appid
	ResultNoUser             = "no_user"              // 找不到对应用户
	ResultTooFrequent        = "too_frequent"         // 推送过于频繁
	ResultSensitiveWord      = "sensitive_word"       // 有敏感词出现
	ResultAppidNotMatch      = "appid_notmatch"       // appid与cid或者appkey不匹配
	ResultNotAuth            = "not_auth"             // 用户没有鉴权
	ResultBlackAppId         = "black_appid"          // 黑名单app
	ResultInvalidParam       = "invalid_param"        // 参数检验不通过
	ResultAliasNotBind       = "alias_notbind"        // 别名没有绑定cid
	ResultTagOverLimit       = "tag_over_limit"       // tag个数超限
	ResultSuccessOnline      = "successed_online"     // 在线下发
	ResultSuccessOffline     = "successed_offline"    // 离线下发
	ResultTagInvalidOrNoAuth = "taginvalid_or_noauth" // tag无效或者没有使用权限
	ResultNoValidPush        = "no_valid_push"        // 没有有效下发
	ResultSuccessIgnore      = "successed_ignore"     // 忽略非活跃用户
	ResultNoTaskId           = "no_taskid"            // 找不到taskid
	ResultOtherError         = "other_error"          // 其他错误
)

type IStyle interface {
}

// 系统样式
type StyleSystem struct {
	Type        int    `json:"type"`          // 必传: 系统样式=0,个推样式=1,纯图样式(背景图样式)=4,展开通知样式=6
	Text        string `json:"text"`          // 必传: 通知内容,
	Title       string `json:"title"`         // 必传: 通知标题,
	Logo        string `json:"logo"`          // 必传: 通知的图标名称，包含后缀名（需要在客户端开发时嵌入），如“push.png”,
	BigStyle    int    `json:"big_style"`     // 通知展示样式,枚举值包括 1,2(style=0时，big_style只能为1或者2),
	BigImageUrl string `json:"big_image_url"` // 知展示文本+大图样式，参数 大图URL地址,
	LogoUrl     string `json:"logourl"`       // "http://xxxx/a.png",
	IsRing      bool   `json:"is_ring"`       // 收到通知是否响铃：true响铃，false不响铃。默认响铃,
	IsVibrate   bool   `json:"is_vibrate"`    // 收到通知是否振动：true振动，false不振动。默认振动,
	IsClearable bool   `json:"is_clearable"`  // 通知是否可清除： true可清除，false不可清除。默认可清除,
	NotifyId    int    `json:"notify_id"`     // 需要被覆盖的消息已经增加了notifyId字段，用于实现下发消息的覆盖。新的消息使用相同的notifyId下发。
	ChannelId   string `json:"channel_id"`    // 通知渠道id,唯一标识,默认Default
	ChannelName string `json:"channel_name"`  // 通知渠道名称，默认Default
	// 该字段代表通知渠道重要性，具体值有0、1、2、3、4；
	//  设置之后不能修改；具体展示形式如下：
	//   0：无声音，无震动，不显示。(不推荐)
	//   1：无声音，无震动，锁屏不显示，通知栏中被折叠显示，导航栏无logo。
	//   2：无声音，无震动，锁屏和通知栏中都显示，通知不唤醒屏幕。
	//   3：有声音，有震动，锁屏和通知栏中都显示，通知唤醒屏幕。（推荐）
	//   4：有声音，有震动，亮屏下通知悬浮展示，锁屏通知以默认形式展示且唤醒屏幕。（推荐）
	ChannelLevel int `json:"channel_level"`
}

func NewStyleSystem() StyleSystem {
	return StyleSystem{
		Type:         0,
		Text:         "",
		Title:        "",
		Logo:         "",
		BigStyle:     1,
		BigImageUrl:  "",
		LogoUrl:      "",
		IsRing:       true,
		IsVibrate:    true,
		IsClearable:  true,
		NotifyId:     0,
		ChannelId:    "Default",
		ChannelName:  "Default",
		ChannelLevel: 0,
	}
}

// 纯图样式(背景图样式)
type StyleImage struct {
	Type        int    `json:"type"`         // 必传: 系统样式=0,个推样式=1,纯图样式(背景图样式)=4,展开通知样式=6
	Logo        string `json:"logo"`         // 必传: 通知的图标名称，包含后缀名（需要在客户端开发时嵌入），如“push.png”,
	BannerUrl   string `json:"banner_url"`   // 必传: 通过url方式指定动态banner图片作为通知背景图
	IsRing      bool   `json:"is_ring"`      // 收到通知是否响铃：true响铃，false不响铃。默认响铃,
	IsVibrate   bool   `json:"is_vibrate"`   // 收到通知是否振动：true振动，false不振动。默认振动,
	IsClearable bool   `json:"is_clearable"` // 通知是否可清除： true可清除，false不可清除。默认可清除,
}

// 个推样式
type StyleGeTui struct {
	Type        int    `json:"type"`         // 必传: 系统样式=0,个推样式=1,纯图样式(背景图样式)=4,展开通知样式=6
	Text        string `json:"text"`         // 必传: 通知内容,
	Title       string `json:"title"`        // 必传: 通知标题,
	Logo        string `json:"logo"`         // 必传: 通知的图标名称，包含后缀名（需要在客户端开发时嵌入），如“push.png”,
	LogoUrl     string `json:"logourl"`      // "http://xxxx/a.png",
	IsRing      bool   `json:"is_ring"`      // 收到通知是否响铃：true响铃，false不响铃。默认响铃,
	IsVibrate   bool   `json:"is_vibrate"`   // 收到通知是否振动：true振动，false不振动。默认振动,
	IsClearable bool   `json:"is_clearable"` // 通知是否可清除： true可清除，false不可清除。默认可清除,
	NotifyId    int    `json:"notify_id"`    // 需要被覆盖的消息已经增加了notifyId字段，用于实现下发消息的覆盖。新的消息使用相同的notifyId下发。
}

// 展开通知样式
type StyleExt struct {
	Type        int    `json:"type"`          // 必传: 系统样式=0,个推样式=1,纯图样式(背景图样式)=4,展开通知样式=6
	Text        string `json:"text"`          // 必传: 通知内容,
	Title       string `json:"title"`         // 必传: 通知标题,
	Logo        string `json:"logo"`          // 必传: 通知的图标名称，包含后缀名（需要在客户端开发时嵌入），如“push.png”,
	LogoUrl     string `json:"logourl"`       // "http://xxxx/a.png",
	BigStyle    int    `json:"big_style"`     // 通知展示样式,枚举值包括 1,2(style=0时，big_style只能为1或者2),
	BigImageUrl string `json:"big_image_url"` // 知展示文本+大图样式，参数 大图URL地址,
	BigText     string `json:"big_text"`      // 通知展示文本+长文本样式，参数是长文本
	BannerUrl   string `json:"banner_url"`    // 必传: 通过url方式指定动态banner图片作为通知背景图
	IsRing      bool   `json:"is_ring"`       // 收到通知是否响铃：true响铃，false不响铃。默认响铃,
	IsVibrate   bool   `json:"is_vibrate"`    // 收到通知是否振动：true振动，false不振动。默认振动,
	IsClearable bool   `json:"is_clearable"`  // 通知是否可清除： true可清除，false不可清除。默认可清除,
	NotifyId    int    `json:"notify_id"`     // 需要被覆盖的消息已经增加了notifyId字段，用于实现下发消息的覆盖。新的消息使用相同的notifyId下发。
	ChannelId   string `json:"channel_id"`    // 通知渠道id,唯一标识,默认Default
	ChannelName string `json:"channel_name"`  // 通知渠道名称，默认Default
	// 该字段代表通知渠道重要性，具体值有0、1、2、3、4；
	//  设置之后不能修改；具体展示形式如下：
	//   0：无声音，无震动，不显示。(不推荐)
	//   1：无声音，无震动，锁屏不显示，通知栏中被折叠显示，导航栏无logo。
	//   2：无声音，无震动，锁屏和通知栏中都显示，通知不唤醒屏幕。
	//   3：有声音，有震动，锁屏和通知栏中都显示，通知唤醒屏幕。（推荐）
	//   4：有声音，有震动，亮屏下通知悬浮展示，锁屏通知以默认形式展示且唤醒屏幕。（推荐）
	ChannelLevel int `json:"channel_level"`
}

// apns推送消息, json串，当手机为ios，并且为离线的时候
type ApnPushInfo struct {
	Aps struct {
		Alert            map[string]interface{} `json:"alert"`
		AutoBadge        string                 `json:"auto_badge"`
		Sound            string                 `json:"sound"`
		ContentAvailable int                    `json:"content-available"`
		Category         string                 `json:"category"`
	} `json:"aps"`
	Payload    string              `json:"payload"`
	Multimedia []map[string]string `json:"multimedia"`
}

// Message 消息
type Message struct {
	AppKey            string `json:"appkey"`              // 注册应用时生成的appkey
	IsOffline         bool   `json:"is_offline"`          // 是否离线推送
	OfflineExpireTime int    `json:"offline_expire_time"` // 消息离线存储有效期，单位：ms
	PushNetworkType   int    `json:"push_network_type"`   // 选择推送消息使用网络类型，0：不限制，1：wifi
	MsgType           string `json:"msgtype"`             // 消息应用类型， 选项：notification、link、notypopload、startactivity, transmission 选择不同消息类型对应不同消息模版，具体消息模版 见详情
}

// 新消息
//  msgType 请使用常量TypeXXX来设置
func NewMessage(msgType string) *Message {
	return &Message{
		MsgType: msgType,
	}
}
