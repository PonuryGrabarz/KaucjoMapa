import ApiClientBase from "../abstractions/api-client-base";

import { JwtToken } from "../jwt-token";

import type { GetMapPointsRequest, LoginRequest, RegisterRequest } from "../requests";
import type { LoginResponse, RegisterResponse } from "../responses";
import type { GetMapPointsResponse, MapPoint } from "../responses/get-map-points-response";

const WAIT_TIME = 2500;
const DUMMY_USER_ID = "00000000-0000-4000-8000-000000000001";
const DUMMY_TIMESTAMP = "2026-07-16 09:15:38.93382 +0000 UTC";

const encodeJwtPart = (value: Record<string, unknown>): string => {
  const bytes = new TextEncoder().encode(JSON.stringify(value));
  const binary = Array.from(bytes, (byte) => String.fromCharCode(byte)).join("");

  return btoa(binary)
    .replaceAll("+", "-")
    .replaceAll("/", "_")
    .replace(/=+$/u, "");
};

const createDummyToken = (username: string, email: string): JwtToken => {
  const issuedAt = Math.floor(Date.now() / 1_000);
  const expiresAt = issuedAt + 3_600;
  const header = encodeJwtPart({ alg: "HS256", typ: "JWT" });
  const claims = encodeJwtPart({
    iss: "KaucjoMapa",
    sub: DUMMY_USER_ID,
    username,
    email,
    iat: issuedAt,
    exp: expiresAt,
  });

  return new JwtToken(`${header}.${claims}.dummy-signature`);
};

const createDummyAuthResponse = <TResponse extends LoginResponse | RegisterResponse>(
  username: string,
  email: string,
): TResponse =>
  ({
    createdAt: DUMMY_TIMESTAMP,
    email,
    id: DUMMY_USER_ID,
    isAdmin: false,
    isBanned: false,
    reputation: 1,
    token: createDummyToken(username, email),
    updatedAt: DUMMY_TIMESTAMP,
    username,
  }) as TResponse;

export class DummyApiClient extends ApiClientBase {
  constructor() {
    super("")
  }

  protected getMapPoints(request: GetMapPointsRequest): Promise<GetMapPointsResponse> {
    const dummyMapPoint1: MapPoint = {
      returnPointId: "1",
      shortName: "GS",
      bgColor: "#193f2a",
      textColor: "#ffffff",
      latitude: 49.6135,
      longitude: 20.7088
    };

    const dummyMapPoint2: MapPoint = {
      returnPointId: "2",
      shortName: "ZK",
      bgColor: "#2f8f46",
      textColor: "#ffffff",
      latitude: 49.6258,
      longitude: 20.6912
    };

    const dummyMapPoint3: MapPoint = {
      returnPointId: "3",
      shortName: "PKP",
      bgColor: "#174ea6",
      textColor: "#ffffff",
      latitude: 49.6189,
      longitude: 20.6878
    };

    const dummyMapPoint4: MapPoint = {
      returnPointId: "4",
      shortName: "S",
      bgColor: "#f0b429",
      textColor: "#3d2b03",
      latitude: 49.6042,
      longitude: 20.7185
    };

    const dummyMapPoint5: MapPoint = {
      returnPointId: "5",
      shortName: "RH",
      bgColor: "#0e2a1b",
      textColor: "#5fd37e",
      latitude: 49.6248,
      longitude: 20.6978
    };

    const dummyMapPoint6: MapPoint = {
      returnPointId: "6",
      shortName: "M",
      bgColor: "#6f4be8 ",
      textColor: "#ffffff",
      latitude: 49.6321,
      longitude: 20.7042
    };

    const dummyMapPoint7: MapPoint = {
      returnPointId: "7",
      shortName: "B",
      bgColor: "#e31b23",
      textColor: "#ffd400",
      latitude: 49.6398,
      longitude: 20.6855
    };

    const dummyResponse: GetMapPointsResponse = {
      points: [dummyMapPoint1, dummyMapPoint2, dummyMapPoint3, dummyMapPoint4, dummyMapPoint5, dummyMapPoint6, dummyMapPoint7]
    };

    return new Promise((resolve) => setTimeout(() => {
      resolve(dummyResponse);
    }, WAIT_TIME))
  }

  protected login(request: LoginRequest): Promise<LoginResponse> {
    const response = createDummyAuthResponse<LoginResponse>(
      request.username,
      `${request.username}@example.local`,
    );

    return new Promise((resolve) => {
      setTimeout(() => resolve(response), WAIT_TIME);
    });
  }

  protected register(request: RegisterRequest): Promise<RegisterResponse> {
    const response = createDummyAuthResponse<RegisterResponse>(
      request.username,
      request.email,
    );

    return new Promise((resolve) => {
      setTimeout(() => resolve(response), WAIT_TIME);
    });
  }
}
