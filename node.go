// additions to the node type

package main

type el int16

const (
	LEFT el = iota
	PROP
	RIGHT
	VALUE
	ID
)

var propKey = []string {"L", "P", "R", "V", "I"}

func ValueKey(element el, value string) string {
	return propKey[element] + value
}

func Key(element el, id int64) string {
	return propKey[element] + string(id)
}

func (n Node) Tokens() []string {
	result := []string{Key(ID, *n.Id), ValueKey(VALUE, *n.Value)}

	if n.Edge != nil {
		result = append(result, Key(LEFT, *n.Edge.Left), Key(PROP, *n.Edge.Prop), Key(RIGHT, *n.Edge.Right))
	}

	return result
}

