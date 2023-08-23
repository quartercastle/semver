# semver

> [!IMPORTANT]
> This is still under active development and should not be
> considered stable. It's not even semantically versioned yet.

Semver is a project about automating semantic versioning by analysing
structural and behavioural aspects of source code.

### Problems with maintaining semver today
TODO: describe why semver can be hard to maintain today.
- Hard to maintain if a project has many contributors or high commit frequency
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
semver --explain path/to/v1.0.0 path/to/v2.0.0
```

Below is an example of the output produced by comparing [v1.4.0](https://github.com/adrianmo/go-nmea/releases/tag/v1.4.0) to the latest
commit [a60cdb4](https://github.com/adrianmo/go-nmea/commit/a60cdb4c706d731910788de3e609e367e8d78400) of the
[nmea](https://github.com/adrianmo/go-nmea) parser module for Go.

```sh
semver --filter major --explain nmea/v1.4.0 nmea/a60cdb4
```

```txt
MAJOR: value spec has changed signature
nmea/v1.4.0/mtk.go:5:2
- TypeMTK = "PMTK"
nmea/a60cdb4/mtk.go:6:2
+ TypeMTK = "MTK001"

MAJOR: type spec has changed signature
nmea/v1.4.0/dbs.go:10:6
- DBS struct {
	DepthFeet       float64
	DepthMeters     float64
	DepthFathoms    float64
}
nmea/a60cdb4/dbs.go:14:6
+ DBS struct {
	DepthFeet       float64
	DepthFeetUnit   string
	DepthMeters     float64
	DepthMeterUnit  string
	DepthFathoms    float64
	DepthFathomUnit string
}

MAJOR 28.070333ms
```

### Next steps
- [ ] Integrate with Git to automatically checkout and cache versions to compare.
- [ ] Extract test cases from previous versions and run them against the latest
      version.
- [ ] Better diffs
- [ ] Better docs explaining why certain changes are breaking.

### License
This project is licensed under the [MIT License](LICENSE).

