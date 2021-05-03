# Copyright 2020 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
#
#
# Locate the Google C++ Mocking Framework.
# Inspired by
# https://github.com/Kitware/CMake/blob/master/Modules/FindGTest.cmake
#
# Result variables
#
#   GMOCK_FOUND - Found the Google Testing framework
#   GMOCK_INCLUDE_DIRS - Include directories
#
# Also defines the library variables below as normal
# variables. These contain debug/optimized keywords when
# a debugging library is found.
#
#   GMOCK_BOTH_LIBRARIES - Both libgmock & libgmock-main
#   GMOCK_LIBRARIES - libgmock
#   GMOCK_MAIN_LIBRARIES - libgmock-main
#
# Accepts the following variables as input:
#
#   GMOCK_ROOT - (as a CMake or environment variable)
#                The root directory of the gmock install prefix
#
#
#-----------------------
# Example Usage:
#
#    find_package(GMock REQUIRED)
#    target_include_directories(foo PRIVATE ${GMOCK_INCLUDE_DIRS})
#
#    add_executable(foo foo.cc)
#    target_link_libraries(foo ${GMOCK_BOTH_LIBRARIES})
#


function(_gmock_append_debugs _endvar _library)
  if(${_library} AND ${_library}_DEBUG)
    set(_output optimized ${${_library}} debug ${${_library}_DEBUG})
  else()
    set(_output ${${_library}})
  endif()
  set(${_endvar} ${_output} PARENT_SCOPE)
endfunction()

function(_gmock_find_library _name)
  find_library(${_name}
      NAMES ${ARGN}
      HINTS
      $ENV{GMOCK_ROOT}
      ${GMOCK_ROOT}
      PATH_SUFFIXES ${_gmock_libpath_suffixes}
      )
  mark_as_advanced(${_name})
endfunction()

set(_gmock_libpath_suffixes lib)

find_path(GMOCK_INCLUDE_DIR gmock/gmock.h
    HINTS
    $ENV{GMOCK_ROOT}/include
    ${GMOCK_ROOT}/include
    )
mark_as_advanced(GMOCK_INCLUDE_DIR)

_gmock_find_library(GMOCK_LIBRARY            gmock)
_gmock_find_library(GMOCK_LIBRARY_DEBUG      gmockd)
_gmock_find_library(GMOCK_MAIN_LIBRARY       gmock_main)
_gmock_find_library(GMOCK_MAIN_LIBRARY_DEBUG gmock_maind)

FIND_PACKAGE_HANDLE_STANDARD_ARGS(GMock DEFAULT_MSG GMOCK_LIBRARY GMOCK_INCLUDE_DIR GMOCK_MAIN_LIBRARY)

if(GMOCK_FOUND)
  set(GMOCK_INCLUDE_DIRS ${GMOCK_INCLUDE_DIR})
  _gmock_append_debugs(GMOCK_LIBRARIES      GMOCK_LIBRARY)
  _gmock_append_debugs(GMOCK_MAIN_LIBRARIES GMOCK_MAIN_LIBRARY)
  set(GMOCK_BOTH_LIBRARIES ${GMOCK_LIBRARIES} ${GMOCK_MAIN_LIBRARIES})
endif()