package epub

import "encoding/xml"

type opfPackage struct {
	XMLName  xml.Name    `xml:"package"`
	Metadata opfMetadata `xml:"metadata"`
	Manifest opfManifest `xml:"manifest"`
	Spine    opfSpine    `xml:"spine"`
}

type opfMetadata struct {
	Title   string   `xml:"title"`
	Creator []string `xml:"creator"`
}

type opfManifest struct {
	Items []manifestItem `xml:"item"`
}

type manifestItem struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type opfSpine struct {
	ItemRefs []spineItemRef `xml:"itemref"`
}

type spineItemRef struct {
	IDRef string `xml:"idref,attr"`
}
