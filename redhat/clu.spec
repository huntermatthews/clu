Name:           clu
Version:        %{_version}
Release:        1%{?dist}
Summary:        System facts collector and analyzer
License:        GPL-2.0-only
URL:            https://github.com/huntermatthews/clu
Source0:        %{name}-%{version}.tar.gz
BuildArch:      x86_64
BuildRequires:  golang >= 1.19

%description
A command-line tool for collecting and analyzing system facts.

%prep
%setup -q

%build
CGO_ENABLED=0 go build -ldflags "-s -w -X github.com/huntermatthews/clu/pkg.Version=%{_version}" -o clu ./cmd/main.go

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_mandir}/man1
mkdir -p %{buildroot}%{_docdir}/%{name}
install -m 755 clu %{buildroot}%{_bindir}/clu
install -m 644 clu.1 %{buildroot}%{_mandir}/man1/clu.1
install -m 644 LICENSE %{buildroot}%{_docdir}/%{name}/LICENSE
install -m 644 README.md %{buildroot}%{_docdir}/%{name}/README.md

%files
%{_bindir}/clu
%{_mandir}/man1/clu.1*
%{_docdir}/%{name}/LICENSE
%{_docdir}/%{name}/README.md

%changelog
* Tue Dec 30 2025 GitHub Actions <noreply@github.com> - %{_version}-1
- Automated build from release %{_version}
