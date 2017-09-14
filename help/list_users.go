package help

func ListUsers() string {
	return `SCIM is the System for Cross-Domain Identity Management. The standard describes
a format, or schema, for representing identity information and specifies the API for
accessing that data.

Of particular interest for this command are the SCIM filters. SCIM filters
provide a limited query interface. See examples below for usage.

Examples:

  - Show all usernames containing gmail.com:
    uaa list-users --filter 'userName co "gmail.com"' --attributes id,emails

  - Show all users from a particular origin (identity provider):
    uaa list-users --filter 'userName eq "bob@example.com" and origin eq "ldap"'

  - Find all unverified users:
    uaa list-users --filter 'verified eq false' --attributes id,userName,name,emails

  - Find users whose username starts with "z":
    uaa list-users --filter 'userName sw "z"'

  - See client approvals for a specific user:
    uaa list-users --filter 'userName eq "bob@example.com' --attributes approvals

  - See everything about a specific user, including group memberships:
    uaa list-users --filter 'userName eq "bob@example.com'

Keep in mind that the UAA must perform SQL joins to determine group membership so
responses will be relatively slow when fetching results for large numbers of users
without filtering to specific attributes of interest with the --attributes flag.
  `
}
