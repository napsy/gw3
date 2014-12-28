package gitw3

import (
	"errors"
	"strings"
	"testing"
)

var showMsg = `commit dafbfbce43a89da9ddb96610e28bad45a770d1e1
Merge: a57a831 d32f052
Author: John Doe <johndoe@example.com>
Date:   Mon Aug 11 12:48:19 2014 +0200

    Merge commit 'd32f05214e94813e831d5587b970695aa9500202' into go1.3

diff --cc src/gojs/context.go
index fcdc6f8,0000000..c25db46
mode 100644,000000..100644
--- a/src/gojs/context.go
+++ b/src/gojs/context.go
@@@ -1,37 -1,0 +1,37 @@@
 +package gojs
 +
 +// #include <stdlib.h>
 +// #include <JavaScriptCore/JSContextRef.h>
`

var logMsg = `f590daf,Mon Oct 27 19:37:13 2014 +0100,Elisa Example,Merge remote-tracking branch 'origin/develop' into develop
bea9c3f,Mon Oct 27 16:57:15 2014 +0100,John Doe,Merge remote-tracking branch 'origin/develop' into develop
df8cdb4,Mon Oct 27 16:56:01 2014 +0100,John Doe,start work on inspector proxy
`

func fakeExec(rootDir string, args ...string) ([]string, error) {
	lines := []string{}
	switch args[0] {
	case "show":
		lines = strings.Split(showMsg, "\n")
	case "log":
		lines = strings.Split(logMsg, "\n")
	default:
		return nil, errors.New("unknown git command")
	}
	return lines, nil
}
func TestLog(t *testing.T) {
	git := GitRoot{exec: fakeExec}
	logs, err := git.Log(3)
	if err != nil {
		t.Fatal(err)
	}
	if logs[2].Message != "start work on inspector proxy" {
		t.Fatalf("invalid message, got '%s'", logs[2].Message)
	}
}
