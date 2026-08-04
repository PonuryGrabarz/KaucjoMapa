import { z } from "zod";

import { Result } from "../result";

import {
  RETURN_POINT_REFUND_METHODS,
  RETURN_POINT_STATUSES,
  RETURN_POINT_TYPES,
} from "../responses";

import type {
  ReturnPointRefundMethod,
  ReturnPointStatus,
  ReturnPointType,
} from "../responses";

export interface GetReturnPointsParams {
  type?: ReturnPointType;
  status?: ReturnPointStatus;
  query?: string;
  min_distance?: number;
  max_distance?: number;
  lat?: number;
  lon?: number;
  pet?: boolean;
  can?: boolean;
  glass?: boolean;
  refund_method?: readonly ReturnPointRefundMethod[];
  page?: number;
  page_size?: number;
}

const FIELD_LABELS: Record<string, string> = {
  type: "Typ punktu",
  status: "Status punktu",
  query: "Fraza wyszukiwania",
  min_distance: "Minimalna odległość",
  max_distance: "Maksymalna odległość",
  lat: "Szerokość geograficzna",
  lon: "Długość geograficzna",
  pet: "Opakowania PET",
  can: "Puszki",
  glass: "Szkło",
  refund_method: "Metody rozliczenia",
  page: "Numer strony",
  page_size: "Liczba elementów na stronie",
};

const listOptions = (values: readonly string[]): string => values.join(", ");

const GET_RETURN_POINTS_SCHEMA = z
  .strictObject({
    type: z
      .enum(RETURN_POINT_TYPES, {
        error: `Typ punktu musi być jedną z wartości: ${listOptions(RETURN_POINT_TYPES)}.`,
      })
      .optional(),

    status: z
      .enum(RETURN_POINT_STATUSES, {
        error: `Status punktu musi być jedną z wartości: ${listOptions(RETURN_POINT_STATUSES)}.`,
      })
      .optional(),

    query: z
      .string({ error: "Fraza wyszukiwania musi być tekstem." })
      .trim()
      .min(1, { error: "Fraza wyszukiwania nie może być pusta." })
      .max(200, { error: "Fraza wyszukiwania może mieć najwyżej 200 znaków." })
      .optional(),

    min_distance: z
      .number({ error: "Minimalna odległość musi być liczbą." })
      .min(0, { error: "Minimalna odległość nie może być ujemna." })
      .optional(),

    max_distance: z
      .number({ error: "Maksymalna odległość musi być liczbą." })
      .min(0, { error: "Maksymalna odległość nie może być ujemna." })
      .optional(),

    lat: z
      .number({ error: "Szerokość geograficzna musi być liczbą." })
      .min(-90, { error: "Szerokość geograficzna nie może być mniejsza niż -90." })
      .max(90, { error: "Szerokość geograficzna nie może być większa niż 90." })
      .optional(),

    lon: z
      .number({ error: "Długość geograficzna musi być liczbą." })
      .min(-180, { error: "Długość geograficzna nie może być mniejsza niż -180." })
      .max(180, { error: "Długość geograficzna nie może być większa niż 180." })
      .optional(),

    pet: z.boolean({ error: "Opakowania PET: oczekiwano wartości logicznej." }).optional(),
    can: z.boolean({ error: "Puszki: oczekiwano wartości logicznej." }).optional(),
    glass: z.boolean({ error: "Szkło: oczekiwano wartości logicznej." }).optional(),

    refund_method: z
      .array(
        z.enum(RETURN_POINT_REFUND_METHODS, {
          error: `Metoda rozliczenia musi być jedną z wartości: ${listOptions(RETURN_POINT_REFUND_METHODS)}.`,
        }),
        { error: "Metody rozliczenia muszą być listą wartości." },
      )
      .min(1, { error: "Lista metod rozliczenia nie może być pusta." })
      .refine((methods) => new Set(methods).size === methods.length, {
        error: "Metody rozliczenia nie mogą się powtarzać.",
      })
      .optional(),

    page: z
      .int({ error: "Numer strony musi być liczbą całkowitą." })
      .min(1, { error: "Numer strony musi wynosić co najmniej 1." })
      .optional(),

    page_size: z
      .int({ error: "Liczba elementów na stronie musi być liczbą całkowitą." })
      .min(1, { error: "Liczba elementów na stronie musi wynosić co najmniej 1." })
      .optional(),
  }, { error: "Żądanie zawiera nieznany parametr." })
  .superRefine((params, context) => {
    if (
      params.min_distance !== undefined
      && params.max_distance !== undefined
      && params.min_distance > params.max_distance
    ) {
      context.addIssue({
        code: "custom",
        path: ["min_distance"],
        message: "Minimalna odległość nie może być większa niż maksymalna.",
      });
    }

    if ((params.lat === undefined) !== (params.lon === undefined)) {
      context.addIssue({
        code: "custom",
        path: [params.lat === undefined ? "lat" : "lon"],
        message: "Szerokość i długość geograficzną trzeba podać razem.",
      });
    }

    const filtersByDistance = params.min_distance !== undefined || params.max_distance !== undefined;

    if (filtersByDistance && (params.lat === undefined || params.lon === undefined)) {
      context.addIssue({
        code: "custom",
        path: ["lat"],
        message: "Filtrowanie po odległości wymaga podania pozycji użytkownika (lat i lon).",
      });
    }
  });

