package response

type DutyChecklistItem struct {
	Date      string `json:"date"`
	Status    string `json:"status"`
	IsChecked bool   `json:"isChecked"`
}
