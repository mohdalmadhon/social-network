-- Create the new invitations table with its own ID
CREATE TABLE group_invitations_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,

    group_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    inviter_id INTEGER NOT NULL,

    status TEXT NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'accepted', 'declined')),

    FOREIGN KEY (group_id)
        REFERENCES groups(id)
        ON DELETE CASCADE,

    FOREIGN KEY (user_id)
        REFERENCES user(id)
        ON DELETE CASCADE,

    FOREIGN KEY (inviter_id)
        REFERENCES user(id)
        ON DELETE CASCADE
);


-- Copy existing invitations
INSERT INTO group_invitations_new (
    group_id,
    user_id,
    inviter_id,
    status
)
SELECT
    group_id,
    user_id,
    inviter_id,
    status
FROM group_invitations;


-- Old invitation notifications used:
-- related_id = group_id
--
-- Change them to:
-- related_id = invitation_id
UPDATE notifications
SET related_id = (
    SELECT gi.id
    FROM group_invitations_new gi
    WHERE gi.group_id = notifications.related_id
      AND gi.user_id = notifications.user_id
      AND gi.inviter_id = notifications.actor_id
)
WHERE category = 'groups'
  AND type = 'invitation'
  AND EXISTS (
      SELECT 1
      FROM group_invitations_new gi
      WHERE gi.group_id = notifications.related_id
        AND gi.user_id = notifications.user_id
        AND gi.inviter_id = notifications.actor_id
  );


-- Remove old table
DROP TABLE group_invitations;


-- Rename new table
ALTER TABLE group_invitations_new
RENAME TO group_invitations;


-- Only one pending invitation can exist
-- between the same group and user
CREATE UNIQUE INDEX group_invitations_one_pending
ON group_invitations(group_id, user_id)
WHERE status = 'pending';


-- If an invitation gets deleted,
-- also delete its notification
CREATE TRIGGER delete_group_invitation_notifications
AFTER DELETE ON group_invitations
FOR EACH ROW
BEGIN
    DELETE FROM notifications
    WHERE category = 'groups'
      AND type = 'invitation'
      AND related_id = OLD.id;
END;


-- Recreate group deletion trigger because
-- invitation related_id now points to invitation.id
DROP TRIGGER IF EXISTS delete_group_notifications;

CREATE TRIGGER delete_group_notifications
BEFORE DELETE ON groups
FOR EACH ROW
BEGIN

    -- Join request notifications
    DELETE FROM notifications
    WHERE category = 'groups'
      AND type = 'join_request'
      AND related_id IN (
          SELECT id
          FROM group_join_requests
          WHERE group_id = OLD.id
      );


    -- Invitation notifications
    DELETE FROM notifications
    WHERE category = 'groups'
      AND type = 'invitation'
      AND related_id IN (
          SELECT id
          FROM group_invitations
          WHERE group_id = OLD.id
      );


    -- Other group notifications whose related_id is still group ID
    DELETE FROM notifications
    WHERE category = 'groups'
      AND type NOT IN ('join_request', 'invitation')
      AND related_id = OLD.id;

END;