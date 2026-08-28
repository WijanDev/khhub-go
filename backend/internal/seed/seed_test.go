package seed

import (
	"strings"
	"testing"
)

func TestResetTruncateListLeavesAuthTables(t *testing.T) {
	t.Parallel()
	if len(demoDataTables) == 0 {
		t.Fatal("demoDataTables must list the congregation tables Reset wipes")
	}
	protected := []string{"users", "sessions"}
	seen := map[string]struct{}{}
	for _, table := range demoDataTables {
		if _, dup := seen[table]; dup {
			t.Fatalf("duplicate truncate table %q", table)
		}
		seen[table] = struct{}{}
		for _, name := range protected {
			if table == name {
				t.Fatalf("Reset must not truncate %s", name)
			}
		}
	}
	sql := truncateDemoSQL()
	for _, name := range protected {
		if strings.Contains(sql, name) {
			t.Fatalf("truncate SQL must not mention %s: %s", name, sql)
		}
	}
}
