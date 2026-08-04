interface Window {
	initialReturnPoints?: import("@shared/api/responses").ReturnPoint[];
}

declare namespace App {
  interface SessionUser {
    id: string;
    username: string;
    email: string;
    isAdmin: boolean;
    reputation: number;
    avatarUrl?: string;
  }

  interface Session {
    user: SessionUser;
    accessToken: string;
  }

	interface Locals {
		session: Session | null;
	}
}
