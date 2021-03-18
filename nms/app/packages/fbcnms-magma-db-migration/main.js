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

const inquirer = require('inquirer');
const {importFromDatabase} = require('@fbcnms/sequelize-models');
import {User} from '@fbcnms/sequelize-models';
import type {Options} from 'sequelize';

const dbQuestions = [
  {
    type: 'input',
    name: 'host',
    message: 'Enter MariaDB host:',
    default: 'mariadb'
  },
  {
    type: 'input',
    name: 'port',
    message: 'Enter MariaDB port:',
    default: 3306
  },
  {
    type: 'input',
    name: 'database',
    message: 'Enter MariaDB NMS database name:',
    default: 'nms'
  },
  {
    type: 'input',
    name: 'username',
    message: 'Enter MariaDB username:',
    default: 'root'
  },
  {
    type: 'input',
    name: 'password',
    message: 'Enter MariaDB password:'
  },
  {
    type: 'input',
    name: 'dialect',
    message: 'Enter MariaDB SQL dialect:',
    default: 'mariadb'
  },
]

async function isMigrationNeeded(): boolean {
  const allUsers = await User.findAll();
  if (allUsers.length > 0) {
    return false;
  }
  return true;
}

function main() {
  (async () => {
    try {
      const willRunMigration = await blah();
			if (!willRunMigration) {
    		console.log('Users found in NMS DB. Skipping migration from source DB');
				return;
			}

		var dbOptions: Options;

    inquirer.prompt(dbQuestions).then(answers => {
      dbOptions = {
        username: answers['username'],
        password: answers['password'],
        database: answers['database'],
        host: answers['host'],
        port: parseInt(answers['port']),
        dialect: answers['dialect'],
        logging: (msg: string) => console.log(msg),
      };

			const notice = `\n` +
                     `MariaDB Connection Config:\n` + 
                     `---------------------------\n` + 
                     `Host: ${answers['host']}:${answers['port']} \n` +
                     `Database: ${answers['database']} \n` +
                     `Username: ${answers['username']} \n` +
                     `Dialect: ${answers['dialect']} \n`
			console.log(notice);

      inquirer.prompt([
        {
          type: 'confirm',
          name: 'willRun',
          message: 'Would you like to run data migration with these settings?:'
        },
      ]).then(confirmation => {
		  	if (confirmation['willRun']) {
          try {
            await importFromDatabase(dbOptions);
            console.log('Completed data migration to current NMS DB');
          } catch (error) {
            console.log(
              `Unable to connect to source database with specified options for migration: ${error}`,
            );
          }
		  		return;
		  	}
		  	console.log('Aborting data migration');
      });;
    });
			
    } catch (e) {
      console.log(
        `Unable to run migration: ${e}`,
      );
    }
  })();
}

main()
