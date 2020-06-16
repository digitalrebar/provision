
VERSION=v4.3.2
CLIVER=v4.3.2
VENDOR=rackn
DESCRIPTION=Digital Rebar Provision platform - $(VERSION)
MAINTAINER=support@rackn.com
URL=https://rackn.com/
LICENSE=RackN
RPMSUM=Digital Rebar Provision platform
ZIPFILE=drp-$(VERSION).zip
PKGS=rpm deb snap
# 'pacman' type fails to build 'services' target correctly
# PKGS=rpm deb snap pacman

FPM=--version "$(VERSION)" --vendor "$(VENDOR)" --maintainer "$(MAINTAINER)" --url '$(URL)' --description "$(DESCRIPTION)" --license "$(LICENSE)" --rpm-summary "$(RPMSUM)"

.PHONY: package

package: services dr-provision dr-client stage-outputs

help:
	@echo ""
	@echo "Make main targets:"
	@echo "	help                   # outputs this help statement"
	@echo "	setup-fpm              # attempt to install/setup FPM"
	@echo "	clean                  # remove packages and staged unzips"
	@echo "	clean-all              # really clean everything"
	@echo "	clean-pkgs             # remove staged packages that might exist"
	@echo "	package	               # make dr-provsion, drpcli, and startup services packages"
	@echo ""
	@echo "Note: package builds: dr-provision dr-client services"
	@echo "      and is equivalent to just 'make' with no targets"
	@echo ""
	@echo "Additional sub-targets:"
	@echo "	dr-provision dr-client services"
	@echo ""

clean:
	rm -rf bin/ usr/ *rpm *deb *snap *pkg.tar.xz

clean-pkgs:
	rm -rf pkgs/

clean-all: clean clean-pkgs
	rm -rf drp-$(VERSION).zip drpcli

# CentOS 7.7 specific
setup-fpm:
	sudo yum install centos-release-scl
	sudo yum install rh-ruby23 rh-ruby23-ruby-devel gcc make rpm-build rubygems
	sudo scl enable rh-ruby23 bash
	sudo gem install fpm
	sudo gem install pleaserun

setup:
	test -f drpcli && echo "drpcli found ... skip download" || curl -fsSl https://rebar-catalog.s3-us-west-2.amazonaws.com/drpcli/$(CLIVER)/amd64/linux/drpcli -o drpcli
	test -f ./drpcli && chmod 755 ./drpcli
	test -f $(ZIPFILE) && echo "dr-provision zip found ... skip download" || ./drpcli catalog item download drp as drp-$(VERSION).zip --version $(VERSION)

setup-drp: setup
	bsdtar -xzvf $(ZIPFILE) bin/linux/amd64/dr-provision
	mkdir -p usr/local/bin/
	mv bin/linux/amd64/dr-provision usr/local/bin/

setup-cli: setup
	bsdtar -xzvf $(ZIPFILE) bin/linux/amd64/drpjoin bin/linux/amd64/drpcli bin/linux/amd64/drbundler
	mkdir -p usr/local/bin/
	mv bin/linux/amd64/drpjoin bin/linux/amd64/drpcli bin/linux/amd64/drbundler usr/local/bin/

dr-provision: setup-drp build-dr-provision

build-dr-provision:
	for PKG in $(PKGS) ; do \
		fpm --input-type dir --output-type $$PKG --name dr-provision --provides dr-provision $(FPM) usr/local/bin/; \
	done
	rm -rf usr/

dr-client: setup-cli build-dr-client

build-dr-client:
	for PKG in $(PKGS) ; do \
		fpm --input-type dir --output-type $$PKG --name drpcli --provides drpcli $(FPM) usr/local/bin/; \
	done
	rm -rf usr/

services: build-services

build-services:
	for PKG in $(PKGS) ; do \
		fpm --input-type pleaserun --output-type $$PKG --name dr-provision-services --provides dr-provision-services $(FPM) /usr/local/bin/dr-provision; \
	done

stage-outputs:
	for PKG in $(PKGS) ; do                                              \
    		mkdir -p pkgs/$$PKG;                                         \
    		DIR=$$PKG;                                                   \
		[[ $$PKG == "pacman" ]] && PKG="pkg.tar.xz";                 \
		ls -1 *\.$$PKG > /dev/null 2>&1 && mv *\.$$PKG pkgs/$$DIR/;  \
	done


