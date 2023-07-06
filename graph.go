package sparqlupdate

import (
	"fmt"
	"strings"

	"github.com/knakk/rdf"
)

type Graph struct {
	nodes        map[string]rdf.Term
	data         map[string]map[string][]string
	blankcounter int
}

func New() *Graph {
	return &Graph{
		nodes:        make(map[string]rdf.Term),
		data:         make(map[string]map[string][]string),
		blankcounter: 0,
	}
}

func (graph *Graph) NewBlank() rdf.Term {
	graph.blankcounter += 1
	b, _ := rdf.NewBlank(fmt.Sprintf("%d", graph.blankcounter))
	return b
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

func (graph *Graph) serializeTerm(term rdf.Term) string {
	str := term.Serialize(rdf.Turtle)
	graph.nodes[str] = term
	return str
}

func (graph *Graph) AddTriple(subj rdf.Term, pred rdf.Term, obj rdf.Term) {
	graph.addInternal(graph.serializeTerm(subj), graph.serializeTerm(pred), graph.serializeTerm(obj))
}

func (graph *Graph) ForEach(handle func(subject string, predicate string, object string)) {
	graph.internalForEach(func(subject, predicate, object string) {
		handle(convertToExternalValue(subject), convertToExternalValue(predicate), convertToExternalValue(object))
	})
}

func convertToExternalValue(term string) string {
	if strings.HasPrefix(term, "<") {
		return strings.Trim(term, "<>")
	}
	iri, err := rdf.NewIRI(term)
	if err == nil {
		return iri.String()
	}
	lit, err := rdf.NewLiteral(term)
	if err == nil {
		return lit.String()
	}

	return term
}

func (graph *Graph) internalForEach(handle func(subject string, predicate string, object string)) {
	for subject, predtoobj := range graph.data {
		for predicate, objects := range predtoobj {
			for _, object := range objects {
				handle(subject, predicate, object)
			}
		}
	}
}

func (graph *Graph) Merge(other *Graph) {
	transformedBlanks := make(map[string]rdf.Term)
	replaceBlank := func(str string, term rdf.Term) rdf.Term {
		if _, isblank := term.(*rdf.Blank); isblank {
			if transformed, has := transformedBlanks[str]; has {
				return transformed
			} else {
				transformed = graph.NewBlank()
				transformedBlanks[str] = transformed
				return transformed
			}
		}
		return term
	}
	other.internalForEach(func(subject, predicate, object string) {
		sterm, ok := graph.nodes[subject]
		if !ok {
			return
		}
		pterm, ok := graph.nodes[predicate]
		if !ok {
			return
		}
		oterm, ok := graph.nodes[object]
		if !ok {
			return
		}
		sterm = replaceBlank(subject, sterm)
		pterm = replaceBlank(predicate, pterm)
		oterm = replaceBlank(object, oterm)
		graph.AddTriple(sterm, pterm, oterm)
	})
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
