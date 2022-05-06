Name:           nmpolicy
Version:        v0.2.1
Release:        1.20220506111911002100.main.8.g2ffc142%{?dist}
Summary:        A simple web app

License:        GPLv3
Source0:        nmpolicy-v0.2.1.tar.gz

BuildRequires:  git make golang

Provides:       %{name} = %{version}

%description
An expressions driven declarative API for dynamic network configuration

%global debug_package %{nil}

%prep
%autosetup -n nmpolicy-v0.2.1


%build
make build


%install
install -Dpm 0755 .out/%{name}ctl %{buildroot}%{_bindir}/%{name}ctl

%check
make unit-test
make integration-test

%files
%{_bindir}/%{name}ctl

%changelog
* Fri May 06 2022 Enrique Llorente <ellorent@redhat.com> - v0.2.1-1.20220506111911002100.main.8.g2ffc142
- resolver: Implement walk with stateVisitor (Quique Llorente)
- resolver: Refactor visitor pattern (Quique Llorente)
- build(deps): bump nokogiri from 1.13.0 to 1.13.3 in /docs (dependabot[bot])
- actions: publish-docs use go 1.16 (Quique Llorente)
- docs: Add new CLI section (Quique Llorente)
- test: Use nmpolicyctl at integration test (Quique Llorente)
- cli: Implement "gen" subcommand (Quique Llorente)
- types: Remove main MetaInfo (Quique Llorente)

* Fri May 06 2022 Enrique Llorente <ellorent@redhat.com> - v0.2.1-1.20220506111647527337.main.8.g2ffc142
- resolver: Implement walk with stateVisitor (Quique Llorente)
- resolver: Refactor visitor pattern (Quique Llorente)
- build(deps): bump nokogiri from 1.13.0 to 1.13.3 in /docs (dependabot[bot])
- actions: publish-docs use go 1.16 (Quique Llorente)
- docs: Add new CLI section (Quique Llorente)
- test: Use nmpolicyctl at integration test (Quique Llorente)
- cli: Implement "gen" subcommand (Quique Llorente)
- types: Remove main MetaInfo (Quique Llorente)

* Fri May 06 2022 Enrique Llorente <ellorent@redhat.com> - v0.2.1-1.20220506111623527932.main.8.g2ffc142
- resolver: Implement walk with stateVisitor (Quique Llorente)
- resolver: Refactor visitor pattern (Quique Llorente)
- build(deps): bump nokogiri from 1.13.0 to 1.13.3 in /docs (dependabot[bot])
- actions: publish-docs use go 1.16 (Quique Llorente)
- docs: Add new CLI section (Quique Llorente)
- test: Use nmpolicyctl at integration test (Quique Llorente)
- cli: Implement "gen" subcommand (Quique Llorente)
- types: Remove main MetaInfo (Quique Llorente)

* Fri May 06 2022 Enrique Llorente <ellorent@redhat.com> - v0.2.1-1.20220506110922624826.main.8.g2ffc142
- resolver: Implement walk with stateVisitor (Quique Llorente)
- resolver: Refactor visitor pattern (Quique Llorente)
- build(deps): bump nokogiri from 1.13.0 to 1.13.3 in /docs (dependabot[bot])
- actions: publish-docs use go 1.16 (Quique Llorente)
- docs: Add new CLI section (Quique Llorente)
- test: Use nmpolicyctl at integration test (Quique Llorente)
- cli: Implement "gen" subcommand (Quique Llorente)
- types: Remove main MetaInfo (Quique Llorente)

* Fri May 06 2022 Enrique Llorente <ellorent@redhat.com> - v0.2.1-1.20220506110837141037.main.8.g2ffc142
- resolver: Implement walk with stateVisitor (Quique Llorente)
- resolver: Refactor visitor pattern (Quique Llorente)
- build(deps): bump nokogiri from 1.13.0 to 1.13.3 in /docs (dependabot[bot])
- actions: publish-docs use go 1.16 (Quique Llorente)
- docs: Add new CLI section (Quique Llorente)
- test: Use nmpolicyctl at integration test (Quique Llorente)
- cli: Implement "gen" subcommand (Quique Llorente)
- types: Remove main MetaInfo (Quique Llorente)

* Fri May 06 2022 Enrique Llorente <ellorent@redhat.com> - v0.2.1-1.20220506110822596585.main.8.g2ffc142
- resolver: Implement walk with stateVisitor (Quique Llorente)
- resolver: Refactor visitor pattern (Quique Llorente)
- build(deps): bump nokogiri from 1.13.0 to 1.13.3 in /docs (dependabot[bot])
- actions: publish-docs use go 1.16 (Quique Llorente)
- docs: Add new CLI section (Quique Llorente)
- test: Use nmpolicyctl at integration test (Quique Llorente)
- cli: Implement "gen" subcommand (Quique Llorente)
- types: Remove main MetaInfo (Quique Llorente)

* Fri May 06 2022 Enrique Llorente <ellorent@redhat.com> - v0.2.1-1.20220506105356372864.main.8.g2ffc142
- resolver: Implement walk with stateVisitor (Quique Llorente)
- resolver: Refactor visitor pattern (Quique Llorente)
- build(deps): bump nokogiri from 1.13.0 to 1.13.3 in /docs (dependabot[bot])
- actions: publish-docs use go 1.16 (Quique Llorente)
- docs: Add new CLI section (Quique Llorente)
- test: Use nmpolicyctl at integration test (Quique Llorente)
- cli: Implement "gen" subcommand (Quique Llorente)
- types: Remove main MetaInfo (Quique Llorente)
