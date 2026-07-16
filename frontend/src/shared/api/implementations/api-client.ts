import ky from 'ky';
import ApiClientBase from '../abstractions/api-client-base';

import { JwtToken } from '../jwt-token';

import type { GetMapPointsRequest, LoginRequest, RegisterRequest } from '../requests';
import type { GetMapPointsResponse, LoginResponse, RegisterResponse } from '../responses';

type MapPointDto = {
  return_point_id: string;
  short_name: string;
  bg_color: string;
  text_color: string;
  latitude: number;
  longitude: number;
};

type AuthResponseDto = {
  created_at: string;
  email: string;
  id: string;
  is_admin: boolean;
  is_banned: boolean;
  reputation: number;
  token: string;
  updated_at: string;
  username: string;
};

const mapAuthResponse = <TResponse extends LoginResponse | RegisterResponse>(
  response: AuthResponseDto,
): TResponse =>
  ({
    createdAt: response.created_at,
    email: response.email,
    id: response.id,
    isAdmin: response.is_admin,
    isBanned: response.is_banned,
    reputation: response.reputation,
    token: new JwtToken(response.token),
    updatedAt: response.updated_at,
    username: response.username,
  }) as TResponse;

export class ApiClient extends ApiClientBase {
  protected async getMapPoints(_request: GetMapPointsRequest): Promise<GetMapPointsResponse> {
    const points = await ky.get('kaucjomaty', { prefix: this.getApiUrl() }).json<MapPointDto[]>();

    return {
      points: points.map((point) => ({
        returnPointId: point.return_point_id,
        shortName: point.short_name,
        bgColor: point.bg_color,
        textColor: point.text_color,
        latitude: point.latitude,
        longitude: point.longitude,
      })),
    };
  }

  protected async login(request: LoginRequest): Promise<LoginResponse> {
    const response = await ky.post('login', {
      prefix: this.getApiUrl(),
      credentials: 'include',
      throwHttpErrors: false,
      json: {
        username: request.username,
        password: request.password,
      },
    });

    if (!response.ok) {
      throw new Error('Logowanie nie powiodło się. Spróbuj ponownie później.');
    }

    return mapAuthResponse<LoginResponse>(await response.json<AuthResponseDto>());
  }

  protected async register(request: RegisterRequest): Promise<RegisterResponse> {
    const response = await ky.post('users', {
      prefix: this.getApiUrl(),
      credentials: 'include',
      throwHttpErrors: false,
      json: {
        username: request.username,
        email: request.email,
        password: request.password,
      },
    });

    if (!response.ok) {
      throw new Error('Rejestracja nie powiodła się. Spróbuj ponownie później.');
    }

    return mapAuthResponse<RegisterResponse>(await response.json<AuthResponseDto>());
  }
}
