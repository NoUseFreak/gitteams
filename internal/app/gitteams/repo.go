package gitteams

func NewRepo(vcsType string, origin *RepoOrigin) Repo {
	return Repo{
		Type:   vcsType,
		Origin: origin,
		Data:   map[string]interface{}{},
	}
}

type Repo struct {
	URL        string
	Origin     *RepoOrigin
	Name       string
	Type       string
	MainBranch string
	TmpDir     string
	Branches   []string
	Data       map[string]interface{}
}

type RepoOrigin struct {
	Name  string
	Short string
}
