package network

import (
	"encoding/xml"
	"github.com/samber/lo"
)

// GrobidTEI represents the root element of the TEI XML document.
type GrobidTEI struct {
	XMLName   xml.Name  `xml:"TEI"`
	XMLSpace  string    `xml:"xml:space,attr"`
	XMLNS     string    `xml:"xmlns,attr"`
	XSI       string    `xml:"xmlns:xsi,attr"`
	SchemaLoc string    `xml:"xsi:schemaLocation,attr"`
	XLink     string    `xml:"xmlns:xlink,attr"`
	Header    TEIHeader `xml:"teiHeader"`
	Text      TEIText   `xml:"text"`
}

func (g GrobidTEI) String() string {
	v := lo.Must(xml.Marshal(&g))
	return string(v)
}
func NewGrobidTEI(str string) (GrobidTEI, error) {
	var response GrobidTEI
	err := xml.Unmarshal([]byte(str), &response)
	if err != nil {
		return GrobidTEI{}, err
	}
	return response, nil
}

// TEIHeader represents the teiHeader element containing metadata.
type TEIHeader struct {
	XMLLang      string       `xml:"xml:lang,attr"`
	FileDesc     FileDesc     `xml:"fileDesc"`
	EncodingDesc EncodingDesc `xml:"encodingDesc"`
	ProfileDesc  ProfileDesc  `xml:"profileDesc"`
}

// FileDesc represents the fileDesc element containing bibliographic information.
type FileDesc struct {
	TitleStmt   TitleStmt   `xml:"titleStmt"`
	Publication Publication `xml:"publicationStmt"`
	SourceDesc  SourceDesc  `xml:"sourceDesc"`
}

// TitleStmt represents the title statement containing title and funders.
type TitleStmt struct {
	Title   Title    `xml:"title"`
	Funders []Funder `xml:"funder"`
}

// Title represents a title element.
type Title struct {
	Level string `xml:"level,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// Funder represents a funding organization.
type Funder struct {
	Ref     string  `xml:"ref,attr,omitempty"`
	OrgName OrgName `xml:"orgName"`
}

// OrgName represents the organization name.
type OrgName struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

// Publication represents the publication statement.
type Publication struct {
	Publisher    string       `xml:"publisher"`
	Availability Availability `xml:"availability"`
}

// Availability represents the availability of the publication.
type Availability struct {
	Status  string `xml:"status,attr"`
	Licence string `xml:"licence"`
}

// SourceDesc represents the source description containing bibliographic structure.
type SourceDesc struct {
	BiblStruct BiblStruct `xml:"biblStruct"`
}

// BiblStruct represents the bibliographic structure.
type BiblStruct struct {
	Analytic Analytic `xml:"analytic"`
	Monogr   Monogr   `xml:"monogr"`
	IDNo     IDNo     `xml:"idno"`
}

// Analytic represents the analytic element containing authors and title.
type Analytic struct {
	Authors []Author `xml:"author"`
	Title   Title    `xml:"title"`
}

// Author represents an author element.
type Author struct {
	PersName    PersName     `xml:"persName"`
	Affiliation *Affiliation `xml:"affiliation,omitempty"`
}

// PersName represents the personal name of the author.
type PersName struct {
	Forename []Forename `xml:"forename"`
	Surname  string     `xml:"surname"`
}

// Forename represents the forename element, which can be first or middle.
type Forename struct {
	Type  string `xml:"type,attr,omitempty"`
	Value string `xml:",chardata"`
}

// Affiliation represents the affiliation of the author.
type Affiliation struct {
	Key      string    `xml:"key,attr,omitempty"`
	OrgNames []OrgName `xml:"orgName"`
}

// Monogr represents the monographic element.
type Monogr struct {
	Imprint Imprint `xml:"imprint"`
}

// Imprint represents the imprint information.
type Imprint struct {
	Date string `xml:"date,omitempty"`
}

// IDNo represents the identifier number.
type IDNo struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

// EncodingDesc represents the encoding description.
type EncodingDesc struct {
	AppInfo AppInfo `xml:"appInfo"`
}

// AppInfo represents application information.
type AppInfo struct {
	Application Application `xml:"application"`
}

// Application represents the application element.
type Application struct {
	Version string `xml:"version,attr"`
	Ident   string `xml:"ident,attr"`
	When    string `xml:"when,attr"`
	Desc    string `xml:"desc"`
	Ref     string `xml:"ref"`
}

// ProfileDesc represents the profile description containing the abstract.
type ProfileDesc struct {
	Abstract Abstract `xml:"abstract"`
}

// Abstract represents the abstract element.
type Abstract struct {
	Div AbstractDiv `xml:"div"`
}

type AbstractDiv struct {
	Paragraphs []Paragraph `xml:"p"`
}

// TEIText represents the text element of the TEI document.
type TEIText struct {
	XMLLang string `xml:"xml:lang,attr"`
	Body    Body   `xml:"body"`
}

// Body represents the body of the text.
type Body struct {
	Divs []Div `xml:"div"`
}

// Div represents a division in the body.
type Div struct {
	XMLName    xml.Name    `xml:"div"`
	Head       Head        `xml:"head"`
	Paragraphs []Paragraph `xml:"p"`
}

// Head represents a heading element.
type Head struct {
	N     string `xml:"n,attr"`
	Value string `xml:",chardata"`
}

// Paragraph represents a paragraph element.
type Paragraph struct {
	Content string `xml:",chardata"`
}
