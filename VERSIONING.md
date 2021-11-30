# Versioning

The NMPolicy versioning follows the [semVer specification](https://semver.org/#semantic-versioning-specification-semver)
with the following meanings for minor, major and patch upgrade levels:

- major: breaking changes at main branch
- minor: non breaking changes at main branch
- patch: patch/bugfix changes at release branch


## Release bugfix procedure

Usually the bugs are fixed at main and then backported to the different 
release branches if needed to do a patch release.

The following are the steps to start bugfixing a major.minor release.

- implement the bugfix PR at main branch
- merge bugfix PR at main branch
- pull tags: `git fetch origin --tags`
- create release branch: `git checkout v1.2.0 -b release-1.2`
- push the release branch: `git push origin release-1.2`
- checkout pr branch: `git checkout origin/release-1.2 -b my-bugfix`
- cherry-pick bugfix PR from main branch 
- merge backport PR at release branch
- pull release branch again: `git fetch origin`
- tag the patch version: `git tag v1.2.1 origin/release-1.2`
- push the tag: `git push origin --tags`
- A v1.2.1 release will appear at github

Follow up bugfixes:

- implement the bugfix PR at main branch
- merge bugfix PR at main branch
- pull tags: `git fetch origin`
- checkout pr branch: `git checkout origin/release-1.2 -b my-bugfix`
- cherry-pick bugfix PR from main branch 
- merge backport PR at release branch
- pull release branch again: `git fetch origin`
- tag the patch version: `git tag v1.2.2 origin/release-1.2`
- push the tag: `git push origin --tags`
- A v1.2.2 release will appear at github
