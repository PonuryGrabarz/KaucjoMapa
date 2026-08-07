import { z } from "astro/zod";
import { isApiError } from "@shared/api/errors";
import { ActionError, defineAction } from "astro:actions";
import { login, logout, register } from "@shared/api/auth";
import { toActionError } from "@shared/helpers/action-errors";

import type { Result } from "@shared/api/result";
import type { ActionAPIContext, ActionErrorCode } from "astro:actions";

const LOGIN_INPUT_SCHEMA = z.object({
	username: z
		.string({ error: "Podaj nazwę użytkownika." })
		.trim()
		.min(1, "Podaj nazwę użytkownika."),
	password: z.string({ error: "Podaj hasło." }).min(1, "Podaj hasło."),
});

const REGISTER_INPUT_SCHEMA = z
	.object({
		username: z
			.string({ error: "Podaj nazwę użytkownika." })
			.trim()
			.min(1, "Podaj nazwę użytkownika."),
		email: z.email({ error: "Podaj poprawny adres e-mail." }).trim(),
		password: z
			.string({ error: "Podaj hasło." })
			.min(8, "Hasło musi mieć co najmniej 8 znaków."),
		confirmPassword: z.string({ error: "Powtórz hasło." }).min(1, "Powtórz hasło."),
	})
	.refine((input) => input.password === input.confirmPassword, {
		error: "Hasła muszą być takie same.",
		path: ["confirmPassword"],
	});

const toSessionContext = (context: ActionAPIContext) => ({
	cookies: context.cookies,
	url: context.url,
});

interface StatusMessage {
	status: number;
	code: ActionErrorCode;
	message: string;
}

const unwrapSession = (
	result: Result<App.Session>,
	overrides: readonly StatusMessage[],
): App.Session => {
	if (!result.isError) return result.value;

	const [error] = result.errors;

	if (isApiError(error) && error.kind === "http") {
		const override = overrides.find((candidate) => candidate.status === error.status);

		if (override) throw new ActionError({ code: override.code, message: override.message });
	}

	throw toActionError(result.errors);
};

export const auth = {
	login: defineAction({
		accept: "form",
		input: LOGIN_INPUT_SCHEMA,
		handler: async (input, context) => {
			const result = await login(toSessionContext(context), input);

			const session = unwrapSession(result, [
				{
					status: 401,
					code: "UNAUTHORIZED",
					message: "Nieprawidłowa nazwa użytkownika lub hasło.",
				},
			]);

			context.locals.session = session;

			return { user: session.user };
		},
	}),

	register: defineAction({
		accept: "form",
		input: REGISTER_INPUT_SCHEMA,
		handler: async ({ confirmPassword: _confirmPassword, ...input }, context) => {
			const result = await register(toSessionContext(context), input);

			const session = unwrapSession(result, [
				{
					status: 409,
					code: "CONFLICT",
					message: "Ta nazwa użytkownika lub adres e-mail jest już zajęty.",
				},
				{
					status: 400,
					code: "BAD_REQUEST",
					message: "Nie udało się utworzyć konta. Sprawdź wprowadzone dane.",
				},
			]);

			context.locals.session = session;

			return { user: session.user };
		},
	}),

	logout: defineAction({
		handler: async (_input, context) => {
			const result = await logout(toSessionContext(context));

			context.locals.session = null;

			if (result.isError) throw toActionError(result.errors);

			return { loggedOut: true };
		},
	}),
};
