import { z } from "astro/zod";
import { ActionError, defineAction } from "astro:actions";

const sessionUserSchema = z.object({
  createdAt: z.string().min(1),
  email: z.email(),
  id: z.string().uuid(),
  isAdmin: z.boolean(),
  isBanned: z.boolean(),
  reputation: z.number().int(),
  token: z.string().min(1),
  updatedAt: z.string().min(1),
  username: z.string().min(1),
});

export const auth = {
  createSession: defineAction({
    input: sessionUserSchema,
    handler: (input, context) => {
      const { session } = context;
      if (!session) {
        throw new ActionError({
          code: "INTERNAL_SERVER_ERROR",
          message: "Obsługa sesji jest niedostępna.",
        });
      }

      const { token, ...userData } = input;

      const user: App.SessionUser = {
        ...userData,
        authenticatedAt: new Date(),
      };

      session.regenerate();

      session.set("user", user);
      session.set("token", token);

      return user;
    },
  }),

  destroySession: defineAction({
    handler: (_input, context) => {
      context.session?.destroy();
      return { success: true };
    },
  }),

  getSession: defineAction({
    handler: async (_input, context) => {
      return (await context.session?.get("user")) ?? null;
    },
  }),
};
