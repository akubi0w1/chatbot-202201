package slack

// channel
const (
	Ch_Admin   string = "C02QM1PR8TF" // manage-butler
	Ch_General string = "C01RCGC1F4L" // general
	Ch_Random  string = "C01R6HEH9T5" // random
)

//color
// ref: https://colorswall.com/palette/3/
const (
	ColorPrimary string = "#0275d8" // 青
	ColorSuccess string = "#5cb85c" // 緑
	ColorInfo    string = "#5bc0de" //水色
	ColorWarning string = "#f0ad4e" // オレンジ
	ColorDanger  string = "#d9534f" // 赤
)

// callback
const (
	CallbackSelectMenu   string = "callbackSelectMenu"
	CallbackFixedContact string = "callbackFixedContact"
)

// action
const (
	ActionSelectContact string = "actionSelectContact"
	ActionSubmitContact string = "actionSubmitContact"
	ActionCancelContact string = "actionCancelContact"
	ActionFixedContact  string = "actionFixedContact"
)
