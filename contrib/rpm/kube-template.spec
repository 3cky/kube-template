%define debug_package %{nil}

Name:           kube-template
Version:        0.1
Release:        1%{?dist}
Summary:        Watches Kubernetes for updates, writing output of a series of templates to files

Group:          System Environment/Daemons
License:        Apache 2.0
URL:            https://github.com/3cky/kube-template
Source0:        https://github.com/3cky/%{name}/archive/v%{version}.tar.gz#/%{name}-%{version}.tar.gz
Source1:        %{name}.service
Source2:        %{name}.config

%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
BuildRequires:  systemd-units
Requires:       systemd
%endif
BuildRequires:  golang >= 1.9

%description
Inspired by HashiCorp's Consul Template, kube-template is a utility that queries a Kubernetes API server
and uses its objects as input data for specified set of templates. Also it can run arbitrary commands for
templates with updated output.

%prep
%setup -c -n %{name}/src/github.com/3cky
cd %{_builddir}/%{name}/src/github.com/3cky
%{__mv} %{name}-%{version} %{name}

%build
export GOPATH=%{_builddir}/%{name}
export PATH=$PATH:"%{_builddir}"/%{name}/bin
cd "$GOPATH/src/github.com/3cky/%{name}"
go get -u github.com/golang/dep/cmd/dep
%make_build

%install
mkdir -p %{buildroot}/%{_bindir}
cp %{_builddir}/%{name}/bin/%{name} %{buildroot}/%{_bindir}
mkdir -p %{buildroot}/%{_sysconfdir}/%{name}

%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
mkdir -p %{buildroot}/%{_unitdir}
cp %{SOURCE1} %{buildroot}/%{_unitdir}/
cp %{SOURCE2} %{buildroot}/%{_sysconfdir}/%{name}/config
%endif

%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
%post
%systemd_post %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun_with_restart %{name}.service
%endif

%clean
rm -rf %{buildroot}


%files
%defattr(-,root,root,-)
%dir %attr(750, root, root) %{_sysconfdir}/%{name}
%if 0%{?fedora} >= 14 || 0%{?rhel} >= 7
%{_sysconfdir}/%{name}/config
%{_unitdir}/%{name}.service
%endif
%attr(755, root, root) %{_bindir}/%{name}

%doc


%changelog
* Wed Dec 19 2018 Victor Antonóvich <v.antonovich@gmail.com> - 0.1-1
- Initial release
