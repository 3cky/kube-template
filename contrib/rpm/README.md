# RPM Spec for kube-template

Spec file structure tries to follow the [packaging guidelines](https://docs.fedoraproject.org/en-US/packaging-guidelines/) from Fedora.

Files provided by built package:

* binary: `/usr/bin/kube-template`
* config: `/etc/kube-template/config`
* systemd service: `/usr/lib/systemd/system/kube-template.service`

# Using

Create the RPMs using the technique outlined in the Build section below, or install Pre-build packages from the Fedora Copr.

## Pre-built packages

Pre-built packages for Fedora/CentOS 7 are available from the [avm/kube-template](https://copr.fedorainfracloud.org/coprs/avm/kube-template/) repository on [Fedora Copr](https://copr.fedoraproject.org/coprs/) system.

# Build

Build the RPM as a non-root user from your home directory:

* Check out this repo.
    ```
    cd $HOME
    git clone https://github.com/3cky/kube-template
    ```

* Install `rpmdevtools` and `mock`.
    ```
    sudo yum install rpmdevtools mock
    ```

* Set up your `rpmbuild` directory tree.
    ```
    rpmdev-setuptree
    ```

* Link the spec file and sources.
    ```
    ln -s $HOME/kube-template/contrib/rpm/kube-template.spec $HOME/rpmbuild/SPECS/
    ln -s $HOME/kube-template/contrib/init/systemd/kube-template.service $HOME/rpmbuild/SOURCES/
    ln -s $HOME/kube-template/contrib/init/systemd/config $HOME/rpmbuild/SOURCES/kube-template.config
    ```

* Download remote source files.
    ```
    spectool -g -R $HOME/rpmbuild/SPECS/kube-template.spec
    ```

* Build the RPM/SRPM packages.
    ```
    rpmbuild -ba $HOME/rpmbuild/SPECS/kube-template.spec
    ```

# Run

* Install the RPM.
* Put config file `kube-template.yaml` in `/etc/kube-template/`.
* Start the service and tail the logs `systemctl start kube-template.service` and `journalctl -f`.
* To enable at reboot `systemctl enable kube-template.service`.
