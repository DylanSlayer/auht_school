package model

type ExamInfo struct {
	Term      string `json:"term"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	Week      string `json:"week"`
	WeekDay   string `json:"week_day"`
	Class     string `json:"class"`
	Teacher   string `json:"teacher"`
	ClassRoom string `json:"class_room"`
	Time      string `json:"time"`
	Type      string `json:"type"`
}

type CommonExamInfo struct {
	Id         string `json:"id"`
	Campus     string `json:"campus"`
	ExamCampus string `json:"exam_campus"`
	Session    string `json:"session"`
	Number     string `json:"number"`
	Name       string `json:"name"`
	Teacher    string `json:"teacher"`
	Time       string `json:"time"`
	Class      string `json:"class"`
	Seat       string `json:"seat"`
}
