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

/**
 * Script for migration of sequelize-models data from one DB to another.
 *
 * Two databases are involved in this script, denoted as the source
 * database, and the target database.
 * After a successful run of this script, both the source and target
 * data should exist on the target database.
 *
 * One DB connection is established automatically using environment
 * variables. The other DB connection must be specified by the user.
 *
 *
 * Script arguments:
 * --username:          DB username
 * --password:          DB password
 * --host:              DB host
 * --port:              DB port
 * --dialect:           DB SQL dialect
 * --export:            export from default sequelize DB
 * --confirm:           skip final confirmation to run migration
 *
 * Example Usage:
 *  $ node -r ../../babelRegister.js main.js
 *  ? Enter DB host: mariadb
 *  ? Enter DB port: 3306
 *  ? Enter DB database name: nms
 *  ? Enter DB username: root
 *  ? Enter DB password: [hidden]
 *  ? Enter DB SQL dialect: mariadb
 *
 *  DB Connection Config:
 *  ---------------------------
 *  Host: mariadb:3306
 *  Database: nms
 *  Username: root
 *  Dialect: mariadb
 *
 *  ? Are you importing from the specified DB, or exporting to it?: import
 *  ? Would you like to run data migration with these settings?: Yes
 *  Completed data migration, importing from specified DB
 */

/* eslint no-console: "off" */

const inquirer = require('inquirer');
const process = require('process');
const argv = require('minimist')(process.argv.slice(2));
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
const {exportToDatabase, importFromDatabase} = require('./index.ts');
import type {Options} from 'sequelize';

const dbQuestions = [
  {
    type: 'input',
    name: 'host',
    message: 'Enter DB host:',
    default: 'mariadb',
  },
  {
    type: 'input',
    name: 'port',
    message: 'Enter DB port:',
    default: 3306,
  },
  {
    type: 'input',
    name: 'database',
    message: 'Enter DB database name:',
    default: 'nms',
  },
  {
    type: 'input',
    name: 'username',
    message: 'Enter DB username:',
    default: 'root',
  },
  {
    type: 'password',
    name: 'password',
    message: 'Enter DB password:',
  },
  {
    type: 'input',
    name: 'dialect',
    message: 'Enter DB SQL dialect:',
    default: 'mariadb',
  },
];

async function getDbOptions(): Promise<Options> {
  let dbOptions: Options = {};

  if (
    argv['username'] &&
    argv['password'] &&
    argv['database'] &&
    argv['host'] &&
    argv['port'] &&
    argv['dialect']
  ) {
    dbOptions = {
      username: argv['username'],
      password: argv['password'],
      database: argv['database'],
      host: argv['host'],
      port: parseInt(argv['port']),
      dialect: argv['dialect'],
      logging: (msg: string) => console.log(msg),
    };
    console.log(argv);
  } else {
    await inquirer.prompt(dbQuestions).then(answers => {
      dbOptions = {
        username: answers['username'],
        password: answers['password'],
        database: answers['database'],
        host: answers['host'],
        port: parseInt(answers['port']),
        dialect: answers['dialect'],
        logging: (msg: string) => console.log(msg),
      };
    });
  }
  return dbOptions;
}

function displayDbOptions(dbOptions: Options) {
  const notice =
    `\n` +
    `DB Connection Config:\n` +
    `---------------------------\n` +
    `Host: ${dbOptions.host || ''}:${dbOptions.port || 0} \n` +
    `Database: ${dbOptions.database || ''} \n` +
    `Username: ${dbOptions.username || ''} \n` +
    `Dialect: ${dbOptions.dialect || ''} \n`;

  console.log(notice);
}

async function confirmAndRunMigration(dbOptions: Options): Promise<void> {
  if (argv['confirm']) {
    if (argv['export']) {
      await runMigration(false, dbOptions);
      return;
    }
    await runMigration(true, dbOptions);
    return;
  }

  await inquirer
    .prompt([
      {
        type: 'rawlist',
        name: 'runType',
        message:
          'Are you importing from the specified DB, or exporting to it?:',
        choices: ['import', 'export'],
      },
      {
        type: 'confirm',
        name: 'willRun',
        message: 'Would you like to run data migration with these settings?:',
      },
    ])
    .then(answers => {
      if (answers['willRun']) {
        const isImport = answers['runType'] === 'import';
        (async () => {
          await runMigration(isImport, dbOptions);
        })();
        return;
      }
      console.log('Aborting data migration');
    });
}

async function runMigration(
  isImport: boolean,
  dbOptions: Options,
): Promise<void> {
  try {
    if (isImport) {
      await importFromDatabase(dbOptions);
      console.log('Completed data migration, imported from specified DB');
    } else {
      await exportToDatabase(dbOptions);
      console.log('Completed data migration, exported to specified DB');
    }
  } catch (error) {
    console.log(
      `Unable to connect to specified database for migration:\n` +
        `--------------------------------------------------------------------------\n` +
        `${error}\n` +
        `--------------------------------------------------------------------------\n`,
    );
    process.exit(1);
  }
}

function main() {
  (async () => {
    const dbOptions: Options = await getDbOptions();
    displayDbOptions(dbOptions);
    await confirmAndRunMigration(dbOptions);
  })();
}

main();
