# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""
Generates a module structure for swagger specifications for python.
The python code expects swagger specifications (.yml) in a module
orc8r.swagger.specs or lte.swagger.specs respectively.
Please note this cannot be achieved with individual rules for each
.yml file. This is due to the python restriction that resource files
require an "__init__.py" and thus all resource files need to be
generated in the same folder.
Due to the specific folder structure required, this implementation
can only be used in "lte/swagger" or "orc8r/swagger".
"""

load("@rules_python//python:defs.bzl", "py_library")

def py_swagger_specs(name, component):
    """Generates swagger specifications.

    Args:
        name: Irrelevant mandatory argument.
        component: The component associated to the specifications.
            This is currently either "lte" or "orc8r".
    """

    SPECS_IMPORT_ROOT = "specs_root"
    SPECS_MODULE_PATH = "{0}/{1}/swagger/specs/".format(SPECS_IMPORT_ROOT, component)
    CREATE_MODULE_NAME = "create_module_{0}_swagger_specs".format(component)
    CREATE_MODULE_NAME_REF = ":{0}".format(CREATE_MODULE_NAME)

    # Wrapper around genrule below in order to add the module root to the python path (see "imports").
    py_library(
        name = "{0}_swagger_specs".format(component),
        srcs = [CREATE_MODULE_NAME_REF],
        data = [CREATE_MODULE_NAME_REF],
        imports = [SPECS_IMPORT_ROOT],
        visibility = ["//visibility:public"],
    )

    # Creates a folder "specs_root/$component/swagger/specs/" and copies all
    # ".yml" files from "$component/swagger" into this folder
    native.genrule(
        name = CREATE_MODULE_NAME,
        srcs = native.glob(["*.yml"]),
        outs = ["{0}__init__.py".format(SPECS_MODULE_PATH)] + [SPECS_MODULE_PATH + spec for spec in native.glob(["*.yml"])],
        cmd = " && ".join([
            "mkdir -p $(RULEDIR)/{0}".format(SPECS_MODULE_PATH),
            "touch $(RULEDIR)/{0}__init__.py".format(SPECS_MODULE_PATH),
            "cp {0}/swagger/*.yml $(RULEDIR)/{1}".format(component, SPECS_MODULE_PATH),
        ]),
    )
