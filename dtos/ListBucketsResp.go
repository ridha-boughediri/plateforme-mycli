package dtos

import (
	"encoding/xml"
	"time"
)

type ListAllMyBucketsResult struct {
	XMLName xml.Name      `xml:"ListAllMyBucketsResult"`
	Xmlns   string        `xml:"xmlns,attr"`
	Buckets []ListBuckets `xml:"Buckets>Bucket"`
}

type ListBuckets struct {
	CreationDate time.Time `xml:"CreationDate"`
	Name         string    `xml:"Name"`
}
