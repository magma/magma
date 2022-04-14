/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import {useEffect, useRef} from 'react';

type RefFunctionType = (HTMLElement | null) => mixed;
type RefType = RefFunctionType | {current: HTMLElement | null};
export type CombinedRefs = Array<RefType>;

const useCombinedRefs = (refs: CombinedRefs) => {
  const targetRef = useRef<?HTMLElement>(null);

  useEffect(() => {
    refs.forEach((ref: ?RefType) => {
      if (ref == null) {
        return;
      }

      if (typeof ref === 'function') {
        const refFunction: RefFunctionType = ref;
        refFunction(targetRef.current || null);
      } else {
        ref.current = targetRef.current || null;
      }
    });
  }, [refs]);

  return targetRef;
};

export default useCombinedRefs;
