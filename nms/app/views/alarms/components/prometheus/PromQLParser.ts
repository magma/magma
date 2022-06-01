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
 */

// @ts-ignore generated
import grammar from './__generated__/PromQLGrammar.js';
import nearley from 'nearley';
import {BinaryOperation} from './PromQL';
import {SyntaxError} from './PromQLTypes';

export function Parser() {
  return new nearley.Parser(nearley.Grammar.fromCompiled(grammar));
}

export function Parse(input: string | undefined | null): BinaryOperation {
  if (!input) {
    throw 'empty input to parser';
  }
  const parser = Parser().feed(input);
  // parser returns array of all possible parsing trees, so access the first
  // element of results since this grammar should only produce 1 for each
  // input
  const ast = parser.results[0] as BinaryOperation;
  if (ast === undefined) {
    throw new SyntaxError('Malformed PromQL expression');
  }
  return ast;
}
