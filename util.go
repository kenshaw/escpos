package escpos

import (
	"errors"

	"github.com/moovweb/gokogiri/xml"
	"github.com/moovweb/gokogiri/xpath"
)

var (
	// ErrBodyElementEmpty is the body element empty error.
	ErrBodyElementEmpty = errors.New("Body element empty")

	// bodyPath is the xpath selector for the
	bodyPath = xpath.Compile("*[local-name()='Body']")
)

// getBodyChildren returns the child nodes contained in the Body element in a XML document.
func getBodyChildren(doc *xml.XmlDocument) ([]xml.Node, error) {
	// grab nodes
	nodes, err := doc.Root().Search(bodyPath)
	if err != nil {
		return nil, err
	}

	// check that the data is present
	if len(nodes) < 1 || nodes[0].CountChildren() < 1 {
		return nil, ErrBodyElementEmpty
	}

	// get body children
	return nodes[0].FirstChild().Search("./*")
}
