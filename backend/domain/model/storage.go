package model

// RADEK: znovu rozdělit struct per soubor, platí všude
type BucketItemUrl struct {
	Name string
	URL  string
}

type BucketItemData struct {
	Name  string
	Data  []byte
	CType string
}
