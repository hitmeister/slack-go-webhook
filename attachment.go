package main

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Attachment struct {
	Fallback  *string  `json:"fallback"`
	Color     *string  `json:"color"`
	PreText   *string  `json:"pretext"`
	Title     *string  `json:"title"`
	TitleLink *string  `json:"title_link"`
	Text      *string  `json:"text"`
	Fields    []*Field `json:"fields"`
}

func (a *Attachment) AddField(field Field) *Attachment {
	a.Fields = append(a.Fields, &field)
	return a
}
