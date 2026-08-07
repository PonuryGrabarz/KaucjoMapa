import { z } from "zod";

import { Result } from "../result";

import {
	RETURN_POINT_TYPES,
	RETURN_POINT_STATUSES,
	RETURN_POINT_REFUND_METHODS,
} from "../responses";

import type { ReturnPointType, ReturnPointStatus, ReturnPointRefundMethod } from "../responses";

export interface CreateReturnPointParams {
	name: string;
	type: ReturnPointType;
	address?: string;
	latitude: number;
	longitude: number;
	pet?: boolean;
	can?: boolean;
	glass?: boolean;
	refund_method: readonly ReturnPointRefundMethod[];
	opening_hours?: Record<string, string>;
	description?: string;
	current_status: ReturnPointStatus;
}

const FIELD_LABELS: Record<string, string> = {
	name: "Nazwa punktu",
	type: "Typ punktu",
	address: "Adres",
	latitude: "Szerokość geograficzna",
	longitude: "Długość geograficzna",
	refund_method: "Formy rozliczenia",
	opening_hours: "Godziny otwarcia",
	description: "Opis",
	current_status: "Status punktu",
};

const listOptions = (values: readonly string[]): string => values.join(", ");

const CREATE_RETURN_POINT_SCHEMA = z
	.strictObject(
		{
			name: z
				.string({ error: "Nazwa punktu musi być tekstem." })
				.trim()
				.min(1, { error: "Nazwa punktu nie może być pusta." })
				.max(255, { error: "Nazwa punktu może mieć najwyżej 255 znaków." }),

			type: z.enum(RETURN_POINT_TYPES, {
				error: `Typ punktu musi być jedną z wartości: ${listOptions(RETURN_POINT_TYPES)}.`,
			}),

			address: z.string({ error: "Adres musi być tekstem." }).trim().default(""),

			latitude: z
				.number({ error: "Szerokość geograficzna musi być liczbą." })
				.min(-90, { error: "Szerokość geograficzna nie może być mniejsza niż -90." })
				.max(90, { error: "Szerokość geograficzna nie może być większa niż 90." }),

			longitude: z
				.number({ error: "Długość geograficzna musi być liczbą." })
				.min(-180, { error: "Długość geograficzna nie może być mniejsza niż -180." })
				.max(180, { error: "Długość geograficzna nie może być większa niż 180." }),

			pet: z.boolean({ error: "Opakowania PET: oczekiwano wartości logicznej." }).default(false),
			can: z.boolean({ error: "Puszki: oczekiwano wartości logicznej." }).default(false),
			glass: z.boolean({ error: "Szkło: oczekiwano wartości logicznej." }).default(false),

			refund_method: z
				.array(
					z.enum(RETURN_POINT_REFUND_METHODS, {
						error: `Metoda rozliczenia musi być jedną z wartości: ${listOptions(RETURN_POINT_REFUND_METHODS)}.`,
					}),
					{ error: "Metody rozliczenia muszą być listą wartości." },
				)
				.min(1, { error: "Wybierz co najmniej jedną formę rozliczenia." })
				.refine((methods) => new Set(methods).size === methods.length, {
					error: "Metody rozliczenia nie mogą się powtarzać.",
				}),

			opening_hours: z.record(z.string(), z.string()).default({}),

			description: z.string({ error: "Opis musi być tekstem." }).trim().default(""),

			current_status: z.enum(RETURN_POINT_STATUSES, {
				error: `Status punktu musi być jedną z wartości: ${listOptions(RETURN_POINT_STATUSES)}.`,
			}),
		},
		{ error: "Żądanie zawiera nieznany parametr." },
	)
	.superRefine((params, context) => {
		if (!params.pet && !params.can && !params.glass) {
			context.addIssue({
				code: "custom",
				path: ["pet"],
				message: "Wybierz co najmniej jeden rodzaj przyjmowanych opakowań.",
			});
		}
	});

type ValidatedParams = z.infer<typeof CREATE_RETURN_POINT_SCHEMA>;

const toErrors = (error: z.ZodError): Error[] =>
	error.issues.map((issue) => {
		const field = issue.path[0];

		const label = typeof field === "string" ? FIELD_LABELS[field] : undefined;

		return new Error(
			label && !issue.message.startsWith(label) ? `${label}: ${issue.message}` : issue.message,
		);
	});

export class CreateReturnPointRequest {
	public readonly kind = "create-return-point" as const;

	private readonly params: Readonly<ValidatedParams>;

	private constructor(params: ValidatedParams) {
		this.params = Object.freeze({ ...params });
	}

	public static create(params: CreateReturnPointParams): Result<CreateReturnPointRequest> {
		const parsed = CREATE_RETURN_POINT_SCHEMA.safeParse(params);

		if (!parsed.success) return Result.failure<CreateReturnPointRequest>(toErrors(parsed.error));

		return Result.success(new CreateReturnPointRequest(parsed.data));
	}

	public toBody(): Record<string, unknown> {
		return { ...this.params };
	}
}
