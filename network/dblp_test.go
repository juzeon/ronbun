package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeDBLPLink(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"https://dblp.org/db/conf/sigmod/sigmod2021.html", "https://dblp.org/db/conf/sigmod/sigmod2021.html"},
		{"https://dblp.uni-trier.de/db/conf/sigmod/sigmod2021.html", "https://dblp.org/db/conf/sigmod/sigmod2021.html"},
		{"https://dblp2.uni-trier.de/db/conf/sigmod/sigmod2021.html", "https://dblp.org/db/conf/sigmod/sigmod2021.html"},
		{"https://dblp.dagstuhl.de/db/conf/sigmod/sigmod2021.html", "https://dblp.org/db/conf/sigmod/sigmod2021.html"},
	}

	for _, test := range tests {
		result := NormalizeDBLPLink(test.input)
		assert.Equal(t, test.expected, result)
	}
}
