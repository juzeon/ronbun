package crawler

type ConferenceInstance struct {
	Slug    string
	Title   string
	Year    int
	TocLink string
}
type Paper struct {
	Title              string
	DBLPLink           string
	DOILink            string
	ConferenceInstance ConferenceInstance
}
