SUBDIRS = src

.PHONY: cover release archive install-services uninstall-services enable-services disable-services

cover:
	@for path in `find . -name cover.out` ; do \
		echo $$path; \
		go tool cover -func=$$path ; \
	done

clean-local:
	-rm -fr ./bin ./pkg ./dist ./build
	-rm -fr ./src/beluganos_fibc.egg-info/
	-find . -name "*~" | xargs rm -f
	-find . -name "*.pyc" | xargs rm -f
	-find . -name cover.out  | xargs rm -f
	-find . -name cover.log  | xargs rm -f

distclean-local:
	-find . -name "Makefile.in" | xargs rm -f
	-rm -f aclocal.m4 config.log config.status
	-rm -fr ./autom4te.cache

release: install
	install -v -C `pwd`/bin/* ${GOPREFIX}/bin/

install-lxcinit:
	cp -v -r `pwd`/etc/lxcinit /etc/

uninstall-lxcinit:
	@echo "*** Remove /etc/lxcinit manually on your needs. ***"

NCMS = ncmd ncmi ncms
CFGS = cfgc cfgd cfgcp cfgbgp cfgbgpc cfgfrr cfglxd cfgnet cfgsysc cfgsysctl cfgvtyc netplan+ lxcinit.sh
install-service: install-lxcinit
	@for ncmname in $(NCMS) ; do \
		install -v -C ${GOPREFIX}/bin/$$ncmname /usr/bin/ ; \
		install -v -C `pwd`/etc/systemd/$$ncmname.service /etc/systemd/system/ ; \
	done

	install -v -C `pwd`/etc/systemd/ncm.target /etc/systemd/system/

	@for cfgname in $(CFGS) ; do \
		install -v -C ${GOPREFIX}/bin/$$cfgname /usr/bin/ ; \
	done

	install -v -C `pwd`/etc/config/ncm* /etc/beluganos/

	systemctl daemon-reload

uninstall-service: uninstall-lxcinit disable-services
	@for ncmname in $(NCMS) ; do \
		rm -v -f /usr/bin/$$ncmname ; \
		rm -v -f /etc/systemd/system/$$ncmname.service; \
	done

	rm -v -f /etc/systemd/system/ncm.target

	@for cfgname in $(CFGS) ; do \
		rm -v -f /usr/bin/$$cfgname ; \
	done

	rm -v -f /etc/beluganos/ncm*

	systemctl daemon-reload

enable-services:
	@for ncmname in $(NCMS) ; do \
		systemctl enable $$ncmname ; \
	done

disable-services:
	@for ncmname in $(NCMS) ; do \
		systemctl disable $$ncmname ; \
	done

archive:
	git archive HEAD --prefix=netconf/ --output=netconf.zip
