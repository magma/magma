---
id: permission-model
title: Permission Model
---

### Entities Definition

* **User**: An identity that is used for authenticating to symphony, contains information about the user
    and also settings that defines his permissions (roles, membership in groups)
* **Group**: Collection of users that are related. Each user can be a member of multiple groups.
* **Policies**(Under active development): Policies are a set of rules to grant access to symphony data
    and operations over it. You can attach policies to a group and the group members "inherit" the group policies

You can manage all these entities inside the [Administrative Tab](/admin/user_management)

### User Roles

Each user has a role that determine access to key parts of Symphony like Settings and User Management
 (it can be further enhanced by policies):
* **User** - Can log in to Symphony desktop and mobile apps. A user has read only access to symphony data
             by default.
* **Admin** - Same as user but can also update settings and manage users and permissions
* **Owner** - Full access over everything, including inventory and workforce data
* **Deactivate** - In Symphony we never delete users since we wish to always be able to manage users history in the system.
     Deactivated user does not have any access or permissions in Symphony

### Moving to new user management
* All existing users are already defined and shown in the "Users" tab
* Users that were defined as "superusers" are now set with the "Owner" role
* Users that were defined as "user" or "read only" are now set with the "User" role
* Pay attention that users with the role "User" have no write permissions (see next section)

### "Write Permission" Group

**Disclaimer**: This a temporary solution until policies are available

If you want to grant a user with write permissions to symphony (create\update\delete)
you need to add it as member to the default "Write Permission" group.

Owners have write permissions even without being members of this group.