type JwtObject = Readonly<Record<string, unknown>>;

type JwtDateClaim = "exp" | "iat" | "nbf";

const decodeJwtObject = (segment: string, partName: string): JwtObject => {
  try {
    const base64 = segment.replaceAll("-", "+").replaceAll("_", "/");

    const paddedBase64 = base64.padEnd(
      base64.length + ((4 - (base64.length % 4)) % 4),
      "=",
    );

    const bytes = Uint8Array.from(atob(paddedBase64), (character) =>
      character.charCodeAt(0),
    );

    const parsedValue: unknown = JSON.parse(new TextDecoder().decode(bytes));

    if (
      typeof parsedValue !== "object" ||
      parsedValue === null ||
      Array.isArray(parsedValue)
    ) {
      throw new TypeError(`${partName} tokena JWT nie jest obiektem JSON.`);
    }

    return Object.freeze(parsedValue as Record<string, unknown>);
  } catch (error) {
    if (
      error instanceof TypeError &&
      error.message.endsWith("nie jest obiektem JSON.")
    ) {
      throw error;
    }

    throw new Error(`Nie udało się odczytać części „${partName}” tokena JWT.`);
  }
};
export class JwtToken {
  public readonly rawValue: string;
  public readonly header: JwtObject;
  public readonly claims: JwtObject;
  public readonly signature: string;

  public constructor(rawValue: string) {
    const segments = rawValue.split(".");

    if (
      segments.length !== 3 ||
      segments.some((segment) => segment.length === 0)
    ) {
      throw new Error("Token JWT musi składać się z trzech niepustych części.");
    }

    this.rawValue = rawValue;
    this.header = decodeJwtObject(segments[0], "nagłówek");
    this.claims = decodeJwtObject(segments[1], "payload");
    this.signature = segments[2];
  }

  public get expiresAt(): Date | undefined {
    return this.getDateClaim("exp");
  }

  public get issuedAt(): Date | undefined {
    return this.getDateClaim("iat");
  }

  public get notBefore(): Date | undefined {
    return this.getDateClaim("nbf");
  }

  public getClaim<TClaim = unknown>(name: string): TClaim | undefined {
    return this.claims[name] as TClaim | undefined;
  }

  public hasClaim(name: string): boolean {
    return Object.hasOwn(this.claims, name);
  }

  public isExpired(at: Date = new Date()): boolean {
    const expiresAt = this.expiresAt;

    return expiresAt !== undefined && expiresAt.getTime() <= at.getTime();
  }

  public toString(): string {
    return this.rawValue;
  }

  private getDateClaim(name: JwtDateClaim): Date | undefined {
    const value = this.getClaim(name);

    if (typeof value !== "number" || !Number.isFinite(value)) {
      return undefined;
    }

    return new Date(value * 1_000);
  }
}
