Name:           clu
Version:        %{_version}
Release:        1%{?dist}
Summary:        System facts collector and analyzer
License:        LicenseRef-NonDistributable
URL:            https://github.com/NHGRI/clu
Source0:        %{name}-%{version}.tar.gz
# BuildArch removed - target architecture specified via rpmbuild --target flag
# BuildRequires:  golang >= 1.20  # Disabled - Go provided by GitHub Actions setup-go

# Disable debuginfo package generation
%global debug_package %{nil}
%define __spec_install_post %{nil}

%description
A command-line tool for collecting and analyzing system facts.

%prep
%setup -q

%build
# Use go from PATH or fallback to Ubuntu 20.04 package location
if command -v go >/dev/null 2>&1; then
    GOCMD=go
elif [ -x /usr/lib/go-1.20/bin/go ]; then
    GOCMD=/usr/lib/go-1.20/bin/go
else
    echo "Error: Go compiler not found"
    exit 1
fi
CGO_ENABLED=0 $GOCMD build -ldflags "-X github.com/NHGRI/clu/pkg/global.Version=%{_version}" -o clu ./cmd/main.go

%install
mkdir -p %{buildroot}%{_bindir}
mkdir -p %{buildroot}%{_mandir}/man1
mkdir -p %{buildroot}%{_docdir}/%{name}
install -m 755 clu %{buildroot}%{_bindir}/clu
install -m 644 clu.1 %{buildroot}%{_mandir}/man1/clu.1
install -m 644 README.md %{buildroot}%{_docdir}/%{name}/README.md

%files
%{_bindir}/clu
%{_mandir}/man1/clu.1*
%{_docdir}/%{name}/README.md

%changelog
* Tue Dec 30 2025 GitHub Actions <noreply@github.com> - %{_version}-1
- Automated build from release %{_version}
