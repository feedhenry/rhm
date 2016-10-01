package commands

// define all the types used in a common file

// Project defines the project response object
type Project struct {
	AuthorEmail    string      `json:"authorEmail"`
	BusinessObject string      `json:"businessObject"`
	Featured       bool        `json:"featured"`
	GUID           string      `json:"guid"`
	JSONTemplateID string      `json:"jsonTemplateId"`
	Migrated       bool        `json:"migrated"`
	Service        bool        `json:"service"`
	SysCreated     int         `json:"sysCreated"`
	SysGroupFlags  int         `json:"sysGroupFlags"`
	SysGroupList   string      `json:"sysGroupList"`
	SysModified    int         `json:"sysModified"`
	SysOwner       string      `json:"sysOwner"`
	SysVersion     int         `json:"sysVersion"`
	Template       interface{} `json:"template"`
	TemplateID     string      `json:"templateId"`
	Title          string      `json:"title"`
	Type           string      `json:"type"`
	Apps           []*App      `json:"apps,omitempty"`
}

// App defines an app request object
type App struct {
	GUID      string `json:"guid"`
	ScmBranch string `json:"scmBranch"`
	ScmCommit string `json:"scmCommit"`
	ScmURL    string `json:"scmUrl"`
	Title     string `json:"title"`
}
