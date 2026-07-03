S001 checks that the AGENTS.md file exists. There is no `invalid.md` fixture
here — the failing case is the *absence* of a file at the expected path,
which can't be represented as a checked-in fixture. The rule's tests cover
it by pointing `CheckFileExists` at a path that is guaranteed not to exist
(e.g. inside a fresh `t.TempDir()`).
