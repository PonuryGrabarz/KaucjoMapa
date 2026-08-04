export class Notification {
  id!: string;
  user_id?: string | null;
  name!: string;
  description!: string;
  seen!: boolean;
  created_at!: string;
}

export type GetNotificationsResponse = Notification[];

export type CreateNotificationResponse = Notification;

export type MarkNotificationSeenResponse = Notification;

export type DeleteNotificationResponse = undefined;
