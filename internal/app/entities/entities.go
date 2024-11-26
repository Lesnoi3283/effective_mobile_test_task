package entities

type Song struct {
	ID          uint64 `gorm:"primary_key"`
	Song        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
