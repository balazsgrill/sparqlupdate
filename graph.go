package sparqlupdate

import (
	"strings"

	"github.com/knakk/rdf"
)

type Graph struct {
	data map[string]map[string][]string
}

func New() *Graph {
	return &Graph{
		data: make(map[string]map[string][]string),
	}
}

func (graph *Graph) addInternal(subject string, predicate string, object string) {
	properties, present := graph.data[subject]
	if !present {
		properties = make(map[string][]string)
		graph.data[subject] = properties
	}
	properties[predicate] = append(properties[predicate], object)
}

func (graph *Graph) Add(triple rdf.Triple) {
	graph.AddTriple(triple.Subj, triple.Pred, triple.Obj)
}

func (graph *Graph) AddTriple(subj rdf.Term, pred rdf.Term, obj rdf.Term) {
	graph.addInternal(subj.Serialize(rdf.Turtle), pred.Serialize(rdf.Turtle), obj.Serialize(rdf.Turtle))
}

func (graph *Graph) UpdateQuery(namedgraph rdf.Term) string {
	var result strings.Builder

	result.WriteString("INSERT DATA\n{")
	if namedgraph != nil {
		result.WriteString("GRAPH ")
		result.WriteString(namedgraph.Serialize(rdf.Turtle))
		result.WriteString(" {\n")
	}
	for subject, properties := range graph.data {
		result.WriteString(subject)
		result.WriteString(" ")
		first := true
		for predicate, objects := range properties {
			if first {
				first = false
			} else {
				result.WriteString(" ;\n\t")
			}
			result.WriteString(predicate)
			result.WriteString(" ")
			result.WriteString(strings.Join(objects, " , "))
		}
		result.WriteString(" .\n")
	}
	if namedgraph != nil {
		result.WriteString("}")
	}
	result.WriteString("}")
	return result.String()
}
