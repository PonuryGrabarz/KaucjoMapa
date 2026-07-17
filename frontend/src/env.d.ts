/// <reference types="astro/client" />

declare namespace App {
  interface SessionUser {
    authenticatedAt: Date;
    createdAt: string;
    email: string;
    id: string;
    isAdmin: boolean;
    isBanned: boolean;
    reputation: number;
    updatedAt: string;
    username: string;
  }

  interface SessionData {
    token: string;
    user: SessionUser;
  }
}
