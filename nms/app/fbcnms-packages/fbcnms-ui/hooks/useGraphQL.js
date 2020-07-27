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

import {Environment, fetchQuery} from 'relay-runtime';
import {useEffect, useState} from 'react';

export default function (
  env: Environment,
  query: any,
  variables: {[string]: mixed},
) {
  const [error, setError] = useState(null);
  const [response, setResponse] = useState(null);
  const [isLoading, setIsLoading] = useState(false);

  const variablesJSON = JSON.stringify(variables);
  useEffect(() => {
    const variables = JSON.parse(variablesJSON);

    setError(null);
    setIsLoading(true);
    fetchQuery(env, query, variables)
      .then(response => {
        setResponse(response);
        setIsLoading(false);
      })
      .catch(error => {
        setError(error);
        setIsLoading(false);
      });
  }, [env, query, variablesJSON]);

  return {error, response, isLoading};
}
