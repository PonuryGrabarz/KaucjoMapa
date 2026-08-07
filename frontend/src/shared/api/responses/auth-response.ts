export class AuthTokens {
  accessToken!: string;
  refreshToken!: string;
  refreshTokenExpiresAt?: string;
}

export class AccessToken {
  accessToken!: string;
}

export type LoginResponse = AuthTokens;

export type RegisterResponse = AuthTokens;

export type RefreshTokenResponse = AccessToken;

export type RevokeTokenResponse = undefined;
