# gostable

gostable is a Golang linter that flags any go driver usage that is outside of the MongoDB Stable API. Detecting "non-stable" usage ranges from easy to difficult. Here are the cases that are handled:

### Function calls
An easy case, this is driven by a [map of maps](https://github.com/fsnow/gostable/blob/0bd607bc7c09485dd59d03e7e50a4a9a00a030c0/common/analyzer.go#L30). Package -> struct -> function name. In the tree descent, this is in the [\*ast.CallExpr case](https://github.com/fsnow/gostable/blob/0bd607bc7c09485dd59d03e7e50a4a9a00a030c0/common/analyzer.go#L134).
### RunCommand
Also handled in the CallExpr case, we flag all calls to RunCommand. The actual command is the first field name of the bson.D passed as the 2nd argument to RunCommand. A future modification is to *not* flag RunCommand were the command is successfully identified and is supported by the Stable API. See [analyzeRunCommand](https://github.com/fsnow/gostable/blob/0bd607bc7c09485dd59d03e7e50a4a9a00a030c0/common/analyzer.go#L266) for our attempt to find the command construction.
### Structs
All references to unsupported struct fields are flagged. This is also [configuration driven](https://github.com/fsnow/gostable/blob/0bd607bc7c09485dd59d03e7e50a4a9a00a030c0/common/analyzer.go#L52). In the tree descent this is the [\*ast.CompositeLit case](https://github.com/fsnow/gostable/blob/0bd607bc7c09485dd59d03e7e50a4a9a00a030c0/common/analyzer.go#L158).
### Aggregation Stages
The unsupported aggregation stages are sufficiently unique that we take a shortcut, flagging string matching the [stage names](https://github.com/fsnow/gostable/blob/0bd607bc7c09485dd59d03e7e50a4a9a00a030c0/common/analyzer.go#L68). The list of stages is here. The [\*ast.BasicLit case](https://github.com/fsnow/gostable/blob/0bd607bc7c09485dd59d03e7e50a4a9a00a030c0/common/analyzer.go#L123) handles this string matching.

# Build
```
./build.sh
```

# Test
```
./test.sh
```

There are two projects under testdata: stable and unstable. The expected output from the linter is in the 2 "golden" files. The test script compares the linter output against these files.
