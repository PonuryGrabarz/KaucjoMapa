import { z } from "zod";

const CLAIMS_SCHEMA = z
  .object({
    user_id: z.uuid(),
    username: z.string().min(1),
    email: z.string(),
    is_admin: z.boolean(),
    reputation: z.number(),
    exp: z.number().optional(),
  })
  .loose();

export interface AccessTokenClaims {
  userId: string;
  username: string;
  email: string;
  isAdmin: boolean;
  reputation: number;
  expiresAtSeconds?: number;
}

const decodeSegment = (segment: string): unknown => {
  const normalized = segment.replaceAll("-", "+").replaceAll("_", "/");
  const padded = normalized.padEnd(Math.ceil(normalized.length / 4) * 4, "=");

  const binary = atob(padded);
  const bytes = Uint8Array.from(binary, (character) => character.charCodeAt(0));

  return JSON.parse(new TextDecoder().decode(bytes));
};

export const decodeAccessToken = (token: string): AccessTokenClaims | undefined => {
  const segments = token.split(".");

  if (segments.length !== 3 || !segments[1]) return undefined;

  let payload: unknown;

  try {
    payload = decodeSegment(segments[1]);
  } catch {
    return undefined;
  }

  const parsed = CLAIMS_SCHEMA.safeParse(payload);

  if (!parsed.success) return undefined;

  return {
    userId: parsed.data.user_id,
    username: parsed.data.username,
    email: parsed.data.email,
    isAdmin: parsed.data.is_admin,
    reputation: parsed.data.reputation,
    expiresAtSeconds: parsed.data.exp,
  };
};

export const EXPIRY_SKEW_SECONDS = 60;

export const isAccessTokenExpired = (
  claims: AccessTokenClaims,
  now: Date = new Date(),
  skewSeconds: number = EXPIRY_SKEW_SECONDS,
): boolean => {
  if (claims.expiresAtSeconds === undefined) return true;

  return claims.expiresAtSeconds - skewSeconds <= Math.floor(now.getTime() / 1000);
};

export const toSessionUser = (claims: AccessTokenClaims): App.SessionUser => ({
  id: claims.userId,
  username: claims.username,
  email: claims.email,
  isAdmin: claims.isAdmin,
  reputation: claims.reputation,
});
