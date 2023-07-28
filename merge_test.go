package sparqlupdate_test

import (
	"testing"

	"github.com/balazsgrill/sparqlupdate"
	"github.com/knakk/rdf"
)

func TestDisjunctMerge(t *testing.T) {
	g1 := sparqlupdate.New()
	g2 := sparqlupdate.New()

	s1, _ := rdf.NewIRI("http://s/1")
	p1, _ := rdf.NewIRI("http://p/1")
	s2, _ := rdf.NewIRI("http://s/2")
	p2, _ := rdf.NewIRI("http://p/2")
	l1, _ := rdf.NewLiteral("v1")
	l2, _ := rdf.NewLiteral("v2")

	g1.AddTriple(s1, p1, l1)
	g2.AddTriple(s2, p2, l2)

	g1.Merge(g2)
	if g1.Length() != 2 {
		t.Fail()
	}
}
