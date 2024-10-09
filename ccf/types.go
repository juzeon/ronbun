package ccf

// Conference represents the validation schema for CCFDDL Conference Specification 3.0.X.
type Conference struct {
	Title       string              `json:"title"`       // Short conference name, without year, uppercase
	Description string              `json:"description"` // Description, or long name, with no session
	Sub         string              `json:"sub"`         // The category that the conference is labeled by CCF.
	Rank        Rank                `json:"rank"`        // The ranking information of the conference
	DBLP        string              `json:"dblp"`        // The suffix in dblp url, e.g., iccv in dblp.uni-trier.de/db/conf/iccv
	Confs       []ConferenceDetails `json:"confs"`       // List of conference details
}

// Rank represents the ranking information of the conference.
type Rank struct {
	CCF   string `json:"ccf"`   // The level that the conference is ranked by CCF, e.g., A, B, C
	CORE  string `json:"core"`  // The level that the conference is ranked by CORE, e.g., A*, A, B, C
	THCPL string `json:"thcpl"` // The level that the conference is ranked by Tsinghua, e.g., A, B, C
}

// ConferenceDetails represents the details of a specific conference.
type ConferenceDetails struct {
	Year     int        `json:"year"`     // Year the conference is happening
	ID       string     `json:"id"`       // Conference name & year, lowercase
	Link     string     `json:"link"`     // URL to the conference home page
	Timeline []Timeline `json:"timeline"` // List of important dates and deadlines
	Timezone string     `json:"timezone"` // Timezone of the conference
	Date     string     `json:"date"`     // When the main conference is happening, e.g., Mar 12-16, 2021
	Place    string     `json:"place"`    // Where the main conference is happening, e.g., city, country
}

// Timeline represents important dates and deadlines for the conference.
type Timeline struct {
	AbstractDeadline string `json:"abstract_deadline"` // Abstract deadline if applicable, optional
	Deadline         string `json:"deadline"`          // Deadline, in the format of yyyy-mm-dd hh:mm:ss or TBD
	Comment          string `json:"comment"`           // Some comments on the conference, optional
}
