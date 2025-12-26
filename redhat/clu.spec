Name:           clu
Version:        0.1.0
Release:        1%{?dist}
Summary:        Minimal build of the clu CLI

License:        GPL-2.0-only
URL:            https://github.com/huntermatthews/clu
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang

%description
Minimal RPM packaging to build the Go-based `clu` binary and install
its README and man page.

%prep
%autosetup -n %{name}-%{version}

%build
go build -o clu ./cmd

%install
install -D -m 0755 clu %{buildroot}%{_bindir}/clu
install -D -m 0644 clu.1 %{buildroot}%{_mandir}/man1/clu.1

%files
%license LICENSE
%doc README.md
%{_bindir}/clu
%{_mandir}/man1/clu.1*

%changelog
* Fri Dec 26 2025 Maintainer <maintainer@example.com> - 0.1.0-1
- Initial minimal spec: build binary, install README and man page
