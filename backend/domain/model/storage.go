package model

type BucketItemUrl struct {
	Name string
	URL  string
}

type BucketItemData struct {
	Name  string
	Data  []byte
	CType string
}
