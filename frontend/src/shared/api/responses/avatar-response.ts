export class AvatarFile {
  data!: ArrayBuffer;
  contentType!: string;
}

export type GetAvatarResponse = AvatarFile;
