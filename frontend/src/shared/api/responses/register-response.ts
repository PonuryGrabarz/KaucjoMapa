import { JwtToken } from "../jwt-token";

export class RegisterResponse {
  createdAt!: string;
  email!: string;
  id!: string;
  isAdmin!: boolean;
  isBanned!: boolean;
  reputation!: number;
  token!: JwtToken;
  updatedAt!: string;
  username!: string;
}
