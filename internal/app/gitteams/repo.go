package gitteams

// NewRepo creates a new instance of a Repo.
func NewRepo(vcsType string, origin *RepoOrigin) Repo {
	return Repo{
		Type:   vcsType,
		Origin: origin,
		Data:   map[string]interface{}{},
	}
}

// Repo represents a repository and also serves as a DTO between processors.
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

// RepoOrigin defines the origin or the repository.
type RepoOrigin struct {
	Name  string
	Short string
}
