# semver

> [!IMPORTANT]
> This is still under active development and should not be
> considered stable. It's not even semantically versioned yet.

Semver is a project about automating semantic versioning by analysing
structural  and behavioural aspects of source code.

Versioning is hard and in many cases ignored as it seems like extra work.
This is unfortunate as semantic versioning has many benefits in terms of
dependency management and it is a good way to communicate change to
dependent projects.

Describe problem of semver today:
- Manuel process relying on reviewers
- Error prone
- No tooling
<!--There are other ways to tackle and maintain semantic versioning, most are
based on formatted commit messages, which can be error prone and the
process is dependent on reviewers too verify that changes are correctly
categoriesed. -->

### How can we improve
There are two aspects in doing semantic versioning the first is structural
and the second is behavioural.

Structural change can be detected by comparing the abstract syntax three (ast)
of an earlier version with latest version. By doing this it is possible to
detect if the latest version contain structural changes and to categorize
them as either a patch, minor or major change.

Behavioural change is done by running previous versions test suite
against the latest version. It is important to mention that
the behavioural part of this project is only as good as its unit tests and
coverage. It is up to the maintainer to ensure that there is enough test
coverage to verify the behavioural aspects hasn't changed between versions.

### Install
Install latest version of semver.
```sh
go install github.com/quartercastle/semver/cmd/semver
```

### Usage
Semver can only do structural change detection at the moment. Checkout two
versions of a project and use semver like below see the structural changes
between the versions.
```sh
semver --explain examples/v1.0.0 examples/v2.0.0
```

### Next steps
- [ ] Integrate with Git to automatically checkout versions to compare.
- [ ] Extract test cases from previous versions and run them against the latest
      version.

### License
This project is licensed under the [MIT License](LICENSE).