type ValidatedParams = z.infer<typeof GET_RETURN_POINTS_SCHEMA>;

const toErrors = (error: z.ZodError): Error[] =>
  error.issues.map((issue) => {
    const field = issue.path[0];

    const label = typeof field === "string" ? FIELD_LABELS[field] : undefined;

    return new Error(
      label && !issue.message.startsWith(label) ? `${label}: ${issue.message}` : issue.message,
    );
  });

export class GetReturnPointsRequest {
  public readonly kind = "get-return-points" as const;

  private readonly params: Readonly<ValidatedParams>;

  private constructor(params: ValidatedParams) {
    this.params = Object.freeze({ ...params });
  }

  public static create(params: GetReturnPointsParams = {}): Result<GetReturnPointsRequest> {
    const parsed = GET_RETURN_POINTS_SCHEMA.safeParse(params);

    if (!parsed.success) return Result.failure<GetReturnPointsRequest>(toErrors(parsed.error));

    return Result.success(new GetReturnPointsRequest(parsed.data));
  }

  public get type(): ReturnPointType | undefined {
    return this.params.type;
  }

  public get status(): ReturnPointStatus | undefined {
    return this.params.status;
  }

  public get query(): string | undefined {
    return this.params.query;
  }

  public get minDistance(): number | undefined {
    return this.params.min_distance;
  }

  public get maxDistance(): number | undefined {
    return this.params.max_distance;
  }

  public get lat(): number | undefined {
    return this.params.lat;
  }

  public get lon(): number | undefined {
    return this.params.lon;
  }

  public get pet(): boolean | undefined {
    return this.params.pet;
  }

  public get can(): boolean | undefined {
    return this.params.can;
  }

  public get glass(): boolean | undefined {
    return this.params.glass;
  }

  public get refundMethods(): readonly ReturnPointRefundMethod[] | undefined {
    return this.params.refund_method;
  }

  public get page(): number | undefined {
    return this.params.page;
  }

  public get pageSize(): number | undefined {
    return this.params.page_size;
  }

  public toSearchParams(): URLSearchParams {
    const searchParams = new URLSearchParams();

    for (const [key, value] of Object.entries(this.params)) {
      if (value === undefined) continue;

      searchParams.set(key, Array.isArray(value) ? value.join(",") : String(value));
    }

    return searchParams;
  }

  public toQueryString(): string {
    return this.toSearchParams().toString();
  }
}
