package gogent

type Gopher struct {
	Role      string
	Goal      string
	Backstory string
}

// Create a new Agent
func NewGopher() *Gopher {
	return &Gopher{}
}

func (g *Gopher) WithRole(role string) *Gopher {
	g.Role = role
	return g
}

func (g *Gopher) WithGoal(goal string) *Gopher {
	g.Goal = goal
	return g
}

func (g *Gopher) WithBackstory(backstory string) *Gopher {
	g.Backstory = backstory
	return g
}
