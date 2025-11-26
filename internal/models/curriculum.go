package models

import "time"

type Curriculum struct {
	ID        int64
	FullName  string
	Email     string
	Phone     string
	BirthDate *time.Time

	Summary    string
	Profession string
	Experience []ExperienceEntry
	Education  []EducationEntry
	Skills     []string
	Languages  []LanguageEntry

	User User
	BaseModel
}

type ExperienceEntry struct {
	Company     string
	Role        string
	Description string
	StartDate   time.Time
	EndDate     *time.Time
}

type EducationEntry struct {
	Institution string
	Degree      EducationDegree
	StartDate   time.Time
	EndDate     *time.Time
}

type LanguageEntry struct {
	Name  string
	Level LanguageLevel
}

type EducationDegree string

const (
	DegreeElementary   EducationDegree = "elementary"
	DegreeHighSchool   EducationDegree = "high_school"
	DegreeTechnical    EducationDegree = "technical"
	DegreeAssociate    EducationDegree = "associate"
	DegreeBachelor     EducationDegree = "bachelor"
	DegreeLicentiate   EducationDegree = "licentiate"
	DegreePostgraduate EducationDegree = "postgraduate"
	DegreeMaster       EducationDegree = "master"
	DegreeDoctorate    EducationDegree = "doctorate"
	DegreeMBA          EducationDegree = "mba"
	DegreeCertificate  EducationDegree = "certificate"
)

type LanguageLevel string

const (
	LanguageBasic        LanguageLevel = "basic"
	LanguageIntermediate LanguageLevel = "intermediate"
	LanguageAdvanced     LanguageLevel = "advanced"
	LanguageFluent       LanguageLevel = "fluent"
	LanguageNative       LanguageLevel = "native"
)
