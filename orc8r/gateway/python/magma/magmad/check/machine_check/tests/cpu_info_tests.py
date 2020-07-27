"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import textwrap
import unittest

from magma.magmad.check.machine_check import cpu_info


class CpuInfoParseTests(unittest.TestCase):
    def test_parse_bad_output(self):
        expected = cpu_info.LscpuCommandResult(
            error='err',
            core_count=None,
            threads_per_core=None,
            architecture=None,
            model_name=None,
        )
        actual = cpu_info.parse_lscpu_output('output', 'err', None)
        self.assertEqual(expected, actual)

        output = 'bad output'.strip().encode('ascii')
        expected = cpu_info.LscpuCommandResult(
            error="Parsing failed: 'NoneType' object has no attribute 'group'\nbad output",
            core_count=None,
            threads_per_core=None,
            architecture=None,
            model_name=None,
        )
        actual = cpu_info.parse_lscpu_output(output, '', None)
        self.assertEqual(expected, actual)

    def test_parse_good_output(self):
        output = textwrap.dedent('''
            Architecture:          x86_64
            CPU op-mode(s):        32-bit, 64-bit
            Byte Order:            Little Endian
            CPU(s):                4
            On-line CPU(s) list:   0-3
            Thread(s) per core:    1
            Core(s) per socket:    4
            Socket(s):             1
            NUMA node(s):          1
            Vendor ID:             GenuineIntel
            CPU family:            6
            Model:                 158
            Model name:            Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz
            Stepping:              10
            CPU MHz:               2903.998
            BogoMIPS:              5807.99
            Hypervisor vendor:     KVM
            Virtualization type:   full
            L1d cache:             32K
            L1i cache:             32K
            L2 cache:              256K
            L3 cache:              12288K
            NUMA node0 CPU(s):     0-3
            Flags:                 fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush mmx fxsr sse sse2 ht syscall nx rdtscp lm constant_tsc rep_good nopl xtopology nonstop_tsc pni pclmulqdq ssse3 cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt aes xsave avx rdrand hypervisor lahf_lm abm 3dnowprefetch invpcid_single kaiser fsgsbase avx2 invpcid rdseed clflushopt
        ''').strip().encode('ascii')

        expected = cpu_info.LscpuCommandResult(
            error=None,
            core_count=4,
            threads_per_core=1,
            architecture='x86_64',
            model_name='Intel(R) Core(TM) i9-8950HK CPU @ 2.90GHz',
        )
        actual = cpu_info.parse_lscpu_output(output, '', None)
        self.assertEqual(expected, actual)


if __name__ == '__main__':
    unittest.main()
