// @flow
import request from 'superagent';

// TODO should we be talking directly to graph container ?
//  or platform-server (reusing the session)
const url = process.env.GRAPH_HOST || 'http://graph/query';

export async function groupsForUser(
  tenant: string,
  userEmail: string,
  role: string,
  sessionId: ?string,
): Promise<string[]> {
  // Sample output
  // {
  //   "data": {
  //     "user": {
  //       "groups": [
  //         {
  //           "id": "171798691840",
  //           "name": "group4",
  //           "members": [
  //             {
  //               "groups": [
  //                 {
  //                   "name": "group4"
  //                 }
  //               ]
  //             }
  //           ]
  //         }
  //       ]
  //     }
  //   }
  // }

  // Add group
  // mutation addUserGroup {
  //   addUsersGroup(input: {name:"group4"}) {
  //     id
  //     name
  //     status
  //     members {
  //       authID
  //     }
  //   }
  // }

  // Add user to group
  // mutation editUsersGroup{
  //   editUsersGroup(input: {
  //     id: 171798691840,
  //     description: "test",
  //     members: [167503724544],
  //   }) {
  //     id
  //     members {
  //       id
  //       email
  //     }
  //   }
  // }

  // Query groups
  // query allGroups {
  //   usersGroups {
  //     edges {
  //       node {
  //         name
  //         id
  //       }
  //     }
  //   }
  // }

  if (sessionId) {
    // TODO cache per session
    // TODO which cache to use ? there is some lru-cache in the packages folder
    // if (cache.get(sessionId)) {
    //   return cache.get(sessionId);
    // }
  }

  return request
    .post(url)
    .send({
      query: `query getUsersGroups {
                user(authID: "${userEmail}") {
                  groups {
                    name
                  }
                }
              }`,
    })
    .set('Accept', 'application/json')
    .set('Content-Type', 'application/json')
    .set('x-auth-organization', tenant)
    .set('x-auth-user-email', userEmail)
    .set('x-auth-user-role', role)
    .then(res => {
      // extract just group names
      return (res.body?.data?.user?.groups ?? []).map(group => group.name);
    })
    .catch(err => {
      // FIXME proper logging and error handling
      console.log('Error retrieving user groups from graphQl');
      console.log(err);
    });
}
