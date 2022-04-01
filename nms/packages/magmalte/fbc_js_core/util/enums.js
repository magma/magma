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
 * @flow
 * @format
 */

// Coerces the type to an enum, or throws
export function assertEnum<T>(
  // We don't use the value, so set to any
  // eslint-disable-next-line flowtype/no-weak-types
  enumObjectOrArray: {[T]: any} | T[],
  maybeEnum: string,
): T {
  const maybe = coerceEnum(enumObjectOrArray, maybeEnum);
  if (maybe != null) {
    return maybe;
  }
  throw new Error('Invalid enum type');
}

// Coerces the type to an enum, or returns null
export function coerceEnum<T>(
  // We don't use the value, so set to any
  // eslint-disable-next-line flowtype/no-weak-types
  enumObjectOrArray: {[T]: any} | T[],
  maybeEnum: string,
): ?T {
  const keys: T[] = Array.isArray(enumObjectOrArray)
    ? enumObjectOrArray
    : Object.keys(enumObjectOrArray);
  // Need to use find because in or contains aren't as typesafe.
  const maybe = keys.find(v => v === maybeEnum);
  if (maybe != null) {
    return maybe;
  }
  // Explicitly return null since otherwise we'd return undefined.
  return null;
}
