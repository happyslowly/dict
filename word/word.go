package word

type Word struct {
	Title       string       `json:"title"`
	ProAmE      string       `json:"pro_ame"`
	ProBrE      string       `json:"pro_bre"`
	Definitions []Definition `json:"definitions"`
	Classes     []string
}

type Definition struct {
	Class        string   `json:"class"`
	Translations []string `json:"translations"`
}
