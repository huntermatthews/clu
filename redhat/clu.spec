Name:           clu
Version:        %{_version}
Release:        1%{?dist}
Summary:        System facts collector and analyzer
License:        LicenseRef-NonDistributable
URL:            https://github.com/NHGRI/clu

# No Source0 - building directly from copied source tree

# Disable debuginfo package generation
%global debug_package %{nil}
%define __spec_install_post %{nil}

%description
A command-line tool for collecting and analyzing system facts.

%prep
# No prep needed - source already copied to BUILD directory by build script

%build
cd %{_builddir}/clu-%{_version}
make build VERSION=%{_version}

%install
cd %{_builddir}/clu-%{_version}
make install VERSION=%{_version} PREFIX=%{buildroot}%{_prefix}

%files
%{_bindir}/clu
%{_mandir}/man1/clu.1*
%{_docdir}/%{name}/README.md

%changelog
* Tue Dec 30 2025 GitHub Actions <noreply@github.com> - %{_version}-1
- Automated build from release %{_version}

# %%global _source_payload    w9.gzdio
# %%global _binary_payload    w9.gzdio
