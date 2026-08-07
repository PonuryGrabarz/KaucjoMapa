export class Bookmark {
  id!: string;
  user_id!: string;
  return_point_id!: string;
  created_at!: string;
}

export type CreateBookmarkResponse = Bookmark;

export type GetBookmarkByReturnPointResponse = Bookmark;

export type DeleteBookmarkByReturnPointResponse = undefined;
