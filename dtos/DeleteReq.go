package dtos

type Delete struct {
	Quiet  bool `xml:"Quiet"`
	Object struct {
		Key string `xml:"Key"`
	} `xml:"Object"`
}

type Deleted struct {
	Key string `xml:"Key"`
}
