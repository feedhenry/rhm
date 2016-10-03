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

// ProjectTemplate defines the project template response object
type ProjectTemplate struct {
	ID       string `json:"id"`
	Title    string `json:"name"`
	Category string `json:"category"`
}

// Environment defines the environments response object
type Environment struct {
	ID      string  `json:"id"`
	Label   string  `json:"label"`
	Token   string  `json:"token"`
	Domain  string  `json:"domain"`
	UID     string  `json:"uid"`
	Enabled bool    `json:"enabled"`
	Target  *Target `json:"target,omitempty"`
}

// MBaaS targets
type Target struct {
	ID          string `json:"id"`
	Domain      string `json:"domain"`
	Owner       string `json:"owner"`
	FhMbaasHost string `json:"fhMbaasHost"`
	Label       string `json:"label"`
	URL         string `json:"url"`
	RouterDNS   string `json:"routerDns"`
	BearerToken string `json:"bearerToken"`
	ServiceKey  string `json:"servicekey"`
	NagiosURL   string `json:"nagiosUrl"`
	Decoupled   bool   `json:"decoupled"`
	Editable    bool   `json:"editable"`
	Enabled     bool   `json:"enabled"`
	Type        string `json:"type"`
	Env         string `json:"_env"`
	Description string `json:"_label"`
}
