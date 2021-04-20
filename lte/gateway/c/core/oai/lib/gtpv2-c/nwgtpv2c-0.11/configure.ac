# this is example-file: configure.ac

# initial information about the project
AC_INIT([NwGtpv2c],[0.11],[])

# check if the source folder is correct
AC_CONFIG_SRCDIR([src/NwGtpv2c.c])

# Checks for programs

# check for C++ preprocessor and compiler and the library compiler
AC_PROG_CXXCPP
AC_PROG_CXX
AC_PROG_RANLIB

# automake initialisation and check for minimal automake API version 1.9
AM_INIT_AUTOMAKE([1.9])

# use the C compiler for the following checks
AC_LANG([C])
AC_C_INLINE

# Checks for header files.
AC_HEADER_STDC
AC_CHECK_HEADERS([fcntl.h])
AC_CHECK_HEADERS([string])
AC_CHECK_HEADERS([iostream])
AC_CHECK_HEADERS([arpa/inet.h])
AC_CHECK_HEADERS([sys/time.h])
AC_CHECK_HEADERS([sys/socket.h])

AC_CHECK_FUNCS([strerror])
AC_CHECK_FUNCS([socket])
AC_CHECK_FUNCS([memset])
AC_CHECK_FUNCS([gettimeofday])
AC_FUNC_MALLOC
AC_PROG_LN_S

AC_CHECK_LIB([event], [event_init])
AM_CONDITIONAL([ENABLE_TESTS_APP], [test x$HAVE_LIBEVENT == xyes])

# Checks for typedefs, structures, and compiler characteristics.
AC_TYPE_SIZE_T

# distribute additional compiler and linker flags
# --> set these variables instead of CXXFLAGS or LDFLAGS
AC_SUBST([AM_CFLAGS])
AC_SUBST([AM_LDFLAGS])

# files to generate via autotools (.am or .in source files)
AC_CONFIG_FILES([Makefile])
AC_CONFIG_FILES([src/Makefile])
AC_CONFIG_FILES([test-app/Makefile])
AC_CONFIG_FILES([test-app/nw-helloworld/Makefile])
AC_CONFIG_FILES([test-app/nw-egtping/Makefile])

# generate the final Makefile etc.
AC_OUTPUT
