all:
	cd ${PKG_BUILD_DIR}; \
	cmake -GNinja -DCMAKE_BUILD_TYPE=release  \
		-DCMAKE_INSTALL_PREFIX=${PKG_INSTALL_DIR} \
		${PKG_REPO_DIR}
	cd ${PKG_BUILD_DIR}; ninja -C .
	cd ${PKG_BUILD_DIR}; ninja -C . install
	ldd ${PKG_INSTALL_DIR}/bin/devmand | tr -s '[:blank:]' '\n' | grep '^/' | \
		xargs -I @ sh -c 'mkdir -p $$(dirname ${PKG_INSTALL_DIR}@); \
		cp @ ${PKG_INSTALL_DIR}@;'
test:
	cd ${PKG_BUILD_DIR}; ctest -T test -V
