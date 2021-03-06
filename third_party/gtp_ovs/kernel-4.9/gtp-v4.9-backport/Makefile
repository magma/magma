PKGNAME=oai-gtp
VERSION=4.9
ITERATION=9

PKGFMT=deb
WORK_DIR=/tmp/build-${PKGNAME}
PKGFILE=${PKGNAME}_${VERSION}-${ITERATION}_${ARCH}.${PKGFMT}
DEPMOD_CONFIG_DIR=/etc/depmod.d/
DEPMOD_CONFIG_FILE=${DEPMOD_CONFIG_DIR}/gtp.conf

OUTPUT_DIR=`pwd`
OUTPUT_PATH=${OUTPUT_DIR}/${PKGFILE}

obj-m += gtp.o

build:
	make -C /lib/modules/$(KVERS)/build M=$(PWD) modules

modules_install: build
	make -C /lib/modules/$(KVERS)/build M=$(PWD) modules_install
	mkdir -p ${DEPMOD_CONFIG_DIR}
	echo "override gtp.ko * extra" >> ${DEPMOD_CONFIG_FILE}
	echo "override gtp.ko * weak-updates" >> ${DEPMOD_CONFIG_FILE}
	depmod -a

package: build 
	rm -rf ${WORK_DIR}
	mkdir ${WORK_DIR}
	make INSTALL_MOD_PATH=${WORK_DIR} -C /lib/modules/$(KVERS)/build M=$(PWD) modules_install
	fpm \
	    -f \
	    -s dir \
	    -t ${PKGFMT} \
	    -a ${ARCH} \
	    -n ${PKGNAME} \
	    -v ${VERSION} \
	    --iteration ${ITERATION} \
	    --provides ${PKGNAME} \
	    --conflicts ${PKGNAME} \
	    --replaces ${PKGNAME} \
	    --package ${OUTPUT_PATH} \
	    --description 'Flow-based GTP kernel module' \
	    --after-install scripts/after-install.sh \
	    -C ${WORK_DIR}

clean:
	make -C /lib/modules/$(KVERS)/build M=$(PWD) clean
