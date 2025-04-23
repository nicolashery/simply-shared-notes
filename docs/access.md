# Access

## Access tokens

- A main value proposition for the application is that it does _not_ require registration or user accounts
- Accessing a Space and its resources such as Notes is done in a "semi-private" way thanks to unguessable **Access Links**
- Access Links are meant to be shared to trusted users that need access to the Space
- Links to the Space are of the form `/s/{token}/*`
  - `s` is short for "space"
  - `{token}` is a valid **Access Token** for that Space
- Each Space has three Access Tokens, one for each **Role** granting different permissions
  - **Admin**: Can do anything the Edit and View roles can do, as well as add/remove Members, rename the Space, regenerate the Access Tokens, delete the Space
  - **Edit**: Can do anything the View role can do, as well as create/edit/delete Notes, edit Members
  - **View**: Can view the Space name, view Notes, view Members
- Spaces are created with an **Email** address
  - Allows for recovery of Access Links if they are lost or forgotten
  - Protects against spam or bots as creating a Space sends an email with the Access Links and does _not_ redirect to the created Space

## Members

- In the absence of proper user accounts, Members help track who created or edited items such as Notes in the Space
- A Member is tied to a Space, and cannot be shared across Spaces
- See also [SQL - Deleting members](sql.md#deleting-members)

## Identity

- When accessing a Space through one of the **Admin** or **Edit** Access Links, a user needs to declare their **Identity** by selecting one of the Space Members
  - The **View** Access Link can be considered as "anonymous"
- The whole system is built on trust, as users can identify as any Member they wish
  - But since the Space is meant to be shared with a small private group of trusted people, this is an acceptable trade-off
- The selected Member ID that the user identified as is stored in a **cookie**
  - This allows the user to make changes to the Space (ex: create a Note) without having to identify at each action
  - The cookie name is `simplysharednotes_session`
  - The cookie is encrypted and cannot be tampered with (the only way a user can identity as another member is through the UI)
  - The cookie value contains the selected Member internal ID (among other session values)
  - The cookie expires after 3 months (90 days)
- The Member selection page is at the URI `/identity`
  - All Space pages `/s/{token}/*` redirect to `/identity` if no valid Identity is found in the cookie
  - Note that a stored Identity can become invalid if a Member is deleted from the Space
