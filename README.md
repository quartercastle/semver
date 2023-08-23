# semver

> [!IMPORTANT]
> This is still under active development and should not be
> considered stable. It's not even semantically versioned yet.

Semver is a project about automating semantic versioning by analysing
structural and behavioural aspects of source code.

### Problems with maintaining semver today
TODO: describe why semver can be hard to maintain today.
- Hard to maintain if a project has many contributers or high commit frequency
- Process is often manuel, relying on reviewers or structured commit messages

### How can we improve
There are two aspects in doing semantic versioning the first is structural
and the second is behavioural.

Structural change can be detected by comparing the abstract syntax three (ast)
of an earlier version with latest version. By doing this it is possible to
detect if the latest version contain structural changes and to categorize
them as either a patch, minor or major change.

Behavioural change can be achieved by running the previous versions test suite
against the latest version. It is important to mention that
the behavioural verification requires a good test suite and coverage.
It is up to the maintainer to ensure that there is enough test
coverage to verify that the behavioural aspects hasn't changed between versions.
And even with the best test suite you wont be able to guarentee that it wont be
breaking be aware of this!

### Install
Install latest version of semver.
```sh
go install github.com/quartercastle/semver/cmd/semver@latest
```

### Usage
Semver can only do structural change detection at the moment. To use it checkout
two versions of a project in different folders and use the semver cli to see
the structural changes between the versions.
```sh
semver --explain examples/v1.0.0 examples/v2.0.0
```

### Next steps
- [ ] Integrate with Git to automatically checkout and cache versions to compare.
- [ ] Extract test cases from previous versions and run them against the latest
      version.

### License
This project is licensed under the [MIT License](LICENSE).

