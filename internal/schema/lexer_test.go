package schema

import (
	"strings"
	"testing"

	"github.com/ory/x/snapshotx"
)

func TestLexer(t *testing.T) {
	t.Run("suite=snapshots", func(t *testing.T) {
		cases := []struct{ name, input string }{
			{"empty", ""},
			{"single class", `
class name implements Namespace {
	metadata = {
		id: "123"
	}
}
`},
			{"comments", `
/**/

/** doc comment
 * content
 * more content
 */
class name implements Namespace {
	// Block comment
	/*
	äny ünicöde characterß???
	*/
}
`},
			{"two classes", `
class user implements Namespace {
	metadata = {
		id: "1"
	}
}

class document implements Namespace {
	metadata = {
		id: "2"
	}

	related: {
		owners: user[]
		editors: user[]
		viewers: user[]
		parent: document[]
	}
}
`},
			{"full class", `
class File implements Namespace {
	metadata = {
		id: "2"
	}
  
	related: {
	  parents: File[]
	  viewers: User[]
	  owners: User[]
	  siblings: File[]
	}
  
	permits = {
	  view: (ctx: Context): boolean =>
		this.related.parents.some(p => p.permits.view(ctx)) ||
		  this.related.viewers.includes(ctx.subject) ||
		  this.related.owners.includes(ctx.subject),
  
	  edit: (ctx: Context) => this.related.owners.includes(ctx.subject),
  
	  rename: (ctx: Context) => this.related.siblings.some(s => s.permits.edit(ctx))
	}
}
`},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				l := Lex(tc.name, tc.input)
				var items []string
				for {
					item := l.nextItem()
					items = append(items, item.String())

					if item.Typ == itemError {
						t.Fail()
						break
					}
					if item.Typ == itemEOF {
						break
					}
				}
				t.Logf("Tokens:\n%s", strings.Join(items, "\n"))
				snapshotx.SnapshotT(t, items)
			})
		}
	})
}
