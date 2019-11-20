package event

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_javaTestsAddedInPatchAddOneJavaTest(t *testing.T) {
	p := `
+    @Test
+    public void Foo() {
+    }
`
	assert.Equal(t, 1, javaTestsAddedInPatch(p))
}

func Test_javaTestsAddedInPatchRemoveOneJavaTest(t *testing.T) {
	p := `
-    @Test
-    public void Foo() {
-    }
`
	assert.Equal(t, -1, javaTestsAddedInPatch(p))
}

func Test_javaTestsAddedInPatchNetZeroJavaTest(t *testing.T) {
	p := `
+    @Test
+    public void Foo() {
+    }
-    @Test
-    public void Bar() {
-    }
`
	assert.Equal(t, 0, javaTestsAddedInPatch(p))
}

func Test_goTestsAddedInPatchOneAdd(t *testing.T) {
	p := `
+func Test_Foo(t *testing.T) {
+}
`
	assert.Equal(t, 1, goTestsAddedInPatch(p))
}

func Test_goTestsAddedInPatchRemoveOne(t *testing.T) {
	p := `
-func Test_Foo(t *testing.T) {
-}
`
	assert.Equal(t, -1, goTestsAddedInPatch(p))
}

func Test_goTestsAddedInPatchNetZero(t *testing.T) {
	p := `
+func Test_Foo(t *testing.T) {
+}
-func TestBar(t *testing.T) {
-}
`
	assert.Equal(t, 0, goTestsAddedInPatch(p))
}
