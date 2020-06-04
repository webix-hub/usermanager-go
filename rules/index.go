package rules

type Type int

type Description struct {
	ID    Type   `json:"id"`
	Short string `json:"short"`
	Long  string `json:"long"`
}

const (
	_ Type = iota
	CanSeeUsers
	CanEditUsers
	CanAdminProjects
	CanAdminTasks
	CanSeeTasks
	CanAdminEmails
	CanSeeEmails
	CanCreateReports
	CanSeeReports
)

func GetDetails() []Description {
	i := make([]Description, 9)

	i[0] = Description{CanSeeUsers, "CanSeeUsers", "Can see user details and access levels"}
	i[1] = Description{CanEditUsers, "CanEditUsers", "Can modify user details and access levels"}
	i[2] = Description{CanAdminProjects, "CanAdminProjects", "Can create projects"}
	i[3] = Description{CanAdminTasks, "CanAdminTasks", "Can create tasks"}
	i[4] = Description{CanSeeTasks, "CanSeeTasks", "Can see tasks"}
	i[5] = Description{CanSeeEmails, "CanSeeEmails", "Can see custom emails"}
	i[6] = Description{CanAdminEmails, "CanAdminEmails", "Can create custom emails"}
	i[7] = Description{CanCreateReports, "CanCreateReports", "Can create reports"}
	i[8] = Description{CanSeeReports, "CanSeeReports", "Can see reports"}

	return i
}

func GetRuleName(ruleIdx int) string {
	rulesDescriptions := GetDetails()
	ruleDescription := Description{}
	for idx, value := range rulesDescriptions {
		if int(value.ID) == ruleIdx {
			ruleDescription = rulesDescriptions[idx]
			break
		}
	}
	return ruleDescription.Short
}
